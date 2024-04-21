package graph

import (
	"RodrigoScola/tracer/pkg/data"
	"errors"
	"fmt"
)

type GraphNode struct {
	Table data.Table   `json:"data"`
	To    []*GraphNode `json:"to"`
	From  []*GraphNode `json:"from"`
}



func (node *GraphNode) HasChild(childNode *GraphNode) *GraphNode {
	for _, cu := range node.To {
		if cu == childNode {
			return cu
		}
	}
	for _, cu := range node.From {
		if cu == childNode {

			return cu
		}
	}
	return nil
}

func (node *GraphNode) PrintColumns() {
	fmt.Println(">> ", node.Table.Name)
	if len(node.To) > 0 {
		fmt.Printf("\tConnects to: \n")
		for i := range node.To {
			fmt.Printf("\t%s\n", node.To[i].Table.Name)
		}
	}
	if len(node.From) == 0 {
		return
	}

	fmt.Println("\n>> is connected from:")
	for _, from_node := range node.From {
		fmt.Printf("\t%s ", from_node.Table.Name)

		col, err := from_node.Table.GetConnectingColumn(node.Table.Name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("using the col %s\n %s ", col.ColumnName, *col.ReferencedColumnName)
	}
}

func (g *GraphNode) IsInSlice(haystack []*GraphNode, needle *GraphNode) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

type Graph struct {
	Nodes []*GraphNode `json:"nodes"`
}

func (g *Graph) FindNode(tableName string) (*GraphNode, error) {

	for _, node := range g.Nodes {
		if node.Table.Name == tableName {
			return node, nil
		}
	}

	return nil, errors.New("could not find table")

}

func MakeGraph(tables []data.Table) (*Graph, error) {
	graph := &Graph{
		Nodes: make([]*GraphNode, 0),
	}

	for i := range tables {
		node := &GraphNode{
			Table: tables[i],
			To:    make([]*GraphNode, 0),
			From:  make([]*GraphNode, 0),
		}

		for j := range tables[i].Columns {
			tablename := tables[i].Columns[j].ReferencedTableName

			if tablename == nil {
				continue
			}

			referencedNode, err := graph.FindNode(*tablename)

			if err != nil {
				continue
			}

			if !node.IsInSlice(node.To, referencedNode) {
				node.To = append(node.To, referencedNode)
			}

			if !referencedNode.IsInSlice(referencedNode.From, node) {
				referencedNode.From = append(referencedNode.From, node)

			}

		}

		graph.Nodes = append(graph.Nodes, node)

	}
	return graph, nil
}

var allPaths [][]*GraphNode = make([][]*GraphNode, 0)

func walk(current *GraphNode, end *GraphNode, path []*GraphNode, seen map[*GraphNode]bool) {
	if current == nil {
		return
	}
	if current == end {
		allPaths = append(allPaths, append(path, current))
		return
	}

	if seen[current] {
		return
	}
	seenCopy := make(map[*GraphNode]bool)
	for k, v := range seen {
		seenCopy[k] = v
	}
	seenCopy[current] = true

	for i := range current.To {
		walk(current.To[i], end, append(path, current), seenCopy)
	}
	for i := range current.From {
		walk(current.From[i], end, append(path, current), seenCopy)
	}
}

func FindPath(start *GraphNode, end *GraphNode) ([]*GraphNode, error) {
	walk(start, end, make([]*GraphNode, 0), make(map[*GraphNode]bool))

	if len(allPaths) == 0 {
		return nil, errors.New("no path found")
	}
	smallest := allPaths[0]

	for i := range allPaths {
		if len(allPaths[i]) < len(smallest) {
			smallest = allPaths[i]
		}
	}

	return smallest, nil
}

func PrintPath(path []*GraphNode) {

	for i, item := range path {
		fmt.Print(item.Table.Name)
		if i >= 0 && i < len(path) {
			fmt.Print(" -> ")
		}

	}
	fmt.Println('\n')

}
