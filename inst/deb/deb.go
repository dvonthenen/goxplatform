package deb

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"

	common "github.com/dvonthenen/goxplatform/inst/common"
	run "github.com/dvonthenen/goxplatform/run"
)

var (
	//ErrExecEmptyOutput failed to generate any output
	ErrExecEmptyOutput = errors.New("Failed to generate any output")
)

//Deb implementation for the Deb package manager
type Deb struct {
	run *run.Run
}

//NewDeb generates a Deb object
func NewDeb() *Deb {
	myRun := run.NewRun()
	myDeb := &Deb{
		run: myRun,
	}
	return myDeb
}

//IsInstalled returns if the package is installed
func (deb *Deb) IsInstalled(packageName string) error {
	log.Debugln("IsInstalled ENTER")
	log.Debugln("packageName:", packageName)

	_, err := deb.GetInstalledVersion(packageName, false)
	if err != nil {
		log.Debugln("Package", packageName, "IS NOT installed")
		log.Debugln("IsInstalled LEAVE")
		return err
	}

	log.Debugln("Package", packageName, "IS installed")
	log.Debugln("IsInstalled LEAVE")
	return nil
}

func correctVersionFromDeb(version string) string {
	if len(version) == 0 {
		return ""
	}

	index := strings.Index(version, "-")
	if index == -1 {
		return version
	}

	fixedVersion := version[:index]
	return fixedVersion
}

//GetInstalledVersion returns the version of the installed package
func (deb *Deb) GetInstalledVersion(packageName string, parseVersion bool) (string, error) {
	log.Debugln("GetInstalledVersion ENTER")
	log.Debugln("packageName:", packageName)

	cmdline := "dpkg -s " + packageName + " | grep Version | sed -n -e 's/^.*Version: //p'"
	output, errCmd := deb.run.CommandOutput(cmdline)
	if errCmd != nil {
		log.Debugln("runCommandOutput Failed:", errCmd)
		log.Debugln("GetInstalledVersion LEAVE")
		return "", errCmd
	}

	if len(output) == 0 {
		log.Debugln("Output length is empty")
		log.Debugln("GetInstalledVersion LEAVE")
		return "", ErrExecEmptyOutput
	}

	if strings.Contains(output, "is not installed") {
		log.Warnln("Package", packageName, "is not installed. Blanking the output.")
		output = ""
	}

	//this is for REX-Ray and DVDCLI that only use the format 0.2.0
	//0.2.0-1 -> 0.2.0
	version := correctVersionFromDeb(output)

	//use the original string but remove anything but the version
	//2.0.10000.2072.Ubuntu.14.04 -> 2.0.10000.2072
	if parseVersion {
		myVersion, errParse := common.ParseVersionFromFilename(output)
		if errParse != nil {
			log.Debugln("ParseVersionFromFilename Failed:", errParse)
			log.Debugln("GetInstalledVersion LEAVE")
			return "", errParse
		}
		version = myVersion
	}

	log.Debugln("GetInstalledVersion Found:", version)
	log.Debugln("GetInstalledVersion LEAVE")

	return version, nil
}
