package main

import (
	"log"
	"net/http"

	"github.com/thenry3/pinstagrad-backend/resources"

	"github.com/gorilla/mux"
)

// Handler Test
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Handler)
	r.HandleFunc("/test", resources.GetAllUsers)

	log.Fatal(http.ListenAndServe(":3000", r))

}
