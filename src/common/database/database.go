package database

import (
	"database/sql"

	. "github.com/go-jet/jet/v2/postgres"

	"github.com/jbenzshawel/go-sandbox/common/cerror"
)

type DbProvider func() (*sql.DB, error)

func Query(dbProvider DbProvider, stmt SelectStatement, result interface{}) (err error) {
	db, err := dbProvider()
	if err != nil {
		return err
	}
	defer func() {
		closeErr := db.Close()
		err = cerror.CombineErrors(err, closeErr)
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
		closeErr := db.Close()
		err = cerror.CombineErrors(err, closeErr)
	}()

	result, err = stmt.Exec(db)
	return
}

func ExecuteUpdate(dbProvider DbProvider, stmt UpdateStatement) (result sql.Result, err error) {
	db, err := dbProvider()
	if err != nil {
		return nil, err
	}
	defer func() {
		closeErr := db.Close()
		err = cerror.CombineErrors(err, closeErr)
	}()

	result, err = stmt.Exec(db)
	return
}
