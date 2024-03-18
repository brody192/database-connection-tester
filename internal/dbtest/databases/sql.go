package databases

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/xo/dburl"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Sql(driverName, mysqlURL string) (time.Duration, error) {
	dsn, err := dburl.Parse(mysqlURL)
	if err != nil {
		return 0, fmt.Errorf("dburl.Parse error: %w", err)
	}

	query := dsn.Query()

	query.Set("connect_timeout", "1")
	query.Set("timeout", "1")

	dsn.RawQuery = query.Encode()

	if driverName == "postgres" {
		driverName = "pgx"
	}

	sT := time.Now()

	db, err := sql.Open(driverName, dsn.DSN)
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

		if err := db.PingContext(ctx); err != nil {
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
