# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

All commands run from the `src/` directory:

```bash
# Install dependencies
cd src && go mod tidy

# Run the bot
cd src && go run ./cmd/discord-bot

# Build
cd src && go build ./cmd/discord-bot

# Run tests
cd src && go test ./...

# Run a single package's tests
cd src && go test ./internal/jukebox/...
```

Docker:
```bash
docker build -f resources/docker/Dockerfile .
```

## Architecture

The bot is a prefix-based Discord bot (e.g. `!stream list`) that acts as an HTTP client to two external LanOps APIs: the **Stream Proxy API** and the **Jukebox API**. All credentials and URLs come from environment variables (see `src/.env.example`).

### Request flow

1. `cmd/discord-bot/main.go` — entry point; loads config, creates `discordgo` session, wires up `bot.Client`
2. `internal/bot/main.go` — registers `OnMessage` and `OnReady` handlers on the Discord session
3. `internal/bot/handlers/router.go` — `OnMessage` strips the command prefix, then does a longest-prefix match against the handler registry (e.g. `"jukebox volume"` matches before `"jukebox"`)
4. `internal/bot/handlers/registry.go` — maps command strings to `HandlerFunc` implementations; add new commands here
5. Handler packages under `internal/bot/handlers/` — each handler checks for admin role membership (`cfg.Discord.AdminRoleId`) before acting, calls the relevant API client, then sends a Discord reply
6. API clients — `internal/jukebox/` and `internal/streams/` — thin HTTP wrappers using basic auth

### Adding a new command

1. Create a new handler package under `internal/bot/handlers/<service>/<command>/`
2. Implement `func Handler(s, m, commandParts, args, cfg, msgCh)` matching `handlers.HandlerFunc`
3. Register it in `internal/bot/handlers/registry.go` `init()`

### Key conventions

- Admin-only commands guard with `slices.Contains(m.Member.Roles, cfg.Discord.AdminRoleId)`
- Async log/error reporting goes through `msgCh chan<- channels.MsgCh` rather than direct logging inside handlers
- Logging uses `zerolog`; structured logging at the `main.go` level via the message channel
- `gorm.io/gorm` is a dependency but not yet used — likely planned for a future feature

## gstack

Use the `/browse` skill from gstack for all web browsing. Never use `mcp__claude-in-chrome__*` tools.

Available gstack skills:
`/office-hours`, `/plan-ceo-review`, `/plan-eng-review`, `/plan-design-review`, `/design-consultation`, `/design-shotgun`, `/design-html`, `/review`, `/ship`, `/land-and-deploy`, `/canary`, `/benchmark`, `/browse`, `/connect-chrome`, `/qa`, `/qa-only`, `/design-review`, `/setup-browser-cookies`, `/setup-deploy`, `/retro`, `/investigate`, `/document-release`, `/codex`, `/cso`, `/autoplan`, `/plan-devex-review`, `/devex-review`, `/careful`, `/freeze`, `/guard`, `/unfreeze`, `/gstack-upgrade`, `/learn`
