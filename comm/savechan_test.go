package comm

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bruceadowns/vropsbot/common"
)

func createDb(t *testing.T) (db *sql.DB, err error) {
	driver := os.Getenv("VROPS_DB_DRIVER")
	filename := os.Getenv("VROPS_DB_FILENAME")

	if driver == "" {
		t.Skip("Skipping. Set VROPS_DB_DRIVER.")
	}
	if filename == "" {
		t.Skip("Skipping. Set VROPS_DB_FILENAME.")
	}

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

func yield(t *testing.T) {
	t.Log("yield to channel goroutine")
	time.Sleep(time.Second)
}

func testSaveKVStore(t *testing.T, db *sql.DB, bufsize int, rows int) {
	t.Logf("Test %dx%d", bufsize, rows)

	chSave := SaveChan(bufsize, db)
	for i := 0; i < rows; i++ {
		chSave <- &common.KVStore{
			Key:   fmt.Sprintf("testSaveKVStore-%d-%d", bufsize, i),
			Value: fmt.Sprintf("%d", rows)}
	}
	yield(t)
}

func TestSaveKVStoreChannel(t *testing.T) {
	db, err := createDb(t)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testSaveKVStore(t, db, 0, 10)
	testSaveKVStore(t, db, 10, 100)
	testSaveKVStore(t, db, 10, 1000)
	testSaveKVStore(t, db, 10, 10000)

	deleteDb(t)
}

/*
func TestSaveSlackItemTypeChannel(t *testing.T) {
	chSave := SaveChan(0, db)

	t.Log("send slackitemtype")
	chSave <- &SlackItemType{
		Name: ""}
	yield()
}

func TestSaveSlackItemChannel(t *testing.T) {
	chSave := SaveChan(0, db)

	t.Log("send slackitem")
	chSave <- &SlackItem{
		Name: "",
		Type: 0}
	yield()
}

func TestSavePullActionsChannel(t *testing.T) {
	chSave := SaveChan(0, db)

	t.Log("send pullactions")
	chSave <- &PullActions{
		When:     "",
		Channel:  "",
		User:     "",
		Request:  "",
		Response: "",
		IsError:  0}
	yield()
}

func TestSavePushActionsChannel(t *testing.T) {
	chSave := SaveChan(0, db)

	t.Log("send pullactions")
	chSave <- &PushActions{
		When:       "",
		PluginName: "",
		Channel:    "",
		Request:    "",
		Response:   "",
		Before:     "",
		After:      "",
		IsError:    0}
	yield()
}
*/
