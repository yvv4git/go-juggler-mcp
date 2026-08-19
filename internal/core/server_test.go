package core

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/mcptest"
	"github.com/yvv4git/go-juggler-mcp/internal/adaptors/fs"
	"github.com/yvv4git/go-juggler/browser"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func newMockBrowserClient(t *testing.T) *MockBrowserClient {
	t.Helper()

	ctrl := gomock.NewController(t)

	return NewMockBrowserClient(ctrl)
}

func newTestHandler(t *testing.T, browser *MockBrowserClient) *Handler {
	t.Helper()

	log := zap.NewNop()

	return New(browser, fs.New(), "test-session", log)
}

// testPNG returns a valid solid-color PNG of the given dimensions.
func testPNG(width, height int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func TestToolsRegistered(t *testing.T) {
	h := newTestHandler(t, newMockBrowserClient(t))

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
	client := newMockBrowserClient(t)
	client.EXPECT().
		OpenTab(gomock.Any(), "test-session", "http://example.com").
		Return(&browser.TabResponse{TabID: "tab-1", URL: "http://example.com"}, nil)

	h := newTestHandler(t, client)

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

func newScreenshotTestServer(t *testing.T, client *MockBrowserClient) *mcptest.Server {
	t.Helper()

	h := newTestHandler(t, client)

	s := mcptest.NewUnstartedServer(t)
	for _, td := range h.tools() {
		if td.tool.Name == "screenshot" {
			s.AddTool(td.tool, td.handler)
		}
	}

	if err := s.Start(context.Background()); err != nil {
		t.Fatalf("start server: %v", err)
	}
	t.Cleanup(s.Close)

	return s
}

func expectScreenshot(client *MockBrowserClient) {
	client.EXPECT().
		Screenshot(gomock.Any(), "tab-1", "test-session", false).
		Return(testPNG(400, 300), nil)
}

func TestScreenshotToolImageContent(t *testing.T) {
	client := newMockBrowserClient(t)
	expectScreenshot(client)

	s := newScreenshotTestServer(t, client)

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

func TestScreenshotToolJPEGContent(t *testing.T) {
	client := newMockBrowserClient(t)
	expectScreenshot(client)

	s := newScreenshotTestServer(t, client)

	res, err := s.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "screenshot",
			Arguments: map[string]any{"tab_id": "tab-1", "format": "jpeg", "quality": 80},
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
		if img, ok := c.(mcp.ImageContent); ok && img.MIMEType == "image/jpeg" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected image/jpeg content in screenshot result")
	}
}

func TestScreenshotToolSaveToDisk(t *testing.T) {
	client := newMockBrowserClient(t)
	expectScreenshot(client)

	s := newScreenshotTestServer(t, client)

	path := filepath.Join(t.TempDir(), "shots", "page.png")
	res, err := s.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "screenshot",
			Arguments: map[string]any{"tab_id": "tab-1", "path": path},
		},
	})
	if err != nil {
		t.Fatalf("call tool: %v", err)
	}

	if res.IsError {
		t.Fatalf("expected success, got error result")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read saved file: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("saved file is empty")
	}

	found := false
	for _, c := range res.Content {
		if txt, ok := c.(mcp.TextContent); ok {
			if txt.Text != "" && txt.Text != "ok" {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("expected text result with save info")
	}
}

func TestProcessImageResize(t *testing.T) {
	src := testPNG(400, 300)

	data, mime, err := processImage(src, formatPNG, 85, 200)
	if err != nil {
		t.Fatalf("process image: %v", err)
	}

	if mime != "image/png" {
		t.Fatalf("mime: got %s, want image/png", mime)
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode resized image: %v", err)
	}
	if got := img.Bounds().Dx(); got != 200 {
		t.Fatalf("width: got %d, want 200", got)
	}
	if got := img.Bounds().Dy(); got != 150 {
		t.Fatalf("height: got %d, want 150", got)
	}
}

func TestProcessImageJPEG(t *testing.T) {
	src := testPNG(400, 300)

	data, mime, err := processImage(src, formatJPEG, 80, 0)
	if err != nil {
		t.Fatalf("process image: %v", err)
	}

	if mime != "image/jpeg" {
		t.Fatalf("mime: got %s, want image/jpeg", mime)
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("decode jpeg: %v", err)
	}
	if format != "jpeg" {
		t.Fatalf("format: got %s, want jpeg", format)
	}
	if got := img.Bounds().Dx(); got != 400 {
		t.Fatalf("width: got %d, want 400", got)
	}
}