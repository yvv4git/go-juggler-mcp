package juggler

import (
	"context"

	"github.com/yvv4git/go-juggler"
	"github.com/yvv4git/go-juggler-mcp/internal/ports"
	"github.com/yvv4git/go-juggler/browser"
)

// Client adapts the juggler HTTP client to the ports.BrowserClient interface.
type Client struct {
	c *juggler.Client
}

// NewClient creates an adapter over the camofox-browser REST API at addr.
func NewClient(addr string) *Client {
	return &Client{c: juggler.NewClient(addr)}
}

func (a *Client) Health(ctx context.Context) (*juggler.HealthResponse, error) {
	return a.c.Health(ctx)
}

func (a *Client) OpenTab(ctx context.Context, sessionKey, url string) (*juggler.TabResponse, error) {
	return a.c.OpenTab(ctx, sessionKey, url)
}

func (a *Client) Navigate(ctx context.Context, tabID, sessionKey, url string) error {
	return a.c.Navigate(ctx, tabID, sessionKey, url)
}

func (a *Client) Snapshot(ctx context.Context, tabID, sessionKey string) (*juggler.SnapshotResponse, error) {
	return a.c.Snapshot(ctx, tabID, sessionKey)
}

func (a *Client) CloseTab(ctx context.Context, tabID, sessionKey string) error {
	return a.c.CloseTab(ctx, tabID, sessionKey)
}

func (a *Client) CloseSession(ctx context.Context, sessionKey string) error {
	return a.c.CloseSession(ctx, sessionKey)
}

func (a *Client) Click(ctx context.Context, tabID, sessionKey, ref, selector string) error {
	return a.c.Click(ctx, tabID, sessionKey, ref, selector)
}

func (a *Client) Type(ctx context.Context, tabID, sessionKey, ref, selector, text string) error {
	return a.c.Type(ctx, tabID, sessionKey, ref, selector, text)
}

func (a *Client) Press(ctx context.Context, tabID, sessionKey, key string) error {
	return a.c.Press(ctx, tabID, sessionKey, key)
}

func (a *Client) Scroll(ctx context.Context, tabID, sessionKey, direction string, amount int) error {
	return a.c.Scroll(ctx, tabID, sessionKey, direction, amount)
}

func (a *Client) Back(ctx context.Context, tabID, sessionKey string) error {
	return a.c.Back(ctx, tabID, sessionKey)
}

func (a *Client) Forward(ctx context.Context, tabID, sessionKey string) error {
	return a.c.Forward(ctx, tabID, sessionKey)
}

func (a *Client) Refresh(ctx context.Context, tabID, sessionKey string) error {
	return a.c.Refresh(ctx, tabID, sessionKey)
}

func (a *Client) Links(ctx context.Context, tabID, sessionKey string, limit, offset int) (*browser.LinksResponse, error) {
	return a.c.Links(ctx, tabID, sessionKey, limit, offset)
}

func (a *Client) Screenshot(ctx context.Context, tabID, sessionKey string, fullPage bool) ([]byte, error) {
	return a.c.Screenshot(ctx, tabID, sessionKey, fullPage)
}

func (a *Client) Evaluate(ctx context.Context, tabID, sessionKey, expression string) (*juggler.EvaluateResponse, error) {
	return a.c.Evaluate(ctx, tabID, sessionKey, expression)
}

func (a *Client) NetworkRequests(ctx context.Context, tabID, sessionKey string) ([]juggler.ResourceEntry, error) {
	return a.c.NetworkRequests(ctx, tabID, sessionKey)
}

func (a *Client) Stats(ctx context.Context, tabID, sessionKey string) (*browser.TabStats, error) {
	return a.c.Stats(ctx, tabID, sessionKey)
}

func (a *Client) ListTabs(ctx context.Context, sessionKey string) (*juggler.ListTabsResponse, error) {
	return a.c.ListTabs(ctx, sessionKey)
}

var _ ports.BrowserClient = (*Client)(nil)
