package dbtest

import (
	"fmt"
	"main/internal/dbtest/databases"
	"net/url"
	"sync"
	"time"
)

type testResults struct {
	Host     string
	Duration time.Duration
	Err      error
}

func RunTests(databaseURLs []string) ([]testResults, time.Duration, error) {
	dburls := []*url.URL{}

	for _, db := range databaseURLs {
		dburl, err := url.Parse(db)
		if err != nil {
			return nil, 0, err
		}

		dburls = append(dburls, dburl)
	}

	sT := time.Now()

	results := []testResults{}

	resultsChan := make(chan testResults, len(dburls))

	defer close(resultsChan)

	ack := make(chan struct{})

	defer close(ack)

	go func() {
		for result := range resultsChan {
			results = append(results, result)
			ack <- struct{}{}
		}
	}()

	wg := sync.WaitGroup{}

	wg.Add(len(dburls))

	for _, dburl := range dburls {
		go func(dburl *url.URL) {
			defer wg.Done()

			result := testResults{}

			switch dburl.Scheme {
			case "postgres", "mysql":
				result.Duration, result.Err = databases.Sql(dburl.Scheme, dburl.String())
			case "redis":
				result.Duration, result.Err = databases.Redis(dburl.String())
			case "mongodb":
				result.Duration, result.Err = databases.Mongo(dburl.String())
			default:
				result.Err = fmt.Errorf(`scheme "%s" not implemented`, dburl.Scheme)
			}

			result.Host = dburl.Hostname()

			resultsChan <- result
			<-ack
		}(dburl)
	}

	wg.Wait()

	return results, time.Since(sT), nil
}
