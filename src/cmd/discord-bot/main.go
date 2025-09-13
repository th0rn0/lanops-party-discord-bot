package main

import (
	"lanops/party-discord-bot/internal/bot"
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"

	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

var (
	logger zerolog.Logger
	cfg    config.Config
	msgCh  = make(chan channels.MsgCh, 20)
)

func main() {
	logger = zerolog.New(
		zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339},
	).Level(zerolog.TraceLevel).With().Timestamp().Caller().Logger()
	logger.Info().Msg("Initializing Party Discord Bot")

	logger.Info().Msg("Loading Config")
	cfg = config.Load()

	logger.Info().Msg("Starting Party Discord Bot")

	// Message Channel
	go func() {
		for msg := range msgCh {
			if msg.Err != nil {
				logger.Error().Err(msg.Err).Msg(msg.Message)
			} else {
				logger.Info().Msg(msg.Message)
			}
		}
	}()

	logger.Info().Msg("Starting Discord Bot")
	discordClient, err := discordgo.New("Bot " + cfg.Discord.Token)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create bot")
	}

	botClient, err := bot.New(cfg, discordClient, msgCh)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to create bot")
	}

	if err := botClient.Run(); err != nil {
		logger.Fatal().Err(err).Msg("Failed to start bot")
	}
	// TODO
	// fix get current track
}
