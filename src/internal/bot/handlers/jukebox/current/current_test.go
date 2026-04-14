package current_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"

	"lanops/party-discord-bot/internal/bot/handlers/jukebox/current"
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"
	"lanops/party-discord-bot/internal/jukebox"
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

func newMsg() *discordgo.MessageCreate {
	return &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author:    &discordgo.User{Bot: false},
			Member:    &discordgo.Member{Roles: []string{}},
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

func TestHandler_Success(t *testing.T) {
	track := jukebox.GetCurrentTrackOutput{} // used just for struct reference
	_ = track
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"name": "Sandstorm", "artist": "Darude"})
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	current.Handler(newMockSession(), newMsg(), []string{"jukebox", "current"}, []string{}, newCfg(srv.URL), msgCh)

	if len(msgCh) == 0 {
		t.Error("expected info message in channel")
	}
}

func TestHandler_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	msgCh := make(chan channels.MsgCh, 10)
	current.Handler(newMockSession(), newMsg(), []string{"jukebox", "current"}, []string{}, newCfg(srv.URL), msgCh)

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
