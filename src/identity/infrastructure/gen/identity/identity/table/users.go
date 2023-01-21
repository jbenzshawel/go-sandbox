//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var Users = newUsersTable("identity", "users", "")

type usersTable struct {
	postgres.Table

	//Columns
	UserID        postgres.ColumnInteger
	UserUUID      postgres.ColumnString
	FirstName     postgres.ColumnString
	LastName      postgres.ColumnString
	Email         postgres.ColumnString
	EmailVerified postgres.ColumnBool
	Enabled       postgres.ColumnBool
	CreatedAt     postgres.ColumnTimestamp
	LastUpdatedAt postgres.ColumnTimestamp

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type UsersTable struct {
	usersTable

	EXCLUDED usersTable
}

// AS creates new UsersTable with assigned alias
func (a UsersTable) AS(alias string) *UsersTable {
	return newUsersTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new UsersTable with assigned schema name
func (a UsersTable) FromSchema(schemaName string) *UsersTable {
	return newUsersTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new UsersTable with assigned table prefix
func (a UsersTable) WithPrefix(prefix string) *UsersTable {
	return newUsersTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new UsersTable with assigned table suffix
func (a UsersTable) WithSuffix(suffix string) *UsersTable {
	return newUsersTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newUsersTable(schemaName, tableName, alias string) *UsersTable {
	return &UsersTable{
		usersTable: newUsersTableImpl(schemaName, tableName, alias),
		EXCLUDED:   newUsersTableImpl("", "excluded", ""),
	}
}

func newUsersTableImpl(schemaName, tableName, alias string) usersTable {
	var (
		UserIDColumn        = postgres.IntegerColumn("user_id")
		UserUUIDColumn      = postgres.StringColumn("user_uuid")
		FirstNameColumn     = postgres.StringColumn("first_name")
		LastNameColumn      = postgres.StringColumn("last_name")
		EmailColumn         = postgres.StringColumn("email")
		EmailVerifiedColumn = postgres.BoolColumn("email_verified")
		EnabledColumn       = postgres.BoolColumn("enabled")
		CreatedAtColumn     = postgres.TimestampColumn("created_at")
		LastUpdatedAtColumn = postgres.TimestampColumn("last_updated_at")
		allColumns          = postgres.ColumnList{UserIDColumn, UserUUIDColumn, FirstNameColumn, LastNameColumn, EmailColumn, EmailVerifiedColumn, EnabledColumn, CreatedAtColumn, LastUpdatedAtColumn}
		mutableColumns      = postgres.ColumnList{UserUUIDColumn, FirstNameColumn, LastNameColumn, EmailColumn, EmailVerifiedColumn, EnabledColumn, CreatedAtColumn, LastUpdatedAtColumn}
	)

	return usersTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UserID:        UserIDColumn,
		UserUUID:      UserUUIDColumn,
		FirstName:     FirstNameColumn,
		LastName:      LastNameColumn,
		Email:         EmailColumn,
		EmailVerified: EmailVerifiedColumn,
		Enabled:       EnabledColumn,
		CreatedAt:     CreatedAtColumn,
		LastUpdatedAt: LastUpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
