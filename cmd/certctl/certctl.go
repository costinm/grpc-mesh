package main

import (
	"context"
	"crypto/tls"
	"log"

	// Import the GCS API definitions and generate using the template grpc/go.
	// buf.build google doesn't include google.security.meshca.v1
	//
	storagev1 "go.buf.build/grpc/go/googleapis/googleapis/google/storage/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Command line tool to generate and refresh mesh certificates, matching GCP file paths.
//
// On GKE, this is done automatically if the Pod is annotated with ...
//
// For go programs, it is recommended to directly call this from uca package.
//
// For other languages, until an native library is available, exec this program periodically.
// Will create /var/run/secrets/....
//
func main() {
	cc, err := grpc.Dial(
		"storage.googleapis.com:443",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})),
	)
	if err != nil {
		log.Fatalf("Failed to dial GCS API: %v", err)
	}
	client := storagev1.NewStorageClient(cc)
	resp, err := client.GetBucket(context.Background(), &storagev1.GetBucketRequest{
		// Public GCS dataset
		// Ref: https://cloud.google.com/healthcare/docs/resources/public-datasets/nih-chest
		Bucket: "gcs-public-data--healthcare-nih-chest-xray",
	})
	if err != nil {
		log.Fatalf("Failed to get bucket: %v", err)
	}
	log.Println(resp)
}
