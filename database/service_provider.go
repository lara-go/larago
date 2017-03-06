package database

import "github.com/lara-go/larago"

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
	application.Bind(func() (*Connection, error) {
		var connection Connection
		application.Make(&connection)

		// Establish database connection.
		err := connection.Connect()
		if err != nil {
			return nil, err
		}

		return &connection, nil
	})
}

func (p *ServiceProvider) registerMigrator(application *larago.Application) {
	application.Bind(&Migrator{})
}
