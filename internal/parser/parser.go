package parser

import (
	"errors"
	"golang.org/x/exp/slices"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// TODO: create a separate package+repo for this

type IParser interface {
	ExtractTags(fileName string) (tags []string, fileType string, err error)
}

type Service struct {
}

func (s Service) ExtractTags(fileName string) (tags []string, fileType string, err error) {
	// get the file extension
	fileExtension := filepath.Ext(fileName)
	// remove the file extension from the file name
	fileName = strings.TrimSuffix(fileName, fileExtension)
	// Define a regular expression to extract tags
	regex := regexp.MustCompile(tagsRegex)
	// find all matches in the file name
	matches := regex.FindAllString(fileName, -1)

	if matchCount := len(matches); matchCount > 0 {
		for _, match := range matches {
			if !isInsignificant(match) {
				tags = append(tags, match)
			}
		}
	} else {
		// no matches
		return nil, fileExtension, errors.New("no tags from filename")
	}
	return tags, fileExtension, nil
}

func isInsignificant(word string) bool {
	insignificantWords := []string{"in", "the", "a", "an"}

	if slices.Contains(insignificantWords, word) || isNumber(word) {
		return true
	} else {
		return false
	}
}

func isNumber(str string) bool {
	if _, err := strconv.Atoi(str); err != nil {
		return false
	} else {
		return true
	}
}

const tagsRegex = `\b\w+\b`
