package streams

import (
	"lanops/party-discord-bot/internal/config"
	"net/http"
)

type Client struct {
	cfg  config.Config
	url  string
	http http.Client
}

type stream struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}
