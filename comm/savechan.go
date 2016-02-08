package comm

import (
	"database/sql"
	"log"
)

// Persist is the interface all persistable object must implement
type Persist interface {
	Save(db *sql.DB) error
}

// SaveChan creates and returns a buffered,
// asynchronous channel used to save structures
func SaveChan(size int, db *sql.DB) (ch chan Persist) {
	ch = make(chan Persist, size)

	go func() {
		for {
			select {
			case iface := <-ch:
				if err := iface.Save(db); err != nil {
					log.Printf("Error occurred saving object: %v", err)
				}
			}
		}
	}()

	return
}
