package database

import (
	"context"
	"fmt"
	"log"

	"firebase.google.com/go/db"

	firebase "firebase.google.com/go"
)

// UploadPhotoToRealtimeDatabase uploads images to realtime DB
func UploadPhotoToRealtimeDatabase(ctx context.Context, currentUser *User, photo *Photo, client *db.Client) {
	ref := client.NewRef("restricted_access/secret_document")
	var data map[string]interface{}
	if err := ref.Get(ctx, &data); err != nil {
		log.Fatalln("Error reading from database:", err)
	}
	fmt.Println(data)
}

// ConnectToRealtimeDatabase connect to Firebase realtime DB
func ConnectToRealtimeDatabase(ctx context.Context, app *firebase.App) *db.Client {
	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}
	return client
}
