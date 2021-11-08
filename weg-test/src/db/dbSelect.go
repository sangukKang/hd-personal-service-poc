package db

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
	kafka "weg-test/src/kafka"
)

type Item struct {
	id 			int64 `json:"id,omitempty"`
	name        string   `json:"name,omitempty"`
	quantity 	int64 `json:"quantity,omitempty"`
}


func (a Item) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *Item) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &a)
}

func SelectFileInfo() string{
	db, err := getDB()
	checkError(err)

	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database")

	//	Read rows from table.
	var b bytes.Buffer
	b.WriteString("SELECT json_agg(inventory) as name from inventory; ")
	b.WriteString(";")

	var str = b.String()
	//	var sql_statement = "SELECT * from inventory" + b.string()

	var res string
	err = db.QueryRow(str).Scan(&res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	return res
}

func SelectFileReq(userId string) string{

	db, err := getDB()
	checkError(err)

	err = db.Ping()
	checkError(err)
	fmt.Println("Successfully created connection to database")

	//	Read rows from table.
	var b bytes.Buffer
	b.WriteString("SELECT json_agg(inventory) as name from inventory ")

	fmt.Println("userId : ",userId)

	if(len(strings.TrimSpace(userId)) > 0){
		b.WriteString("where id = '")
		b.WriteString(userId)
		b.WriteString("';")
	}else{
		b.WriteString(";")
	}
	var str = b.String()
	//	var sql_statement = "SELECT * from inventory" + b.string()

	var res string
	err = db.QueryRow(str).Scan(&res)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	return res
}

func TestKafka() string{
	kafka.Producer()
	go kafka.Consumer()
	return ""
}

func SQLToJSON(rows *sql.Rows) ([]byte, error) {
    columns, err := rows.Columns()
	fmt.Println("columns : " , columns)
    if err != nil {
        return nil, fmt.Errorf("Column error: %v", err)
    }

    tt, err := rows.ColumnTypes()
	fmt.Println("tt : " , tt)
    if err != nil {
        return nil, fmt.Errorf("Column type error: %v", err)
    }

    values := make([]interface{}, len(tt))

	fmt.Println("values : " , values)

    data := make(map[string][]interface{})

	fmt.Println("data : " , data)

    for rows.Next() {
		
		fmt.Println("data : ???????")

        for i := range values {
            values[i] = ""
        }
        err = rows.Scan(values...)
        if err != nil {
            return nil, fmt.Errorf("Failed to scan values: %v", err)
        }
        for i, v := range values {
            data[columns[i]] = append(data[columns[i]], v)
        }
    }

	fmt.Println("data : " , data)

    return json.Marshal(data)
}
