package database

import (
	"context"
	"log"

	"firebase.google.com/go/db"

	firebase "firebase.google.com/go"
)

// RetrieveAllPhotos retrieves images in realtime DB
func RetrieveAllPhotos(ctx context.Context, reference *db.Ref) Photo {
	var photo Photo
	if err := reference.Get(ctx, &photo); err != nil {
		log.Fatalln("Error reading value:", err)
	}
	return photo
}

// UpdatePhotoInRealtimeDatabase update images in realtime DB
func UpdatePhotoInRealtimeDatabase(ctx context.Context, updatedPhoto *Photo, reference *db.Ref) {
	currentPhoto := reference.Child("photos").Child(updatedPhoto.Pointer)
	err := currentPhoto.Update(ctx, map[string]interface{}{
		"UserID":       updatedPhoto.UserID,
		"Pointer":      updatedPhoto.Pointer,
		"Tags":         updatedPhoto.Tags,
		"Uploadtime":   updatedPhoto.Uploadtime,
		"Photographer": updatedPhoto.Photographer,
	})
	if err != nil {
		log.Fatalln("Error setting value:", err)
	}
}

// UploadPhotoToRealtimeDatabase uploads images to realtime DB
func UploadPhotoToRealtimeDatabase(ctx context.Context, photo Photo, reference *db.Ref) {
	photoRef := reference.Child("all").Child(photo.UUID)

	err := photoRef.Set(ctx, photo)
	if err != nil {
		log.Fatalln("Error setting value:", err)
	}
}

// ConnectToReference uploads to document containing photos
func ConnectToReference(ctx context.Context, client *db.Client) *db.Ref {
	ref := client.NewRef("photos/uploads")
	var data map[string]interface{}
	if err := ref.Get(ctx, &data); err != nil {
		log.Fatalln("Error reading from database:", err)
	}
	return ref
}

// ConnectToRealtimeDatabase connect to Firebase realtime DB
func ConnectToRealtimeDatabase(ctx context.Context, app *firebase.App) *db.Client {
	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}
	return client
}
