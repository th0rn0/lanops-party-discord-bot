package help

import (
	"fmt"
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"
	"slices"
	"sort"

	"github.com/bwmarrin/discordgo"
)

func Handler(s *discordgo.Session, m *discordgo.MessageCreate, commandParts []string, args []string, cfg config.Config, msgCh chan<- channels.MsgCh) {
	if slices.Contains(m.Member.Roles, cfg.Discord.AdminRoleId) {
		helpMsg := "**Available Commands:**\n"
		helpMsg += formatCommands(map[string]string{
			cfg.Discord.CommandPrefix + "stream enable <stream name>":  "Enable live stream",
			cfg.Discord.CommandPrefix + "stream disable <stream name>": "Disable live stream",
			cfg.Discord.CommandPrefix + "stream list":                  "List available scenes",
			cfg.Discord.CommandPrefix + "jukebox play":                 "Play a track",
			cfg.Discord.CommandPrefix + "jukebox pause":                "Pause the current track",
			cfg.Discord.CommandPrefix + "jukebox skip":                 "Skip the current track",
			cfg.Discord.CommandPrefix + "jukebox queue":                "Show the current queue",
			cfg.Discord.CommandPrefix + "jukebox stop":                 "Stop playback and clear the queue",
		})
		s.ChannelMessageSend(m.ChannelID, helpMsg)
	}
}

func formatCommands(cmds map[string]string) string {
	keys := make([]string, 0, len(cmds))
	for k := range cmds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var out string
	for _, k := range keys {
		out += fmt.Sprintf("`%s` - %s\n", k, cmds[k])
	}
	return out
}
