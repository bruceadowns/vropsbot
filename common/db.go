package common

import (
	"database/sql"
	"log"
	"os"
	"time"

	// importing sqlite3 driver for init side effect
	_ "github.com/mattn/go-sqlite3"
)

// KVStore holds the KVStore db struct
type KVStore struct {
	Key   string
	Value string
}

// SlackItem holds the SlackItem db struct
type SlackItem struct {
	ID   string
	Name string
	Type int
}

// PullActions stores push actions
type PullActions struct {
	When     string
	Channel  string
	User     string
	Request  string
	Response string
	IsError  int
}

// PushActions stores push actions
type PushActions struct {
	When       string
	PluginName string
	Channel    string
	Request    string
	Response   string
	Before     string
	After      string
	IsError    int
}

// Save persists the KVStore structure
func (kvs *KVStore) Save(db *sql.DB) (err error) {
	_, err = db.Exec(
		"INSERT INTO 'KVStore' ('Key', 'Value') VALUES (?, ?)",
		kvs.Key, kvs.Value)
	return
}

// Save persists the SlackItem structure
func (si *SlackItem) Save(db *sql.DB) (err error) {
	_, err = db.Exec(
		"INSERT OR IGNORE INTO 'SlackItem' ('ID', 'Name', 'Type') VALUES (?, ?, ?)",
		si.ID, si.Name, si.Type)
	return
}

// Save persists the PushActions structure
func (pa *PushActions) Save(db *sql.DB) (err error) {
	_, err = db.Exec(
		`INSERT INTO 'PushActions' (
			'When'
			, 'PluginName'
			, 'Channel'
			, 'Request'
			, 'Response'
			, 'Before'
			, 'After'
			, 'IsError')
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		time.Now(),
		pa.PluginName,
		pa.Channel,
		pa.Request,
		pa.Response,
		pa.Before,
		pa.After,
		pa.IsError)

	return
}

// InitDatabase creates/opens and initializes the sqlite3 database
func InitDatabase() (db *sql.DB, err error) {
	// check for db env variables
	// fall back to config file

	driver := os.Getenv("VROPS_DB_DRIVER")
	if driver == "" {
		config, _ := NewConfig()
		driver = config.DB.Driver
	}

	filename := os.Getenv("VROPS_DB_FILENAME")
	if filename == "" {
		config, _ := NewConfig()
		filename = config.DB.Filename
	}

	log.Printf("sql driver: %s", driver)
	log.Printf("sql filename: %s", filename)

	if _, err = os.Stat(filename); err == nil {
		db, err = sql.Open(driver, filename)
	} else {
		db, err = sql.Open(driver, filename)
		if err == nil {
			_, err = db.Exec(createDbSQL)
		}
	}

	return
}
