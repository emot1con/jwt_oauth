package helper

import (
	"database/sql"

	"github.com/sirupsen/logrus"
)

func CommitOrRollback(tx *sql.Tx) {
	logrus.Info("rollback panic")
	if err := recover(); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			panic(errRollback)
		}
	}
}

func BeginTransaction(db *sql.DB) (*sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}
