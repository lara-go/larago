package bootstrappers

import "github.com/lara-go/larago"

// LoadConfig imports config from the lazy loader.
func LoadConfig(application *larago.Application) error {
	application.ImportConfig()

	return nil
}
