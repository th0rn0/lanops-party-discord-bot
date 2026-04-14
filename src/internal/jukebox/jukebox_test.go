package jukebox

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"lanops/party-discord-bot/internal/config"
)

// errorTransport always returns an error, simulating a network failure.
type errorTransport struct{}

func (e *errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("mock transport error")
}

func makeConfig(url string) config.Config {
	return config.Config{
		Lanops: config.Lanops{
			JukeboxApiUsername: "user",
			JukeboxApiPassword: "pass",
			JukeboxApiUrl:      url,
		},
	}
}

func clientWithError() Client {
	return Client{
		username: "user",
		password: "pass",
		url:      "http://example.com",
		http:     http.Client{Transport: &errorTransport{}},
	}
}

func TestNew(t *testing.T) {
	cfg := makeConfig("http://localhost:9999")
	c := New(cfg)
	if c.username != "user" || c.password != "pass" || c.url != "http://localhost:9999" {
		t.Errorf("New() fields incorrect: username=%q password=%q url=%q", c.username, c.password, c.url)
	}
}

// invalidURL triggers an http.NewRequest error (unparseable host).
const invalidURL = "http://[invalid"

// -- Start --

func TestStart_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/player/start" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	if err := New(makeConfig(srv.URL)).Start(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStart_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	if err := New(makeConfig(srv.URL)).Start(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestStart_RequestFails(t *testing.T) {
	if err := clientWithError().Start(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestStart_NewRequestError(t *testing.T) {
	if err := (Client{url: invalidURL}).Start(); err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

// -- Stop --

func TestStop_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/player/stop" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	if err := New(makeConfig(srv.URL)).Stop(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStop_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500) }))
	defer srv.Close()
	if err := New(makeConfig(srv.URL)).Stop(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestStop_RequestFails(t *testing.T) {
	if err := clientWithError().Stop(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestStop_NewRequestError(t *testing.T) {
	if err := (Client{url: invalidURL}).Stop(); err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

// -- Pause --

func TestPause_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/player/pause" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	if err := New(makeConfig(srv.URL)).Pause(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestPause_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500) }))
	defer srv.Close()
	if err := New(makeConfig(srv.URL)).Pause(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestPause_RequestFails(t *testing.T) {
	if err := clientWithError().Pause(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestPause_NewRequestError(t *testing.T) {
	if err := (Client{url: invalidURL}).Pause(); err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

// -- Skip --

func TestSkip_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/player/skip" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	if err := New(makeConfig(srv.URL)).Skip(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSkip_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500) }))
	defer srv.Close()
	if err := New(makeConfig(srv.URL)).Skip(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSkip_RequestFails(t *testing.T) {
	if err := clientWithError().Skip(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSkip_NewRequestError(t *testing.T) {
	if err := (Client{url: invalidURL}).Skip(); err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

// -- SetVolume --

func TestSetVolume_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/player/volume" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	if err := New(makeConfig(srv.URL)).SetVolume(50); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSetVolume_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500) }))
	defer srv.Close()
	if err := New(makeConfig(srv.URL)).SetVolume(50); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSetVolume_RequestFails(t *testing.T) {
	if err := clientWithError().SetVolume(50); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestSetVolume_NewRequestError(t *testing.T) {
	if err := (Client{url: invalidURL}).SetVolume(50); err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

// -- GetVolume --

func TestGetVolume_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/player/volume" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(GetVolumeOutput{Volume: 75})
	}))
	defer srv.Close()

	vol, err := New(makeConfig(srv.URL)).GetVolume()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if vol != 75 {
		t.Errorf("expected volume 75, got %d", vol)
	}
}

func TestGetVolume_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500) }))
	defer srv.Close()
	if _, err := New(makeConfig(srv.URL)).GetVolume(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetVolume_RequestFails(t *testing.T) {
	if _, err := clientWithError().GetVolume(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetVolume_NewRequestError(t *testing.T) {
	if _, err := (Client{url: invalidURL}).GetVolume(); err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestGetVolume_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	if _, err := New(makeConfig(srv.URL)).GetVolume(); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

// -- GetCurrentTrack --

func TestGetCurrentTrack_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tracks/current" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(GetCurrentTrackOutput{Name: "Sandstorm", Artist: "Darude"})
	}))
	defer srv.Close()

	result, err := New(makeConfig(srv.URL)).GetCurrentTrack()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Currently Playing: Sandstorm - Darude" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestGetCurrentTrack_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(500) }))
	defer srv.Close()
	if _, err := New(makeConfig(srv.URL)).GetCurrentTrack(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetCurrentTrack_RequestFails(t *testing.T) {
	if _, err := clientWithError().GetCurrentTrack(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetCurrentTrack_NewRequestError(t *testing.T) {
	if _, err := (Client{url: invalidURL}).GetCurrentTrack(); err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestGetCurrentTrack_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	if _, err := New(makeConfig(srv.URL)).GetCurrentTrack(); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
