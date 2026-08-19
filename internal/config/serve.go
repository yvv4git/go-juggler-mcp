package config

// Transport is the MCP server transport type.
type Transport string

const (
	// TransportStdio speaks MCP over stdin/stdout. Compatible with the
	// majority of MCP clients (opencode, Claude Desktop, Cursor, etc.).
	TransportStdio Transport = "stdio"

	// TransportHTTP speaks MCP over the streamable HTTP transport.
	TransportHTTP Transport = "http"
)

// Serve holds the MCP server and juggler browser settings.
type Serve struct {
	// Addr is the camofox-browser REST API endpoint.
	Addr string `toml:"addr"`

	// Session is the juggler session key. When empty, a unique one is
	// generated at startup.
	Session string `toml:"session"`

	// Transport selects the MCP transport: stdio (default) or http.
	Transport Transport `toml:"transport"`

	// HTTPAddr is the listen address used when Transport == http.
	HTTPAddr string `toml:"http_addr"`
}
