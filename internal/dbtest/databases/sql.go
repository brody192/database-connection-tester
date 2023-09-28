package databases

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/xo/dburl"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func Sql(driverName, mysqlURL string) (time.Duration, error) {
	dsn, err := dburl.Parse(mysqlURL)
	if err != nil {
		return 0, fmt.Errorf("dburl.Parse error: %w", err)
	}

	sT := time.Now()

	db, err := sql.Open(driverName, dsn.DSN)
	if err != nil {
		return 0, fmt.Errorf("sql.Open error: %w", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		return 0, fmt.Errorf("db.Ping error: %w", err)
	}

	return time.Since(sT), nil
}
