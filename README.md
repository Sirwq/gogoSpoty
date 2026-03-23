# gogoSpoty

![gogoSpoty widget](docs/demo.gif)

Spotify "Now Playing" OBS widget with Twitch chat song requests.

Shows the current track, artist, album art, and progress bar as a browser source in OBS. Twitch viewers can request songs via `!sr` command in chat — requests are queued in Redis and automatically added to Spotify playback.

## Features

- Real-time OBS widget with track info, album cover, and progress bar
- Twitch chat integration (`!sr <song name>` to request tracks)
- Redis-backed song queue
- Per-user cooldowns on song requests
- Graceful shutdown (SIGINT/SIGTERM)

## Requirements

- Docker and Docker Compose **or** Go 1.26+ with Redis
- Spotify Premium account
- Twitch account
- [Spotify Developer App](https://developer.spotify.com/dashboard) (Client ID + Secret)
- [Twitch Developer App](https://dev.twitch.tv/console) (Client ID + Secret)

## Project Structure

```
gogoSpoty/
├── cmd/gogoSpoty/       — entry point
├── internal/
│   ├── app/             — application assembly and lifecycle
│   ├── bot/             — Twitch bot, song queue, cooldowns, auth
│   ├── config/          — configuration (all env vars)
│   ├── crypto/          — random state generation for OAuth
│   ├── poller/          — Spotify playback polling and queue processing
│   └── widget/          — Spotify OAuth, HTTP server, track state, OBS widget handlers
├── static/              — widget HTML/CSS/JS and placeholder image
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

## Setup

### 1. Clone

```bash
git clone https://github.com/Sirwq/gogoSpoty.git
cd gogoSpoty
```

### 2. Configure

```bash
cp .env.example .env
```

Open `.env` and fill in your Spotify and Twitch credentials. Redirect URLs in `.env` must match those in your Spotify Dashboard and Twitch Developer Console.

### 3. Run

#### Docker Compose (recommended)

```bash
docker compose up --build
```

This starts both the app and Redis. No need to install Go or Redis locally.

To stop:

```bash
docker compose down
```

#### Without Docker

Requires Go 1.26+ and a running Redis instance. Set `REDIS_ADDR=localhost:6379` in `.env`.

```bash
make run
```

### 4. First Launch

On first launch, the app will print OAuth URLs for Spotify and Twitch — open them in a browser to authorize. Tokens are saved locally and reused on next start.

## Usage

### OBS Widget

Add a **Browser Source** in OBS with URL:

```
http://localhost:5111/widget
```

Recommended size: 500×150. Set background to transparent.

### Song Requests

Viewers type in Twitch chat:

```
!sr never gonna give you up
```

The bot searches Spotify, adds the first result to the queue, and confirms in chat. Songs are queued to Spotify playback automatically.

## Makefile

| Command | Description |
|---|---|
| `make build` | Build the binary |
| `make run` | Build and run |
| `make clean` | Remove the binary |
| `make docker` | Start with Docker Compose |
| `make docker-down` | Stop Docker Compose |

## API

| Endpoint | Description |
|---|---|
| `GET /widget` | OBS browser source |
| `GET /api/current` | Current track JSON |
| `GET /static/*` | Static assets |
