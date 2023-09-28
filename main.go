package main

import (
	"fmt"
	"main/internal/dbtest"
	"main/internal/tools"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/olekukonko/tablewriter"
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

	fmt.Printf("Starting %d Tests... \n\n", len(dbs))

	results, duration, err := dbtest.RunTests(dbs)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	table := tablewriter.NewWriter(os.Stdout)

	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{"Host", "Status", "Error", "duration"})
	table.SetFooter([]string{"", "", "Total", duration.String()})

	for _, result := range results {
		if result.Err != nil {
			table.Append([]string{result.Host, "Error", result.Err.Error(), result.Duration.String()})
			continue
		}

		table.Append([]string{result.Host, "Success", "N/A", result.Duration.String()})
	}

	table.Render()
}
