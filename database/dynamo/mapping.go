package dynamo

import (
	"regexp"
	"strings"
)

func mapExpressionNames(expression string) map[string]*string {
	re := regexp.MustCompile("#\\S+")

	matches := re.FindAllString(expression, -1)

	matchingFields := map[string]*string{}
	for _, val := range matches {
		originalFieldName := strings.Replace(val, "#", "", 1)
		matchingFields[val] = &originalFieldName
	}

	return matchingFields
}
