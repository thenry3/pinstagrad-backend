package gcloud

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"config"

	"cloud.google.com/go/storage"
)

// CloudClient creates GCloud client
func CloudClient() *storage.Client {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	return client
}

//NewBucket gets or creates new bucket with bucket name passed as a param
func NewBucket(client *storage.Client, bucketName string) *storage.BucketHandle {
	conf := config.New()

	ctx := context.Background()

	projectID := conf.CloudStorage.Project

	bucket := client.Bucket(bucketName)
	exists, err := bucket.Attrs(ctx)
	if err != nil {
		log.Printf("Bucket under creation: %s", bucketName)
	}

	if exists == nil {
		if err := bucket.Create(ctx, projectID, nil); err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
		fmt.Printf("Bucket %v created.\n", bucketName)
	}

	return bucket
}

// GetBucket gets or creates and gets Google Cloud bucket
func GetBucket(client *storage.Client) *storage.BucketHandle {
	conf := config.New()

	ctx := context.Background()

	projectID := conf.CloudStorage.Project

	bucketName := conf.CloudStorage.Bucket

	bucket := client.Bucket(bucketName)
	exists, err := bucket.Attrs(ctx)
	if err != nil {
		log.Printf("Bucket under creation: %s", bucketName)
	}

	if exists == nil {
		if err := bucket.Create(ctx, projectID, nil); err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}
		fmt.Printf("Bucket %v created.\n", bucketName)
	}

	return bucket
}

// Upload uploads image to GCloud
func Upload(client *storage.Client, signedURL string, bucketname string, image string, uuid string, contentType string) {
	ctx := context.Background()

	b, err := ioutil.ReadFile(image)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("PUT", signedURL, bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", contentType)
	RESTclient := new(http.Client)
	resp, err := RESTclient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	acl := client.Bucket(bucketname).Object(uuid).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		log.Print(err)
	}

	fmt.Println(resp)
}

// URI returns uri of bucket
func URI(attrs storage.ObjectAttrs) string {
	return attrs.Metadata["selfLink"]
}
