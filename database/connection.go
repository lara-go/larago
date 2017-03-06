package database

import (
	"github.com/lara-go/larago/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // Import sqlite3 dialect
)

// Connection to the database.
type Connection struct {
	Driver string `di:"Config.Database.Driver"`
	DSN    string `di:"Config.Database.DSN"`
	Debug  bool   `di:"Config.App.Debug"`
	Logger *logger.Logger

	connection *gorm.DB
}

// Connect to the database.
func (c *Connection) Connect() error {
	// Open connection.
	db, err := gorm.Open(c.Driver, c.DSN)
	if err != nil {
		return err
	}

	// Check if connection is active.
	if err = db.DB().Ping(); err != nil {
		return err
	}

	c.Logger.Debug("Connected to %s via %s", c.DSN, c.Driver)

	if c.Debug {
		db.LogMode(true)
	}

	c.connection = db

	return nil
}

// Disconnect from the database.
func (c *Connection) Disconnect() {
	if c.connection != nil {
		c.connection.Close()
	}
}

// GetConnection to db.
func (c *Connection) GetConnection() *gorm.DB {
	return c.connection
}
