package enable

import (
	"errors"
	"fmt"
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"
	"lanops/party-discord-bot/internal/streams"
	"slices"

	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate, commandParts []string, args []string, cfg config.Config, msgCh chan<- channels.MsgCh) {
	if slices.Contains(m.Member.Roles, cfg.Discord.AdminRoleId) {
		var returnString string
		msgCh <- channels.MsgCh{Err: nil, Message: "Message Event - Stream - Triggered", Level: "INFO"}
		streamClient := streams.New(cfg)
		enable := true
		returnString = fmt.Sprintf("Enabled Stream: %s", args[0])
		if commandParts[1] == "disable" {
			enable = false
			returnString = fmt.Sprintf("Enabled Stream: %s", args[0])
		}
		_, err := streamClient.EnableStreamByName(args[0], enable)
		if err != nil {
			if errors.Is(err, streams.ErrNotFound) {
				returnString = fmt.Sprintf("Stream %s not found", args[0])
			} else {
				msgCh <- channels.MsgCh{Err: err, Message: "Something went wrong", Level: "ERROR"}
				returnString = "There was a error connecting to the API"
			}
		}
		s.ChannelMessageSend(m.ChannelID, returnString)
	}
}
