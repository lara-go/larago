package stubs

// MigrationStub template.
const MigrationStub = `
package migrations

import (
	"github.com/jinzhu/gorm"
)

// {{.Name}} migration.
type {{.Name}} struct{}

// Migrate runs migrations.
func (m *{{.Name}}) Migrate(tx *gorm.DB) error {
	return nil
}

// Rollback changes.
func (m *{{.Name}}) Rollback(tx *gorm.DB) error {
	return nil
}
`
