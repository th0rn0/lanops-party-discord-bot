package enable_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"

	"lanops/party-discord-bot/internal/bot/handlers/streams/enable"
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

func newCfg(streamsURL string) config.Config {
	return config.Config{
		Discord: config.Discord{AdminRoleId: "admin-role"},
		Lanops: config.Lanops{
			StreamProxyApiAddress:  streamsURL,
			StreamProxyApiUsername: "user",
			StreamProxyApiPassword: "pass",
		},
	}
}

func TestHandler_NotAdmin(t *testing.T) {
	msgCh := make(chan channels.MsgCh, 10)
	enable.Handler(newMockSession(), newMsg(false), []string{"stream", "enable"}, []string{"mystream"}, newCfg("http://localhost"), msgCh)
	if len(msgCh) != 0 {
		t.Error("expected no messages for non-admin user")
	}
}

func TestHandler_Admin_Enable_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"name": "mystream", "enabled": true})
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	enable.Handler(newMockSession(), newMsg(true), []string{"stream", "enable"}, []string{"mystream"}, newCfg(srv.URL), msgCh)
	if len(msgCh) == 0 {
		t.Error("expected info message in channel")
	}
}

func TestHandler_Admin_Disable_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{"name": "mystream", "enabled": false})
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	// commandParts[1] == "disable" triggers the disable branch
	enable.Handler(newMockSession(), newMsg(true), []string{"stream", "disable"}, []string{"mystream"}, newCfg(srv.URL), msgCh)
	if len(msgCh) == 0 {
		t.Error("expected info message in channel")
	}
}

func TestHandler_Admin_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	enable.Handler(newMockSession(), newMsg(true), []string{"stream", "enable"}, []string{"nostream"}, newCfg(srv.URL), msgCh)
	if len(msgCh) == 0 {
		t.Error("expected info message in channel")
	}
}

func TestHandler_Admin_APIError(t *testing.T) {
	// Close the server before the handler makes its request to trigger a transport error.
	// A 500 response does NOT produce an error from EnableStreamByName due to a bug in
	// the streams client, so a connection failure is required to exercise the error branch.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	enable.Handler(newMockSession(), newMsg(true), []string{"stream", "enable"}, []string{"mystream"}, newCfg(srv.URL), msgCh)

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
