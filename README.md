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

| Variable                           | Description                                            |
|------------------------------------|--------------------------------------------------------|
| `DISCORD_TOKEN`                    | Token for authenticating the Discord bot.              |
| `DISCORD_SERVER_ID`                | ID of the Discord server (guild) to connect to.        |
| `DISCORD_ADMIN_ROLE_ID`            | Role ID that grants admin permissions in Discord.      |
| `DISCORD_COMMAND_PREFIX`           | Prefix used for bot commands in Discord (e.g., `!`).   |
| `LANOPS_STREAM_PROXY_API_USERNAME` | Username for authenticating with the Stream Proxy API. |
| `LANOPS_STREAM_PROXY_API_PASSWORD` | Password for authenticating with the Stream Proxy API. |
| `LANOPS_STREAM_PROXY_API_ADDRESS`  | Address of the LANOPS Stream Proxy API.                |
| `LANOPS_JUKEBOX_API_USERNAME`      | Username for authenticating with the jukebox API.      |
| `LANOPS_JUKEBOX_API_PASSWORD`      | Password for authenticating with the jukebox API.      |
| `LANOPS_JUKEBOX_API_URL`           | Base URL for the LANOPS jukebox API.                   |

## Docker

```docker build -f resources/docker/Dockerfile .```

```
docker run -d \
  --name party-discord-bot \
  --restart unless-stopped \
  -e DISCORD_TOKEN= \
  -e DISCORD_SERVER_ID= \
  -e DISCORD_ADMIN_ROLE_ID= \
  -e DISCORD_COMMAND_PREFIX=! \
  -e LANOPS_STREAM_PROXY_API_USERNAME= \
  -e LANOPS_STREAM_PROXY_API_PASSWORD= \
  -e LANOPS_STREAM_PROXY_API_ADDRESS=http://localhost:8080 \
  -e LANOPS_JUKEBOX_API_USERNAME= \
  -e LANOPS_JUKEBOX_API_PASSWORD= \
  -e LANOPS_JUKEBOX_API_URL=http://localhost:9999 \
  -p 8888:8888 \
  -v /mnt/servdata/lanops/jukebox/db:/db \
  th0rn0/lanops-party-discord-bot:latest
```

```
  party-discord-bot:
    image: th0rn0/lanops-party-discord-bot:latest
    container_name: party-discord-bot
    restart: unless-stopped
    environment:
      DISCORD_TOKEN: 
      DISCORD_SERVER_ID: 
      DISCORD_ADMIN_ROLE_ID: 
      DISCORD_COMMAND_PREFIX: "!"
      LANOPS_STREAM_PROXY_API_USERNAME: 
      LANOPS_STREAM_PROXY_API_PASSWORD:
      LANOPS_STREAM_PROXY_API_ADDRESS: "http://localhost:8080"
      LANOPS_JUKEBOX_API_USERNAME: 
      LANOPS_JUKEBOX_API_PASSWORD: 
      LANOPS_JUKEBOX_API_URL: "http://localhost:9999"
```