package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	// Core AWS types, interfaces, config
	"github.com/aws/aws-sdk-go-v2/config"     // Configuration loading
	"github.com/aws/aws-sdk-go-v2/service/s3" // S3 Service Client
	// S3 Specific types (like BucketLocationConstraint)
)

var (
	dirName = "../../tmp"
	wg      sync.WaitGroup
)

func main() {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:", currDir)

	bucketName := os.Getenv("BUCKET_NAME")
	if bucketName == "" {
		log.Fatalf("BUCKET_NAME environment variable is not set")
	}
	// Load the AWS configuration (credentials, region, etc.)
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("sa-east-1"))
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
		if *bucket.Name == bucketName {
			fmt.Println("Bucket found")
		}
	}

	// This serves as a semaphore to control the number of threads
	uploadControl := make(chan int, 3)

	// This serves to retry sending files that failed
	// I found it a bit bad, because it can fall into an infinite loop
	retryUpload := make(chan string, 2)
	go func() {
		for fileName := range retryUpload {
			fmt.Println("Retrying upload for file:", fileName)
			// Increment WaitGroup and respect semaphore for retries
			wg.Add(1)
			uploadControl <- 1
			go uploadFile(client, bucketName, fileName, uploadControl, retryUpload)
		}
	}()

	dir, err := os.ReadDir(dirName)
	if err != nil {
		log.Fatalf("failed to read directory, %v", err)
	}
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}
		fmt.Println(entry.Name())
		// Sequential upload
		// if err := uploadFile(client, bucketName, entry.Name()); err != nil {
		// 	log.Fatalf("failed to upload file, %v", err)
		// }
		// Parallel upload
		wg.Add(1)
		// adds 1 to the channel
		uploadControl <- 1
		go uploadFile(client, bucketName, entry.Name(), uploadControl, retryUpload)
	}
	// close the channel
	close(uploadControl)

	// Wait for	all uploads to complete
	wg.Wait()
}

func uploadFile(client *s3.Client, bucketName string, filePath string, uploadControl <-chan int, retryUpload chan<- string) error {
	// necessary for parallelism
	defer wg.Done()
	defer func() { <-uploadControl }()

	file, err := os.Open(dirName + "/" + filePath)
	if err != nil {
		retryUpload <- filePath
		return fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer file.Close()

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &filePath,
		Body:   file,
	})
	if err != nil {
		retryUpload <- filePath
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}
