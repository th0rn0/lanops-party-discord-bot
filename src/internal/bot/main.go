package bot

import (
	"lanops/party-discord-bot/internal/bot/handlers"
	"lanops/party-discord-bot/internal/channels"
	"lanops/party-discord-bot/internal/config"

	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func New(cfg config.Config, discordClient *discordgo.Session, msgCh chan<- channels.MsgCh) (*Client, error) {
	client := &Client{
		cfg:   cfg,
		dg:    discordClient,
		msgCh: msgCh,
	}

	// Register the Events for Discord Go
	client.dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		handlers.OnMessage(s, m, cfg, msgCh)
	})
	client.dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		handlers.OnReady(s, m, cfg, msgCh)
	})

	// Set the Intents
	client.dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	return client, nil
}

func (client *Client) Run() error {
	// Open the websocket
	if err := client.dg.Open(); err != nil {
		return err
	}
	defer client.dg.Close()

	// Wait for CTRL+C
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	return nil
}
