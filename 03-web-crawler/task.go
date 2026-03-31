package main

import "sync"

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) ([]string, error) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	if depth <= 0 {
		return nil, nil
	}
	var mu sync.RWMutex

	urlVisited := map[string]struct{}{}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		return nil, err
	}
	result := []string{body}

	for _, u := range urls {
		go func() {
			if res, err := Crawl(u, depth-1, fetcher); err == nil {
				result = append(result, res...)
			}
		}()
	}
	return result, nil
}
