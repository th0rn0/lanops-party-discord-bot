package handlers

import (
	jukeboxCurrent "lanops/party-discord-bot/internal/bot/handlers/jukebox/current"
	jukeboxPause "lanops/party-discord-bot/internal/bot/handlers/jukebox/pause"
	jukeboxSkip "lanops/party-discord-bot/internal/bot/handlers/jukebox/skip"
	jukeboxStart "lanops/party-discord-bot/internal/bot/handlers/jukebox/start"
	jukeboxStop "lanops/party-discord-bot/internal/bot/handlers/jukebox/stop"
	jukeboxVolume "lanops/party-discord-bot/internal/bot/handlers/jukebox/volume"
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
	// Streams
	Register("stream list", streamList.Handler)
	Register("stream enable", streamEnable.Handler)
	Register("stream disable", streamEnable.Handler)
	// Jukebox
	Register("jukebox start", jukeboxStart.Handler)
	Register("jukebox stop", jukeboxStop.Handler)
	Register("jukebox pause", jukeboxPause.Handler)
	Register("jukebox skip", jukeboxSkip.Handler)
	Register("jukebox volume", jukeboxVolume.Handler)
	Register("jukebox current", jukeboxCurrent.Handler)

}
