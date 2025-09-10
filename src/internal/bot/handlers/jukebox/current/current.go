package current

import (
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"
	"lanops/party-discord-bot/internal/jukebox"

	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate, commandParts []string, args []string, cfg config.Config, msgCh chan<- channels.MsgCh) {
	var returnString string
	msgCh <- channels.MsgCh{Err: nil, Message: "Message Event - Jukebox - Triggered", Level: "INFO"}
	jukeboxClient := jukebox.New(cfg)
	returnString, err := jukeboxClient.GetCurrentTrack()
	if err != nil {
		msgCh <- channels.MsgCh{Err: err, Message: "Something went wrong", Level: "ERROR"}
		returnString = "There was a error connecting to the API"
	}
	s.ChannelMessageSend(m.ChannelID, returnString)
}
