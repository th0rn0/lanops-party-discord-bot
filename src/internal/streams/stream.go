package streams

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"lanops/party-discord-bot/internal/config"
	"net/http"
)

var (
	client      http.Client
	ErrNotFound = errors.New("Stream not found")
)

func New(cfg config.Config) Client {
	c := Client{
		url: cfg.LanopsStreamProxyApiAddress,
		cfg: cfg,
	}
	return c
}

func (c Client) GetStreams() (streams []stream, err error) {
	req, err := http.NewRequest("GET", c.url+"/streams/", nil)
	if err != nil {
		return streams, err
	}
	// req.SetBasicAuth(username, password)
	resp, err := client.Do(req)
	if err != nil {
		return streams, err
	}
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return streams, err
	}
	err = json.Unmarshal(body, &streams)
	if err != nil {
		return streams, err
	}
	return streams, nil
}

func (c Client) EnableStreamByName(name string, enabled bool) (stream stream, err error) {
	var jsonData []byte
	params := map[string]interface{}{
		"enabled": enabled,
	}
	jsonData, _ = json.Marshal(params)
	req, err := http.NewRequest("POST", c.url+"/streams/"+name+"/enable", bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return stream, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return stream, err
	}
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return stream, ErrNotFound
		} else {
			return stream, err
		}
	}
	err = json.Unmarshal(body, &stream)
	if err != nil {
		return stream, err
	}
	return stream, nil
}

func (c Client) GetStreamByName(name string) (stream stream, err error) {
	req, err := http.NewRequest("GET", c.url+"/streams/"+name, nil)
	if err != nil {
		return stream, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return stream, err
	}
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusFound {
		if resp.StatusCode == http.StatusNotFound {
			return stream, ErrNotFound
		} else {
			return stream, err
		}
	}
	err = json.Unmarshal(body, &stream)
	if err != nil {
		return stream, err
	}
	return stream, nil
}
