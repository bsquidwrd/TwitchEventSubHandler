package routes

import (
	"fmt"
	"html"
	"log/slog"
	"net/http"
)

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	slog.Info(
		"Healthcheck called",
		"endpoint", html.EscapeString(r.URL.Path),
	)

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
