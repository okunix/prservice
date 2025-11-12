package models

import (
	"fmt"
	"strings"
)

type ValidationError map[string]error

func (v ValidationError) Error() string {
	errorLines := []string{}
	for k, v := range v {
		errorLines = append(errorLines, fmt.Sprintf("%s: %s", k, v.Error()))
	}
	return strings.Join(errorLines, "\n")
}
