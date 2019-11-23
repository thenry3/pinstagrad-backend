package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler Test
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

/**
API Endpoints
**/

func uploadPicture(w http.ResponseWriter, r *http.Request) {

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
