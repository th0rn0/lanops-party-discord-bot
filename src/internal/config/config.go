package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() Config {
	godotenv.Load()
	// Discord
	discordToken := os.Getenv("DISCORD_TOKEN")
	if discordToken == "" {
		log.Fatal("❌ DISCORD_TOKEN not set in environment")
	}

	discordGuildId := os.Getenv("DISCORD_SERVER_ID")
	if discordGuildId == "" {
		log.Fatal("❌ DISCORD_SERVER_ID not set in environment")
	}

	discordAdminRoleId := os.Getenv("DISCORD_ADMIN_ROLE_ID")
	if discordAdminRoleId == "" {
		log.Fatal("❌ DISCORD_ADMIN_ROLE_ID not set in environment")
	}

	discordCommandPrefix := os.Getenv("DISCORD_COMMAND_PREFIX")
	if discordCommandPrefix == "" {
		log.Fatal("❌ DISCORD_COMMAND_PREFIX not set in environment")
	}

	// Stream Proxy
	lanopsStreamProxyApiAddress := os.Getenv("LANOPS_STREAM_PROXY_API_ADDRESS")
	if lanopsStreamProxyApiAddress == "" {
		log.Fatal("❌ LANOPS_STREAM_PROXY_API_ADDRESS not set in environment")
	}

	return Config{
		DiscordCommandPrefix:        discordCommandPrefix,
		DiscordToken:                discordToken,
		DiscordGuildId:              discordGuildId,
		DiscordAdminRoleId:          discordAdminRoleId,
		LanopsStreamProxyApiAddress: lanopsStreamProxyApiAddress,
	}
}
