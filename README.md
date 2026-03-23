# gogoSpoty

Spotify "Now Playing" OBS widget with Twitch chat song requests.

Shows the current track, artist, album art, and progress bar as a browser source in OBS. Twitch viewers can request songs via `!sr` command in chat — requests are queued and automatically added to Spotify playback.

## Demo

![gogoSpoty widget](docs/demo.gif)

## Features

- Real-time OBS widget with track info, album cover, and progress bar
- Twitch chat integration (`!sr <song name>` to request tracks)
- Song request queue with per-user cooldowns
- Graceful shutdown (SIGINT/SIGTERM)
- Standalone mode — single binary, no external dependencies

## Requirements

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
├── static/              — widget HTML/CSS/JS (embedded into binary)
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── go.mod
```

## Installation

### Option 1: Download Binary (easiest)

No dependencies required. Everything is bundled into a single executable.

1. Download the binary for your platform from [Releases](https://github.com/Sirwq/gogoSpoty/releases)
2. Download [.env.example](https://github.com/Sirwq/gogoSpoty/blob/main/.env.example) and rename it to `.env`
3. Fill in your Spotify and Twitch credentials
4. Run:

```bash
# Linux / macOS
chmod +x gogoSpoty-linux   # or gogoSpoty-macos
./gogoSpoty-linux

# Windows
gogoSpoty.exe
```

That's it. No Redis, no Docker, no Go installation needed.

### Option 2: Docker Compose

Uses Redis for persistent song queue (survives restarts).

```bash
git clone https://github.com/Sirwq/gogoSpoty.git
cd gogoSpoty
cp .env.example .env
# Edit .env with your credentials, set REDIS_ADDR=redis:6379
docker compose up --build
```

To stop:

```bash
docker compose down
```

### Option 3: Build from Source

```bash
git clone https://github.com/Sirwq/gogoSpoty.git
cd gogoSpoty
cp .env.example .env
# Edit .env with your credentials

# Standalone (in-memory queue, no Redis)
go build -tags standalone -o gogoSpoty ./cmd/gogoSpoty/

# With Redis
# Set REDIS_ADDR=localhost:6379 in .env, start Redis, then:
make run
```

## Configuration

Edit `.env` with your credentials. See [.env.example](https://github.com/Sirwq/gogoSpoty/blob/main/.env.example) for all available options.

Redirect URLs in `.env` must match those configured in your Spotify Dashboard and Twitch Developer Console.

On first launch, the app will print OAuth URLs for Spotify and Twitch — open them in a browser to authorize. Tokens are saved locally and reused on next start.

## Usage

### OBS Widget

Add a **Browser Source** in OBS with URL:

```
http://localhost:5111/static/widget.html
```

Recommended size: 500×150. Set background to transparent.

### Song Requests

Viewers type in Twitch chat:

```
!sr never gonna give you up
```

The bot searches Spotify, adds the first result to the queue, and confirms in chat. Songs are queued to Spotify playback automatically.

## Standalone vs Docker

| | Standalone binary | Docker Compose |
|---|---|---|
| Setup | Download + .env | Clone + .env + docker compose up |
| Dependencies | None | Docker |
| Song queue | In-memory (lost on restart) | Redis (persists across restarts) |
| Best for | Quick setup, personal use | Persistent queue, production use |

## Makefile

| Command | Description |
|---|---|
| `make build` | Build with Redis support |
| `make run` | Build and run with Redis |
| `make release` | Build standalone binaries for Linux, macOS, Windows |
| `make docker` | Start with Docker Compose |
| `make docker-down` | Stop Docker Compose |
| `make clean` | Remove all binaries |

## API

| Endpoint | Description |
|---|---|
| `GET /static/widget.html` | OBS browser source |
| `GET /api/current` | Current track JSON |
| `GET /static/*` | Static assets |

## License

MIT