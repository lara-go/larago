package database

import (
	"github.com/jinzhu/gorm"
	"github.com/lara-go/larago"
)

// ServiceProvider struct.
type ServiceProvider struct{}

// Register service.
func (p *ServiceProvider) Register(application *larago.Application) {
	p.registerDatabaseConnection(application)
	p.registerMigrator(application)

	application.Commands(
		&CommandDBSeed{},
		&CommandMakeMigration{},
		&CommandMigrate{},
		&CommandMigrateRollback{},
		&CommandMigrateReset{},
	)
}

func (p *ServiceProvider) registerDatabaseConnection(application *larago.Application) {
	application.Bind(&Manager{}, "db")

	application.Bind(func() (*gorm.DB, error) {
		var manager Manager
		application.Make(&manager)

		return manager.GetConnection()
	}, "db.connection")
}

func (p *ServiceProvider) registerMigrator(application *larago.Application) {
	application.Bind(&Migrator{})
}
