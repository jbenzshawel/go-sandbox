package database

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"
)

type DbProvider func() (*sql.DB, error)

func Query(dbProvider DbProvider, stmt SelectStatement, result interface{}) (err error) {
	db, err := dbProvider()
	if err != nil {
		return err
	}
	defer func() {
		err = db.Close()
	}()

	err = stmt.Query(db, result)
	return
}

func ExecuteInsert(dbProvider DbProvider, stmt InsertStatement) (result sql.Result, err error) {
	db, err := dbProvider()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = db.Close()
	}()

	result, err = stmt.Exec(db)
	return
}
