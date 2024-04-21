package main

import (
	"RodrigoScola/tracer/pkg/data"
	"RodrigoScola/tracer/pkg/data/graph"
	"RodrigoScola/tracer/pkg/query"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

// db, err := sqlx.Connect("mysql", "root@tcp(127.0.0.1:3306)/tracer")
// if err != nil {
// 	panic(err)
// }

type Input struct {
	InputTable  string
	OutputTable string
	OnColumn    string
	Value       any
}

func main() {

	db, err := sqlx.Connect("mysql", "root:password@(127.0.0.1:3306)/tracer")
	if err != nil {
		panic(err)
	}
	input := Input{
		InputTable:  "sku",
		OutputTable: "products",
		OnColumn:    "id",
		Value:       22,
	}

	_ = input

	schema, err := data.GetSchemaFromDb(db, "tracer")

	// schema, err := GetSchemaFromFile("create_tables.json")

	if err != nil {
		panic(err)
	}

	g, err := graph.MakeGraph(schema.Tables)

	if err != nil {
		panic(err)
	}

	adspath, err := g.FindNode(input.InputTable)
	if err != nil {
		panic(err)
	}

	catpath, err := g.FindNode(input.OutputTable)
	if err != nil {
		panic(err)
	}
	path, err := graph.FindPath(adspath, catpath)
	if err != nil {
		panic(err)
	}

	// var completeQuery strings.Builder

	var S = query.NewQuery()

	graph.PrintPath(path)

	skip_next := false

	for i := range path {
		if i == 0 {
			primary, err:= path[i].Table.GetPrimary()
			if err != nil {
			  continue
			}

			S.AddSelect(path[i].Table.Name, primary.ColumnName)
			continue
		}
		if skip_next {
			skip_next = false
			continue
		}
		nextNode := i - 1

		if nextNode > len(path)-1 {
			nextNode = i + 1
		}
		connectingNode := path[nextNode]

		fmt.Printf("path[i].Table.Name: %v\n", path[i].Table.Name)

		conn, err := path[i].Table.GetConnectingColumn(connectingNode.Table.Name)

		if err != nil {
			conn, err = connectingNode.Table.GetConnectingColumn(path[i].Table.Name)
			if err != nil {
				panic(err)
			}

			if strings.Compare(input.InputTable, *conn.ReferencedTableName) == 0 {
				continue
			}

			S.AddSelect(*conn.ReferencedTableName, *conn.ReferencedColumnName)

			S.AddJoin(
				"inner join",
				*conn.ReferencedTableName,
				*conn.ReferencedColumnName,
				conn.TableName,
				conn.ColumnName,
			)
		}

		if strings.Compare(input.InputTable, conn.TableName) == 0 {
			continue
		}

		S.AddSelect(conn.TableName, conn.ColumnName)
		S.AddJoin(
			"inner join",
			conn.TableName,
			conn.ColumnName,
			*conn.ReferencedTableName,
			*conn.ReferencedColumnName,
		)


	}
	var builder strings.Builder

	// selectstr:= ""


	builder.WriteString(fmt.Sprintf("select %s from %s ",
		strings.Join(S.Rows, ","),

		input.InputTable))

	for _, v := range S.Joins {
		builder.WriteString(fmt.Sprintf(" %s %s on %s.%s = %s.%s",
			v.Type,
			v.Column.TableName,
			v.Column.TableName,
			v.Column.ColumnName,
			v.Column.ReferenceTableName,
			v.Column.ReferenceColumnName,
		))
	}

	builder.WriteString(fmt.Sprintf(" where %s.%s = %v", input.InputTable, input.OnColumn, input.Value))
	fmt.Println(builder.String())


	rows, err := db.Queryx(builder.String())
	if err != nil {
		panic(err)
	}

	items := make(map[string][]string)
	_ = items

	for rows.Next() {

		results := make(map[string]interface{})
		err = rows.MapScan(results)
		if err != nil {
			panic(err)
		}

		for key, value := range results {
			if _, ok := items[key]; !ok {
				items[key] = make([]string, 0)
			}

			strvalue, ok := value.(string)
			if ok {
				value = strvalue
			}

			items[key] = append(items[key], fmt.Sprintf("%v", value))
		}
	}
	tableItems := make([][]string, 0)

	for key, value := range items {
		curr := make([]string, 0)
		curr = append(curr, key)
		curr = append(curr, value...)

		tableItems = append(tableItems, curr)

	}
	fmt.Println(tableItems)

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("99"))).
		Rows(tableItems...)

	fmt.Println(t)

}
