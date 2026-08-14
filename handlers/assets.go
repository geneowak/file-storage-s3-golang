package handlers

import (
	"os"
)

func (cfg ApiConfig) ensureAssetsDir() error {
	if _, err := os.Stat(cfg.AssetsRoot); os.IsNotExist(err) {
		return os.Mkdir(cfg.AssetsRoot, 0755)
	}
	return nil
}
