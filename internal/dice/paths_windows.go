//go:build windows

package dice

import (
	"path/filepath"

	"golang.org/x/sys/windows"
)

func HistoryDir() string {
	doc, err := windows.KnownFolderPath(windows.FOLDERID_Documents, 0)
	if err != nil {
		return filepath.Join(".", "history")
	}
	return filepath.Join(doc, "dice-roller", "history")
}
