package data

import (
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Table struct {
	Name    string       `json:"name"`
	Columns []SchemaInfo `json:"columns"`
}

func (t *Table) GetConnectingColumn(tableName string) (*SchemaInfo, error) {
	for _, col := range t.Columns {
		if col.ReferencedTableName != nil && strings.Compare(*col.ReferencedTableName, tableName) == 0 {
		return &col, nil
		}
	}
	return nil, errors.New("could not find connecting column")
}

type Schema struct {
	Name   string  `json:"name"`
	Tables []Table `json:"tables"`

	hashmap map[string]*Table
}

func GetSchemaFromFile() (*Schema, error) {
	file, err := os.ReadFile("create_tables.json")

	if err != nil {
		return nil, err
	}

	var schema Schema

	err = json.Unmarshal(file, &schema)

	if err != nil {
		return nil, err
	}

	return &schema, nil

}

func GetSchemaFromDb(db *sqlx.DB, name string) (*Schema, error) {
	var info []SchemaInfo

	err := db.Select(&info, getSchemaQuery())
	if err != nil {
		return nil, err
	}

	schema := &Schema{
		Name:    name,
		Tables:  make([]Table, 0),
		hashmap: make(map[string]*Table),
	}

	for _, i := range info {
		foundTable := false
		for index, table := range schema.Tables {
			if i.TableName == table.Name {
				foundTable = true
				schema.Tables[index].Columns = append(table.Columns, i)
			}
		}
		if !foundTable {
			column := make([]SchemaInfo, 0)
			column = append(column, i)
			table := Table{
				Name:    i.TableName,
				Columns: column,
			}
			schema.hashmap[table.Name] = &table
			schema.Tables = append(schema.Tables, table)
		}
	}
	return schema, nil
}
