package main

import (
	"RodrigoScola/tracer/pkg/data"
	"encoding/json"
	"os"
)

func GetSchemaFromFile(filename string) (*data.Schema, error) {
	file, err := os.ReadFile(filename)

	if err != nil {
	  return nil, err
	}

	var schema data.Schema

	err =json.Unmarshal(file, &schema)

	if err != nil {
	  return nil, err
	}
	
	return &schema, nil

}