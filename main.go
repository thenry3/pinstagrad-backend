package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	email := strings.Join(r.URL.Query()["email"], "")
	phoneNum := "+" + strings.Join(r.URL.Query()["phonenum"], "")
	log.Print(phoneNum)
	pswrd := strings.Join(r.URL.Query()["password"], "")
	name := strings.Join(r.URL.Query()["name"], "")
	photourl := strings.Join(r.URL.Query()["photourl"], "")

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
	log.Print(dbRef)
	database.UploadPhotoToUserDatabase(ctx, newUser, dbRef)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID = strings.Join(r.URL.Query()["uid"], "")
}

func sendLoginCredentials(w http.ResponseWriter, r *http.Request) {

}

func uploadPictureHandler(w http.ResponseWriter, r *http.Request) {
	contentType := strings.Join(r.URL.Query()["content"], "")
	id, err := strconv.Atoi(strings.Join(r.URL.Query()["userid"], ""))
	log.Print(id)
	if err != nil {
		log.Fatalln("Error converting userid:", err)
	}
	point := strings.Join(r.URL.Query()["point"], "")
	log.Print(point)
	location := strings.Join(r.URL.Query()["location"], "")
	timeTaken, err := strconv.Atoi(strings.Join(r.URL.Query()["time"], ""))
	if err != nil {
		log.Fatalln("Error converting time:", err)
	}
	season, err := strconv.Atoi(strings.Join(r.URL.Query()["season"], ""))
	if err != nil {
		log.Fatalln("Error converting season:", err)
	}
	photographer := strings.Join(r.URL.Query()["photographer"], "")
	log.Print(photographer)

	signedURL, key, err := SignedURL.SignedURLoptions("config/pinstagrad-back-7.json", "PUT", conf.CloudStorage.Bucket, contentType)

	log.Print(signedURL)

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
		Location:     location,
		Time:         database.TimeDay(timeTaken),
		Season:       database.Season(season),
		Uploadtime:   time.Now(),
		Photographer: photographer,
		UUID:         key,
	}

	log.Print(photo)

	ctx := context.Background()
	dbRef := database.ConnectToReference(ctx, database.ConnectToRealtimeDatabase(ctx, app), "photos/uploads")
	log.Print(dbRef)
	database.UploadPhotoToRealtimeDatabase(ctx, photo, dbRef)
}

func getAllPicturesWithFilter(w http.ResponseWriter, r *http.Request) {

}

func retrievePicturesHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	dbRef := database.ConnectToReference(ctx, database.ConnectToRealtimeDatabase(ctx, app), "photos/uploads")
	log.Print(dbRef)
	log.Print(database.RetrieveAllPhotos(ctx, dbRef))
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
	r.HandleFunc("/signup", createNewUserHandler)
	r.HandleFunc("/uploadpicture", uploadPictureHandler)
	r.HandleFunc("/uploadpicture/{userid}/{point}/{}", uploadPictureHandler)
	r.HandleFunc("/retrievepictures", retrievePicturesHandler)
	r.HandleFunc("/retrievepicture/{userid}/{picture}", retrievePicture)

	log.Printf("Listening on %s%s", os.Getenv("HOST"), os.Getenv("PORT"))

	log.Fatal(http.ListenAndServe(":3000", r))

}
