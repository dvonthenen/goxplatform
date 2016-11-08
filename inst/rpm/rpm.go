package rpm

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

//Rpm implementation for the Rpm package manager
type Rpm struct {
	run *run.Run
}

//NewRpm generates a Rpm object
func NewRpm() *Rpm {
	myRun := run.NewRun()
	myRpm := &Rpm{
		run: myRun,
	}
	return myRpm
}

//IsInstalled returns if the package is installed
func (rpm *Rpm) IsInstalled(packageName string) error {
	log.Debugln("IsInstalled ENTER")
	log.Debugln("packageName:", packageName)

	_, err := rpm.GetInstalledVersion(packageName, false)
	if err != nil {
		log.Debugln("Package", packageName, "IS NOT installed")
		log.Debugln("IsInstalled LEAVE")
		return err
	}

	log.Debugln("Package", packageName, "IS installed")
	log.Debugln("IsInstalled LEAVE")
	return nil
}

func correctVersionFromRpm(version string) string {
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
func (rpm *Rpm) GetInstalledVersion(packageName string, parseVersion bool) (string, error) {
	log.Debugln("GetInstalledVersion ENTER")
	log.Debugln("packageName:", packageName)

	cmdline1 := "rpm -qi " + packageName + " | grep Version | sed -n -e 's/^.*Version.*: //p'"
	output1, errCmd1 := rpm.run.CommandOutput(cmdline1)
	if errCmd1 != nil {
		log.Debugln("runCommandOutput Failed:", errCmd1)
		log.Debugln("GetInstalledVersion LEAVE")
		return "", errCmd1
	}
	if len(output1) == 0 {
		log.Debugln("Output1 length is empty")
		log.Debugln("GetInstalledVersion LEAVE")
		return "", ErrExecEmptyOutput
	}

	cmdline2 := "rpm -qi " + packageName + " | grep Release | sed -n -e 's/^.*Release.*: //p'"
	output2, errCmd2 := rpm.run.CommandOutput(cmdline2)
	if errCmd2 != nil {
		log.Debugln("runCommandOutput Failed:", errCmd2)
		log.Debugln("GetInstalledVersion LEAVE")
		return "", errCmd2
	}
	if len(output2) == 0 {
		log.Debugln("Output2 length is empty")
		log.Debugln("GetInstalledVersion LEAVE")
		return "", ErrExecEmptyOutput
	}

	output := output1 + "-" + output2

	if strings.Contains(output1, "is not installed") ||
		strings.Contains(output2, "is not installed") {
		log.Warnln("Package", packageName, "is not installed. Blanking the output.")
		output = ""
	}

	//this is for REX-Ray and DVDCLI that only use the format 0.2.0
	//0.2.0-1 -> 0.2.0
	version := correctVersionFromRpm(output)

	//use the original string but remove anything but the version
	//2.0-10000.2072.el7 -> 2.0-10000.2072
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
