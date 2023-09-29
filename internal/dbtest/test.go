package dbtest

import (
	"fmt"
	"main/internal/dbtest/databases"
	"net/url"
	"sync"
	"time"
)

type testMeta struct {
	URL    *url.URL
	Host   string
	Scheme string
}

type testResults struct {
	Meta     testMeta
	Duration time.Duration
	Err      error
}

func RunTests(databaseURLs []*url.URL) ([]*testResults, time.Duration, error) {
	sT := time.Now()

	results := []*testResults{}

	resultsChan := make(chan *testResults, len(databaseURLs))

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

	wg.Add(len(databaseURLs))

	for _, dburl := range databaseURLs {
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

			result.Meta = testMeta{
				URL:    dburl,
				Host:   dburl.Hostname(),
				Scheme: dburl.Scheme,
			}

			resultsChan <- &result
			<-ack
		}(dburl)
	}

	wg.Wait()

	return results, time.Since(sT), nil
}
