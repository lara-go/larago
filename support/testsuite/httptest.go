package testsuite

import (
	net_http "net/http"
	"testing"

	"github.com/gavv/httpexpect"
)

// NewHTTPExpect constructor.
func NewHTTPExpect(handler net_http.Handler, t *testing.T) *httpexpect.Expect {
	testConfiguration := httpexpect.Config{
		Client: &net_http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	}

	testConfiguration.Printers = []httpexpect.Printer{
		httpexpect.NewDebugPrinter(t, true),
	}

	return httpexpect.WithConfig(testConfiguration)
}
