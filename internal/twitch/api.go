package twitch

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
)

var (
	clientID     = os.Getenv("CLIENTID")
	clientSecret = os.Getenv("CLIENTSECRET")
)

type clientCredentials struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func fakeDBCall() (string, error) { return "", nil }

func getAuthKey(dbServices *database.Services) string {
	redisKey := "twitch:api:authkey"
	existingAuthKey := dbServices.Redis.GetString(redisKey)

	if existingAuthKey != "" {
		return existingAuthKey
	}

	dbKey, err := fakeDBCall()
	if err != nil {
		slog.Error("Error getting Auth key from DB", err)
	}
	if dbKey != "" {
		dbServices.Redis.SetString(redisKey, dbKey, 5*time.Minute)
		return dbKey
	}

	gotAuthLock := dbServices.Twitch.AuthLock.TryLock()

	if gotAuthLock {
		defer dbServices.Twitch.AuthLock.Unlock()
		newAuthKey, err := getNewAuthKey()
		if err != nil {
			slog.Error("Error getting new Auth Key", err)
		}
		if newAuthKey.ExpiresIn > 0 {
			dbServices.Redis.SetString(redisKey, newAuthKey.AccessToken, time.Duration(newAuthKey.ExpiresIn))
			// set value in db too
			return newAuthKey.AccessToken
		}
	} else {
		dbServices.Twitch.AuthLock.Lock()
		dbServices.Twitch.AuthLock.Unlock()
		return getAuthKey(dbServices)
	}

	return ""
}

func getNewAuthKey() (*clientCredentials, error) {
	requestUrl := "https://id.twitch.tv/oauth2/token"
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")

	request, err := http.NewRequest(http.MethodPost, requestUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return &clientCredentials{}, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		slog.Error("Error calling API", err)
		return nil, err
	}
	defer response.Body.Close()

	rawBody, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("Error reading API body", err)
		return &clientCredentials{}, err
	}

	var credentials clientCredentials
	json.Unmarshal(rawBody, &credentials)

	return &credentials, nil
}

func CallApi(dbServices *database.Services, endpoint string, data string, parameters *url.Values) ([]byte, error) {
	requestUrl, _ := url.ParseRequestURI("https://api.twitch.tv/")
	requestUrl.Path = endpoint
	requestUrl.RawQuery = parameters.Encode()

	request, err := http.NewRequest(http.MethodGet, requestUrl.String(), strings.NewReader(data))
	authKey := getAuthKey(dbServices)
	request.Header.Add("Client-ID", clientID)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authKey))

	if err != nil {
		slog.Error("Error assembling API request", err)
		return nil, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		slog.Error("Error calling API", err)
		return nil, err
	}
	defer response.Body.Close()

	rawBody, err := io.ReadAll(response.Body)
	if err != nil {
		slog.Error("Error reading API body", err)
		return nil, err
	}

	return rawBody, nil
}
