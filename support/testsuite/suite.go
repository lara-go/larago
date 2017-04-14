package testsuite

import (
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/lara-go/larago"
	"github.com/lara-go/larago/foundation/bootstrappers"
	"github.com/lara-go/larago/http"
)

// ApplicationSuite interface.
type ApplicationSuite interface {
	// MakeApplication factory.
	MakeApplication() *larago.Application
}

// LaragoSuite prepares to run your tests
type LaragoSuite struct {
	ApplicationSuite
	Application *larago.Application
}

// BootstrapApplication to have new instance every test.
func (s *LaragoSuite) BootstrapApplication(application *larago.Application) {
	// Run default bootstrappers.
	err := application.BootstrapWith(
		bootstrappers.LoadConfig,
		bootstrappers.BootProviders,
	)

	if err != nil {
		panic(err)
	}

	// Change default config values.
	application.Config().Set("App.Env", "testing")
	application.Config().Set("HTTP.LogRequests", false)

	// Bootstrap router.
	router := application.Get("router").(*http.Router)
	router.Bootstrap()
}

// HTTP returns instance of httpexpect.
func (s *LaragoSuite) HTTP(t *testing.T) *httpexpect.Expect {
	router := s.Application.Get("router").(*http.Router)

	return NewHTTPExpect(router.GetHTTPRouter(), t)
}
