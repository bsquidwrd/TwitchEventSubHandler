package routes

import (
	"fmt"
	"html"
	"log/slog"
	"net/http"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
)

func HandleHealthCheck(dbServices *database.Services) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info(
			"Healthcheck called",
			"endpoint", html.EscapeString(r.URL.Path),
		)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}
}
