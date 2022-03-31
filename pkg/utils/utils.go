package utils

import (
	"regexp"
	"strings"
)

func SnakeCasePath(path string) string {
	re := regexp.MustCompile(`[[:punct:]](\w)`)
	p := re.ReplaceAllStringFunc(path, strings.Title)
	re = regexp.MustCompile(`[[:punct:]]`)
	return re.ReplaceAllString(p, "")
}
