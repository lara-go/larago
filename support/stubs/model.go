package stubs

// ModelStub template.
const ModelStub = `
package models

import "github.com/jinzhu/gorm"

// {{.Name}} model.
type {{.Name}} struct {
	gorm.Model
}
`
