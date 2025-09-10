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

	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	content := strings.TrimPrefix(m.Content, "!")
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

	// if strings.HasPrefix(m.Content, commandPrefix+"jukebox") {
	// 	logger.Info().Msg("Message Create Event - Jukebox - Triggered")
	// 	returnString = "we are testing pal jukebox"
	// 	sendMessage = true
	// }

	// if m.Content == commandPrefix+"jb current" {
	// 	logger.Info().Msg("Message Create Event - Jukebox Currently playing - Triggered")
	// 	returnString = jukeboxAPI.GetCurrentTrack()
	// 	sendMessage = true
	// } else if slices.Contains(m.Member.Roles, discordJukeBoxControlRoleID) {
	// 	if strings.HasPrefix(m.Content, commandPrefix+"jb") {
	// 		logger.Info().Msg("Message Create Event - Jukebox Control - Triggered")
	// 		jukeboxCommand := strings.Split(m.Content, " ")
	// 		returnString = jukeboxAPI.Control(jukeboxCommand[1])
	// 		sendMessage = true
	// 	}
}

func OnReady(s *discordgo.Session, m *discordgo.MessageCreate, cfg config.Config, msgCh chan<- channels.MsgCh) {
	// Set the playing status.
	s.UpdateGameStatus(0, "Lan Partying!")
}
