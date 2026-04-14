package handlers_test

import (
	"testing"

	"github.com/bwmarrin/discordgo"

	"lanops/party-discord-bot/internal/bot/handlers"
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"
)

func TestRegister(t *testing.T) {
	called := false
	handlers.Register("regtest", func(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, args []string, cfg config.Config, msgCh chan<- channels.MsgCh) {
		called = true
	})
	defer delete(handlers.Registry, "regtest")

	fn, ok := handlers.Registry["regtest"]
	if !ok {
		t.Fatal("handler not found in registry after Register()")
	}
	fn(nil, nil, nil, nil, config.Config{}, make(chan channels.MsgCh, 1))
	if !called {
		t.Error("registered handler was not called")
	}
}

func TestRegistry_ExpectedCommands(t *testing.T) {
	expected := []string{
		"help",
		"stream list",
		"stream enable",
		"stream disable",
		"jukebox start",
		"jukebox stop",
		"jukebox pause",
		"jukebox skip",
		"jukebox volume",
		"jukebox current",
	}
	for _, cmd := range expected {
		if _, ok := handlers.Registry[cmd]; !ok {
			t.Errorf("command %q not registered", cmd)
		}
	}
}
