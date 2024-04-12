package databases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/xo/dburl"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func Sql(driverName, sqlURL string) (time.Duration, error) {
	var dsn string

	url, err := dburl.Parse(sqlURL)
	if err != nil && err == dburl.ErrUnknownDatabaseScheme {
		// if dburl doesn't recognize the database schema, use the url as-is anyway
		dsn = sqlURL
	} else if err != nil {
		return 0, fmt.Errorf("dburl.Parse error: %w", err)
	} else {
		query := url.Query()

		query.Set("connect_timeout", "1")
		query.Set("timeout", "1")

		url.RawQuery = query.Encode()

		dsn = url.DSN
	}

	if driverName == "postgres" {
		driverName = "pgx"
	}

	sT := time.Now()

	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return time.Since(sT), fmt.Errorf("sql.Open error: %w", err)
	}

	defer db.Close()

	var lastErr error

	maxTestTime := time.Now().Add(8 * time.Second)

	for range 10 {
		if time.Now().After(maxTestTime) {
			if lastErr == nil {
				lastErr = context.DeadlineExceeded
			}
			break
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		_, err := db.QueryContext(ctx, "SELECT 1")
		if err != nil {
			lastErr = err
		} else {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	if lastErr != nil {
		return time.Since(sT), fmt.Errorf("db.Ping error: %w", lastErr)
	}

	return time.Since(sT), nil
}
