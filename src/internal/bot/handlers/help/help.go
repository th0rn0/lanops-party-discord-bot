package help

import (
	"fmt"
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"
	"slices"

	"github.com/bwmarrin/discordgo"
)

var commands = map[string]string{
	"!help":                         "Show available commands",
	"!stream enable <stream name>":  "Enable live stream",
	"!stream disable <stream name>": "Disable live stream",
	"!stream list":                  "List available scenes",
	"!jukebox play":                 "Play a track",
	"!jukebox pause":                "Pause the current track",
	"!jukebox skip":                 "Skip the current track",
	"!jukebox queue":                "Show the current queue",
	"!jukebox stop":                 "Stop playback and clear the queue",
}

func Handler(s *discordgo.Session, m *discordgo.MessageCreate, commandParts []string, args []string, cfg config.Config, msgCh chan<- channels.MsgCh) {
	if slices.Contains(m.Member.Roles, cfg.Discord.AdminRoleId) {
		helpMsg := "**Available Commands:**\n"
		for cmd, desc := range commands {
			helpMsg += fmt.Sprintf("`%s` - %s\n", cmd, desc)
		}
		s.ChannelMessageSend(m.ChannelID, helpMsg)
	}
}
