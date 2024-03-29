package controller

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"

	"config"

	GCloud "gcloud"

	SignedURL "gcloud/SignedURL"
)

// CreateFirebaseUser Creates Firebase user for DB interactions
func CreateFirebaseUser(conf *config.Config, contentType string, gCloudClient *storage.Client, client *auth.Client, email string, emailVerified bool, phoneNumber string, password string, name string, photoURL string, status bool) *auth.UserRecord {
	ctx := context.Background()

	signedURL, key, err := SignedURL.SignedURLoptions("config/pinstagrad-back-7.json", "PUT", conf.CloudSettings.ProfileBucket, contentType)
	if err != nil {
		log.Fatalf("Error creating auth key: %v", err)
	}

	GCloud.Upload(gCloudClient, signedURL, conf.CloudSettings.ProfileBucket, photoURL, key, contentType)

	// REGEX match ^\+[1-9]\d{1,14}$ for valid phonenumber
	params := (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(emailVerified).
		PhoneNumber(phoneNumber).
		Password(password).
		DisplayName(name).
		PhotoURL(conf.CloudSettings.URI + key).
		Disabled(status)
	usr, err := client.CreateUser(ctx, params)
	if err != nil {
		log.Fatalf("error creating user: %v\n", err)
	}
	log.Printf("Successfully created user: %v\n", usr)
	return usr
}

// FirebaseCustomToken Intializes Firebase custom token for auth
func FirebaseCustomToken(app *firebase.App) string {
	ctx := context.Background()

	client, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.CustomToken(ctx, "some-uid")
	if err != nil {
		log.Fatalf("error minting custom token: %v\n", err)
	}

	log.Printf("Got custom token: %v\n", token)

	return token
}

// VerifyFirebaseToken Verifies Firebase custom token for auth
func VerifyFirebaseToken(app *firebase.App, idToken string) *auth.Token {
	ctx := context.Background()

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Fatalf("error verifying ID token: %v\n", err)
	}

	log.Printf("Verified ID token: %v\n", token)

	return token
}
