package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"./database"
	firebaseController "./firebase"

	"github.com/gorilla/mux"
)

//Helper functions

func createTagList(connectedTags string) []string {
	var tags []string
	return tags
}

func createUploadTime() time.Time {
	start := time.Now()
	return start
}

// Handler Test
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

/**
API Endpoints
**/

func createNewUser(w http.ResponseWriter, r *http.Request) {

}

func sendLoginCredentials(w http.ResponseWriter, r *http.Request) {

}

func uploadPicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["userid"])
	if err != nil {
		log.Fatalln("Error converting userid:", err)
	}
	point := vars["pointer"]
	tags := createTagList(vars["tags"])

	photo := database.Photo{UserID: id,
		Pointer:      point,
		Tags:         tags,
		Uploadtime:   createUploadTime(),
		Photographer: vars["photographer"]}

	ctx := context.Background()
	app := firebaseController.FirebaseSDK()
	dbRef := database.ConnectToReference(ctx, database.ConnectToRealtimeDatabase(ctx, app))
	database.UploadPhotoToRealtimeDatabase(ctx, photo, dbRef)
}

func retrievePicture(w http.ResponseWriter, r *http.Request) {

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Handler)
	r.HandleFunc("/uploadpicture/{userid}/{picture}", uploadPicture)
	r.HandleFunc("/retrievepicture/{userid}/{picture}", retrievePicture)

	log.Fatal(http.ListenAndServe(":3000", r))

}
