package host

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/MeguruMacabre/meguru-pack-compact/internal/config"
	"github.com/MeguruMacabre/meguru-pack-compact/internal/instance"
)

type InstancePaths struct {
	InstanceRoot string
	GameRoot     string
}

// prepareInstance проверяет директорию сборки и определяет пути инстанса и игры.
func prepareInstance(
	appRoot string,
	cfg config.HostConfig,
) (InstancePaths, error) {
	instanceExists, err := dirExists(appRoot, cfg.InstanceDirectory)
	if err != nil {
		return InstancePaths{}, err
	}

	if !instanceExists {
		return InstancePaths{}, fmt.Errorf(
			"instance directory does not exist: %s",
			cfg.InstanceDirectory,
		)
	}

	instanceRoot := filepath.Join(appRoot, cfg.InstanceDirectory)

	gameRoot, err := instance.FindGameRoot(instanceRoot)
	if err != nil {
		return InstancePaths{}, err
	}

	return InstancePaths{
		InstanceRoot: instanceRoot,
		GameRoot:     gameRoot,
	}, nil
}

// dirExists проверяет существование директории внутри указанного корня.
func dirExists(root string, name string) (bool, error) {
	dirPath := filepath.Join(root, name)

	info, err := os.Stat(dirPath)
	if err == nil {
		return info.IsDir(), nil
	}

	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}

	return false, err
}
