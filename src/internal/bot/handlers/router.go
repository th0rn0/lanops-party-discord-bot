package handlers

import (
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"

	"strings"

	"github.com/bwmarrin/discordgo"
)

// func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
func OnMessage(s *discordgo.Session, m *discordgo.MessageCreate, cfg config.Config, msgCh chan<- channels.MsgCh) {
	// var returnString = "Default Message. If you are seeing this, Corey, Trevor... You fucked up!"
	// var sendMessage = true
	// var sendMessage = false

	// Ignore all messages created by the bot itself
	// if m.Author.ID == s.State.User.ID {
	// 	return
	// }

	if m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, "!") {
		return
	}
	// s.ChannelMessageSend(m.ChannelID, returnString)

	content := strings.TrimPrefix(m.Content, "!")
	parts := strings.Fields(content)
	if len(parts) == 0 {
		// s.ChannelMessageSend(m.ChannelID, returnString)

		return
	}

	for i := len(parts); i > 0; i-- {
		key := strings.Join(parts[:i], " ")
		if handler, ok := Registry[key]; ok {
			handler(s, m, parts[:i], parts[i:], cfg, msgCh)

			// s.ChannelMessageSend(m.ChannelID, returnString)

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

	// Return the Message
	// if sendMessage {
	// 	s.ChannelMessageSend(m.ChannelID, returnString)
	// }
}

func OnReady(s *discordgo.Session, event *discordgo.Ready, msgCh chan<- channels.MsgCh) {
	// Set the playing status.
	s.UpdateGameStatus(0, "Lan Partying!")
}
