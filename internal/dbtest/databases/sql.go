package databases

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/xo/dburl"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func Sql(driverName, mysqlURL string) (time.Duration, error) {
	dsn, err := dburl.Parse(mysqlURL)
	if err != nil {
		return 0, fmt.Errorf("dburl.Parse error: %w", err)
	}

	query := dsn.Query()

	query.Set("connect_timeout", "3")
	query.Set("timeout", "3")

	dsn.RawQuery = query.Encode()

	sT := time.Now()

	db, err := sql.Open(driverName, dsn.DSN)
	if err != nil {
		return time.Since(sT), fmt.Errorf("sql.Open error: %w", err)
	}

	defer db.Close()

	maxTestTime := time.Now().Add(8 * time.Second)

	if err := retry.Do(func() error {
		return db.Ping()
	},
		retry.LastErrorOnly(true),
		retry.Attempts(300),
		retry.Delay(10),
		retry.RetryIf(func(err error) bool {
			return !time.Now().After(maxTestTime)
		}),
	); err != nil {
		return time.Since(sT), fmt.Errorf("db.Ping error: %w", err)
	}

	return time.Since(sT), nil
}
