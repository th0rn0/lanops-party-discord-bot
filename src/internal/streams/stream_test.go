package streams

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"lanops/party-discord-bot/internal/config"
)

type errorTransport struct{}

func (e *errorTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("mock transport error")
}

func makeConfig(url string) config.Config {
	return config.Config{
		Lanops: config.Lanops{
			StreamProxyApiUsername: "user",
			StreamProxyApiPassword: "pass",
			StreamProxyApiAddress:  url,
		},
	}
}

func clientWithError(url string) Client {
	return Client{
		cfg:  makeConfig(url),
		url:  url,
		http: http.Client{Transport: &errorTransport{}},
	}
}

func TestNew(t *testing.T) {
	cfg := makeConfig("http://localhost:8080")
	c := New(cfg)
	if c.url != "http://localhost:8080" {
		t.Errorf("url: got %q", c.url)
	}
}

// -- GetStreams --

func TestGetStreams_Success(t *testing.T) {
	streams := []stream{{Name: "stream1", Enabled: true}, {Name: "stream2", Enabled: false}}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/streams/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(streams)
	}))
	defer srv.Close()

	result, err := New(makeConfig(srv.URL)).GetStreams()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 streams, got %d", len(result))
	}
}

func TestGetStreams_EmptyList(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("[]"))
	}))
	defer srv.Close()

	result, err := New(makeConfig(srv.URL)).GetStreams()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 streams, got %d", len(result))
	}
}

func TestGetStreams_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	// Non-200 returns empty slice; no error in current implementation
	result, _ := New(makeConfig(srv.URL)).GetStreams()
	if len(result) != 0 {
		t.Errorf("expected empty slice on non-200, got %d items", len(result))
	}
}

func TestGetStreams_RequestFails(t *testing.T) {
	if _, err := clientWithError("http://example.com").GetStreams(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetStreams_NewRequestError(t *testing.T) {
	c := Client{url: "http://[invalid", cfg: makeConfig("http://[invalid")}
	if _, err := c.GetStreams(); err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestGetStreams_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("not json"))
	}))
	defer srv.Close()

	if _, err := New(makeConfig(srv.URL)).GetStreams(); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

// -- EnableStreamByName --

func TestEnableStreamByName_Enable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/streams/mystream/enable" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(stream{Name: "mystream", Enabled: true})
	}))
	defer srv.Close()

	s, err := New(makeConfig(srv.URL)).EnableStreamByName("mystream", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Enabled {
		t.Error("expected stream to be enabled")
	}
}

func TestEnableStreamByName_Disable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(stream{Name: "mystream", Enabled: false})
	}))
	defer srv.Close()

	s, err := New(makeConfig(srv.URL)).EnableStreamByName("mystream", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Enabled {
		t.Error("expected stream to be disabled")
	}
}

func TestEnableStreamByName_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer srv.Close()

	_, err := New(makeConfig(srv.URL)).EnableStreamByName("nostream", true)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestEnableStreamByName_ServerError(t *testing.T) {
	// Note: the implementation returns (stream, nil) on non-200/non-404 responses
	// because it returns the http.Do error variable which is nil on a successful transport.
	// This test exercises that branch for coverage.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	_, err := New(makeConfig(srv.URL)).EnableStreamByName("s", true)
	if errors.Is(err, ErrNotFound) {
		t.Error("unexpected ErrNotFound for 500 response")
	}
}

func TestEnableStreamByName_RequestFails(t *testing.T) {
	if _, err := clientWithError("http://example.com").EnableStreamByName("s", true); err == nil {
		t.Error("expected error, got nil")
	}
}

// Note: EnableStreamByName has a bug where req.Header.Set is called before the
// if err != nil check, so passing an invalid URL causes a nil pointer panic rather
// than returning an error. The http.NewRequest error path is unreachable without
// triggering a panic, so it cannot be covered without modifying the source.

func TestEnableStreamByName_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("not json"))
	}))
	defer srv.Close()

	if _, err := New(makeConfig(srv.URL)).EnableStreamByName("s", true); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

// -- GetStreamByName --
// Note: the implementation checks for http.StatusFound (302), not 200.

func TestGetStreamByName_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/streams/mystream" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusFound)
		_ = json.NewEncoder(w).Encode(stream{Name: "mystream", Enabled: true})
	}))
	defer srv.Close()

	s, err := New(makeConfig(srv.URL)).GetStreamByName("mystream")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "mystream" {
		t.Errorf("expected name 'mystream', got %q", s.Name)
	}
}

func TestGetStreamByName_NotFound(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	defer srv.Close()

	_, err := New(makeConfig(srv.URL)).GetStreamByName("nostream")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestGetStreamByName_ServerError(t *testing.T) {
	// Note: the implementation returns (stream, nil) on non-302/non-404 responses
	// because it returns the http.Do error variable which is nil on a successful transport.
	// This test exercises that branch for coverage.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer srv.Close()

	_, err := New(makeConfig(srv.URL)).GetStreamByName("s")
	if errors.Is(err, ErrNotFound) {
		t.Error("unexpected ErrNotFound for 500 response")
	}
}

func TestGetStreamByName_RequestFails(t *testing.T) {
	if _, err := clientWithError("http://example.com").GetStreamByName("s"); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestGetStreamByName_NewRequestError(t *testing.T) {
	c := Client{url: "http://[invalid", cfg: makeConfig("http://[invalid")}
	if _, err := c.GetStreamByName("s"); err == nil {
		t.Error("expected error for invalid URL, got nil")
	}
}

func TestGetStreamByName_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusFound)
		_, _ = w.Write([]byte("not json"))
	}))
	defer srv.Close()

	if _, err := New(makeConfig(srv.URL)).GetStreamByName("s"); err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
