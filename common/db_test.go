package common

import (
	"database/sql"
	"os"
	"testing"
)

func createDb(t *testing.T) (db *sql.DB, err error) {
	if os.Getenv("VROPS_DB_DRIVER") == "" {
		t.Skip("Skipping. Set VROPS_DB_DRIVER.")
	}
	if os.Getenv("VROPS_DB_FILENAME") == "" {
		t.Skip("Skipping. Set VROPS_DB_FILENAME.")
	}

	driver := os.Getenv("VROPS_DB_DRIVER")
	filename := os.Getenv("VROPS_DB_FILENAME")

	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}

	db, err = sql.Open(driver, filename)
	if err == nil {
		var createKvstoreSQL = `
			CREATE TABLE IF NOT EXISTS 'KVStore' (
		  	'ID' INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE
		  	, 'Key' TEXT
		  	, 'Value' TEXT
				);
		`
		_, err = db.Exec(createKvstoreSQL)
	}

	return
}

func deleteDb(t *testing.T) (err error) {
	filename := os.Getenv("VROPS_DB_FILENAME")

	if _, err = os.Stat(filename); err == nil {
		os.Remove(filename)
	}

	return
}

func TestSaveKVStoreManual(t *testing.T) {
	db, err := createDb(t)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	kvs := &KVStore{
		Key:   "TestSaveKVStoreManualOne",
		Value: "True"}
	if err := kvs.Save(db); err != nil {
		t.Fatal(err)
	}

	deleteDb(t)
}
