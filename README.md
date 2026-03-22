# gogoSpoty

Spotify "Now Playing" OBS widget with Twitch chat song requests.

Shows the current track, artist, album art, and progress bar as a browser source in OBS. Twitch viewers can request songs via `!sr` command in chat — requests are queued in Redis and automatically added to Spotify playback.

## Features

- Real-time OBS widget with track info, album cover, and progress bar
- Twitch chat integration (`!sr <song name>` to request tracks)
- Redis-backed song queue
- Per-user cooldowns on song requests
- Graceful shutdown (SIGINT/SIGTERM)

## Requirements

- Go 1.22+
- Redis
- Spotify Premium account
- Twitch account
- [Spotify Developer App](https://developer.spotify.com/dashboard) (Client ID + Secret)
- [Twitch Developer App](https://dev.twitch.tv/console/apps) (Client ID + Secret)

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

Create a `.env` file in the project root:

```env
# Spotify
CLIENT_ID_SPOTY=your_spotify_client_id
CLIENT_SECRET_SPOTY=your_spotify_client_secret
REDIRECT_URL_SPOTY=http://127.0.0.1:5111/callback

# Twitch
TWITCH_USERNAME=your_bot_username
TWITCH_CHANNEL=your_channel
TWITCH_CLIENT_ID=your_twitch_client_id
TWITCH_CLIENT_SECRET=your_twitch_client_secret
TWITCH_REDIRECT_URL=http://localhost:6111/callback

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=your_redis_password
```

### 3. Start Redis

```bash
# Docker
docker run -d -p 6379:6379 redis

# Or locally
redis-server
```

### 4. Build & Run

```bash
make run
```

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

## API

| Endpoint | Description |
|---|---|
| `GET /widget` | OBS browser source |
| `GET /api/current` | Current track JSON |
| `GET /static/*` | Static assets |