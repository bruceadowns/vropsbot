package comm

import (
	"database/sql"
	"log"
	"time"
)

// HeartbeatChan persists a heartbeat
func HeartbeatChan(size int, db *sql.DB) (ch chan string) {
	ch = make(chan string, size)

	go func() {
		for {
			select {
			case name := <-ch:
				_, err := db.Exec(
					`INSERT OR REPLACE INTO 'Heartbeat' (
							'Name'
							, 'When')
						VALUES (?, ?)`,
					name,
					time.Now())
				if err != nil {
					log.Printf("Error updating heartbeat for %s [%s]", name, err)
				}
			}
		}
	}()

	return
}
