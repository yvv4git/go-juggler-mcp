package core

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"
)

// --- health ---

func (h *Handler) healthTool() toolDef {
	tool := mcp.NewTool("health",
		mcp.WithDescription("Check the browser status: engine, connection and memory usage"),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		h.log.Debug("health", zap.String("session", h.sessionOf(req)))

		res, err := h.client.Health(ctx)
		if err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultJSON(res)
	}}
}

// --- open_tab ---

func (h *Handler) openTabTool() toolDef {
	tool := mcp.NewTool("open_tab",
		mcp.WithDescription("Open a new tab and navigate to the given URL"),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL to navigate to"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		url, err := req.RequireString("url")
		if err != nil {
			return wrapErr(err)
		}

		h.log.Debug("open_tab", zap.String("session", h.sessionOf(req)), zap.String("url", url))

		res, err := h.client.OpenTab(ctx, h.sessionOf(req), url)
		if err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultJSON(res)
	}}
}

// --- navigate ---

func (h *Handler) navigateTool() toolDef {
	tool := mcp.NewTool("navigate",
		mcp.WithDescription("Load a URL in an existing tab"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("URL to load"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		url, err := req.RequireString("url")
		if err != nil {
			return wrapErr(err)
		}

		h.log.Debug("navigate", zap.String("tab", id), zap.String("url", url))

		if err := h.client.Navigate(ctx, id, h.sessionOf(req), url); err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultText("ok"), nil
	}}
}

// --- snapshot ---

func (h *Handler) snapshotTool() toolDef {
	tool := mcp.NewTool("snapshot",
		mcp.WithDescription("Get the ARIA tree of the page with element refs. Use the refs to act on elements"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		h.log.Debug("snapshot", zap.String("tab", id))

		res, err := h.client.Snapshot(ctx, id, h.sessionOf(req))
		if err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultJSON(res)
	}}
}

// --- click ---

func (h *Handler) clickTool() toolDef {
	tool := mcp.NewTool("click",
		mcp.WithDescription("Click an element by ref or CSS selector"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("ref",
			mcp.Description("Element ref from the snapshot (e.g. e1)"),
		),
		mcp.WithString("selector",
			mcp.Description("CSS selector to locate the element"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		ref := req.GetString("ref", "")
		selector := req.GetString("selector", "")

		h.log.Debug("click", zap.String("tab", id), zap.String("ref", ref), zap.String("selector", selector))

		if err := h.client.Click(ctx, id, h.sessionOf(req), ref, selector); err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultText("ok"), nil
	}}
}

// --- type ---

func (h *Handler) typeTool() toolDef {
	tool := mcp.NewTool("type",
		mcp.WithDescription("Fill an input field by ref or selector with text"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("Text to type"),
		),
		mcp.WithString("ref",
			mcp.Description("Element ref from the snapshot (e.g. e1)"),
		),
		mcp.WithString("selector",
			mcp.Description("CSS selector to locate the input"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		text, err := req.RequireString("text")
		if err != nil {
			return wrapErr(err)
		}

		ref := req.GetString("ref", "")
		selector := req.GetString("selector", "")

		h.log.Debug("type", zap.String("tab", id), zap.String("ref", ref), zap.String("selector", selector))

		if err := h.client.Type(ctx, id, h.sessionOf(req), ref, selector, text); err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultText("ok"), nil
	}}
}

// --- press ---

func (h *Handler) pressTool() toolDef {
	tool := mcp.NewTool("press",
		mcp.WithDescription("Press a keyboard key (Enter, Tab, Escape, etc.)"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("key",
			mcp.Required(),
			mcp.Description("Key name, e.g. Enter, Tab, Escape"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		key, err := req.RequireString("key")
		if err != nil {
			return wrapErr(err)
		}

		h.log.Debug("press", zap.String("tab", id), zap.String("key", key))

		if err := h.client.Press(ctx, id, h.sessionOf(req), key); err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultText("ok"), nil
	}}
}

// --- scroll ---

func (h *Handler) scrollTool() toolDef {
	tool := mcp.NewTool("scroll",
		mcp.WithDescription("Scroll the page up or down by N pixels"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("direction",
			mcp.Required(),
			mcp.Description("Scroll direction"),
			mcp.Enum("up", "down"),
		),
		mcp.WithInteger("amount",
			mcp.Required(),
			mcp.Description("Pixels to scroll"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		direction, err := req.RequireString("direction")
		if err != nil {
			return wrapErr(err)
		}

		amount, err := req.RequireInt("amount")
		if err != nil {
			return wrapErr(err)
		}

		h.log.Debug("scroll", zap.String("tab", id), zap.String("direction", direction), zap.Int("amount", amount))

		if err := h.client.Scroll(ctx, id, h.sessionOf(req), direction, amount); err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultText("ok"), nil
	}}
}

// --- back ---

func (h *Handler) backTool() toolDef {
	tool := mcp.NewTool("back",
		mcp.WithDescription("Navigate back in the tab history"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		h.log.Debug("back", zap.String("tab", id))

		if err := h.client.Back(ctx, id, h.sessionOf(req)); err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultText("ok"), nil
	}}
}

// --- forward ---

func (h *Handler) forwardTool() toolDef {
	tool := mcp.NewTool("forward",
		mcp.WithDescription("Navigate forward in the tab history"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		h.log.Debug("forward", zap.String("tab", id))

		if err := h.client.Forward(ctx, id, h.sessionOf(req)); err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultText("ok"), nil
	}}
}

// --- refresh ---

func (h *Handler) refreshTool() toolDef {
	tool := mcp.NewTool("refresh",
		mcp.WithDescription("Reload the current page"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		h.log.Debug("refresh", zap.String("tab", id))

		if err := h.client.Refresh(ctx, id, h.sessionOf(req)); err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultText("ok"), nil
	}}
}

// --- links ---

func (h *Handler) linksTool() toolDef {
	tool := mcp.NewTool("links",
		mcp.WithDescription("List all links on the page with pagination"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithInteger("limit",
			mcp.Description("Maximum number of links to return"),
			mcp.DefaultNumber(50),
		),
		mcp.WithInteger("offset",
			mcp.Description("Pagination offset"),
			mcp.DefaultNumber(0),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		limit := req.GetInt("limit", 50)
		offset := req.GetInt("offset", 0)

		h.log.Debug("links", zap.String("tab", id), zap.Int("limit", limit), zap.Int("offset", offset))

		res, err := h.client.Links(ctx, id, h.sessionOf(req), limit, offset)
		if err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultJSON(res)
	}}
}

// --- screenshot ---

func (h *Handler) screenshotTool() toolDef {
	tool := mcp.NewTool("screenshot",
		mcp.WithDescription("Take a PNG screenshot of the page or full page"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithBoolean("full_page",
			mcp.Description("Capture the full page instead of the viewport"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		fullPage := req.GetBool("full_page", false)

		h.log.Debug("screenshot", zap.String("tab", id), zap.Bool("full_page", fullPage))

		png, err := h.client.Screenshot(ctx, id, h.sessionOf(req), fullPage)
		if err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultImage(
			fmt.Sprintf("Screenshot of tab %s (%d bytes)", id, len(png)),
			base64.StdEncoding.EncodeToString(png),
			"image/png",
		), nil
	}}
}

// --- evaluate ---

func (h *Handler) evaluateTool() toolDef {
	tool := mcp.NewTool("evaluate",
		mcp.WithDescription("Run arbitrary JavaScript in the page context"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("expression",
			mcp.Required(),
			mcp.Description("JavaScript expression to evaluate"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		expression, err := req.RequireString("expression")
		if err != nil {
			return wrapErr(err)
		}

		h.log.Debug("evaluate", zap.String("tab", id))

		res, err := h.client.Evaluate(ctx, id, h.sessionOf(req), expression)
		if err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultJSON(res)
	}}
}

// --- network_requests ---

func (h *Handler) networkRequestsTool() toolDef {
	tool := mcp.NewTool("network_requests",
		mcp.WithDescription("Get all resources loaded by the page (navigation + subresources)"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		h.log.Debug("network_requests", zap.String("tab", id))

		res, err := h.client.NetworkRequests(ctx, id, h.sessionOf(req))
		if err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultJSON(map[string]any{"requests": res})
	}}
}

// --- stats ---

func (h *Handler) statsTool() toolDef {
	tool := mcp.NewTool("stats",
		mcp.WithDescription("Get tab state: URL, visited URLs, refs count"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		h.log.Debug("stats", zap.String("tab", id))

		res, err := h.client.Stats(ctx, id, h.sessionOf(req))
		if err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultJSON(res)
	}}
}

// --- list_tabs ---

func (h *Handler) listTabsTool() toolDef {
	tool := mcp.NewTool("list_tabs",
		mcp.WithDescription("List all open tabs in a session"),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		h.log.Debug("list_tabs", zap.String("session", h.sessionOf(req)))

		res, err := h.client.ListTabs(ctx, h.sessionOf(req))
		if err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultJSON(res)
	}}
}

// --- close_tab ---

func (h *Handler) closeTabTool() toolDef {
	tool := mcp.NewTool("close_tab",
		mcp.WithDescription("Close a tab"),
		mcp.WithString("tab_id",
			mcp.Required(),
			mcp.Description("Tab identifier returned by open_tab"),
		),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := tabID(req)
		if err != nil {
			return wrapErr(err)
		}

		h.log.Debug("close_tab", zap.String("tab", id))

		if err := h.client.CloseTab(ctx, id, h.sessionOf(req)); err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultText("ok"), nil
	}}
}

// --- close_session ---

func (h *Handler) closeSessionTool() toolDef {
	tool := mcp.NewTool("close_session",
		mcp.WithDescription("Destroy a session and all its tabs. Defaults to the server session"),
		mcp.WithString("session",
			mcp.Description("Juggler session key. Defaults to the server session"),
		),
	)

	return toolDef{tool: tool, handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		h.log.Debug("close_session", zap.String("session", h.sessionOf(req)))

		if err := h.client.CloseSession(ctx, h.sessionOf(req)); err != nil {
			return wrapErr(err)
		}

		return mcp.NewToolResultText("ok"), nil
	}}
}
