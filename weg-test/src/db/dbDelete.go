package db

import (
	"fmt"

	_ "github.com/lib/pq"
)

func Delete(id string) {
	// Initialize connection object.
	db, err := getDB()
	checkError(err)

	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database")

	// Delete some data from table.
	sql_statement := "DELETE FROM inventory WHERE id = $1;"
	_, err = db.Exec(sql_statement, id)
	checkError(err)
	fmt.Println("Deleted 1 row of data")
}
