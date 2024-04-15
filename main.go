package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

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

type Table struct {
	Name    string       `json:"name"`
	Columns []SchemaInfo `json:"columns"`
}

type Schema struct {
	Name    string  `json:"name"`
	Tables  []Table `json:"tables"`
	hashmap map[string]*Table
}

type Input struct {
	InputTable  string
	OutputTable string
	OnColumn    string
	Value       any
}

func main() {

	input := Input{
		InputTable:  "ads",
		OutputTable: "products",
		OnColumn:    "id",
		Value:       957,
	}

	_ = input

	db, err := sqlx.Connect("mysql", "root:password@tcp(127.0.0.1:3306)/tracer")
	if err != nil {
		panic(err)
	}

	query := getSchemaQuery()
	// tables, err := showCreateTables(db)

	// if err != nil {
	// 	panic(err)
	// }

	var info []SchemaInfo
	err = db.Select(&info, query)
	if err != nil {
		panic(err)
	}

	schema := Schema{
		Name:    "tracer",
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

	// for _, table := range schema.Tables {
	// 	for _, row := range table.Columns {
	// 		fmt.Println("> COL", row.ColumnName)
	// 		fmt.Println(row.ColumnType)
	// 		// fmt.Println(row)
	// 		// if row.ReferencedColumnName != nil {
	// 		// 	fmt.Print("\n > References ", *row.ReferencedColumnName)
	// 		// 	fmt.Print(" of ", *row.ReferencedTableName)
	// 		// 	fmt.Print("\n")
	// 		// }
	// 	}
	// }

	jsonTable, err := json.Marshal(schema)
	if err != nil {
		panic(err)
	}
	os.WriteFile("create_tables.json", jsonTable, os.ModeAppend)

	graph := MakeGraph(schema.Tables)

	graph.PrintGraph()
}
