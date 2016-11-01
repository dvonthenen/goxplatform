package fs

import (
	"testing"

	log "github.com/Sirupsen/logrus"
	assert "github.com/stretchr/testify/assert"
)

var fs *Fs

func TestMain(m *testing.M) {
	log.SetLevel(log.InfoLevel)
	log.Debugln("Start tests")
	fs = NewFs()
	m.Run()
}

func TestGetFilename(t *testing.T) {
	URI := "http://127.0.0.1:8080/dir/myfile.deb"
	file := fs.GetFilenameFromURIOrFullPath(URI)
	assert.Equal(t, file, "myfile.deb")
}

func TestFilenameOnly(t *testing.T) {
	URI := "myfile.deb"
	file := fs.GetFilenameFromURIOrFullPath(URI)
	assert.Equal(t, file, "myfile.deb")
}

func TestPathFromFullPath(t *testing.T) {
	path := "/tmp/dir/myfile.deb"
	dir := fs.GetPathFileFullFilename(path)
	assert.Equal(t, dir, "/tmp/dir")
}
