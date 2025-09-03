package handlers

import (
	streamEnable "lanops/party-discord-bot/internal/bot/handlers/streams/enable"
	streamList "lanops/party-discord-bot/internal/bot/handlers/streams/list"
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"

	"github.com/bwmarrin/discordgo"
)

type HandlerFunc func(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	commandParts []string,
	args []string,
	cfg config.Config,
	msgCh chan<- channels.MsgCh,
)

var Registry = map[string]HandlerFunc{}

func Register(command string, handler HandlerFunc) {
	Registry[command] = handler
}

func init() {
	Register("stream list", streamList.Handler)
	Register("stream enable", streamEnable.Handler)
	Register("stream disable", streamEnable.Handler)
}
