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
				tags = append(tags, strings.ToUpper(match))
			}
		}
	} else {
		// no matches
		return nil, fileExtension, errors.New("no tags from filename")
	}
	return tags, fileExtension, nil
}

func isInsignificant(word string) bool {
	if slices.Contains(combineInsignificantWords(), strings.ToUpper(word)) || isNumber(word) {
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

func combineInsignificantWords() (combinedList []string) {
	articles := []string{"A", "AN", "THE"}
	conjunctions := []string{"AND", "OR", "BUT"}
	prepositions := []string{"AT", "IN", "ON", "TO", "FROM", "WITH"}
	pronouns := []string{"I", "YOU", "HE", "SHE", "IT", "WE", "THEY"}
	possessivePronouns := []string{"MY", "YOUR", "HIS", "HER", "ITS", "OUR", "THEIR"}
	personalPronouns := []string{"ME", "HIM", "HER", "US", "THEM"}
	demonstrativePronouns := []string{"THIS", "THAT", "THESE", "THOSE"}
	quantifiers := []string{"SOME", "ANY", "MANY", "FEW", "SEVERAL", "ALL"}
	auxiliaryVerbs := []string{"IS", "ARE", "WAS", "WERE", "BE", "BEEN", "AM", "HAVE", "HAS", "HAD", "DO", "DOES", "DID"}
	commonVerbs := []string{"DO", "MAKE", "CREATE", "BUILD", "DEVELOP", "WRITE", "PRODUCE", "GENERATE"}

	combinedList = append(combinedList, articles...)
	combinedList = append(combinedList, conjunctions...)
	combinedList = append(combinedList, prepositions...)
	combinedList = append(combinedList, pronouns...)
	combinedList = append(combinedList, possessivePronouns...)
	combinedList = append(combinedList, personalPronouns...)
	combinedList = append(combinedList, demonstrativePronouns...)
	combinedList = append(combinedList, quantifiers...)
	combinedList = append(combinedList, auxiliaryVerbs...)
	combinedList = append(combinedList, commonVerbs...)

	return combinedList
}
