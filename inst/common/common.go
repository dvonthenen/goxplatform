package common

import (
	"errors"
	"regexp"

	log "github.com/Sirupsen/logrus"
)

var (
	//ErrParseVersionFailed failed to parse version from filename
	ErrParseVersionFailed = errors.New("Failed to parse version from filename")
)

//ParseVersionFromFilename this parses the version string out of the filename
func ParseVersionFromFilename(filename string) (string, error) {
	log.Debugln("ParseVersionFromFilename ENTER")
	log.Debugln("filename:", filename)

	r, err := regexp.Compile(".*([0-9]+\\.[0-9]+[\\.\\-][0-9]+\\.[0-9]+).*")
	if err != nil {
		log.Debugln("regexp is invalid")
		log.Debugln("ParseVersionFromFilename LEAVE")
		return "", err
	}
	strings := r.FindStringSubmatch(filename)
	if strings == nil || len(strings) < 2 {
		log.Debugln("Unable to find version from string")
		log.Debugln("ParseVersionFromFilename LEAVE")
		return "", ErrParseVersionFailed
	}

	version := strings[1]

	log.Debugln("Found:", version)
	log.Debugln("ParseVersionFromFilename LEAVE")

	return version, nil
}
