package scanner

import (
	"path/filepath"
	"strings"
)

var excludedDirectories = map[string]bool{
	".git":         true,
	"node_modules": true,
	".next":        true,
	"dist":         true,
	"build":        true,
	"__pycache__":  true,
	".venv":        true,
	"venv":         true,
	".tmp":         true,
	".idea":        true,
	".vscode":      true,
	"coverage":     true,
}

var excludedFiles = map[string]bool{
	"package-lock.json": true,
	"pnpm-lock.yaml":    true,
	"yarn.lock":         true,
	"bun.lock":          true,
	"bun.lockb":         true,
	".DS_Store":         true,
	"context":           true,
	"context.txt":       true,
}

var binaryExtensions = map[string]bool{
	".png":      true,
	".jpg":      true,
	".jpeg":     true,
	".gif":      true,
	".webp":     true,
	".ico":      true,
	".mp3":      true,
	".wav":      true,
	".mp4":      true,
	".mov":      true,
	".zip":      true,
	".tar":      true,
	".gz":       true,
	".7z":       true,
	".pdf":      true,
	".exe":      true,
	".bin":      true,
	".iso":      true,
	".woff":     true,
	".woff2":    true,
	".ttf":      true,
	".eot":      true,
	".pyc":      true,
	".db":       true,
	".sqlite":   true,
	".sqlite3":  true,
	".dmg":      true,
	".appimage": true,
}

func isEnvironmentFile(name string) bool {
	return name == ".env" ||
		strings.HasPrefix(name, ".env.")
}

func isSafeEnvironmentFile(name string) bool {
	switch name {
	case ".env.example", ".env.template", ".env.sample":
		return true
	default:
		return false
	}
}

func isBinaryFile(path string) bool {
	extension := strings.ToLower(
		filepath.Ext(path),
	)

	return binaryExtensions[extension]
}
