package helper

import (
	"database/sql"
	"log"
)

func CommitOrRollback(tx *sql.Tx) {
	err := recover()
	if err != nil {
		errorRollback := tx.Rollback()
		if errorRollback != nil {
			log.Println(errorRollback)
		}
	} else {
		errorCommit := tx.Commit()
		if errorCommit != nil {
			log.Println(errorCommit)
		}
	}
}
