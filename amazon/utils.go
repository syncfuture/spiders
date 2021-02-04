package amazon

import (
	"regexp"
	"strings"
)

var (
	_spaceLineRegex = regexp.MustCompile(_spaceLineRegexStr)
)

func TrimSpaceAndLines(str string) string {
	str = CompactStr(str)
	str = strings.TrimSpace(str)
	return str
}

func CompactStr(str string) string {
	str = _spaceLineRegex.ReplaceAllString(str, "")
	return str
}
