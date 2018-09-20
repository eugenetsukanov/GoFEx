package main

import (
	"errors"
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

//type SafeCounter struct
var fetched = struct {
	v   map[string]error
	mux sync.Mutex
}{v: make(map[string]error)}

var loading = errors.New("url load in progress") // sentinel value

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {	
	// This implementation doesn't do either:
	if depth <= 0 {
		fmt.Printf("<- Done with %v, depth 0.\n", url)
		return
	}

	fetched.mux.Lock()

	if _, ok := fetched.v[url]; ok {
		fetched.mux.Unlock()
		fmt.Printf("<- Done with %v, already fetched.\n", url)
		return
	}

	fetched.v[url] = loading
	fetched.mux.Unlock()

	body, urls, err := fetcher.Fetch(url)

	fetched.mux.Lock()
	fetched.v[url] = err // if error
	fetched.mux.Unlock()

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)

	done := make(chan bool)
	for i, u := range urls {
		fmt.Printf("-> Crawling child %v/%v of %v : %v.\n", i, len(urls), url, u)
		go func(url string) {
			Crawl(url, depth-1, fetcher)
			done <- true
		}(u)
	}

	for i, u := range urls {
		fmt.Printf("<- [%v] %v/%v Waiting for child %v.\n", url, i, len(urls), u)
		<-done
	}

	return
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)

	fmt.Println("\n\nFetching stats\n--------------")
	for url, error := range fetched.v {
		if error != nil {
			fmt.Printf("%v failed: %v\n", url, error)
		} else {
			fmt.Printf("%v was fetched\n", url)
		}
	}
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
