package trino

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tekyin/go-libs/errors"
	_ "github.com/trinodb/trino-go-client/trino"
)

var Conn *Connection

type Connection struct {
	DB *sql.DB
}

func InitConnection(user string, host string, port int) error {
	connStr := fmt.Sprintf("http://%s@%s:%d", user, host, port)
	db, err := sql.Open("trino", connStr)
	errors.CheckError(err)

	Conn = &Connection{DB: db}
	result := RunQuery(context.Background(), "Select version() as version")
	fmt.Println("Connected to Trino at", host, "on port", port, " Trino version: ", result[0]["version"])
	return nil
}

func RunQuery(ctx context.Context, query string, args ...any) []map[string]interface{} {
	rows, err := Conn.DB.QueryContext(ctx, query, args...)
	errors.CheckError(err)
	defer func(rows *sql.Rows) {
		err := rows.Close()
		errors.CheckError(err)
	}(rows)

	var result = make([]map[string]interface{}, 0)

	columns, err := rows.Columns()
	errors.CheckError(err)

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for rows.Next() {
		for i := 0; i < len(columns); i++ {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		errors.CheckError(err)
		row := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			row[col] = v
		}
		result = append(result, row)
	}
	return result
}
