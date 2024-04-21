package routes

import (
	"fmt"
	"net/http"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
)

func HandleHealthCheck(dbServices *database.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "OK")
	}
}
