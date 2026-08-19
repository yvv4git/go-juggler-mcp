package commands

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"
	"github.com/yvv4git/go-juggler-mcp/internal/adaptors/fs"
	"github.com/yvv4git/go-juggler-mcp/internal/adaptors/juggler"
	"github.com/yvv4git/go-juggler-mcp/internal/config"
	mcpconfig "github.com/yvv4git/go-juggler-mcp/internal/config/mcp"
	"github.com/yvv4git/go-juggler-mcp/internal/core"
	"github.com/yvv4git/go-juggler-mcp/internal/infra"
	"go.uber.org/zap"
)

func ServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the MCP server for browser automation",
		Long: `Start the MCP server that exposes browser control through the
Juggler protocol.

Examples:
  go-juggler-mcp serve
  go-juggler-mcp serve --transport http --http-addr :8080
  go-juggler-mcp serve --addr http://localhost:9377 --session demo`,
		RunE: serve,
	}

	cmd.Flags().String("addr", "http://localhost:9377", "camofox-browser REST API endpoint (overrides config)")
	cmd.Flags().String("session", "", "juggler session key (overrides config)")
	cmd.Flags().String("transport", "stdio", "MCP transport: stdio or http (overrides config)")
	cmd.Flags().String("http-addr", ":8080", "listen address for the http transport (overrides config)")
	cmd.Flags().String("log-level", "info", "log level (overrides config)")

	return cmd
}

func serve(cmd *cobra.Command, _ []string) error {
	cfgPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("get config flag: %w", err)
	}

	var cfg mcpconfig.Config
	if err := mcpconfig.Load(cfgPath, &cfg); err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Precedence: explicit flag > config file > cobra default.
	for _, f := range []struct {
		name  string
		value string
	}{
		{name: "addr", value: cfg.Serve.Addr},
		{name: "session", value: cfg.Serve.Session},
		{name: "transport", value: string(cfg.Serve.Transport)},
		{name: "http-addr", value: cfg.Serve.HTTPAddr},
		{name: "log-level", value: cfg.Log.Level},
	} {
		if f.value != "" && !cmd.Flags().Changed(f.name) {
			if err := cmd.Flags().Set(f.name, f.value); err != nil {
				return fmt.Errorf("set %s flag: %w", f.name, err)
			}
		}
	}

	addr, err := cmd.Flags().GetString("addr")
	if err != nil {
		return fmt.Errorf("get addr flag: %w", err)
	}

	session, err := cmd.Flags().GetString("session")
	if err != nil {
		return fmt.Errorf("get session flag: %w", err)
	}

	transport, err := cmd.Flags().GetString("transport")
	if err != nil {
		return fmt.Errorf("get transport flag: %w", err)
	}

	httpAddr, err := cmd.Flags().GetString("http-addr")
	if err != nil {
		return fmt.Errorf("get http-addr flag: %w", err)
	}

	logLevel, err := cmd.Flags().GetString("log-level")
	if err != nil {
		return fmt.Errorf("get log-level flag: %w", err)
	}

	log, err := infra.NewWithLogLevel(logLevel)
	if err != nil {
		return fmt.Errorf("setup logger: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Info("Config",
		zap.String("config path", cfgPath),
		zap.String("juggler addr", addr),
		zap.String("transport", transport),
		zap.String("session", session),
	)

	client := juggler.NewClient(addr)
	filesystem := fs.New()
	handler := core.New(client, filesystem, session, log)
	mcpServer := handler.MCPServer()

	switch config.Transport(transport) {
	case config.TransportStdio:
		return serveStdio(mcpServer, log)
	case config.TransportHTTP:
		return serveHTTP(ctx, mcpServer, httpAddr, log)
	default:
		return fmt.Errorf("unsupported transport %q", transport)
	}
}

func serveStdio(s *server.MCPServer, log *zap.Logger) error {
	log.Info("Serving MCP over stdio")

	return server.ServeStdio(s)
}

func serveHTTP(ctx context.Context, s *server.MCPServer, addr string, log *zap.Logger) error {
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           server.NewStreamableHTTPServer(s),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = httpServer.Shutdown(shutdownCtx)
	}()

	log.Info("Serving MCP over streamable HTTP", zap.String("addr", addr))

	if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
