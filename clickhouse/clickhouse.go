package clickhouse

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/tekyin/go-libs/errors"
)

var Conn *ChConnection

type ChConnection struct {
	DB driver.Conn
}

func Init(host string, port int, user string, password string, secure bool) (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Protocol: clickhouse.HTTP,
			Addr:     []string{fmt.Sprintf("%s:%d", host, port)},
			Auth: clickhouse.Auth{
				Database: "default",
				Username: user,
				Password: password,
			},
			ClientInfo: clickhouse.ClientInfo{
				Products: []struct {
					Name    string
					Version string
				}{
					{Name: "an-example-go-client", Version: "0.1"},
				},
			},
			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
			TLS: nil,
		})
	)

	if err != nil {
		fmt.Printf("open connection error: %s\n", err)
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		fmt.Println("ping error:", err)
		return nil, err
	}

	Conn = &ChConnection{DB: conn}
	res := RunQuery(ctx, "select version()")
	fmt.Println("ClickHouse version:", res)
	fmt.Println("Connected to ClickHouse at", host, "on port", port)
	return conn, nil
}

func RunQuery(ctx context.Context, query string, args ...any) []map[string]interface{} {
	rows, err := Conn.DB.Query(ctx, query, args...)
	errors.CheckError(err)
	defer func(rows driver.Rows) {
		err := rows.Close()
		errors.CheckError(err)
	}(rows)

	var result = make([]map[string]interface{}, 0)
	columns := rows.Columns()
	columnTypes := rows.ColumnTypes()

	for rows.Next() {
		valuePtrs := make([]interface{}, len(columns))

		// Create appropriate pointers based on column types
		for i, colType := range columnTypes {
			switch colType.DatabaseTypeName() {
			case "String", "FixedString":
				var s string
				valuePtrs[i] = &s
			case "Int8", "Int16", "Int32", "Int64":
				var n int64
				valuePtrs[i] = &n
			case "UInt8", "UInt16", "UInt32", "UInt64":
				var n uint64
				valuePtrs[i] = &n
			case "Float32", "Float64":
				var f float64
				valuePtrs[i] = &f
			case "Date", "DateTime", "DateTime64":
				var t interface{}
				valuePtrs[i] = &t
			default:
				var v interface{}
				valuePtrs[i] = &v
			}
		}

		err := rows.Scan(valuePtrs...)
		errors.CheckError(err)

		row := make(map[string]interface{})
		for i, col := range columns {
			switch v := valuePtrs[i].(type) {
			case *string:
				row[col] = *v
			case *int64:
				row[col] = *v
			case *uint64:
				row[col] = *v
			case *float64:
				row[col] = *v
			case *interface{}:
				row[col] = *v
			default:
				row[col] = v
			}
		}
		result = append(result, row)
	}

	return result
}
