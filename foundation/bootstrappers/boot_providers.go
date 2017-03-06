package bootstrappers

import "github.com/lara-go/larago"

// BootProviders lounches application boot process.
func BootProviders(application *larago.Application) error {
	application.Boot()

	return nil
}
