package core

import (
	"os"
	"path/filepath"
	"strings"
)

func ExpandHome(p string) string {
	if p == "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return p
	}
	if p == "~" {
		return home
	}
	if strings.HasPrefix(p, "~/") {
		return filepath.Join(home, p[2:])
	}
	return p
}

func IsSubPath(targetPath, basePath string) bool {
	target := filepath.Clean(targetPath)
	base := filepath.Clean(basePath)
	if target == base {
		return true
	}
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}
	return rel != "." && !strings.HasPrefix(rel, "..")
}

func PathExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

func DirSize(startPath string) int64 {
	var total int64 = 0
	_ = filepath.WalkDir(startPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			return nil
		}
		info, statErr := d.Info()
		if statErr != nil {
			return nil
		}
		total += info.Size()
		return nil
	})
	return total
}
