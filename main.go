// package main

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	_ "github.com/alexbrainman/odbc"
// )

// func main() {
// 	connString := "DSN=demodata;UID=demo;PWD=demo"

// 	db, err := sql.Open("odbc", connString)
// 	if err != nil {
// 		log.Fatal("Error connecting to the database: ", err)
// 	}
// 	defer db.Close()

// 	rows, err := db.Query("SELECT * FROM Student")
// 	if err != nil {
// 		log.Fatal("Error executing query: ", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var (
// 			column1 string
// 			column2 string
// 			column3 string
// 			column4 string
// 			column5 string
// 			column6 string
// 			column7 string
// 			column8 string
//
// 		)

//
// 		err := rows.Scan(&column1, &column2, &column3, &column4, &column5, &column6, &column7, &column8)
// 		if err != nil {
// 			log.Fatal("Error scanning row: ", err)
// 		}

// 		fmt.Printf("Column1: %s, Column2: %s, Column3: %s, Column4: %s, Column5: %s, Column6: %s, Column7: %s, Column8: %s\n", column1, column2, column3, column4, column5, column6, column7, column8)
// 	}

//		if err = rows.Err(); err != nil {
//			log.Fatal("Error iterating rows: ", err)
//		}
//	}
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

	_ "github.com/alexbrainman/odbc"
)

type Record struct {
	Field1 string `json:"field1"`
	Field2 string `json:"field2"`
	Field3 string `json:"field3"`
	Field4 string `json:"field4"`
	Field5 string `json:"field5"`
	Field6 string `json:"field6"`
	Field7 string `json:"field7"`
	Field8 string `json:"field8"`
}

// struct stores request body
type RequestBody struct {
	DSN   string `json:"dsn"`
	UID   string `json:"uid"`
	PWD   string `json:"pwd"`
	Query string `json:"query"`
}

// / opens connection and fetches table
func fetchRecords(w http.ResponseWriter, r *http.Request) {

	var reqBody RequestBody
	// Decode the JSON body of the request
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// tableName := r.URL.Query().Get("table")
	// if tableName == "" {
	// 	http.Error(w, "Table name is required", http.StatusBadRequest)
	// 	return
	// }

	db, err := sql.Open("odbc", fmt.Sprintf("DSN=%s;UID=%s;PWD=%s", reqBody.DSN, reqBody.UID, reqBody.PWD))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Printf("tableName :: %s", reqBody.Query)
	rows, err := db.Query(reqBody.Query) //"SELECT * FROM %s", reqBody.TableName
	if err != nil {
		log.Fatal(err)
	}
	// defer rows.Close()

	// var records []Record
	// for rows.Next() {
	// 	var rec Record
	// 	if err := rows.Scan(&rec.Field1, &rec.Field2, &rec.Field3, &rec.Field4, &rec.Field5, &rec.Field6, &rec.Field7, &rec.Field8); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	records = append(records, rec)
	// }
	// if err := rows.Err(); err != nil {
	// 	log.Fatal(err)
	// }

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	values := make([]interface{}, len(columns))
	record := make(map[string]interface{})

	for i := range values {
		values[i] = new(interface{})
	}

	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			log.Fatal(err)
		}

		for i, colName := range columns {
			val := reflect.ValueOf(values[i]).Elem().Interface()
			record[colName] = val
		}

		json.NewEncoder(w).Encode(record)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/fetch", fetchRecords)
	fmt.Println("Server is running on port 8082...")
	log.Fatal(http.ListenAndServe(":8082", nil))
}
