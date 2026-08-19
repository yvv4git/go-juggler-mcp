# go-juggler-mcp

MCP-compatible server that lets you control a browser programmatically via the
[Juggler](https://github.com/yvv4git/go-juggler) protocol. It exposes browser
automation as MCP tools, so any MCP client (opencode, Claude Desktop, Cursor,
etc.) can open tabs, inspect pages, click, type and take screenshots.

## Requirements

- A running [camofox-browser](https://github.com/daijro/camofox) endpoint
  (default `http://localhost:9377`).

## Install

```sh
go build -o go-juggler-mcp .
```

## Usage

```sh
./go-juggler-mcp serve
./go-juggler-mcp serve --addr http://localhost:9377 --session demo
./go-juggler-mcp serve --transport http --http-addr :8080
```

Default transport is `stdio`, compatible with the majority of MCP clients.
For remote access use `--transport http` (streamable HTTP).

### Connect from opencode

Add to your `opencode.json`:

```json
{
  "mcpServers": {
    "juggler": {
      "command": "/path/to/go-juggler-mcp",
      "args": ["serve"]
    }
  }
}
```

### Use in opencode

Replace `-m "opencode/deepseek-v4-flash-free"` with your model. The word
`juggler` in the prompt tells the model which MCP server to use.

```bash
# Health and raw result
opencode run -m "opencode/deepseek-v4-flash-free" "Check juggler health and print the raw JSON response from the tool"

# Open a page and show the resulting tab id / title
opencode run -m "opencode/deepseek-v4-flash-free" "Open https://rutube.ru with juggler and show the raw open_tab result (tab_id, title, url)"

# Snapshot: page structure with element refs
opencode run -m "opencode/deepseek-v4-flash-free" "Open https://rutube.ru with juggler, take a snapshot, and summarize the page structure: main blocks, headings, navigation, and the first 10 clickable elements with their refs"

# Navigation: history back / forward / refresh
opencode run -m "opencode/deepseek-v4-flash-free" "Open https://rutube.ru with juggler, navigate to https://google.com, go back, go forward, refresh, then show the current URL and the stats for that tab"

# Tabs: list and manage
opencode run -m "opencode/deepseek-v4-flash-free" "Use juggler to open https://rutube.ru and https://www.google.com, list all tabs, then close the google tab"

# Interaction: click + type + press
opencode run -m "opencode/deepseek-v4-flash-free" "Open https://www.google.com with juggler, snapshot, find the search input, type 'golang juggler' into it, press Enter, wait, then show the final URL and list the top 5 result links"

# Click an element by ref
opencode run -m "opencode/deepseek-v4-flash-free" "Open https://rutube.ru with juggler, snapshot, click the first video element by its ref, wait, then evaluate document.title"

# Scroll and collect links with pagination
opencode run -m "opencode/deepseek-v4-flash-free" "Open https://rutube.ru with juggler, scroll down 3000px, then list all links on the page with limit 20"

# Extract data with evaluate (no screenshots needed)
opencode run -m "opencode/deepseek-v4-flash-free" "Open https://rutube.ru with juggler, then use evaluate to return a JSON object with: page title, meta description, number of images, number of links, and performance metrics (page load time)"

# Extract structured data from the DOM
opencode run -m "opencode/deepseek-v4-flash-free" "Open https://rutube.ru with juggler, use evaluate to extract all video titles and thumbnails visible on the page as a JSON array"

# Network analysis
opencode run -m "opencode/deepseek-v4-flash-free" "Open https://pixelscan.net/bot-check with juggler, then use network_requests to list all API calls (XHR/fetch), their status codes and response sizes"

# End-to-end scenario
opencode run -m "opencode/deepseek-v4-flash-free" "Use juggler to: open https://rutube.ru, snapshot the page, click the search icon, type 'news' in the search box, press Enter, wait for results, then evaluate document.title and extract the top 5 result titles"
```

Note: `screenshot` returns a base64 image inside the MCP tool result; it is not
saved to disk. Only use it with a vision-capable model, otherwise the image is
discarded after the response.

## Configuration

Configuration is loaded from a TOML file (`--config`/`-c`, default
`config.toml`). See [config.example.toml](config.example.toml):

```toml
[log]
level = "info"

[serve]
addr = "http://localhost:9377"
session = ""
transport = "stdio"
http_addr = ":8080"
```

Every setting can be overridden with a CLI flag:
`--addr`, `--session`, `--transport`, `--http-addr`, `--log-level`.

The session key is auto-generated at startup when left empty. All tools use
this session by default; each tool accepts an optional `session` argument to
override it.

## Tools

| Tool                  | Description                                          |
|-----------------------|------------------------------------------------------|
| `health`              | Check browser status (engine, connection)            |
| `open_tab`            | Open a new tab and navigate to a URL                 |
| `navigate`            | Load a URL in an existing tab                        |
| `snapshot`            | Get the ARIA tree of the page (element refs)         |
| `click`               | Click an element by ref or CSS selector              |
| `type`                | Fill an input field by ref or selector               |
| `press`               | Press a keyboard key (Enter, Tab, Escape, etc.)      |
| `scroll`              | Scroll the page up/down by N pixels                  |
| `back`                | Navigate back in history                             |
| `forward`             | Navigate forward in history                          |
| `refresh`             | Reload the current page                              |
| `links`               | List all links on the page with pagination           |
| `screenshot`          | Take a PNG screenshot (page or viewport)             |
| `evaluate`            | Run arbitrary JavaScript in the page context         |
| `network_requests`    | Get all loaded resources                             |
| `stats`               | Get tab state (URL, visited URLs, refs)              |
| `list_tabs`           | List all tabs in a session                           |
| `close_tab`           | Close a tab                                          |
| `close_session`       | Destroy a session and all its tabs                   |

## Layout

```text
go-juggler-mcp/
├── main.go                   # cobra entry point
├── cmd/serve.go              # serve subcommand
├── internal/
│   ├── config/               # TOML configuration (cleanenv)
│   │   └── mcp/              # root Config struct and Load
│   ├── ports/                # BrowserClient interface
│   ├── adaptors/juggler/     # juggler client adapter
│   ├── core/                 # MCP server, tools and handlers
│   ├── domain/               # domain errors
│   └── infra/                # zap logger
└── config.example.toml
```

## License

MIT, see [LICENSE](LICENSE). See [NOTICE](NOTICE) for attribution.
