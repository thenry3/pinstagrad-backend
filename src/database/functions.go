package database

import (
	"context"
	"log"

	"firebase.google.com/go/db"

	firebase "firebase.google.com/go"
)

//AddLikes add liked photo to user's liked list
func AddLikes(ctx context.Context, reference *db.Ref, currentUserUID string, photoPointer string) string {
	currentUser := User{}
	if err := reference.Child(currentUserUID).Get(ctx, &currentUser); err != nil {
		log.Fatalln("Error reading from database:", err)
	}
	liked := append(currentUser.Liked, photoPointer)
	currentUser.Liked = liked
	err := reference.Child(currentUserUID).Update(ctx, map[string]interface{}{
		"liked": currentUser.Liked,
	})
	if err != nil {
		log.Fatalln("Error setting value:", err)
	}
	return photoPointer
}

// RetrieveIndividualPhoto retrieves singular image by uuid
func RetrieveIndividualPhoto(ctx context.Context, reference *db.Ref, uuid string) Photo {
	v := Photo{}
	if err := reference.Child("all").Child(uuid).Get(ctx, &v); err != nil {
		log.Fatalln("Error reading value:", err)
	}
	return v
}

// RetrieveAllPhotos retrieves images in realtime DB
func RetrieveAllPhotos(ctx context.Context, reference *db.Ref) map[string]Photo {
	v := map[string]Photo{}
	if err := reference.Child("all").Get(ctx, &v); err != nil {
		log.Fatalln("Error reading value:", err)
	}
	return v
}

// UploadPhotoToUserDatabase uploads images to realtime DB
func UploadPhotoToUserDatabase(ctx context.Context, user User, reference *db.Ref) {
	photoRef := reference.Child(user.UID)

	err := photoRef.Set(ctx, user)
	if err != nil {
		log.Fatalln("Error setting value:", err)
	}
}

// UpdatePhotoInRealtimeDatabase update images in realtime DB
func UpdatePhotoInRealtimeDatabase(ctx context.Context, updatedPhoto *Photo, reference *db.Ref) {
	currentPhoto := reference.Child("photos").Child(updatedPhoto.Pointer)
	err := currentPhoto.Update(ctx, map[string]interface{}{
		"UserID":       updatedPhoto.UserID,
		"Pointer":      updatedPhoto.Pointer,
		"Location":     updatedPhoto.Location,
		"Season":       updatedPhoto.Season,
		"Time":         updatedPhoto.Time,
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
func ConnectToReference(ctx context.Context, client *db.Client, path string) *db.Ref {
	ref := client.NewRef(path)
	var data map[string]interface{}
	if err := ref.Get(ctx, &data); err != nil {
		log.Fatalln("Error reading from database:", err)
	}
	return ref
}

// ConnectToRealtimeDatabase connect to Firebase realtime database
func ConnectToRealtimeDatabase(ctx context.Context, app *firebase.App) *db.Client {
	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("Error initializing database client:", err)
	}
	return client
}
