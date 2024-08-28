package recipeschema

import (
	"regexp"
	"strings"
)

func normalizeLine(line string) string {
	line = strings.Replace(line, "\n", " ", -1)
	line = strings.TrimSpace(line)
	line = regexp.MustCompile(`\s+`).ReplaceAllString(line, " ")
	return line
}
