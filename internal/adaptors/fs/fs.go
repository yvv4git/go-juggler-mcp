package fs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yvv4git/go-juggler-mcp/internal/ports"
)

// FS adapts the filesystem to the ports.Filesystem interface.
type FS struct{}

// New creates a filesystem adapter.
func New() *FS {
	return &FS{}
}

// Save writes data to path, creating parent directories if needed.
func (f *FS) Save(path string, data []byte) error {
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create parent dir %q: %w", dir, err)
		}
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write file %q: %w", path, err)
	}

	return nil
}

var _ ports.Filesystem = (*FS)(nil)