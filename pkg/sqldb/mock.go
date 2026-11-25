package sqldb

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitMockDatabase opens an in-memory SQLite database that is closed automatically when the test ends.
func InitMockDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := ":memory:"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql db: %v", err)
	}
	t.Cleanup(func() {
		if err := sqlDB.Close(); err != nil {
			t.Errorf("failed to close database: %v", err)
		}
	})
	return db
}
