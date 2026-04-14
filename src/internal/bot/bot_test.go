package bot_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/bwmarrin/discordgo"

	"lanops/party-discord-bot/internal/bot"
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"
)

type errorTransport struct{}

func (e *errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("mock transport error")
}


func newTestCfg() config.Config {
	return config.Config{
		Discord: config.Discord{
			Token:         "fake-token",
			AdminRoleId:   "admin-role",
			CommandPrefix: "!",
		},
	}
}

func TestNew(t *testing.T) {
	s, err := discordgo.New("Bot fake-token")
	if err != nil {
		t.Fatalf("failed to create discordgo session: %v", err)
	}
	msgCh := make(chan channels.MsgCh, 10)

	client, err := bot.New(newTestCfg(), s, msgCh)
	if err != nil {
		t.Fatalf("bot.New() returned unexpected error: %v", err)
	}
	if client == nil {
		t.Fatal("bot.New() returned nil client")
	}
}

func TestRun_OpenError(t *testing.T) {
	s, _ := discordgo.New("Bot fake-token")
	s.Client = &http.Client{Transport: &errorTransport{}}

	msgCh := make(chan channels.MsgCh, 10)
	client, _ := bot.New(newTestCfg(), s, msgCh)

	err := client.Run()
	if err == nil {
		t.Error("expected error from Run() when Open() fails, got nil")
	}
}
