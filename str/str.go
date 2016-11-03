package str

import (
	"errors"
	"regexp"
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
