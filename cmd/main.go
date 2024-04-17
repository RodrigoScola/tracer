package main

import (
	"RodrigoScola/tracer/pkg/data/graph"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// db, err := sqlx.Connect("mysql", "root@tcp(127.0.0.1:3306)/tracer")
// if err != nil {
// 	panic(err)
// }

// schema, err := data.GetSchemaFromDb(db, "tracer")

// if err != nil {
//   panic(err)
// }

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

	schema, err := GetSchemaFromFile("create_tables.json")

	if err != nil {
		panic(err)
	}

	g, err := graph.MakeGraph(schema.Tables)

	if err != nil {
		panic(err)
	}

	adspath, err := g.FindNode("ads")
	if err != nil {
		panic(err)
	}

	catpath, err := g.FindNode("category_specification")
	if err != nil {
		panic(err)
	}

	path, err := graph.FindPath(adspath, catpath)
	if err != nil {
		panic(err)
	}

	var builder strings.Builder

	for i, node := range path {
		if i == len(path)-1 {

			fmt.Println("last line")
			continue
		}

		nextNode := path[i+1]

		fmt.Println(node.Table.Name, "->", nextNode.Table.Name)

		connecting, err := node.Table.GetConnectingColumn(nextNode.Table.Name)
		if err != nil {
			connecting, err = nextNode.Table.GetConnectingColumn(node.Table.Name)
			if err != nil {
				panic(err)
			}
		}
		if i == 0 {
			builder.WriteString("select * from ")
			builder.WriteString(node.Table.Name)
		}

		addJoin(&connectTable{
			Name: node.Table.Name,
			Col:  connecting.ColumnName,
		},
			&connectTable{
				Name: nextNode.Table.Name,
				Col:  *connecting.ReferencedColumnName,
			},
			&builder,
		)
		continue

	}
	fmt.Println(builder.String())

	_ = g

}

type connectTable struct {
	Name string
	Col  string
}

func addJoin(mainTable *connectTable, connectingTable *connectTable, builder *strings.Builder) {
	builder.WriteString(" inner join ")
	builder.WriteString(connectingTable.Name)
	builder.WriteString(" on ")
	builder.WriteString(mainTable.Name + "." + mainTable.Col)
	builder.WriteString(" = ")
	builder.WriteString(connectingTable.Name + "." + connectingTable.Col)
}
