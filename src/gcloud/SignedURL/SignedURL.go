package signedurl

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"cloud.google.com/go/storage"
	uuid "github.com/google/uuid"
	"golang.org/x/oauth2/google"
)

// SignedURLoptions Returns signed url for users to interact with Gcloud without having credentials
func SignedURLoptions(serviceAccount string, option string, bucket string) (string, string, error) {
	jsonKey, err := ioutil.ReadFile(serviceAccount)
	if err != nil {
		return "", "", fmt.Errorf("cannot read the JSON key file, err: %v", err)
	}
	conf, err := google.JWTConfigFromJSON(jsonKey)
	if err != nil {
		return "", "", fmt.Errorf("google.JWTConfigFromJSON: %v", err)
	}
	key := uuid.New().String()

	expires := time.Now().Add(time.Minute * 50)
	r, err := storage.SignedURL(bucket, key, &storage.SignedURLOptions{
		GoogleAccessID: conf.Email,
		PrivateKey:     conf.PrivateKey,
		Method:         option,
		Expires:        expires,
		ContentType:    "image/png",
	})

	log.Print(r)
	return r, key, err
}
