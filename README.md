# TRexT 🦖

A Terminal UI based REST API Client built with Go and [tview](https://github.com/rivo/tview). Similar to Postman or Insomnia, but right in your terminal!

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## Features

- 🚀 **HTTP Methods**: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- 📝 **Request Builder**: URL input, headers editor, body editor
- 📊 **Response Viewer**: Status, headers, formatted JSON body, timing
- 💾 **Persistence**: Save requests to SQLite database
- 📁 **Collections**: Organize requests in collections
- ⌨️ **Keyboard-driven**: Full keyboard navigation
- 🎨 **Theming**: Customizable color themes

## Installation

### From Source

```bash
git clone https://github.com/YashIIT0909/TRexT.git
cd TRexT
go build -o trext ./cmd/trext
./trext
```

### Go Install

```bash
go install github.com/YashIIT0909/TRexT/cmd/trext@latest
```

## Usage

Run the application:

```bash
./trext
```

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Navigate between panels |
| `Ctrl+Enter` | Send request |
| `Ctrl+N` | New request |
| `Ctrl+S` | Save request |
| `Ctrl+U` | Focus URL input |
| `Ctrl+H` | Focus collections (left) |
| `Ctrl+L` | Focus response (right) |
| `Ctrl+Q` | Quit |
| `n` | New request (in collections list) |
| `d` | Delete request (in collections list) |

### Layout

```
┌─────────────────┬──────────────────────────┬──────────────────────────┐
│  Collections    │       Request            │       Response           │
│                 │  ┌──────────────────┐    │  Status: 200 OK | 45ms   │
│  + New Request  │  │ GET ▼ │ URL...   │    │  ──────────────────────  │
│                 │  └──────────────────┘    │  {                       │
│  > GET /users   │  ┌─ Headers ────────┐    │    "data": "..."         │
│  > POST /login  │  │ Content-Type:... │    │  }                       │
│                 │  └──────────────────┘    │                          │
│                 │  ┌─ Body ───────────┐    │                          │
│                 │  │ {"key": "value"} │    │                          │
│                 │  └──────────────────┘    │                          │
│                 │     [Send Request]       │                          │
└─────────────────┴──────────────────────────┴──────────────────────────┘
│ Ctrl+Enter: Send | Ctrl+S: Save | Ctrl+N: New | Tab: Navigate         │
└───────────────────────────────────────────────────────────────────────┘
```

## Project Structure

```
TRexT/
├── cmd/trext/main.go           # Entry point
├── internal/
│   ├── app/
│   │   ├── app.go              # Main application logic
│   │   └── theme.go            # Color theming
│   ├── components/
│   │   ├── request_panel.go    # Request builder UI
│   │   ├── response_view.go    # Response display UI
│   │   ├── collections_list.go # Sidebar collections
│   │   └── dialogs.go          # Modal dialogs
│   ├── http/
│   │   ├── client.go           # HTTP client wrapper
│   │   ├── request.go          # Request model
│   │   └── response.go         # Response model
│   ├── storage/
│   │   ├── db.go               # SQLite database
│   │   ├── config.go           # YAML configuration
│   │   └── models.go           # Data models
│   └── utils/
│       └── json.go             # JSON utilities
├── configs/default.yaml        # Default configuration
├── go.mod
└── README.md
```

## Configuration

Config is stored at `~/.config/trext/config.yaml`:

```yaml
theme: default          # or "dracula"
defaultTimeout: 30      # seconds
sslVerify: true
history:
  maxItems: 100
  enabled: true
```

## Data Storage

Requests and history are stored in SQLite at `~/.config/trext/data.db`.

## Roadmap

- [ ] cURL import/export
- [ ] Environment variables
- [ ] Authentication helpers (Basic, Bearer, OAuth)
- [ ] Response syntax highlighting
- [ ] Request templates
- [ ] Export/Import collections
- [ ] WebSocket support

## License

MIT License
