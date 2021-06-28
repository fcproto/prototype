package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/fcproto/prototype/pkg/api"
	"github.com/fcproto/prototype/pkg/logger"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

var firestoreClient *firestore.Client
var log *logrus.Logger

type clientInfo struct {
	ClientID   string
	LastUpdate time.Time
	UpdateSize int
	LastSpeed  float64
}

var clientStatusLock sync.Mutex
var clientStatus = make([]*clientInfo, 0)

func updateClient(id string, updateSize int, lastSpeed float64) {
	clientStatusLock.Lock()
	defer clientStatusLock.Unlock()

	var info *clientInfo
	for _, cInfo := range clientStatus {
		if cInfo.ClientID == id {
			info = cInfo
			break
		}
	}

	if info == nil {
		info = &clientInfo{
			ClientID: id,
		}
		clientStatus = append(clientStatus, info)
	}
	info.LastUpdate = time.Now()
	info.UpdateSize = updateSize
	info.LastSpeed = lastSpeed
}

func createClient() *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	var projectId string

	// Sets thj Google Cloud Platform project ID.
	if os.Getenv("APP_USER") == "air" {
		// Server is running locally, so we look for the project id in the credential file
		jsonFile, err := os.Open("fcproto-credentials.json")
		if err != nil {
			log.Fatal(err)
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		defer jsonFile.Close()

		type ProjectID struct {
			ProjectID string `json:"project_id"`
		}
		var projectIdStruct ProjectID
		err = json.Unmarshal(byteValue, &projectIdStruct)
		if err != nil {
			log.Fatalf("The fcproto-credentials.json file seems to be malformed: %s", err)
		}
		projectId = projectIdStruct.ProjectID
	} else if len(os.Getenv("K_SERVICE")) > 0 {
		// Server runs on Cloud Run so we get the project id from the google metadata endpoint
		url := "http://metadata.google.internal/computeMetadata/v1/project/project-id"
		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Metadata-Flavor", "Google")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		projectId = string(body)
	} else {
		log.Fatal("Cannot determine if the server runs locally or on Cloud Run. " +
			"If you are trying to run the binary locally please use air " +
			"as stated in the documentation")
	}
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectId)
	if err != nil {
		// panic if client cannot be created
		log.Fatal(err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}

func main() {
	log = logger.New()
	log.Info("starting server...")

	firestoreClient = createClient()
	defer firestoreClient.Close()

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/status", Status)
	router.POST("/", StoreData)
	router.GET("/near/:client-id", GetNearCars)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Start HTTP server.
	go func() {
		log.Printf("listening on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	log.Println("stopping server...")

	if err := srv.Close(); err != nil {
		log.Error(err)
	}
}

func Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	status := table.NewWriter()
	status.SetStyle(table.StyleLight)
	status.AppendHeader(table.Row{"#", "Client ID", "Last Update", "Update Size", "Average Speed"})
	status.AppendSeparator()

	now := time.Now()
	clientStatusLock.Lock()
	defer clientStatusLock.Unlock()
	for i, info := range clientStatus {
		lastUpdate := fmt.Sprintf("%.0fs ago", now.Sub(info.LastUpdate).Seconds())
		speed := fmt.Sprintf("%.2fm/s", info.LastSpeed)
		status.AppendRow(table.Row{1 + i, info.ClientID[:8], lastUpdate, info.UpdateSize, speed})
		status.AppendSeparator()
	}

	w.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprintf(w, "\n\n\n%s\n", status.Render())
	if err != nil {
		log.Error(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	iter := firestoreClient.Collection("sensor-data").Documents(r.Context())
	data := make([]*api.SensorData, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf("An error has occurred: %s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		var el api.SensorData
		err = doc.DataTo(&el)
		if err != nil {
			log.Errorf("An error has occurred: %s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		data = append(data, &el)
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Errorf("An error has occurred: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	log.Infof("Sent %d documents", len(data))
}

func StoreData(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	data := make([]*api.SensorData, 0)
	batch := firestoreClient.Batch()

	// Populate the user data
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		log.Errorf("An error has occurred: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}

	clientId := ""
	lastSpeed := 0.0
	for _, el := range data {
		ref := firestoreClient.Collection("sensor-data").NewDoc()
		clientId = el.ClientID
		lastSpeed += el.Sensors["gps"]["speed"]
		batch.Set(ref, el)
	}

	_, err = batch.Commit(r.Context())

	if err != nil {
		log.Errorf("An error has occurred: %s", err)
		http.Error(w, err.Error(), 500)
	} else {
		// Write content-type, statuscode, payload
		size := len(data)
		updateClient(clientId, size, lastSpeed/float64(size))
		log.Infof("Stored %d documents for client %s", size, clientId[:8])
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		// json.NewEncoder(w).Encode(data)
	}
}

func GetNearCars(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	clientID := params.ByName("client-id")

	iter := firestoreClient.Collection("sensor-data").
		OrderBy("Timestamp", firestore.Desc).
		Documents(r.Context())

	data := make([]*api.SensorData, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Errorf("An error has occurred: %s", err)
			http.Error(w, err.Error(), 500)
			return
		}

		d := api.SensorData{}
		err = doc.DataTo(&d)
		if err != nil {
			log.Errorf("An error has occurred: %s", err)
			http.Error(w, err.Error(), 500)
			return
		}
		data = append(data, &d)
	}
	carIds := map[string]bool{clientID: true}
	nearCars := make([]*api.SensorData, 0)
	for _, el := range data {
		if !carIds[el.ClientID] {
			carIds[el.ClientID] = true
			nearCars = append(nearCars, el)
		}
		if len(nearCars) > 2 {
			break
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	err := json.NewEncoder(w).Encode(nearCars)
	if err != nil {
		log.Errorf("An error has occurred: %s", err)
		http.Error(w, err.Error(), 500)
		return
	}
	log.Infof("Sent %d documents", len(nearCars))
}
