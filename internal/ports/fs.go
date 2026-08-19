package ports

// Filesystem is the minimal surface the MCP tools need to persist files.
type Filesystem interface {
	// Save writes data to path, creating parent directories if needed.
	Save(path string, data []byte) error
}