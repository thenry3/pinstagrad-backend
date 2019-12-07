package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"firebase.google.com/go/db"

	"github.com/Rahul12344/pinstagrad-backend/config"

	"firebase.google.com/go/auth"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	firebaseController "github.com/Rahul12344/pinstagrad-backend/src/firebase"
	"github.com/joho/godotenv"

	"github.com/Rahul12344/pinstagrad-backend/src/database"

	"github.com/gorilla/mux"

	GCloud "github.com/Rahul12344/pinstagrad-backend/src/gcloud"
	SignedURL "github.com/Rahul12344/pinstagrad-backend/src/gcloud/SignedURL"
)

// Global variables
var (
	cloudClient *storage.Client

	app *firebase.App

	authClient *auth.Client

	conf *config.Config

	currentUserID string
)

//UploadData data necessary to upload photos
type UploadData struct {
	Point        string
	Location     string
	TimeTaken    int
	Season       int
	Photographer string
	UserID       string
}

//Helper functions

type filtererType func([]database.Photo) []database.Photo
type updates func(context.Context, *db.Ref, string, string) string

func filters(filterer filtererType, photos []database.Photo) ([]database.Photo, error) {
	return filterer(photos), nil
}

func update(ctx context.Context, updater updates, ref *db.Ref, id string, fields string) {
	updater(ctx, ref, id, fields)
}

func createTagList(connectedTags string) [5]string {
	var tags [5]string
	tags[0] = "target"
	return tags
}

// Handler Test
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

/**
API Endpoints
**/

func createNewUserHandler(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	phoneNum := r.FormValue("phonenum")
	pswrd := r.FormValue("password")
	name := r.FormValue("name")
	photourl := r.FormValue("photourl")

	user := firebaseController.CreateFirebaseUser(conf, "image/png", cloudClient, authClient, email, false, phoneNum, pswrd,
		name, photourl, false)

	newUser := database.User{
		UID:        user.UID,
		Email:      user.Email,
		Name:       name,
		Password:   pswrd,
		ProfilePic: user.PhotoURL,
	}
	ctx := context.Background()
	dbRef := database.ConnectToReference(ctx, database.ConnectToRealtimeDatabase(ctx, app), "users")
	database.UploadPhotoToUserDatabase(ctx, newUser, dbRef)

	JSONuser, err := json.Marshal(newUser)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.Write(JSONuser)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID = strings.Join(r.URL.Query()["uid"], "")
	w.Write([]byte(currentUserID))
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID = ""
	w.Write([]byte(currentUserID))
}

func sendLoginCredentials(w http.ResponseWriter, r *http.Request) {

}

func uploadPictureHandler(w http.ResponseWriter, r *http.Request) {
	if currentUserID == "" {
		http.Error(w, "Anon User must be logged in to upload", 500)
		return
	}
	contentType := r.FormValue("content")
	id, err := strconv.Atoi(r.FormValue("userid"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	point := r.FormValue("point")
	location := r.FormValue("location")
	timeTaken, err := strconv.Atoi(r.FormValue("time"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	season, err := strconv.Atoi(r.FormValue("season"))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	photographer := r.FormValue("photographer")

	signedURL, key, err := SignedURL.SignedURLoptions("config/pinstagrad-back-7.json", "PUT", conf.CloudStorage.Bucket, contentType)

	GCloud.Upload(cloudClient, signedURL, conf.CloudStorage.Bucket, point, key, contentType)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	log.Printf("Created signed URL: %s", signedURL)
	log.Printf("PUT request using HTTP with URL")

	photo := database.Photo{
		UserID:       id,
		Pointer:      conf.CloudStorage.URI + key,
		Location:     location,
		Time:         database.TimeDay(timeTaken),
		Season:       database.Season(season),
		Uploadtime:   time.Now(),
		Photographer: photographer,
		UUID:         key,
	}

	ctx := context.Background()
	dbRef := database.ConnectToReference(ctx, database.ConnectToRealtimeDatabase(ctx, app), "photos/uploads")
	database.UploadPhotoToRealtimeDatabase(ctx, photo, dbRef)

	JSONphoto, err := json.Marshal(photo)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	w.Write(JSONphoto)
}

func getAllPicturesWithFilter(w http.ResponseWriter, r *http.Request) {

}

func updateUserFieldsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	field := r.FormValue("field")
	switch field {
	case "likes":
		dbRef := database.ConnectToReference(ctx, database.ConnectToRealtimeDatabase(ctx, app), "users")
		update(ctx, database.AddLikes, dbRef, currentUserID, r.FormValue("photourl"))
	}
}

func retrievePicturesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	dbRef := database.ConnectToReference(ctx, database.ConnectToRealtimeDatabase(ctx, app), "photos/uploads")
	photos := database.RetrieveAllPhotos(ctx, dbRef)

	JSONphotos, err := json.Marshal(photos)
	if err != nil {
		w.Write([]byte("Error converting photos to JSON"))
	}

	w.Write(JSONphotos)
}

func retrievePicture(w http.ResponseWriter, r *http.Request) {

}

func init() {
	currentUserID = ""

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "/Users/rahulnatarajan/go/src/github.com/Rahul12344/pinstagrad-backend/config/pinstagrad-back-7.json")

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	log.Print("Setting environment variables...")

	cloudClient = GCloud.CloudClient()
	GCloud.GetBucket(cloudClient)
	GCloud.NewBucket(cloudClient, os.Getenv("PROFILE_BUCKET_NAME"))

	log.Print("Intiliazing google cloud client and bucket...")

	app = firebaseController.FirebaseSDK()
	authClient = firebaseController.AuthClient(app)
	log.Print("Intiliazing firebase client and admin...")
}

func main() {
	conf = config.New()

	r := mux.NewRouter()
	log.Print("Intiliazing Mux router...")

	r.HandleFunc("/", Handler)
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/logout", logoutHandler)
	r.HandleFunc("/signup", createNewUserHandler)
	r.HandleFunc("/uploadpicture", uploadPictureHandler)
	r.HandleFunc("/uploadpicture/{userid}/{point}/{}", uploadPictureHandler)
	r.HandleFunc("/retrievepictures", retrievePicturesHandler)
	r.HandleFunc("/retrievepicture/{userid}/{picture}", retrievePicture)

	log.Printf("Listening on %s%s", os.Getenv("HOST"), os.Getenv("PORT"))

	log.Fatal(http.ListenAndServe(":3000", r))

}
