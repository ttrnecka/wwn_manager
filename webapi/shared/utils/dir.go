package utils

import (
	"os"
	"path/filepath"
)

func BinaryOrBuildDir() string {
	// 1. Check if running via `go run` (detect temp binary)
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)

		// On Windows, go run temp binaries are usually in Temp directories
		tempDir := os.TempDir()
		rel, err := filepath.Rel(tempDir, exeDir)
		if err == nil && rel != exeDir {
			// likely a go run temporary binary
			if cwd, err := os.Getwd(); err == nil {
				return cwd
			}
		}

		// Otherwise, return executable directory
		return exeDir
	}

	panic("unable to determine binary/build location")
}
