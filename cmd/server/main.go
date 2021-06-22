// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"context"

	"cloud.google.com/go/firestore"
	"github.com/julienschmidt/httprouter"
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

	// Start HTTP server.
	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	iter := firestoreClient.Collection("sensor-data").Documents(context.Background())
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
