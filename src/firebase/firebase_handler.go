package controller

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

// CreateFirebaseUser Creates Firebase user for DB interactions
func CreateFirebaseUser(ctx context.Context, client *auth.Client, email string, emailVerified bool, phoneNumber string, password string, name string, photoURL string, status bool) *auth.UserRecord {
	params := (&auth.UserToCreate{}).
		Email(email).
		EmailVerified(emailVerified).
		PhoneNumber(phoneNumber).
		Password(password).
		DisplayName(name).
		PhotoURL(photoURL).
		Disabled(status)
	usr, err := client.CreateUser(ctx, params)
	if err != nil {
		log.Fatalf("error creating user: %v\n", err)
	}
	log.Printf("Successfully created user: %v\n", usr)
	return usr
}

// FirebaseCustomToken Intializes Firebase custom token for auth
func FirebaseCustomToken(ctx context.Context, app *firebase.App) string {
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
func VerifyFirebaseToken(ctx context.Context, app *firebase.App, idToken string) *auth.Token {
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
