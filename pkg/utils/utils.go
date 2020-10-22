package utils

import (
	"encoding/csv"
	"path/filepath"
	"regexp"
	"strings"
)

func SedReplaceGlobal(subject string, search string, replace string) []string {
	var linesTransformed []string

	lines := strings.Split(subject, "\n")

	for _, line := range lines {
		lineTransformed := regexp.MustCompile(search).ReplaceAllString(line, replace)

		linesTransformed = append(linesTransformed, lineTransformed)
	}

	return linesTransformed
}

func RemoveEmptyLines(lines []string) []string {
	var linesFiltered []string

	for i := range lines {
		if lines[i] != "" {
			linesFiltered = append(linesFiltered, lines[i])
		}
	}

	return linesFiltered
}

func SplitPreservingQuotes(s string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(s))

	r.Comma = ' ' // space

	fields, err := r.Read()

	if err != nil {
		return nil, err
	}

	return fields, nil
}

func FilterOut(files []string, patterns []string) []string {
	if len(patterns) == 0 {
		return files
	}

	var filesFiltered []string

	pattern := patterns[0]

	for i := range files {
		file := files[i]

		match, _ := filepath.Match(pattern, file)

		if !match {
			filesFiltered = append(filesFiltered, file)
		}
	}

	return FilterOut(filesFiltered, patterns[1:])
}
