package twitch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
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

func getAuthKey(dbServices *database.Service) string {
	cachedKey := "twitch:api:authkey"
	existingAuthKey := dbServices.Cache.GetString(cachedKey)

	if existingAuthKey != "" {
		return existingAuthKey
	}

	dbAuthData := dbServices.Database.QueryRow(
		context.Background(),
		"select access_token from public.twitch_auth where expired is not true and expires_at >= current_timestamp order by expires_at desc limit 1",
	)
	var dbKey string
	dbAuthData.Scan(&dbKey)
	if dbKey != "" {
		dbServices.Cache.SetString(cachedKey, dbKey, 5*time.Minute)
		return dbKey
	}

	if dbServices.Twitch.AuthLock.TryLock() {
		defer dbServices.Twitch.AuthLock.Unlock()
		newAuthKey, err := getNewAuthKey()
		if err != nil {
			slog.Error("Error getting new Auth Key", err)
		}
		if newAuthKey.ExpiresIn > 0 {
			expirationDuration := time.Duration(newAuthKey.ExpiresIn) * time.Second
			dbServices.Cache.SetString(cachedKey, newAuthKey.AccessToken, expirationDuration)

			_, err := dbServices.Database.Exec(
				context.Background(),
				"update public.twitch_auth set expired = true where expired is not true",
			)
			if err != nil {
				slog.Error("Error invalidating old access tokens in db", "error", err)
			}

			_, err = dbServices.Database.Exec(
				context.Background(),
				`INSERT INTO public.twitch_auth
				(client_id, access_token, expires_at)
				VALUES($1, $2, $3);`,
				clientID,
				newAuthKey.AccessToken,
				time.Now().UTC().Add(expirationDuration),
			)
			if err != nil {
				slog.Error("Error inserting new access token into db", "error", err)
			}

			return newAuthKey.AccessToken
		}
	} else {
		// If we couldn't get a lock on it, wait until we can
		// Not being able to get a lock means someone else is refreshing the auth key
		dbServices.Twitch.AuthLock.Lock()
		dbServices.Twitch.AuthLock.Unlock()
		return getAuthKey(dbServices)
	}

	return ""
}

func getNewAuthKey() (*clientCredentials, error) {
	tokenUrl := "https://id.twitch.tv/oauth2/token"
	if os.Getenv("API_URL") != "" {
		tokenUrl = fmt.Sprintf("%s/auth/token", os.Getenv("API_URL"))
	}

	requestUrl, _ := url.ParseRequestURI(tokenUrl)
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "client_credentials")

	request, err := http.NewRequest(http.MethodPost, requestUrl.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return &clientCredentials{}, err
	}
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	var response *http.Response

	// Retry the request 3 times
	for i := 0; i < 3; i++ {
		err = nil
		response, err = client.Do(request)
		if err != nil {
			time.After(1 * time.Second)
			continue
		} else {
			break
		}
	}
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

func CallApi(dbServices *database.Service, method string, endpoint string, data string, parameters *url.Values) (int, []byte, error) {
	baseUrl := "https://api.twitch.tv/helix/"
	if os.Getenv("API_URL") != "" {
		baseUrl = fmt.Sprintf("%s/mock/", os.Getenv("API_URL"))
	}

	requestUrl, _ := url.ParseRequestURI(baseUrl)
	requestUrl.Path += endpoint

	if parameters == nil {
		parameters = &url.Values{}
	}
	requestUrl.RawQuery = parameters.Encode()

	if method == "" {
		method = http.MethodGet
	}

	request, err := http.NewRequest(method, requestUrl.String(), strings.NewReader(data))
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Client-ID", clientID)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", getAuthKey(dbServices)))

	if err != nil {
		slog.Error("Error assembling API request", err)
		return 0, nil, err
	}

	// Make sure there's no existing rate limit in place
	if dbServices.Twitch.RatelimitLock.TryLock() {
		dbServices.Twitch.RatelimitLock.Unlock()
	} else {
		// If we couldn't get a lock on it, wait until we can
		// Not being able to get a lock means a rate limit is in effect
		dbServices.Twitch.AuthLock.Lock()
		dbServices.Twitch.AuthLock.Unlock()
	}

	client := &http.Client{}
	var response *http.Response

	// Retry the request 3 times
	for i := 0; i < 3; i++ {
		err = nil
		response, err = client.Do(request)
		if err != nil {
			time.After(1 * time.Second)
			continue
		} else {
			break
		}
	}
	if err != nil {
		slog.Error("Error calling API", err)
		return 0, nil, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusTooManyRequests {
		ratelimitResetValue, err := strconv.ParseInt(response.Header.Get("Ratelimit-Reset"), 10, 64)
		if err != nil {
			return 0, nil, err
		}
		slog.Info("Twitch API rate limit hit, waiting it out", "reset", ratelimitResetValue)
		ratelimitReset := time.Unix(ratelimitResetValue, 0)
		dbServices.Twitch.RatelimitLock.Lock()
		time.After(ratelimitReset.Sub(time.Now().UTC()))
		dbServices.Twitch.RatelimitLock.Unlock()
		return CallApi(dbServices, method, endpoint, data, parameters)
	} else {
		rawBody, err := io.ReadAll(response.Body)
		if err != nil {
			slog.Error("Error reading API body", err)
			return 0, nil, err
		}

		return response.StatusCode, rawBody, nil
	}
}

func DeleteSubscription(dbServices *database.Service, id string) (int, []byte, error) {
	parameters := &url.Values{}
	parameters.Add("id", id)
	return CallApi(dbServices, http.MethodDelete, "eventsub/subscriptions", "", parameters)
}
