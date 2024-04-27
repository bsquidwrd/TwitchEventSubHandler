package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
)

func HandleHealthCheck(dbServices *database.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := dbServices.Database.Ping(context.Background())
		if err != nil {
			http.Error(w, "ERROR", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}
}
