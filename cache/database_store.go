package cache

import (
	"encoding/base64"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/uniplaces/carbon"
)

// DatabaseItem to store.
type DatabaseItem struct {
	Key        string `gorm:"not null;unique_index"`
	Value      string `gorm:"type:text"`
	Expiration time.Time

	tableName string `gorm:"-"`
}

// TableName getter.
func (s *DatabaseItem) TableName() string {
	if s.tableName != "" {
		return s.tableName
	}

	return "cache"
}

// DatabaseStore .
type DatabaseStore struct {
	DB *gorm.DB

	table string
}

// NewDatabaseStore constructor.
func NewDatabaseStore(table string) *DatabaseStore {
	return &DatabaseStore{
		table: table,
	}
}

// Has checks if there is such item.
func (s *DatabaseStore) Has(key string) bool {
	return s.findItem(key) != nil
}

// Put value in cache by key.
func (s *DatabaseStore) Put(key string, value interface{}, duration time.Duration) error {
	// Serialize item in string.
	serialized, err := serializeValue(value)
	if err != nil {
		return ErrorSerialize
	}

	// Create item or update.
	item := s.makeItem()
	item.Key = key
	item.Value = base64.StdEncoding.EncodeToString(serialized)
	item.Expiration = carbon.NewCarbon(time.Now().Add(duration)).Time

	if err := s.DB.Create(item).Error; err != nil {
		s.DB.Where("key = ?", key).Update(item)
	}

	return nil
}

// Get saved value by the key.
func (s *DatabaseStore) Get(key string, target interface{}) error {
	item := s.findItem(key)
	if item == nil {
		return ErrorMissed
	}

	// Base64 decode item.
	by, err := base64.StdEncoding.DecodeString(item.Value)
	if err != nil {
		return ErrorUnserialize
	}

	// Unserialize value and set it to the target
	err = unserializeValue(by, target)
	if err != nil {
		s.Forget(key)

		return ErrorUnserialize
	}

	return nil
}

func (s *DatabaseStore) findItem(key string) *DatabaseItem {
	// Check if is still alive and if not, forget it.
	item := s.makeItem()
	if s.DB.Where("key = ?", key).First(&item).RecordNotFound() {
		return nil
	}

	// Check expiration.
	if carbon.NewCarbon(item.Expiration).IsPast() {
		s.Forget(key)

		return nil
	}

	return item
}

// Forever put value in store by key forever.
func (s *DatabaseStore) Forever(key string, value interface{}) error {
	return s.Put(key, value, time.Minute*5256000)
}

// Forget the value.
func (s *DatabaseStore) Forget(key string) {
	item := s.makeItem()

	s.DB.Where("key = ?", key).Delete(item)
}

// Clear storage.
func (s *DatabaseStore) Clear() {
	s.DB.Delete(s.makeItem())
}

func (s *DatabaseStore) makeItem() *DatabaseItem {
	return &DatabaseItem{
		tableName: s.table,
	}
}
