package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func BinaryOrBuildDir() string {
	// 1. Check if running via `go run` (detect temp binary)
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)

		// On Windows, go run temp binaries are usually in Temp directories
		tempDir := os.TempDir()
		isSub, err := isSubDir(tempDir, exeDir)
		if err != nil {
			panic("unable to determine binary/build location")
		}
		if isSub {
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

func isSubDir(parent, child string) (bool, error) {
	// Get absolute, cleaned paths
	parentAbs, err := filepath.Abs(parent)
	if err != nil {
		return false, err
	}
	childAbs, err := filepath.Abs(child)
	if err != nil {
		return false, err
	}

	// Add a trailing separator to avoid false positives
	// e.g., "/tmp/foo" should not match "/tmp/foobar"
	parentAbs = filepath.Clean(parentAbs) + string(filepath.Separator)
	childAbs = filepath.Clean(childAbs)

	return strings.HasPrefix(childAbs, parentAbs), nil
}
