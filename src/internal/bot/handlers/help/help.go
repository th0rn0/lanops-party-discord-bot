package help

import (
	"fmt"
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"
	"slices"

	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate, commandParts []string, args []string, cfg config.Config, msgCh chan<- channels.MsgCh) {
	if slices.Contains(m.Member.Roles, cfg.Discord.AdminRoleId) {
		var commands = map[string]string{
			cfg.Discord.CommandPrefix + "help":                         "Show available commands",
			cfg.Discord.CommandPrefix + "stream enable <stream name>":  "Enable live stream",
			cfg.Discord.CommandPrefix + "stream disable <stream name>": "Disable live stream",
			cfg.Discord.CommandPrefix + "stream list":                  "List available scenes",
			cfg.Discord.CommandPrefix + "jukebox play":                 "Play a track",
			cfg.Discord.CommandPrefix + "jukebox pause":                "Pause the current track",
			cfg.Discord.CommandPrefix + "jukebox skip":                 "Skip the current track",
			cfg.Discord.CommandPrefix + "jukebox queue":                "Show the current queue",
			cfg.Discord.CommandPrefix + "jukebox stop":                 "Stop playback and clear the queue",
		}
		helpMsg := "**Available Commands:**\n"
		for cmd, desc := range commands {
			helpMsg += fmt.Sprintf("`%s` - %s\n", cmd, desc)
		}
		s.ChannelMessageSend(m.ChannelID, helpMsg)
	}
}
