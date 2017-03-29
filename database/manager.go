package database

import (
	"github.com/jinzhu/gorm"
	"github.com/lara-go/larago/logger"
)

// Manager to the database.
type Manager struct {
	Driver string `di:"Config.Database.Driver"`
	DSN    string `di:"Config.Database.DSN"`
	Debug  bool   `di:"Config.App.Debug"`
	Logger *logger.Logger

	connection *gorm.DB
}

// Connect to the database.
func (m *Manager) Connect() error {
	// Open connection.
	db, err := gorm.Open(m.Driver, m.DSN)
	if err != nil {
		return err
	}

	// Check if connection is active.
	if err = db.DB().Ping(); err != nil {
		return err
	}

	m.Logger.Debug("Connected to %s via %s", m.DSN, m.Driver)

	if m.Debug {
		db.LogMode(true)
	}

	m.connection = db

	return nil
}

// Disconnect from the database.
func (m *Manager) Disconnect() {
	if m.connection != nil {
		m.connection.Close()
	}
}

// GetConnection to db.
func (m *Manager) GetConnection() (*gorm.DB, error) {
	var err error
	if m.connection == nil {
		err = m.Connect()
	}

	return m.connection, err
}
