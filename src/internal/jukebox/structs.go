package jukebox

type API struct {
	URL string
}

type GetCurrentTrackOutput struct {
	Artists []struct {
		Name         string `json:"name"`
		ID           string `json:"id"`
		URI          string `json:"uri"`
		Href         string `json:"href"`
		ExternalUrls struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
	} `json:"artists"`
	Name  string `json:"name"`
	Album struct {
		Images []struct {
			Height int    `json:"height"`
			Width  int    `json:"width"`
			URL    string `json:"url"`
		} `json:"images"`
	} `json:"album"`
}
