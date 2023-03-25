package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"sort"
)

func main() {
	var urls []string
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		urls = append(urls, sc.Text())
	}
	sort.Strings(urls)

	uniqueUrl := &url.URL{}
	uniqueQuery := url.Values{}

	for _, val := range urls {
		currentUrl, err := url.Parse(val)
		if err != nil {
			continue
		}

		currentQuery := currentUrl.Query()
		currentUrl.RawQuery = ""

		if uniqueUrl.String() != currentUrl.String() {
			if uniqueUrl.String() != "" {
				uniqueUrl.RawQuery = uniqueQuery.Encode()
				fmt.Println(uniqueUrl)
			}

			uniqueUrl = currentUrl
			uniqueQuery = currentQuery
			continue
		}

		for queryKey, queryValues := range currentQuery {
			uniqueQuery[queryKey] = queryValues
		}
	}

	uniqueUrl.RawQuery = uniqueQuery.Encode()
	fmt.Println(uniqueUrl)
}
