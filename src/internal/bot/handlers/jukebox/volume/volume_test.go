package volume_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"

	"lanops/party-discord-bot/internal/bot/handlers/jukebox/volume"
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
	volume.Handler(newMockSession(), newMsg(false), []string{"jukebox", "volume"}, []string{}, newCfg("http://localhost"), msgCh)
	if len(msgCh) != 0 {
		t.Error("expected no messages for non-admin user")
	}
}

func TestHandler_Admin_NoArgs_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(map[string]int{"volume": 42})
			return
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	volume.Handler(newMockSession(), newMsg(true), []string{"jukebox", "volume"}, []string{}, newCfg(srv.URL), msgCh)
	if len(msgCh) == 0 {
		t.Error("expected info message in channel")
	}
}

func TestHandler_Admin_NoArgs_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	volume.Handler(newMockSession(), newMsg(true), []string{"jukebox", "volume"}, []string{}, newCfg(srv.URL), msgCh)

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

func TestHandler_Admin_WithArg_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	volume.Handler(newMockSession(), newMsg(true), []string{"jukebox", "volume"}, []string{"50"}, newCfg(srv.URL), msgCh)
	if len(msgCh) == 0 {
		t.Error("expected info message in channel")
	}
}

func TestHandler_Admin_WithArg_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	volume.Handler(newMockSession(), newMsg(true), []string{"jukebox", "volume"}, []string{"50"}, newCfg(srv.URL), msgCh)

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

// TestHandler_Admin_InvalidArg covers the strconv.Atoi error branch ("I can't do that brian").
// Due to a bug in the handler, SetVolume(0) is still called after the parse error,
// so the final Discord message is determined by SetVolume's outcome.
func TestHandler_Admin_InvalidArg_APISuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	volume.Handler(newMockSession(), newMsg(true), []string{"jukebox", "volume"}, []string{"notanumber"}, newCfg(srv.URL), msgCh)
	if len(msgCh) == 0 {
		t.Error("expected info message in channel")
	}
}

func TestHandler_Admin_InvalidArg_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	volume.Handler(newMockSession(), newMsg(true), []string{"jukebox", "volume"}, []string{"notanumber"}, newCfg(srv.URL), msgCh)

	var hasError bool
	for len(msgCh) > 0 {
		msg := <-msgCh
		if msg.Err != nil {
			hasError = true
		}
	}
	if !hasError {
		t.Error("expected error message in channel for API failure after invalid arg")
	}
}
