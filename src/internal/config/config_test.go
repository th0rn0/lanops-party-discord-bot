package config_test

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	"lanops/party-discord-bot/internal/config"
)

var envFixture = map[string]string{
	"DISCORD_TOKEN":                    "tok",
	"DISCORD_ADMIN_ROLE_ID":            "role",
	"DISCORD_COMMAND_PREFIX":           "!",
	"LANOPS_STREAM_PROXY_API_USERNAME": "su",
	"LANOPS_STREAM_PROXY_API_PASSWORD": "sp",
	"LANOPS_STREAM_PROXY_API_ADDRESS":  "http://stream",
	"LANOPS_JUKEBOX_API_USERNAME":      "ju",
	"LANOPS_JUKEBOX_API_PASSWORD":      "jp",
	"LANOPS_JUKEBOX_API_URL":           "http://jukebox",
}

func TestLoad_AllVarsSet(t *testing.T) {
	for k, v := range envFixture {
		t.Setenv(k, v)
	}

	cfg := config.Load()

	if cfg.Discord.Token != "tok" {
		t.Errorf("Discord.Token: got %q, want %q", cfg.Discord.Token, "tok")
	}
	if cfg.Discord.AdminRoleId != "role" {
		t.Errorf("Discord.AdminRoleId: got %q, want %q", cfg.Discord.AdminRoleId, "role")
	}
	if cfg.Discord.CommandPrefix != "!" {
		t.Errorf("Discord.CommandPrefix: got %q, want %q", cfg.Discord.CommandPrefix, "!")
	}
	if cfg.Lanops.StreamProxyApiUsername != "su" {
		t.Errorf("StreamProxyApiUsername: got %q", cfg.Lanops.StreamProxyApiUsername)
	}
	if cfg.Lanops.StreamProxyApiPassword != "sp" {
		t.Errorf("StreamProxyApiPassword: got %q", cfg.Lanops.StreamProxyApiPassword)
	}
	if cfg.Lanops.StreamProxyApiAddress != "http://stream" {
		t.Errorf("StreamProxyApiAddress: got %q", cfg.Lanops.StreamProxyApiAddress)
	}
	if cfg.Lanops.JukeboxApiUsername != "ju" {
		t.Errorf("JukeboxApiUsername: got %q", cfg.Lanops.JukeboxApiUsername)
	}
	if cfg.Lanops.JukeboxApiPassword != "jp" {
		t.Errorf("JukeboxApiPassword: got %q", cfg.Lanops.JukeboxApiPassword)
	}
	if cfg.Lanops.JukeboxApiUrl != "http://jukebox" {
		t.Errorf("JukeboxApiUrl: got %q", cfg.Lanops.JukeboxApiUrl)
	}
}

// TestLoad_MissingVar_subprocess is the subprocess entry point called by TestLoad_MissingVars.
func TestLoad_MissingVar_subprocess(t *testing.T) {
	varName := os.Getenv("TEST_MISSING_VAR")
	if varName == "" {
		t.Skip("subprocess only")
	}
	config.Load()
}

// baseEnv returns a clean environment with all fixture vars except skipVar.
func baseEnv(skipVar string) []string {
	var env []string
	for _, e := range os.Environ() {
		key := strings.SplitN(e, "=", 2)[0]
		if strings.HasPrefix(key, "DISCORD_") || strings.HasPrefix(key, "LANOPS_") {
			continue
		}
		env = append(env, e)
	}
	for k, v := range envFixture {
		if k != skipVar {
			env = append(env, k+"="+v)
		}
	}
	return env
}

func TestLoad_MissingVars(t *testing.T) {
	for varName := range envFixture {
		varName := varName
		t.Run(varName, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=TestLoad_MissingVar_subprocess")
			cmd.Env = append(baseEnv(varName), "TEST_MISSING_VAR="+varName)
			err := cmd.Run()
			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				return
			}
			t.Fatalf("expected non-zero exit for missing %s, got %v", varName, err)
		})
	}
}
