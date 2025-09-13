package volume

import (
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"
	"lanops/party-discord-bot/internal/jukebox"
	"slices"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate, commandParts []string, args []string, cfg config.Config, msgCh chan<- channels.MsgCh) {
	if slices.Contains(m.Member.Roles, cfg.Discord.AdminRoleId) {
		var returnString string
		msgCh <- channels.MsgCh{Err: nil, Message: "Message Event - Jukebox Volume - Triggered", Level: "INFO"}
		jukeboxClient := jukebox.New(cfg)
		convertedInt, err := strconv.Atoi(args[0])
		if err != nil {
			returnString = "I can't do that brian"
		}
		if err := jukeboxClient.Volume(convertedInt); err != nil {
			msgCh <- channels.MsgCh{Err: err, Message: "Something went wrong", Level: "ERROR"}
			returnString = "There was a error connecting to the API"
		} else {
			returnString = "Changing Volume on Jukebox"
		}
		s.ChannelMessageSend(m.ChannelID, returnString)
	}
}
