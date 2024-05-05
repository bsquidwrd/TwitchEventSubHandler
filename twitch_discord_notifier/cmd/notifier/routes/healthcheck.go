package routes

import (
	"fmt"
	"net/http"

	"github.com/bsquidwrd/TwitchEventSubHandler/twitch_discord_notifier/internal/service"
)

func HandleHealthCheck(dbServices *service.DiscordNotifierService) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := dbServices.HealthCheck()
		if err != nil {
			http.Error(w, "ERROR", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}
}
