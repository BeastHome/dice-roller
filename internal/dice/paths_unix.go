//go:build !windows

package dice

import (
	"os"
	"path/filepath"
)

func HistoryDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "./history"
	}
	return filepath.Join(home, ".local", "share", "dice-roller")
}
