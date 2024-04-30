package discordnotifierhandlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/discord"
)

// Returned values are Request Status Code, Response Body, Error
func sendWebhook(url string, method string, body *discord.WebhookBody) (int, discord.WebhookBody, error) {
	data, err := json.Marshal(body)
	if err != nil {
		slog.Warn("Could not marshal body", err)
		return 0, discord.WebhookBody{}, err
	}

	request, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		slog.Warn("Could not assemble request", err)
		return 0, discord.WebhookBody{}, err
	}

	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	var response *http.Response

	// Retry the request 3 times
	for i := 0; i < 3; i++ {
		response, err = client.Do(request)
		if err != nil {
			time.After(1 * time.Second)
			continue
		} else {
			break
		}
	}
	if err != nil {
		slog.Error("Error sending webhook", err)
		return 0, discord.WebhookBody{}, err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("Unable to parse body from response", err)
		return 0, discord.WebhookBody{}, err
	}

	var discordResponse discord.WebhookBody
	err = json.Unmarshal(responseBody, &discordResponse)
	if err != nil {
		slog.Error("Could not unmarshal response body")
		return 0, discord.WebhookBody{}, err
	}

	return response.StatusCode, discordResponse, nil
}
