package handlers_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"

	"lanops/party-discord-bot/internal/bot/handlers"
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

func newMsg(content string, isBot bool, roles []string) *discordgo.MessageCreate {
	return &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author:    &discordgo.User{Bot: isBot},
			Member:    &discordgo.Member{Roles: roles},
			Content:   content,
			ChannelID: "test-channel",
		},
	}
}

func newTestCfg() config.Config {
	return config.Config{
		Discord: config.Discord{
			CommandPrefix: "!",
			AdminRoleId:   "admin-role",
		},
	}
}

func newMsgCh() chan channels.MsgCh {
	return make(chan channels.MsgCh, 10)
}

func TestOnMessage_BotAuthor(t *testing.T) {
	called := false
	handlers.Register("bottest", func(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, args []string, cfg config.Config, msgCh chan<- channels.MsgCh) {
		called = true
	})
	defer delete(handlers.Registry, "bottest")

	handlers.OnMessage(newMockSession(), newMsg("!bottest", true, nil), newTestCfg(), newMsgCh())
	if called {
		t.Error("handler should not be called for bot author")
	}
}

func TestOnMessage_NoPrefix(t *testing.T) {
	called := false
	handlers.Register("noprefix", func(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, args []string, cfg config.Config, msgCh chan<- channels.MsgCh) {
		called = true
	})
	defer delete(handlers.Registry, "noprefix")

	handlers.OnMessage(newMockSession(), newMsg("noprefix", false, nil), newTestCfg(), newMsgCh())
	if called {
		t.Error("handler should not be called without prefix")
	}
}

func TestOnMessage_EmptyContent(t *testing.T) {
	// Prefix only, no command
	handlers.OnMessage(newMockSession(), newMsg("!", false, nil), newTestCfg(), newMsgCh())
	// Should not panic; nothing to assert other than no crash
}

func TestOnMessage_MatchingCommand(t *testing.T) {
	called := false
	var gotParts, gotArgs []string
	handlers.Register("testroute arg", func(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, args []string, cfg config.Config, msgCh chan<- channels.MsgCh) {
		called = true
		gotParts = parts
		gotArgs = args
	})
	defer delete(handlers.Registry, "testroute arg")

	handlers.OnMessage(newMockSession(), newMsg("!testroute arg extra", false, nil), newTestCfg(), newMsgCh())
	if !called {
		t.Error("handler was not called")
	}
	if len(gotParts) != 2 || gotParts[0] != "testroute" || gotParts[1] != "arg" {
		t.Errorf("unexpected parts: %v", gotParts)
	}
	if len(gotArgs) != 1 || gotArgs[0] != "extra" {
		t.Errorf("unexpected args: %v", gotArgs)
	}
}

func TestOnMessage_NoMatchingCommand(t *testing.T) {
	// A command with no registered handler should simply do nothing.
	handlers.OnMessage(newMockSession(), newMsg("!unknowncmd", false, nil), newTestCfg(), newMsgCh())
}

func TestOnReady(t *testing.T) {
	// OnReady calls s.UpdateGameStatus which returns ErrWSNotFound on a disconnected session.
	// The function ignores that error, so this should not panic.
	handlers.OnReady(newMockSession(), newMsg("", false, nil), newTestCfg(), newMsgCh())
}
