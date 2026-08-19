package ports

import (
	"context"

	"github.com/yvv4git/go-juggler/browser"
)

// BrowserClient is the minimal surface of the camofox-browser REST API
// needed by the MCP tools.
type BrowserClient interface {
	Health(ctx context.Context) (*browser.HealthResponse, error)

	OpenTab(ctx context.Context, sessionKey, url string) (*browser.TabResponse, error)
	Navigate(ctx context.Context, tabID, sessionKey, url string) error
	Snapshot(ctx context.Context, tabID, sessionKey string) (*browser.SnapshotResponse, error)
	CloseTab(ctx context.Context, tabID, sessionKey string) error
	CloseSession(ctx context.Context, sessionKey string) error

	Click(ctx context.Context, tabID, sessionKey, ref, selector string) error
	Type(ctx context.Context, tabID, sessionKey, ref, selector, text string) error
	Press(ctx context.Context, tabID, sessionKey, key string) error
	Scroll(ctx context.Context, tabID, sessionKey, direction string, amount int) error
	Back(ctx context.Context, tabID, sessionKey string) error
	Forward(ctx context.Context, tabID, sessionKey string) error
	Refresh(ctx context.Context, tabID, sessionKey string) error

	Links(ctx context.Context, tabID, sessionKey string, limit, offset int) (*browser.LinksResponse, error)
	Screenshot(ctx context.Context, tabID, sessionKey string, fullPage bool) ([]byte, error)
	Evaluate(ctx context.Context, tabID, sessionKey, expression string) (*browser.EvaluateResponse, error)
	NetworkRequests(ctx context.Context, tabID, sessionKey string) ([]browser.ResourceEntry, error)
	Stats(ctx context.Context, tabID, sessionKey string) (*browser.TabStats, error)
	ListTabs(ctx context.Context, sessionKey string) (*browser.ListTabsResponse, error)
}
