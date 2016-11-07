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
	strings, err := str.RegexMatch("http://10.0.0.1:9000", ".+//(.*):[0-9]+")
	assert.Equal(t, nil, err)
	assert.NotEqual(t, nil, strings)
	assert.Equal(t, 2, len(strings))
	assert.Equal(t, "10.0.0.1", strings[1])
}

func TestTrim(t *testing.T) {
	string := str.Trim(" !test!!!", " !")
	assert.Equal(t, "test", string)
}
