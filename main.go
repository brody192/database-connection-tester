package main

import (
	"fmt"
	"main/internal/dbtest"
	"main/internal/tools"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func main() {
	if os.Getenv("SLEEP") == "true" {
		fmt.Println("sleeping for 5 seconds")

		time.Sleep(5 * time.Second)
	}

	dbs, err := tools.GetURLsFromEnvironment("TEST_")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Println("running tests")

	results, duration, err := dbtest.RunTests(dbs)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	for _, result := range results {
		if result.Err != nil {
			fmt.Printf("error testing %s: %s\n", result.Host, result.Err)
			continue
		}

		fmt.Printf("tested %s successfully in %v\n", result.Host, result.Duration)
	}

	fmt.Printf("testing finished in %v", duration)
}
