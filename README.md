# LanOps Discord Party Bot

Discord bot for interacting with various on site services at events.

Intended for use at [LanOps Events](https://www.lanops.co.uk)

## Prerequisites

- ```cp src/.env.example src/.env``` and fill it in

#### Install Dependencies

```bash
cd src
go mod tidy
```

## Usage

Entry Point:
```bash
go run ./cmd/discord-bot
```

## Bot Commands

| Command           | Input               | Description                        |
|-------------------|---------------------|------------------------------------|
| `stream list`     | None                | Lists all available streams.       |
| `stream enable`   | stream name         | Enables a specific stream.         |
| `stream disable`  | stream name         | Disables a specific stream.        |
| `jukebox start`   | None                | Starts the jukebox playback.       |
| `jukebox stop`    | None                | Stops the jukebox playback.        |
| `jukebox pause`   | None                | Pauses the current playback.       |
| `jukebox skip`    | None                | Skips to the next track.           |
| `jukebox volume`  | number from 0 - 100 | Adjusts the jukebox volume.        |
| `jukebox current` | None                | Shows the currently playing track. |

## Env

| Variable                          | Description                                          |
|-----------------------------------|------------------------------------------------------|
| `DISCORD_TOKEN`                   | Token for authenticating the Discord bot.            |
| `DISCORD_SERVER_ID`               | ID of the Discord server (guild) to connect to.      |
| `DISCORD_ADMIN_ROLE_ID`           | Role ID that grants admin permissions in Discord.    |
| `DISCORD_COMMAND_PREFIX`          | Prefix used for bot commands in Discord (e.g., `!`). |
| `LANOPS_STREAM_PROXY_API_ADDRESS` | Address of the LANOPS Stream Proxy API.              |
| `LANOPS_JUKEBOX_API_USERNAME`     | Username for authenticating with the jukebox API.    |
| `LANOPS_JUKEBOX_API_PASSWORD`     | Password for authenticating with the jukebox API.    |
| `LANOPS_JUKEBOX_API_URL`          | Base URL for the LANOPS jukebox API.                 |

## Docker