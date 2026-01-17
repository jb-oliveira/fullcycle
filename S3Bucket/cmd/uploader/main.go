package main

import (
	"context"
	"fmt"
	"log"
	"os"

	// Core AWS types, interfaces, config
	"github.com/aws/aws-sdk-go-v2/config"     // Configuration loading
	"github.com/aws/aws-sdk-go-v2/service/s3" // S3 Service Client
	// S3 Specific types (like BucketLocationConstraint)
)

func main() {

	// Load the AWS configuration (credentials, region, etc.)
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	// Create an S3 client
	client := s3.NewFromConfig(cfg)

	// Example: List Buckets
	listBucketsInput := &s3.ListBucketsInput{}
	buckets, err := client.ListBuckets(context.Background(), listBucketsInput)
	if err != nil {
		log.Fatalf("failed to list buckets, %v", err)
	}
	fmt.Println("Buckets:")
	for _, bucket := range buckets.Buckets {
		fmt.Printf(" - %s\n", *bucket.Name)
	}

	// Example: Using the Manager for Upload
	// uploader := manager.NewUploader(client) // Requires client.ConfigProvider, which cfg satisfies
	// _, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{ ... })

	dir, err := os.ReadDir("./tmp")
	if err != nil {
		log.Fatalf("failed to read directory, %v", err)
	}
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		fmt.Println(entry.Name())
		// if err := uploadFile(client, "first-bucket-s3", entry.Name()); err != nil {
		// 	log.Fatalf("failed to upload file, %v", err)
		// }
	}

}

func uploadFile(client *s3.Client, bucketName string, filePath string) error {
	file, err := os.Open("./tmp/" + filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer file.Close()

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &filePath,
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}
