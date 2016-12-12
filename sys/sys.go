package sys

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	uuid "github.com/twinj/uuid"

	common "github.com/dvonthenen/goxplatform/common"
	fs "github.com/dvonthenen/goxplatform/fs"
	run "github.com/dvonthenen/goxplatform/run"
	str "github.com/dvonthenen/goxplatform/str"
)

const (
	//OsUnknown unknown OS
	OsUnknown = 0

	//OsRhel is RHEL
	OsRhel = 1

	//OsSuse is SuSE
	OsSuse = 2

	//OsUbuntu is Ubuntu
	OsUbuntu = 3

	//OsCoreOs is CoreOS
	OsCoreOs = 4

	//OsMac is OSX
	OsMac = 5
)

var (
	//ErrSrcNotExist src file doesnt exist
	ErrSrcNotExist = errors.New("Source file does not exist")

	//ErrSrcNotRegularFile src file is not a regular file
	ErrSrcNotRegularFile = errors.New("Source file is not a regular file")

	//ErrDstNotRegularFile dst file is not a regular file
	ErrDstNotRegularFile = errors.New("Destination file is not a regular file")

	//ErrUnknownOsVersion unable to determine OS version
	ErrUnknownOsVersion = errors.New("Unknown OS version")
)

//Sys is a static class that provides System related functions
type Sys struct {
	run *run.Run
	fs  *fs.Fs
	str *str.Str
}

//NewSys generates a Sys object
func NewSys() *Sys {
	myRun := run.NewRun()
	myFs := fs.NewFs()
	myStr := str.NewStr()
	mySys := &Sys{
		run: myRun,
		fs:  myFs,
		str: myStr,
	}
	return mySys
}

//GetUUID generates a UUID
func (sys *Sys) GetUUID() []byte {
	myUUID := uuid.NewV1()
	log.Debugln("UUID Generated:", myUUID.String())
	return myUUID.Bytes()
}

//GetUUIDStr generates a UUID
func (sys *Sys) GetUUIDStr() string {
	myUUID := uuid.NewV1()
	log.Debugln("UUID Generated:", myUUID.String())
	return myUUID.String()
}

//GetOsType gets the OS type
func (sys *Sys) GetOsType() int {
	log.Debugln("GetOsType ENTER")

	osType := OsUnknown
	if sys.fs.DoesFileExist("/etc/redhat-release") {
		osType = OsRhel
	} else if sys.fs.DoesFileExist("/etc/SuSE-release") {
		osType = OsSuse
	} else if sys.fs.DoesFileExist("/etc/lsb-release") {
		osType = OsUbuntu
		//	} else if sys.fs.DoesFileExist("/etc/release") {
		//		return OsCoreOs
	} else {
		out, err := sys.run.CommandOutput("uname -s")
		if err == nil && strings.EqualFold(out, "Darwin") {
			osType = OsMac
		} else {
			log.Warnln("Unable to determine OS type")
		}
	}

	log.Debugln("GetOsType =", osType)
	log.Debugln("GetOsType LEAVE")
	return osType
}

func parseVersionFromString(str string) (int, int, error) {
	strs := strings.Split(str, ".")
	imajor, errMajor := strconv.Atoi(strs[0])
	if errMajor != nil {
		return 0, 0, errMajor
	}
	iminor, errMinor := strconv.Atoi(strs[1])
	if errMinor != nil {
		return 0, 0, errMinor
	}
	return imajor, iminor, nil
}

//GetOsVersion returns the major minor version of the OS
func (sys *Sys) GetOsVersion() (int, int, error) {
	log.Debugln("GetOsVersion ENTER")

	itype := sys.GetOsType()
	switch itype {
	case OsRhel:
		data, errRead := ioutil.ReadFile("/etc/redhat-release")
		if errRead != nil {
			return 0, 0, errRead
		}
		log.Debugln(string(data))
		needles, errRegex := sys.str.RegexMatch(string(data), " ([0-9]+\\.[0-9]+[\\.]*[0-9]*) ")
		if errRegex != nil {
			return 0, 0, errRegex
		}
		return parseVersionFromString(needles[0])

	case OsSuse:
		return 0, 0, common.ErrNotImplemented //TODO

	case OsUbuntu:
		data, errRead := ioutil.ReadFile("/etc/lsb-release")
		if errRead != nil {
			return 0, 0, errRead
		}
		log.Debugln(string(data))
		needles, errRegex := sys.str.RegexMatch(string(data), "DISTRIB_RELEASE=([0-9]+\\.[0-9]+[\\.]*[0-9]*) ")
		if errRegex != nil {
			return 0, 0, errRegex
		}
		return parseVersionFromString(needles[0])

	case OsCoreOs:
		return 0, 0, common.ErrNotImplemented //TODO

	case OsMac:
		return 0, 0, common.ErrNotImplemented //TODO
	}

	log.Debugln("GetOsVersion LEAVE")

	return 0, 0, ErrUnknownOsVersion
}

//GetOsStrByType gets the OS string
func (sys *Sys) GetOsStrByType(iType int) string {
	log.Debugln("GetOsStrByType ENTER")

	osStr := "Unknown"
	switch iType {
	case OsRhel:
		osStr = "RHEL"
	case OsSuse:
		osStr = "SUSE"
	case OsUbuntu:
		osStr = "Ubuntu"
	case OsCoreOs:
		osStr = "CoreOS"
	case OsMac:
		osStr = "OSX"
	}

	log.Debugln("GetOsStrByType =", osStr)
	log.Debugln("GetOsStrByType LEAVE")
	return osStr
}

//GetRunningKernelVersion returns the running kernel version
func (sys *Sys) GetRunningKernelVersion() (string, error) {
	log.Debugln("GetRunningKernelVersion ENTER")

	cmdline := "uname -r"
	output, err := sys.run.CommandOutput(cmdline)
	if err != nil {
		log.Debugln("runCommandOutput Failed:", err)
		log.Debugln("GetRunningKernelVersion LEAVE")
		return "", err
	}

	version := output

	log.Debugln("GetRunningKernelVersion Kernel:", version)
	log.Debugln("GetRunningKernelVersion LEAVE")

	return version, nil
}

//GetDeviceList returns the list of all devices on the system
func (sys *Sys) GetDeviceList() ([]string, error) {
	log.Debugln("GetDeviceList ENTER")

	list := []string{}

	outputCmd := "fdisk -l | grep \\/dev\\/"
	output, err := sys.run.CommandOutput(outputCmd)
	if err != nil {
		log.Errorln("Failed to get device list. Err:", err)
		log.Debugln("GetDeviceList LEAVE")
		return list, err
	}

	buffer := bytes.NewBufferString(output)
	for {
		str, err := buffer.ReadString('\n')
		if err == io.EOF {
			break
		}

		needles, errRegex := sys.str.RegexMatch(str, "Disk (/dev/.*): ")
		if errRegex != nil {
			log.Errorln("RegexMatch Failed. Err:", err)
			log.Debugln("GetDeviceList LEAVE")
			return list, err
		}
		device := needles[0]
		log.Debugln("Device Found:", device)

		list = append(list, device)
	}

	log.Debugln("GetDeviceList Succeeded. Device Count:", len(list))
	log.Debugln("GetDeviceList LEAVE")

	return list, nil
}

//GetInUseDeviceList returns the list of all devices on the system
func (sys *Sys) GetInUseDeviceList() ([]string, error) {
	log.Debugln("GetInUseDeviceList ENTER")

	list := []string{}

	outputCmd := "blkid"
	output, err := sys.run.CommandOutput(outputCmd)
	if err != nil {
		log.Errorln("Failed to get blkid list. Err:", err)
		log.Debugln("GetInUseDeviceList LEAVE")
		return list, err
	}

	buffer := bytes.NewBufferString(output)
	for {
		str, err := buffer.ReadString('\n')
		if err == io.EOF {
			break
		}

		needles, errRegex := sys.str.RegexMatch(str, "(/dev/.*):")
		if errRegex != nil {
			log.Errorln("RegexMatch Failed. Err:", err)
			log.Debugln("GetInUseDeviceList LEAVE")
			return list, err
		}
		device := needles[0]
		log.Debugln("Device Found:", device)

		list = append(list, device)
	}

	log.Debugln("GetInUseDeviceList Succeeded. Device Count:", len(list))
	log.Debugln("GetInUseDeviceList LEAVE")

	return list, nil
}
