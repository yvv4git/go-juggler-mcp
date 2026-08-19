package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/yvv4git/go-juggler-mcp/cmd"
)

func main() {
	rootCommand := &cobra.Command{
		Use:   "go-juggler-mcp",
		Short: "MCP server for browser automation via the Juggler protocol",
	}
	rootCommand.PersistentFlags().StringP("config", "c", "config.toml", "Path to config file")

	rootCommand.AddCommand(commands.ServeCommand())

	if err := rootCommand.Execute(); err != nil {
		os.Exit(1)
	}
}
