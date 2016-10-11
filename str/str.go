package str

import (
	"errors"
	"regexp"

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

	return strings, nil
}
