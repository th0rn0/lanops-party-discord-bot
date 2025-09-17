package jukebox

import (
	"net/http"
)

type Client struct {
	username string
	password string
	url      string
	http     http.Client
}

type GetCurrentTrackOutput struct {
	Name   string `json:"name"`
	Artist string `json:"artist"`
}

type GetVolumeOutput struct {
	Volume int `json:"volume"`
}
