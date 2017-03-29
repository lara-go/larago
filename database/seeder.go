package database

import "github.com/jinzhu/gorm"

// Seeder engine to populate DB with inital info.
type Seeder struct {
	Connection *gorm.DB
}
