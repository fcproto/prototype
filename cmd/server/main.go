// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

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
	http.HandleFunc("/", handler)

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	client := createClient()

	ctx := context.Background()
	iter := client.Collection("sensor-data").Documents(ctx)
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
		fmt.Println(doc.Data())
		data = append(data, doc.Data())
	}
	jsonData, _ := json.Marshal(data)
	fmt.Fprint(w, string(jsonData))
}
