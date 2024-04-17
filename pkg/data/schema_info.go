package data

import (
	"errors"
	"fmt"
)

type SchemaInfo struct {
	ColumnName           string  `db:"COLUMN_NAME" json:"name"`
	ColumnType           string  `db:"COLUMN_TYPE" json:"type"`
	TableName            string  `db:"TABLE_NAME" json:"tableName"`
	IsPrimary            bool    `db:"is_primary"`
	TableSchema          string  `db:"TABLE_SCHEMA" json:"-"`
	ReferencedTableName  *string `db:"REFERENCED_TABLE_NAME" json:"referencedTable,omitempty"`
	ReferencedColumnName *string `db:"REFERENCED_COLUMN_NAME" json:"referencedColumn,omitempty"`
}



func (s *SchemaInfo) GetReferenceTable(tables *[]Table) (*Table, error) {
	if s.ReferencedTableName == nil {
		return nil, errors.New("need to reference a table")
	}
	fmt.Println(*s.ReferencedTableName)

	for _, table := range *tables {
		if table.Name == *s.ReferencedTableName {
			return &table, nil
		}
	}
	return nil, errors.New("no table found")
}

func (s *SchemaInfo) GetTable(schema *Schema) (*Table, error) {

	if schema.hashmap[s.TableName] == nil {
		return nil, errors.New("table not found")
	}
	return schema.hashmap[s.TableName], nil
}