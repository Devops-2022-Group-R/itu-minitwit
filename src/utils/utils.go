package utils

import (
	"database/sql"
	"log"

)

type Row = map[string]interface{}

// Queries the database and returns a list of dictionaries.
func QueryDb(db *sql.DB, query string, args ...interface{}) []Row {
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatalf("queryDb failure: %v", err)
	}
	defer rows.Close()

	columnNames, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	results := make([]Row, 0)
	for rows.Next() {
		row := make(Row)

		values := make([]interface{}, len(columnNames))
		valuesRef := make([]interface{}, len(columnNames))
		for i := 0; i < len(columnNames); i++ {
			valuesRef[i] = &values[i]
		}
		err := rows.Scan(valuesRef...)
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < len(columnNames); i++ {
			row[columnNames[i]] = values[i]
		}
		results = append(results, row)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return results
}
