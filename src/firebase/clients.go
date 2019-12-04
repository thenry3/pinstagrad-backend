package controller

import (
	"context"
	"io/ioutil"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/zabawaba99/firego"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// AuthClient Auth client initialization
func AuthClient(app *firebase.App) *auth.Client {
	ctx := context.Background()
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
	return client
}

// FirebaseClient Intializes Firebase client for DB interactions
func FirebaseClient(credentialsFile string) (*firego.Firebase, error) {
	b, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(b, "https://www.googleapis.com/auth/firebase.database", "https://www.googleapis.com/auth/userinfo.email")
	if err != nil {
		return nil, err
	}

	client := conf.Client(oauth2.NoContext)
	dbinteract := firego.New("https://my-sample-app.firebaseio.com", client)
	return dbinteract, nil
}

// FirebaseSDK Intializes Firebase SDK
func FirebaseSDK() *firebase.App {
	config := &firebase.Config{
		DatabaseURL: "https://pinstagrad-7f307.firebaseio.com/",
	}
	opt := option.WithCredentialsFile("config/serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	return app
}
