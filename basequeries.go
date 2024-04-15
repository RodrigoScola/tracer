package main

import (
	"github.com/jmoiron/sqlx"
)

func createTables(db *sqlx.DB) {
	db.MustExec(`CREATE TABLE if not exists author_specification (
		cityId INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100),
		state INT
	) `)

	db.MustExec(`CREATE TABLE if not exists author (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100),
		lastName varchar(100), 
		age INT
	) `)

	db.Exec(`create table if not exists category (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100)
		)`)

	db.Exec(`create table if not exists genre (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100)
		)`)
	db.MustExec(`CREATE TABLE if not exists book (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100),
		color VARCHAR(7),
		categoryId INT,
		authorId INT,
		genreId INT

	) `)
}

func showCreateTables(db *sqlx.DB) ([]string, error) {

	query := `
	select TABLE_NAME from
	INFORMATION_SCHEMA.TABLES
	where
	table_schema = DATABASE()

	 `

	var tables []string
	createTables := []string{}

	err := db.Select(&tables, query)
	if err != nil {
		return []string{}, err
	}

	for _, table := range tables {
		var tableName, createTable string
		err := db.QueryRow(`show create table `+table).Scan(&tableName, &createTable)

		if err != nil {
			return []string{}, err
		}
		createTables = append(createTables, createTable)
	}
	return createTables, nil

}

func getSchemaQuery() string {
	return `SELECT 
        c.column_name,
        c.column_type,
        t.table_name,
        t.table_schema,
        kcu.referenced_table_name,
        kcu.referenced_column_name,
        CASE WHEN kcu.constraint_name = 'PRIMARY' THEN true ELSE false END AS is_primary
    FROM 
        information_schema.tables t 
    JOIN 
        information_schema.columns c 
    ON 
        t.table_name = c.table_name 
        AND t.table_catalog = c.table_catalog 
    LEFT JOIN 
        information_schema.key_column_usage kcu 
    ON 
        t.table_name = kcu.table_name 
        AND t.table_schema = kcu.table_schema 
        AND c.column_name = kcu.column_name
    WHERE 
        t.table_schema = 'tracer'`
}
