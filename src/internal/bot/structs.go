package bot

import (
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"

	"github.com/bwmarrin/discordgo"
)

type Client struct {
	cfg config.Config
	dg  *discordgo.Session
	// streams streams.Client
	msgCh chan<- channels.MsgCh
}
