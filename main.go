package main

import (
	"fmt"
	"main/internal/dbtest"
	"main/internal/tools"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/olekukonko/tablewriter"
)

func main() {
	sleepTime := 3 * time.Second

	if b, _ := strconv.ParseBool(os.Getenv("SLEEP")); b {
		fmt.Printf("sleeping for %v\n", sleepTime)

		time.Sleep(sleepTime)
	}

	dbs, dbmap, err := tools.GetURLsFromEnvironment("TEST_")
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

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Duration < results[j].Duration
	})

	table := tablewriter.NewWriter(os.Stdout)

	table.SetAutoFormatHeaders(false)
	table.SetHeader([]string{"ENV", "Host", "Scheme", "Status", "Error", "Duration"})
	table.SetFooter([]string{"", "", "", "", "Total", duration.String()})

	for _, result := range results {
		if result.Err != nil {
			table.Append([]string{
				dbmap[result.Meta.URL],
				result.Meta.Host, result.Meta.Scheme,
				"Error",
				result.Err.Error(),
				result.Duration.String(),
			})
			continue
		}

		table.Append([]string{
			dbmap[result.Meta.URL],
			result.Meta.Host,
			result.Meta.Scheme,
			"Success",
			"N/A",
			result.Duration.String(),
		})
	}

	table.Render()
}
