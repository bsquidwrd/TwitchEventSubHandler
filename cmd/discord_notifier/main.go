package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bsquidwrd/TwitchEventSubHandler/cmd/discord_notifier/routes"
	"github.com/bsquidwrd/TwitchEventSubHandler/internal/database"
	"github.com/bsquidwrd/TwitchEventSubHandler/pkg/models/twitch"
	amqp "github.com/rabbitmq/amqp091-go"
)

const queueName = "discord_notifier"

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Discord Notifier Starting Up...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = "8081"

	dbServices := database.NewDiscordNotifierService()
	defer dbServices.Cleanup()

	topics := []string{
		"stream.online",
		"stream.offline",
		"channel.update",
		"user.update",
	}

	dbServices.Queue.StartConsuming(queueName, topics, func(msg amqp.Delivery) {
		testFunc(dbServices, msg)
	})

	http.HandleFunc("/", routes.HandleRoot)
	http.HandleFunc("GET /healthcheck", routes.HandleHealthCheck(dbServices))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Running server on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func testFunc(dbServices *database.DiscordNotifierService, msg amqp.Delivery) {
	var userId string
	switch msg.RoutingKey {
	case "channel.update":
		var event twitch.ChannelUpdateEventSubEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			slog.Warn("Could not parse message", "topic", msg.RoutingKey, err)
			return
		}
		userId = event.BroadcasterUserID
	case "user.update":
		var event twitch.UserUpdateEventSubEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			slog.Warn("Could not parse message", "topic", msg.RoutingKey, err)
			return
		}
		userId = event.UserID
	case "stream.online":
		var event twitch.StreamUpEventSubEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			slog.Warn("Could not parse message", "topic", msg.RoutingKey, err)
			return
		}
		userId = event.BroadcasterUserID
	case "stream.offline":
		var event twitch.StreamDownEventSubEvent
		err := json.Unmarshal(msg.Body, &event)
		if err != nil {
			slog.Warn("Could not parse message", "topic", msg.RoutingKey, err)
			return
		}
		userId = event.BroadcasterUserID
	}

	if userId == "" {
		return
	}

	dbUser := dbServices.Database.QueryRow(
		context.Background(),
		`
			select id, "name", login, description, title, language, category_id, category_name, last_online_at, last_offline_at, live
			from public.twitch_user
			where id=$1
		`,
		userId,
	)

	var user twitch.DatabaseUser
	err := dbUser.Scan(&user.Id, &user.Name, &user.Login, &user.Description, &user.Title, &user.Language, &user.CategoryID, &user.CategoryName, &user.LastOnlineAt, &user.LastOfflineAt, &user.Live)
	if err != nil {
		slog.Warn("Could not retrieve user from database", err)
		return
	}

	slog.Info("Got user info!", user)
}
