package start_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"

	"lanops/party-discord-bot/internal/bot/handlers/jukebox/start"
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"
)

type mockDiscordTransport struct{}

func (t *mockDiscordTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"id":"1","channel_id":"ch","content":"ok","author":{"id":"1","username":"bot","discriminator":"0000"}}`)),
		Header:     make(http.Header),
	}, nil
}

func newMockSession() *discordgo.Session {
	s, _ := discordgo.New("Bot fake-token")
	s.Client = &http.Client{Transport: &mockDiscordTransport{}}
	return s
}

func newMsg(isAdmin bool) *discordgo.MessageCreate {
	roles := []string{}
	if isAdmin {
		roles = []string{"admin-role"}
	}
	return &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author:    &discordgo.User{Bot: false},
			Member:    &discordgo.Member{Roles: roles},
			ChannelID: "test-channel",
		},
	}
}

func newCfg(jukeboxURL string) config.Config {
	return config.Config{
		Discord: config.Discord{AdminRoleId: "admin-role"},
		Lanops: config.Lanops{
			JukeboxApiUrl:      jukeboxURL,
			JukeboxApiUsername: "user",
			JukeboxApiPassword: "pass",
		},
	}
}

func TestHandler_NotAdmin(t *testing.T) {
	msgCh := make(chan channels.MsgCh, 10)
	start.Handler(newMockSession(), newMsg(false), []string{"jukebox", "start"}, []string{}, newCfg("http://localhost"), msgCh)
	if len(msgCh) != 0 {
		t.Error("expected no messages for non-admin user")
	}
}

func TestHandler_Admin_APISuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	start.Handler(newMockSession(), newMsg(true), []string{"jukebox", "start"}, []string{}, newCfg(srv.URL), msgCh)
	if len(msgCh) == 0 {
		t.Error("expected info message in channel")
	}
}

func TestHandler_Admin_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	start.Handler(newMockSession(), newMsg(true), []string{"jukebox", "start"}, []string{}, newCfg(srv.URL), msgCh)

	// Should have an error message in channel
	var hasError bool
	for len(msgCh) > 0 {
		msg := <-msgCh
		if msg.Err != nil {
			hasError = true
		}
	}
	if !hasError {
		t.Error("expected error message in channel for API failure")
	}
}
