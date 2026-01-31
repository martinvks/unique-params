package main

import (
	"bufio"
	"flag"
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

const pathParamPlaceholder = "{ID}"

var pathParamRegex = regexp.MustCompile(`^\d+$`)
var uuidParamRegex = regexp.MustCompile(`^(?i)[\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12}$`)

func main() {
	max := flag.Int(
		"max",
		0,
		"maximum number of URLs to output (0 = unlimited)",
	)
	maxQuery := flag.Int(
		"max-query",
		0,
		"maximum number of query parameters per URL (0 = unlimited)",
	)
	flag.Parse()

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

		if len(paramUrl.query) == 0 && !strings.Contains(paramUrl.pathPattern, pathParamPlaceholder) {
			continue
		}

		result, ok := results[paramUrl.pathPattern]
		if !ok {
			results[paramUrl.pathPattern] = paramUrl
			continue
		}

		for qk, qv := range paramUrl.query {
			if *maxQuery > 0 && len(result.query) >= *maxQuery {
				break
			}
			_, ok := result.query[qk]
			if !ok {
				result.query[qk] = qv
			}
		}
	}

	i := 0
	for _, val := range results {
		if *max > 0 && i >= *max {
			break
		}
		target := val.pathUrl
		target.RawQuery = val.query.Encode()
		fmt.Println(target)
		i++
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
		if pathParamRegex.MatchString(segment) || uuidParamRegex.MatchString(segment) {
			patternSegments[index] = pathParamPlaceholder
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
