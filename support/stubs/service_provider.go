package stubs

// ServiceProviderStub template.
const ServiceProviderStub = `
package providers

import (
	"github.com/lara-go/larago"
)

// {{.Name}} .
type {{.Name}} struct{}

// Register service.
func (p *{{.Name}}) Register(application *larago.Application) {

}
`
