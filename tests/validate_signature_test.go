package tests

import (
	"testing"

	"github.com/bsquidwrd/TwitchEventSubHandler/internal/utils"
)

func TestValidateSignature(t *testing.T) {
	secret := []byte("0123456789")
	messageID := "c6f33bb6-1c72-46ff-81f0-6c73e2cf78ba"
	messageTimestamp := "2024-04-19T13:23:54-08:00"
	body := []byte("{\"subscription\":{\"id\":\"c6f33bb6-1c72-46ff-81f0-6c73e2cf78ba\",\"type\":\"stream.offline\",\"version\":\"1\",\"status\":\"enabled\",\"cost\":0,\"condition\":{\"broadcaster_user_id\":\"22812120\"},\"created_at\":\"2024-04-19T20:23:54.634234626Z\",\"transport\":{\"method\":\"webhook\",\"callback\":\"https://example.com/webhooks/callback\"}},\"event\":{\"broadcaster_user_id\":\"22812120\",\"broadcaster_user_login\":\"bsquidwrd\",\"broadcaster_user_name\":\"bsquidwrd\"}}")
	messageSignature := "9098685d5aebcafdcc958dcb12fa3d44f3d2f00f02f0b53f4c13b096f582e6cf"

	signatureMatches := utils.ValidateSignature(secret, messageID, messageTimestamp, body, messageSignature)
	if !signatureMatches {
		t.Error("Signature of message did not match expectation")
	}
}
