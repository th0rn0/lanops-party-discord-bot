package jukebox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lanops/party-discord-bot/internal/config"
	"net/http"
)

func New(cfg config.Config) Client {
	return Client{
		username: cfg.Lanops.JukeboxApiUsername,
		password: cfg.Lanops.JukeboxApiPassword,
		url:      cfg.Lanops.JukeboxApiUrl,
	}
}

func (c Client) Start() error {
	req, err := http.NewRequest("POST", c.url+"/player/start", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status code %d", resp.StatusCode))
	}
	return nil
}

func (c Client) Stop() error {
	req, err := http.NewRequest("POST", c.url+"/player/stop", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status code %d", resp.StatusCode))
	}
	return nil
}

func (c Client) Pause() error {
	req, err := http.NewRequest("POST", c.url+"/player/pause", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status code %d", resp.StatusCode))
	}
	return nil
}

func (c Client) Skip() error {
	req, err := http.NewRequest("POST", c.url+"/player/skip", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status code %d", resp.StatusCode))
	}
	return nil
}

func (c Client) Volume(volume int) error {
	var jsonData []byte
	params := map[string]interface{}{
		"volume": volume,
	}
	jsonData, _ = json.Marshal(params)
	req, err := http.NewRequest("POST", c.url+"/player/volume", bytes.NewReader(jsonData))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Status code %d", resp.StatusCode))
	}
	return nil
}

func (c Client) GetCurrentTrack() (returnString string, err error) {
	var getCurrentTrackOutput GetCurrentTrackOutput
	req, err := http.NewRequest("GET", c.url+"/tracks/current", nil)
	if err != nil {
		return returnString, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return returnString, err
	}
	if resp.StatusCode != 200 {
		return returnString, errors.New(fmt.Sprintf("Status code %d", resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &getCurrentTrackOutput)
	if err != nil {
		return returnString, err
	}
	return fmt.Sprintf("Currently Playing: %s - %s", getCurrentTrackOutput.Name, getCurrentTrackOutput.Artists[0].Name), nil
}
