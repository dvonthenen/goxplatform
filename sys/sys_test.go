package sys

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	assert "github.com/stretchr/testify/assert"
)

var sys *Sys

func TestMain(m *testing.M) {
	log.SetLevel(log.InfoLevel)
	log.Debugln("Start tests")
	sys = NewSys()
	m.Run()
}

func TestGetUUID(t *testing.T) {
	uuid := sys.GetUUID()
	assert.NotEqual(t, nil, uuid)
}

func TestGetUUIDStr(t *testing.T) {
	uuid := sys.GetUUIDStr()
	assert.NotEqual(t, "", uuid)
}

func TestGetOsType(t *testing.T) {
	itype := sys.GetOsType()
	assert.NotEqual(t, OsUnknown, itype)
}
