package str

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var (
	//ErrRegexpFailed Failed to parse tokens from regexp pattern
	ErrRegexpFailed = errors.New("Failed to parse tokens from regexp pattern")
)

//Str is a static class that captures string related functions
type Str struct{}

//NewStr generates a String object
func NewStr() *Str {
	myStr := &Str{}
	return myStr
}

//RegexMatch performs a regex tokenize match
func (str *Str) RegexMatch(haystack string, regex string) ([]string, error) {
	log.Debugln("RegexMatch ENTER")
	log.Debugln("haystack:", haystack)
	log.Debugln("regexp:", regex)

	r, err := regexp.Compile(regex)
	if err != nil {
		log.Errorln("Rexexp Failed:", err)
		log.Debugln("RegexMatch LEAVE")
		return nil, err
	}

	strings := r.FindStringSubmatch(haystack)
	if strings == nil {
		log.Errorln("Rexexp Failed:", err)
		log.Debugln("RegexMatch LEAVE")
		return nil, ErrRegexpFailed
	}

	log.Debugln("RegexMatch:", strings)
	log.Debugln("RegexMatch LEAVE")
	return strings, nil
}

//Trim all chars in col from front and end of the provided line
func (str *Str) Trim(line string, col string) string {
	log.Debugln("Trim ENTER")
	log.Debugln("line:", line)
	log.Debugln("col:", col)

	tmp := strings.TrimLeft(line, col)
	tmp = strings.TrimRight(tmp, col)

	log.Debugln("Trim:", tmp)
	log.Debugln("Trim LEAVE")

	return tmp
}

//IsNumeric returns true if numbers only
func (str *Str) IsNumeric(line string) bool {
	log.Debugln("IsNumeric ENTER")
	log.Debugln("line:", line)

	if _, err := strconv.Atoi(line); err == nil {
		log.Debugln("IsNumeric = TRUE")
		log.Debugln("IsNumeric LEAVE")
		return true
	}

	log.Debugln("IsNumeric = FALSE")
	log.Debugln("IsNumeric LEAVE")
	return false
}

//IsAlpha returns true if is all letter
func (str *Str) IsAlpha(line string) bool {
	log.Debugln("IsAlpha ENTER")
	log.Debugln("line:", line)

	for i := range line {
		if line[i] < 'A' || line[i] > 'z' {
			log.Debugln("IsAlpha = FALSE")
			log.Debugln("IsAlpha LEAVE")
			return false
		} else if line[i] > 'Z' && line[i] < 'a' {
			log.Debugln("IsAlpha = FALSE")
			log.Debugln("IsAlpha LEAVE")
			return false
		}
	}

	log.Debugln("IsAlpha = TRUE")
	log.Debugln("IsAlpha LEAVE")
	return true
}
