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

var RolePermissions = newRolePermissionsTable("identity", "role_permissions", "")

type rolePermissionsTable struct {
	postgres.Table

	//Columns
	RoleID       postgres.ColumnInteger
	PermissionID postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type RolePermissionsTable struct {
	rolePermissionsTable

	EXCLUDED rolePermissionsTable
}

// AS creates new RolePermissionsTable with assigned alias
func (a RolePermissionsTable) AS(alias string) *RolePermissionsTable {
	return newRolePermissionsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new RolePermissionsTable with assigned schema name
func (a RolePermissionsTable) FromSchema(schemaName string) *RolePermissionsTable {
	return newRolePermissionsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new RolePermissionsTable with assigned table prefix
func (a RolePermissionsTable) WithPrefix(prefix string) *RolePermissionsTable {
	return newRolePermissionsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new RolePermissionsTable with assigned table suffix
func (a RolePermissionsTable) WithSuffix(suffix string) *RolePermissionsTable {
	return newRolePermissionsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newRolePermissionsTable(schemaName, tableName, alias string) *RolePermissionsTable {
	return &RolePermissionsTable{
		rolePermissionsTable: newRolePermissionsTableImpl(schemaName, tableName, alias),
		EXCLUDED:             newRolePermissionsTableImpl("", "excluded", ""),
	}
}

func newRolePermissionsTableImpl(schemaName, tableName, alias string) rolePermissionsTable {
	var (
		RoleIDColumn       = postgres.IntegerColumn("role_id")
		PermissionIDColumn = postgres.IntegerColumn("permission_id")
		allColumns         = postgres.ColumnList{RoleIDColumn, PermissionIDColumn}
		mutableColumns     = postgres.ColumnList{}
	)

	return rolePermissionsTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		RoleID:       RoleIDColumn,
		PermissionID: PermissionIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
