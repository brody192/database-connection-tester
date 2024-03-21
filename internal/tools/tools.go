package tools

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

func GetURLsFromEnvironment(prefix string) ([]*url.URL, map[*url.URL]string, error) {
	urls := []*url.URL{}
	urlToEnv := make(map[*url.URL]string)

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			continue
		}

		kv := strings.SplitN(env, "=", 2)

		if len(kv) != 2 {
			return nil, nil, fmt.Errorf("malformed kv found: %s", env)
		}

		key := kv[0]
		value := kv[1]

		if strings.TrimSpace(value) == "" {
			return nil, nil, fmt.Errorf("empty database url found for key %s", key)
		}

		urlParsed, err := url.Parse(value)
		if err != nil {
			return nil, nil, fmt.Errorf("key %s is not a valid URL: %s", key, value)
		}

		if urlParsed.Scheme == "" {
			return nil, nil, fmt.Errorf("url %s missing scheme", value)
		}

		urls = append(urls, urlParsed)
		urlToEnv[urlParsed] = strings.SplitN(key, prefix, 2)[1]
	}

	if len(urls) == 0 {
		return nil, nil, fmt.Errorf("no urls in environment where found with given prefix: %s", prefix)
	}

	return urls, urlToEnv, nil
}
