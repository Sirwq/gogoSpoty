# gogoSpoty

A lightweight Go backend that polls the Spotify API for the currently playing track and exposes it as a browser overlay widget — useful for streaming or screen recording.

## How it works

1. On first launch, the app starts an OAuth flow and prints an authorization URL in the terminal.
2. After you authorize in the browser, the token is saved locally — subsequent runs skip the login step.
3. A background goroutine polls `PlayerCurrentlyPlaying` every 5 seconds and keeps an in-memory `Track` struct up to date.
4. A small HTTP server serves the widget HTML and a JSON API endpoint.

## Endpoints

|        Route        |             Description           |
|---|---|
| `GET /widget`       | Serves the HTML overlay widget    |
| `GET /api/current`  | Returns the current track as JSON |
| `GET /callback`     | OAuth redirect handler            |
| `/static/`          | Static assets (JS, CSS, images)   |

## Project structure

```
gogoSpoty/
├── main.go          # Entry point, server setup, polling loop
├── spoty/
│   ├── auth.go      # OAuth flow, token save/load
│   ├── handlers.go  # HTTP handlers
│   ├── track.go     # Track struct with mutex, update logic
│   └── *_test.go    # Tests
├── botik/           # (in progress)
├── static/
│   ├── widget.html  # Overlay UI
│   ├── script.js    # Fetches /api/current and updates the DOM
│   ├── styles.css
│   └── placeholder.png
├── go.mod
└── go.sum
```

## Setup

**Prerequisites:** Go 1.21+, a Spotify Developer account.

1. Create an app at [developer.spotify.com](https://developer.spotify.com/dashboard) and set the redirect URI to `http://127.0.0.1:5111/callback`.

2. Create a `.env` file in the project root:
   ```
   SPOTIFY_ID=your_client_id
   SPOTIFY_SECRET=your_client_secret
   ```

3. Run:
   ```bash
   go run main.go
   ```

4. Open the printed URL in your browser to authorize. After that, the widget is available at `http://127.0.0.1:5111/widget`.

## Dependencies

- [`zmb3/spotify`](https://github.com/zmb3/spotify) — Spotify Web API client
- [`joho/godotenv`](https://github.com/joho/godotenv) — `.env` loading
- `golang.org/x/oauth2` — OAuth2 token management
