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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Handler)

	log.Fatal(http.ListenAndServe(":3000", r))

}
