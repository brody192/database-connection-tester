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
	results := []*testResults{}

	resultsChan := make(chan *testResults, len(databaseURLs))
	defer close(resultsChan)

	wg := sync.WaitGroup{}

	go func() {
		for result := range resultsChan {
			results = append(results, result)
			wg.Done()
		}
	}()

	wg.Add(len(databaseURLs))

	sT := time.Now()

	for _, dburl := range databaseURLs {
		go func(dburl *url.URL) {
			defer wg.Done()

			result := testResults{}

			if dburl.Scheme == "postgresql" {
				dburl.Scheme = "postgres"
			}

			switch dburl.Scheme {
			case "postgres", "mysql", "libsql":
				result.Duration, result.Err = databases.Sql(dburl.Scheme, dburl.String())
			case "mariadb":
				result.Duration, result.Err = databases.Sql("mysql", dburl.String())
			case "redis":
				result.Duration, result.Err = databases.Redis(dburl.String())
			case "mongodb":
				result.Duration, result.Err = databases.Mongo(dburl.String())
			case "edgedb":
				result.Duration, result.Err = databases.EdgeDB(dburl.String())
			case "memcache", "memcached":
				result.Duration, result.Err = databases.Memcache(dburl.Host)
			default:
				result.Err = fmt.Errorf(`scheme "%s" not implemented`, dburl.Scheme)
			}

			result.Meta = testMeta{
				URL:    dburl,
				Host:   dburl.Hostname(),
				Scheme: dburl.Scheme,
			}

			resultsChan <- &result
			wg.Add(1)
		}(dburl)
	}

	wg.Wait()

	return results, time.Since(sT), nil
}
