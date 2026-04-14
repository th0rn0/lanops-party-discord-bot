package help

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"

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
			Content:   "!help",
			ChannelID: "test-channel",
		},
	}
}

func newCfg() config.Config {
	return config.Config{
		Discord: config.Discord{
			CommandPrefix: "!",
			AdminRoleId:   "admin-role",
		},
	}
}

func TestHandler_Admin(t *testing.T) {
	msgCh := make(chan channels.MsgCh, 10)
	// Should not panic; Discord message is sent via mock transport.
	Handler(newMockSession(), newMsg(true), []string{"help"}, []string{}, newCfg(), msgCh)
}

func TestHandler_NotAdmin(t *testing.T) {
	msgCh := make(chan channels.MsgCh, 10)
	// Non-admin: handler returns early, no Discord send.
	Handler(newMockSession(), newMsg(false), []string{"help"}, []string{}, newCfg(), msgCh)
}

func TestFormatCommands_Empty(t *testing.T) {
	result := formatCommands(map[string]string{})
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestFormatCommands_Single(t *testing.T) {
	result := formatCommands(map[string]string{"!cmd": "does stuff"})
	if !strings.Contains(result, "`!cmd`") {
		t.Errorf("expected formatted command in output, got: %q", result)
	}
	if !strings.Contains(result, "does stuff") {
		t.Errorf("expected description in output, got: %q", result)
	}
}

func TestFormatCommands_Sorted(t *testing.T) {
	result := formatCommands(map[string]string{
		"!zzz": "last",
		"!aaa": "first",
	})
	idxA := strings.Index(result, "!aaa")
	idxZ := strings.Index(result, "!zzz")
	if idxA > idxZ {
		t.Error("expected commands to be sorted alphabetically")
	}
}
