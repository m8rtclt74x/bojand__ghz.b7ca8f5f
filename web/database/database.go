package database

import (
	"os"
	"path/filepath"

	"github.com/bojand/ghz/web/model"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"    // enable the mysql dialect
	_ "github.com/jinzhu/gorm/dialects/postgres" // enable the postgres dialect
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // enable the sqlite3 dialect
)

const dbName = "../test/test.db"

// New creates a new wrapper for the gorm database framework.
func New(dialect, connection string, log bool) (*Database, error) {
	db, err := gorm.Open(dialect, connection)
	if err != nil {
		return nil, err
	}

	if err := createDirectoryIfSqlite(dialect, connection); err != nil {
		return nil, err
	}

	db.LogMode(log)

	db.AutoMigrate(
		new(model.Project),
		new(model.Report),
		new(model.Options),
		new(model.Detail),
		new(model.Histogram),
	)

	db.DB().SetMaxOpenConns(1)

	if dialect == "sqlite3" {
		db.DB().SetMaxOpenConns(10)

		db.Exec("PRAGMA foreign_keys = ON;")
	}

	return &Database{DB: db}, nil
}

func createDirectoryIfSqlite(dialect string, connection string) error {
	if dialect == "sqlite3" {
		if _, err := os.Stat(filepath.Dir(connection)); os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(connection), 0777); err != nil {
				return err
			}
		}
	}
	return nil
}

// Database is a wrapper for the gorm framework.
type Database struct {
	DB *gorm.DB
}

// Close closes the gorm database connection.
func (d *Database) Close() error {
	return d.DB.Close()
}
