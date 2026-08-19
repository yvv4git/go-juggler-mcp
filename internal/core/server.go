package core

import (
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/yvv4git/go-juggler-mcp/internal/ports"
	"go.uber.org/zap"
)

// toolDef couples a tool definition with its handler.
type toolDef struct {
	tool    mcp.Tool
	handler server.ToolHandlerFunc
}

// Handler wires the juggler browser client into an MCP server.
type Handler struct {
	client   ports.BrowserClient
	fs       ports.Filesystem
	session  string
	log      *zap.Logger
}

// New creates a Handler. An empty session generates a unique key at startup.
func New(client ports.BrowserClient, fs ports.Filesystem, session string, log *zap.Logger) *Handler {
	if session == "" {
		session = fmt.Sprintf("mcp-%d", time.Now().UnixNano())
	}

	return &Handler{
		client:  client,
		fs:      fs,
		session: session,
		log:     log,
	}
}

// MCPServer builds the MCP server with all tools registered.
func (h *Handler) MCPServer() *server.MCPServer {
	s := server.NewMCPServer(
		"go-juggler-mcp",
		"1.0.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
		server.WithInstructions(
			"Browser automation over the Juggler protocol. Start by calling "+
				"open_tab, then inspect the page with snapshot and act on elements "+
				"via click/type/press/scroll. Use list_tabs and stats to see open "+
				"tabs and their state. Screenshots are returned as images.",
		),
	)

	h.register(s, h.tools()...)

	return s
}

// tools returns every registered tool definition.
func (h *Handler) tools() []toolDef {
	return []toolDef{
		h.healthTool(), h.openTabTool(), h.navigateTool(), h.snapshotTool(),
		h.clickTool(), h.typeTool(), h.pressTool(), h.scrollTool(),
		h.backTool(), h.forwardTool(), h.refreshTool(), h.linksTool(),
		h.screenshotTool(), h.evaluateTool(), h.networkRequestsTool(),
		h.statsTool(), h.listTabsTool(), h.closeTabTool(), h.closeSessionTool(),
	}
}

func (h *Handler) register(s *server.MCPServer, tools ...toolDef) {
	for _, t := range tools {
		s.AddTool(t.tool, t.handler)
	}
}

// sessionOf returns the session override from the request, or the default.
func (h *Handler) sessionOf(req mcp.CallToolRequest) string {
	if s := req.GetString("session", ""); s != "" {
		return s
	}

	return h.session
}

func tabID(req mcp.CallToolRequest) (string, error) {
	return req.RequireString("tab_id")
}

func wrapErr(err error) (*mcp.CallToolResult, error) {
	return mcp.NewToolResultError(err.Error()), nil
}
