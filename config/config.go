package config

import (
	"os"
	"path/filepath"
)

func WriteState(path string, content []byte, mode os.FileMode) error {
	if !FileExists(path) {
		if err := os.MkdirAll(filepath.Dir(path), mode); err != nil {
			return err
		}
		os.Create(path)
	}

	return os.WriteFile(path, content, mode)
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
