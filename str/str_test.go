package str

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	assert "github.com/stretchr/testify/assert"
)

var str *Str

func TestMain(m *testing.M) {
	log.SetLevel(log.InfoLevel)
	log.Debugln("Start tests")
	str = NewStr()
	m.Run()
}

func TestRegexMatch(t *testing.T) {
	strings, err := str.RegexMatch("hello", "el+")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, strings)
	assert.Equal(t, 1, len(strings))
	assert.Equal(t, "ell", strings[0])
}

func TestTrim(t *testing.T) {
	string := str.Trim(" !test!!!", " !")
	assert.Equal(t, "test", string)
}
