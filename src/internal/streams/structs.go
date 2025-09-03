package streams

import "lanops/party-discord-bot/internal/config"

type Client struct {
	cfg config.Config
	url string
}

type stream struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
