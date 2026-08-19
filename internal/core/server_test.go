package core

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/mcptest"
	"github.com/yvv4git/go-juggler/browser"
	"go.uber.org/zap"
)

type mockClient struct{}

func (mockClient) Health(ctx context.Context) (*browser.HealthResponse, error) {
	return &browser.HealthResponse{OK: true, Engine: "camoufox", Sessions: 1, BrowserConnected: true}, nil
}

func (mockClient) OpenTab(ctx context.Context, sessionKey, url string) (*browser.TabResponse, error) {
	return &browser.TabResponse{TabID: "tab-1", URL: url}, nil
}

func (mockClient) Navigate(ctx context.Context, tabID, sessionKey, url string) error { return nil }
func (mockClient) Snapshot(ctx context.Context, tabID, sessionKey string) (*browser.SnapshotResponse, error) {
	return &browser.SnapshotResponse{URL: "http://example.com", Snapshot: "snapshot", RefsCount: 1}, nil
}
func (mockClient) CloseTab(ctx context.Context, tabID, sessionKey string) error { return nil }
func (mockClient) CloseSession(ctx context.Context, sessionKey string) error    { return nil }
func (mockClient) Click(ctx context.Context, tabID, sessionKey, ref, selector string) error {
	return nil
}
func (mockClient) Type(ctx context.Context, tabID, sessionKey, ref, selector, text string) error {
	return nil
}
func (mockClient) Press(ctx context.Context, tabID, sessionKey, key string) error { return nil }
func (mockClient) Scroll(ctx context.Context, tabID, sessionKey, direction string, amount int) error {
	return nil
}
func (mockClient) Back(ctx context.Context, tabID, sessionKey string) error    { return nil }
func (mockClient) Forward(ctx context.Context, tabID, sessionKey string) error { return nil }
func (mockClient) Refresh(ctx context.Context, tabID, sessionKey string) error { return nil }
func (mockClient) Links(ctx context.Context, tabID, sessionKey string, limit, offset int) (*browser.LinksResponse, error) {
	return &browser.LinksResponse{}, nil
}
func (mockClient) Screenshot(ctx context.Context, tabID, sessionKey string, fullPage bool) ([]byte, error) {
	return []byte("png-data"), nil
}
func (mockClient) Evaluate(ctx context.Context, tabID, sessionKey, expression string) (*browser.EvaluateResponse, error) {
	return &browser.EvaluateResponse{OK: true, Result: "ok"}, nil
}
func (mockClient) NetworkRequests(ctx context.Context, tabID, sessionKey string) ([]browser.ResourceEntry, error) {
	return []browser.ResourceEntry{{Name: "http://example.com", Type: "navigation", Status: 200}}, nil
}
func (mockClient) Stats(ctx context.Context, tabID, sessionKey string) (*browser.TabStats, error) {
	return &browser.TabStats{TabID: tabID, URL: "http://example.com"}, nil
}
func (mockClient) ListTabs(ctx context.Context, sessionKey string) (*browser.ListTabsResponse, error) {
	return &browser.ListTabsResponse{Running: true, Tabs: []browser.TabInfo{{TabID: "tab-1", URL: "http://example.com"}}}, nil
}

func newTestHandler(t *testing.T) *Handler {
	t.Helper()

	log := zap.NewNop()

	return New(mockClient{}, "test-session", log)
}

func TestToolsRegistered(t *testing.T) {
	h := newTestHandler(t)

	s := mcptest.NewUnstartedServer(t)
	for _, td := range h.tools() {
		s.AddTool(td.tool, td.handler)
	}

	if err := s.Start(context.Background()); err != nil {
		t.Fatalf("start server: %v", err)
	}
	defer s.Close()

	got, err := s.Client().ListTools(context.Background(), mcp.ListToolsRequest{})
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	if len(got.Tools) != len(h.tools()) {
		t.Fatalf("tools count: got %d, want %d", len(got.Tools), len(h.tools()))
	}
}

func TestOpenTabTool(t *testing.T) {
	h := newTestHandler(t)

	s := mcptest.NewUnstartedServer(t)
	s.AddTool(h.tools()[1].tool, h.tools()[1].handler)

	if err := s.Start(context.Background()); err != nil {
		t.Fatalf("start server: %v", err)
	}
	defer s.Close()

	res, err := s.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "open_tab",
			Arguments: map[string]any{"url": "http://example.com"},
		},
	})
	if err != nil {
		t.Fatalf("call tool: %v", err)
	}

	if res.IsError {
		t.Fatalf("expected success, got error result")
	}
	if len(res.Content) == 0 {
		t.Fatal("expected content in result")
	}
}

func TestScreenshotToolImageContent(t *testing.T) {
	h := newTestHandler(t)

	s := mcptest.NewUnstartedServer(t)
	for _, td := range h.tools() {
		if td.tool.Name == "screenshot" {
			s.AddTool(td.tool, td.handler)
		}
	}

	if err := s.Start(context.Background()); err != nil {
		t.Fatalf("start server: %v", err)
	}
	defer s.Close()

	res, err := s.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "screenshot",
			Arguments: map[string]any{"tab_id": "tab-1"},
		},
	})
	if err != nil {
		t.Fatalf("call tool: %v", err)
	}

	if res.IsError {
		t.Fatalf("expected success, got error result")
	}

	found := false
	for _, c := range res.Content {
		if img, ok := c.(mcp.ImageContent); ok && img.MIMEType == "image/png" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected image/png content in screenshot result")
	}
}
