package main

import (
	"errors"
	"fmt"
	"strings"
)

type GraphNode struct {
	Data Table        `json:"data"`
	To   []*GraphNode `json:"to"`
	From []*GraphNode `json:"from"`
}

func (node *GraphNode) PrintColumns() {
	fmt.Println(">> NODE: ", node.Data.Name)
	fmt.Println(">> TO:")
	for i := range node.To {
		fmt.Println(node.To[i].Data.Name)
	}
	fmt.Println(">> FROM:")
	for i := range node.From {
		fmt.Println(node.From[i].Data.Name)
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
	Nodes []GraphNode `json:"nodes"`
}

func (g *Graph) FindNode(tableName string) (*GraphNode, error) {

	for i := range g.Nodes {
		if g.Nodes[i].Data.Name == tableName {
			return &g.Nodes[i], nil
		}
	}

	return nil, errors.New("could not find table")

}

func MakeGraph(tables []Table) *Graph {
	graph := &Graph{
		Nodes: make([]GraphNode, 0),
	}

	for i := range tables {
		node := GraphNode{
			Data: tables[i],
			To:   make([]*GraphNode, 0),
			From: make([]*GraphNode, 0),
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

			if !referencedNode.IsInSlice(referencedNode.From, &node) {
				referencedNode.From = append(referencedNode.From, &node)

			}

		}

		graph.Nodes = append(graph.Nodes, node)

	}

	// first pass to add everything

	return graph
}

func walk(current *GraphNode, end *GraphNode, seen map[*GraphNode]bool, path []*GraphNode) ([]*GraphNode, error) {
	newPath := make([]*GraphNode, len(path))
	copy(newPath, path)
	newPath = append(newPath, current)

	if seen[current] {
		return nil, nil
	}
	seen[current] = true

	if current == end {
		return newPath, nil
	}

	for _, node := range append(current.From, current.To...) {
		if foundPath, _ := walk(node, end, seen, newPath); foundPath != nil {
			return foundPath, nil
		}
	}

	return nil, nil
}

func (g *Graph) FindPath(start *GraphNode, end *GraphNode) ([]*GraphNode, error) {

	found, err := walk(start, end, make(map[*GraphNode]bool), make([]*GraphNode, 0))

	if err != nil {
		return nil, nil
	}

	return found, nil

}

func (g *Graph) PrintGraph(path []*GraphNode, level int) {
	for _, node := range g.Nodes {
		fmt.Print("-> ", node.Data.Name, "\n")
		strings.Repeat(" -> ", level)

	}
}

func (g *Graph) PrintPath(path []*GraphNode) {

	for i, item := range path {
		fmt.Print(item.Data.Name)
		if i >= 0 && i < len(path) {
			fmt.Print(" -> ")
		}

	}
	fmt.Println('\n')

}
