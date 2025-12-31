# TRexT 🦖

A Terminal UI based REST API Client built with Go and [tview](https://github.com/rivo/tview). Similar to Postman or Insomnia, but right in your terminal!

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/License-MIT-green.svg)

## Features

- 🚀 **HTTP Methods**: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS
- 📝 **Request Builder**: URL input, headers editor, body editor
- 📊 **Response Viewer**: Status, headers, formatted JSON body, timing
- 💾 **Persistence**: Save requests to PostgreSQL database
- 📁 **Collections**: Organize requests in collections
- ⌨️ **Keyboard-driven**: Full keyboard navigation
- 🎨 **Theming**: Customizable color themes
- 🔄 **Type-safe SQL**: Uses [sqlc](https://sqlc.dev/) for generated database code
- 📦 **Migrations**: Database migrations with [goose](https://github.com/pressly/goose)

## Installation

### Prerequisites

- Go 1.21+
- PostgreSQL database

### From Source

```bash
git clone https://github.com/YashIIT0909/TRexT.git
cd TRexT

# Set up PostgreSQL database
cp .env.example .env
# Edit .env with your PostgreSQL connection string

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
│   │   ├── database.go         # PostgreSQL database connection
│   │   ├── config.go           # YAML configuration
│   │   ├── models.go           # Data models
│   │   └── db/                 # sqlc generated code
│   │       ├── db.go
│   │       ├── models.go
│   │       ├── collections.sql.go
│   │       ├── history.sql.go
│   │       └── requests.sql.go
│   └── utils/
│       └── json.go             # JSON utilities
├── sql/
│   ├── queries/                # SQL queries for sqlc
│   │   ├── collections.sql
│   │   ├── history.sql
│   │   └── requests.sql
│   └── schemas/                # Goose migrations
│       ├── 001_initial_schema.sql
│       └── embed.go
├── configs/default.yaml        # Default configuration
├── .env.example                # Example environment file
├── sqlc.yaml                   # sqlc configuration
├── Makefile                    # Build and database commands
├── go.mod
└── README.md
```

## Configuration

Config is stored at `~/.config/trext/config.yaml`:

```yaml
theme: default          # or "dracula"
defaultTimeout: 30      # seconds
sslVerify: true
proxy: ""
history:
  maxItems: 100
  enabled: true
keybindings:
  sendRequest: Ctrl+Enter
  newRequest: Ctrl+N
  saveRequest: Ctrl+S
  focusURL: Ctrl+U
```

## Database Setup

TRexT uses PostgreSQL for data storage. Set up your database connection using environment variables:

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your PostgreSQL connection string:
   ```
   DATABASE_URL=postgres://user:password@localhost:5432/trext?sslmode=disable
   ```

3. Migrations run automatically on startup using embedded goose migrations.

## Development

### Makefile Commands

| Command | Description |
|---------|-------------|
| `make build` | Build the application |
| `make run` | Build and run the application |
| `make sqlc` | Generate Go code from SQL queries |
| `make goose-up` | Run all pending migrations |
| `make goose-down` | Rollback the last migration |
| `make goose-status` | Show migration status |
| `make goose-create` | Create a new migration file |
| `make generate` | Generate all code (sqlc) |

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
