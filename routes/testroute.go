package routes

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/twitch"
)

func HandleTestRoute(dbServices *database.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parameters := r.URL.Query()
		endpoint := parameters.Get("endpoint")
		parameters.Del("endpoint")

		if endpoint == "" {
			endpoint = "helix/users"
		}
		if len(parameters) == 0 {
			parameters.Set("login", "bsquidwrd")
		}

		result, err := twitch.CallApi(dbServices, http.MethodGet, endpoint, "", &parameters)
		if err != nil {
			slog.Error("Error calling API", err)
		}

		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(result))
	}
}
