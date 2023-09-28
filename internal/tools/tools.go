package tools

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

func GetURLsFromEnvironment(prefix string) ([]string, error) {
	urls := []string{}

	for _, env := range os.Environ() {
		if !strings.HasPrefix(env, prefix) {
			continue
		}

		kv := strings.SplitN(env, "=", 2)

		if len(kv) != 2 {
			return nil, fmt.Errorf("malformed kv found: %s", env)
		}

		key := kv[0]
		value := kv[1]

		if _, err := url.Parse(value); err != nil {
			return nil, fmt.Errorf("key %s is not a valid URL: %s", key, value)
		}

		urls = append(urls, value)
	}

	if len(urls) == 0 {
		return nil, fmt.Errorf("no urls in environment where found with given prefix")
	}

	return urls, nil
}
