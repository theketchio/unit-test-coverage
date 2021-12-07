package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
  a short program that asserts that unit test coverage is not below established lower limits.
*/

type coverage map[string]float64

var (
	update       = flag.Bool("update", false, "set --update to increase limits file (e.g. limits.json) to any current higher unit test coverages")
	limitsFile   = flag.String("limits", "limits.json", "path to the file containing lower bounds for unit test coverage")
	coverageFile = flag.String("coverage", "coverage.txt", "path the to output of go test coverage (e.g. go test ./... -coverprofile cover.out > coverage.txt)")
	bypass       = flag.Bool("bypass", false, "set --bypass to avoid exit status 1 on insufficient coverage")
)

func main() {
	flag.Parse()

	limitsMap, err := getLimits(*limitsFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	coverageMap, err := getCoverage(*coverageFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if *update {
		err = updateCoverage(coverageMap, limitsMap, *limitsFile)
	} else {
		err = assertCoverage(coverageMap, limitsMap)
	}
	if err != nil {
		fmt.Println(err.Error())
		if !*bypass {
			os.Exit(1)
		}
	}
	os.Exit(0)
}

// getLimits returns a map of limits created from limits.json
func getLimits(filename string) (coverage, error) {
	limitsMap := make(coverage)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&limitsMap)
	return limitsMap, err
}

// getCoverage returns a map of current coverage from coverage.txt
func getCoverage(filename string) (coverage, error) {
	coverageMap := make(coverage)
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		vals := strings.Split(scanner.Text(), "\t")
		percentage := 0.0
		if len(vals) >= 4 {
			percentage, err = strconv.ParseFloat(strings.TrimLeft(strings.TrimRight(vals[3], "% of statements"), "coverage: "), 64)
			if err != nil {
				return nil, err
			}
			coverageMap[vals[1]] = percentage
		}
	}

	return coverageMap, err
}

// assertCoverage returns aggregated errors for coverages below limits
func assertCoverage(coverageMap, limitsMap coverage) error {
	var errs []string
	for path, percentage := range coverageMap {
		limit, ok := limitsMap[path]
		if !ok {
			continue
		}
		if percentage < limit {
			errs = append(errs, fmt.Sprintf("coverage for %s is %.2f but expected to be >= %.2f", path, percentage, limit))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("coverage errors:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}

// updateCoverage updates and writes the limits map with any new, higher coverages
func updateCoverage(coverageMap, limitsMap coverage, filename string) error {
	var wasUpdated bool
	for path, percentage := range coverageMap {
		limit, ok := limitsMap[path]
		if !ok || limit < percentage {
			limitsMap[path] = percentage
			wasUpdated = true
		}
	}
	if wasUpdated {
		f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()
		j, err := json.MarshalIndent(limitsMap, "", "\t")
		if err != nil {
			return err
		}
		_, err = f.Write(j)
		if err != nil {
			return err
		}
	}
	return nil
}
