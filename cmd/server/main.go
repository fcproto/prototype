// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/fcproto/prototype/pkg/api"
	"github.com/julienschmidt/httprouter"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/iterator"
)

var firestoreClient *firestore.Client

func createClient() *firestore.Client {
	// Sets your Google Cloud Platform project ID.
	projectID := "fcproto"
	ctx := context.Background()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	// Close client when done with
	// defer client.Close()
	return client
}

func main() {
	log.Print("starting server...")

	firestoreClient = createClient()

	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/", StoreData)
	router.GET("/near/:client-id", GetNearCars)

	// Start HTTP server.
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	iter := firestoreClient.Collection("sensor-data").Documents(r.Context())
	type keyvalue map[string]interface{}
	data := make([]keyvalue, 0)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		data = append(data, doc.Data())
	}
	json.NewEncoder(w).Encode(data)
}

func StoreData(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	data := make([]*api.SensorData, 0)
	batch := firestoreClient.Batch()

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&data)

	for _, el := range data {
		ref := firestoreClient.Collection("sensor-data").NewDoc()
		batch.Set(ref, el)
	}

	_, err := batch.Commit(r.Context())

	if err != nil {
		log.Printf("An error has occurred: %s", err)
		http.Error(w, err.Error(), 500)
	} else {
		// Write content-type, statuscode, payload
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
			log.Printf("An error has occurred: %s", err)
			http.Error(w, err.Error(), 500)
		}

		d := api.SensorData{}
		err = mapstructure.Decode(doc.Data(), &d)
		if err != nil {
			log.Printf("An error has occurred: %s", err)
			http.Error(w, err.Error(), 500)
		}
		data = append(data, &d)
	}
	carIds := map[string]bool{clientID: true}
	nearCars := []*api.SensorData{}
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
	json.NewEncoder(w).Encode(nearCars)
}
