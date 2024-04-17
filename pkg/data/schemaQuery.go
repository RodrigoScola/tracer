package data

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