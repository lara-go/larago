package cache_test

import (
	"testing"

	"github.com/fatih/structs"
	"github.com/jinzhu/gorm"
	"github.com/lara-go/larago/cache"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func databaseStoreFactory() *cache.DatabaseStore {
	store := cache.NewDatabaseStore("cache")

	db, _ := gorm.Open("sqlite3", "file::memory:?mode=memory&cache=shared")
	db.LogMode(true)
	db.DB().SetMaxIdleConns(10)

	db.AutoMigrate(cache.DatabaseItem{})

	s := structs.New(store)
	dbField := s.Field("DB")
	dbField.Set(db)

	return store
}

func TestDatabaseStore_Expire(t *testing.T) {
	testExpiration(t, databaseStoreFactory())
}

func TestDatabaseStore_Forget(t *testing.T) {
	testForget(t, databaseStoreFactory())
}

func TestDatabaseStore_Clear(t *testing.T) {
	testClear(t, databaseStoreFactory())
}
