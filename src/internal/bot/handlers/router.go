package handlers

import (
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"

	"strings"

	"github.com/bwmarrin/discordgo"
)

func OnMessage(s *discordgo.Session, m *discordgo.MessageCreate, cfg config.Config, msgCh chan<- channels.MsgCh) {
	if m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, cfg.Discord.CommandPrefix) {
		return
	}

	content := strings.TrimPrefix(m.Content, cfg.Discord.CommandPrefix)
	parts := strings.Fields(content)
	if len(parts) == 0 {
		return
	}

	for i := len(parts); i > 0; i-- {
		key := strings.Join(parts[:i], " ")
		if handler, ok := Registry[key]; ok {
			handler(s, m, parts[:i], parts[i:], cfg, msgCh)
			return
		}
	}
}

func OnReady(s *discordgo.Session, m *discordgo.MessageCreate, cfg config.Config, msgCh chan<- channels.MsgCh) {
	// Set the playing status.
	s.UpdateGameStatus(0, "Lan Partying!")
}
