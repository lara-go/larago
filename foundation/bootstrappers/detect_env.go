package bootstrappers

import (
	"fmt"
	"os"
	"path"

	"github.com/joho/godotenv"
	"github.com/lara-go/larago"
)

// DetectEnv loads environment variables from .env file.
func DetectEnv(application *larago.Application) error {
	home := getHomeDirectory(application.HomeDirectory)
	dotEnv := getConfigFile(path.Join(home, ".env"))

	// Load env variables.
	err := godotenv.Load(dotEnv)

	switch err {
	case os.ErrPermission:
		return fmt.Errorf("Can't load %s file. Check permissions", dotEnv)
	case os.ErrNotExist:
		return nil
	default:
		return err
	}
}

// Get home directory path.
// TODO: fix this dirty hack!
func getHomeDirectory(defaultHome string) string {
	for i, val := range os.Args {
		if (val == "-H" || val == "--home") && len(os.Args) >= i+2 {
			return os.Args[i+1]
		}
	}

	return defaultHome
}

// Get config path.
// TODO: fix this dirty hack!
func getConfigFile(defaultConfig string) string {
	for i, val := range os.Args {
		if (val == "-c" || val == "--config") && len(os.Args) >= i+2 {
			return os.Args[i+1]
		}
	}

	return defaultConfig
}
