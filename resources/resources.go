package resources

import (
	"net/http"
	"github.com/thenry3/pinstagrad-backend/firebase"
	"context"
)

ctx := context.Background()
config := &firebase.Config{
	DatabaseURL: "https://database-name.firebaseio.com",
  }

// GetAllUsers lol
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HELLO\n"))
}
