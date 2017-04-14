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
	if err := godotenv.Load(dotEnv); err != nil {
		return fmt.Errorf("Can't load .env file in %s", home)
	}

	return nil
}

// Get home directory path.
// TODO: fix this dirty hack!
func getHomeDirectory(defaultHome string) string {
	for i, val := range os.Args {
		if (val == "-r" || val == "--home") && len(os.Args) >= i+2 {
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
