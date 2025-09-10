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

	// Jukebox
	lanopsJukeboxApiUsername := os.Getenv("LANOPS_JUKEBOX_API_USERNAME")
	if lanopsJukeboxApiUsername == "" {
		log.Fatal("❌ LANOPS_JUKEBOX_API_USERNAME not set in environment")
	}
	lanopsJukeboxApiPassword := os.Getenv("LANOPS_JUKEBOX_API_PASSWORD")
	if lanopsJukeboxApiPassword == "" {
		log.Fatal("❌ LANOPS_JUKEBOX_API_PASSWORD not set in environment")
	}
	lanopsJukeboxApiUrl := os.Getenv("LANOPS_JUKEBOX_API_URL")
	if lanopsJukeboxApiUrl == "" {
		log.Fatal("❌ LANOPS_JUKEBOX_API_URL not set in environment")
	}
	discord := Discord{
		CommandPrefix: discordCommandPrefix,
		Token:         discordToken,
		AdminRoleId:   discordAdminRoleId,
	}
	lanops := Lanops{
		StreamProxyApiAddress: lanopsStreamProxyApiAddress,
		JukeboxApiUsername:    lanopsJukeboxApiUsername,
		JukeboxApiPassword:    lanopsJukeboxApiPassword,
		JukeboxApiUrl:         lanopsJukeboxApiUrl,
	}
	return Config{
		Discord: discord,
		Lanops:  lanops,
	}
}
