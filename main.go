package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type ParamUrl struct {
	pathPattern string
	pathUrl     *url.URL
	query       url.Values
}

var pathParamRegex = regexp.MustCompile(`^\d+$`)

func main() {
	var urls []string
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		urls = append(urls, sc.Text())
	}

	if err := sc.Err(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read input: %v\n", err)
	}

	results := map[string]ParamUrl{}

	for _, val := range urls {
		paramUrl, err := getParamUrl(val)
		if err != nil {
			continue
		}

		result, ok := results[paramUrl.pathPattern]
		if !ok {
			results[paramUrl.pathPattern] = paramUrl
			continue
		}

		for qk, qv := range paramUrl.query {
			_, ok := result.query[qk]
			if !ok {
				result.query[qk] = qv
			}
		}
	}

	for _, val := range results {
		target := val.pathUrl
		target.RawQuery = val.query.Encode()
		fmt.Println(target)
	}
}

func getParamUrl(val string) (ParamUrl, error) {
	target, err := url.Parse(val)
	if err != nil {
		return ParamUrl{}, err
	}

	segments := strings.Split(target.Path, "/")

	patternSegments := make([]string, len(segments))
	for index, segment := range segments {
		if pathParamRegex.MatchString(segment) {
			patternSegments[index] = "ID"
		} else {
			patternSegments[index] = segment
		}
	}

	pathPattern := strings.Join(patternSegments, "/")
	query := target.Query()
	target.RawQuery = ""

	return ParamUrl{
		pathPattern: pathPattern,
		pathUrl:     target,
		query:       query,
	}, nil
}
