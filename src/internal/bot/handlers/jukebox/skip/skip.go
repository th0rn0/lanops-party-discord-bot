package skip

import (
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"
	"lanops/party-discord-bot/internal/jukebox"
	"slices"

	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate, commandParts []string, args []string, cfg config.Config, msgCh chan<- channels.MsgCh) {
	if slices.Contains(m.Member.Roles, cfg.Discord.AdminRoleId) {
		var returnString string
		msgCh <- channels.MsgCh{Err: nil, Message: "Message Event - Jukebox Skip - Triggered", Level: "INFO"}
		jukeboxClient := jukebox.New(cfg)
		if err := jukeboxClient.Skip(); err != nil {
			msgCh <- channels.MsgCh{Err: err, Message: "Something went wrong", Level: "ERROR"}
			returnString = "There was a error connecting to the API"
		} else {
			returnString = "Skipping Track on Jukebox"
		}
		s.ChannelMessageSend(m.ChannelID, returnString)
	}
}
