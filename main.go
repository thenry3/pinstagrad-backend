package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"firebase.google.com/go/auth"

	firebaseController "firebase"

	"gcloud/config"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/joho/godotenv"

	"database"

	"github.com/gorilla/mux"

	GCloud "gcloud"

	SignedURL "gcloud/SignedURL"
)

// Global variables
var (
	cloudClient *storage.Client

	app *firebase.App

	authClient *auth.Client

	conf *config.Config
)

//Helper functions

type filtererType func([]database.Photo) []database.Photo

func filters(filterer filtererType, photos []database.Photo) ([]database.Photo, error) {
	return filterer(photos), nil
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
	vars := mux.Vars(r)

	firebaseController.CreateFirebaseUser(authClient, vars["email"], false, vars["phonenumber"], vars["password"],
		vars["name"], vars["photourl"], false)
}

func sendLoginCredentials(w http.ResponseWriter, r *http.Request) {

}

func uploadPictureHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(strings.Join(r.URL.Query()["userid"], ""))
	log.Print(id)
	if err != nil {
		log.Fatalln("Error converting userid:", err)
	}
	point := strings.Join(r.URL.Query()["point"], "")
	log.Print(point)
	tags := createTagList(vars["tags"])
	log.Print(tags)
	photographer := strings.Join(r.URL.Query()["photographer"], "")
	log.Print(photographer)

	signedURL, key, err := SignedURL.SignedURLoptions("./gcloud/config/pinstagrad-back-7.json", "PUT", conf.CloudStorage.Bucket)
	GCloud.Upload(cloudClient, signedURL, conf.CloudStorage.Bucket, point, key)
	if err != nil {
		log.Fatalf("Failed to create signed URL: %v", err)
	}
	log.Printf("Created signed URL: %s", signedURL)
	log.Printf("PUT request using HTTP with URL")

	w.Write([]byte("Handled!\n"))
	photo := database.Photo{
		UserID:       id,
		Pointer:      conf.CloudStorage.URI + key,
		Tags:         tags,
		Uploadtime:   time.Now(),
		Photographer: photographer,
		UUID:         key,
	}

	log.Print(photo)

	ctx := context.Background()
	dbRef := database.ConnectToReference(ctx, database.ConnectToRealtimeDatabase(ctx, app))
	log.Print(dbRef)
	database.UploadPhotoToRealtimeDatabase(ctx, photo, dbRef)
}

func getAllPicturesWithFilter(w http.ResponseWriter, r *http.Request) {

}

func getAllPictures(w http.ResponseWriter, r *http.Request) {

}

func retrievePicture(w http.ResponseWriter, r *http.Request) {

}

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./vendor/gcloud/config/pinstagrad-back-7.json")

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	cloudClient = GCloud.CloudClient()
	GCloud.GetBucket(cloudClient)

	app = firebaseController.FirebaseSDK()
	authClient = firebaseController.AuthClient(app)
}

func main() {
	conf = config.New()

	r := mux.NewRouter()
	r.HandleFunc("/", Handler)
	r.HandleFunc("/signup", createNewUserHandler)
	r.HandleFunc("/uploadpicture", uploadPictureHandler)
	r.HandleFunc("/uploadpicture/{userid}/{point}/{}", uploadPictureHandler)
	r.HandleFunc("/retrievepicture/{userid}/{picture}", retrievePicture)
	log.Fatal(http.ListenAndServe(":3000", r))

}
