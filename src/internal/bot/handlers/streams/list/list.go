package List

import (
	"fmt"
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"
	"lanops/party-discord-bot/internal/streams"

	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate, commandParts []string, args []string, cfg config.Config, msgCh chan<- channels.MsgCh) {
	var returnString string
	msgCh <- channels.MsgCh{Err: nil, Message: "Message Event - Stream - Triggered", Level: "INFO"}
	streamsClient := streams.New(cfg)
	streams, err := streamsClient.GetStreams()
	if err != nil {
		msgCh <- channels.MsgCh{Err: err, Message: "Something went wrong", Level: "ERROR"}
		returnString = "There was a error connecting to the API"
	} else {
		returnString = "Available Streams:\n"
		for _, stream := range streams {
			returnString += fmt.Sprintf("Name: %s - Enabled: %t \n", stream.Name, stream.Enabled)
		}
		if len(streams) == 0 {
			returnString += "NO STREAMS AVAILABLE"
		}
	}
	s.ChannelMessageSend(m.ChannelID, returnString)
}
