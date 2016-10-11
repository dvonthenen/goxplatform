package sys

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	uuid "github.com/twinj/uuid"

	fs "github.com/dvonthenen/goxplatform/fs"
	run "github.com/dvonthenen/goxplatform/run"
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
)

var (
	//ErrSrcNotExist src file doesnt exist
	ErrSrcNotExist = errors.New("Source file does not exist")

	//ErrSrcNotRegularFile src file is not a regular file
	ErrSrcNotRegularFile = errors.New("Source file is not a regular file")

	//ErrDstNotRegularFile dst file is not a regular file
	ErrDstNotRegularFile = errors.New("Destination file is not a regular file")
)

//Sys is a static class that provides System related functions
type Sys struct {
	run *run.Run
	fs  *fs.Fs
}

//NewSys generates a Sys object
func NewSys() *Sys {
	myRun := run.NewRun()
	myFs := fs.NewFs()
	mySys := &Sys{
		run: myRun,
		fs:  myFs,
	}
	return mySys
}

//GetUUID generates a UUID
func (sys *Sys) GetUUID() []byte {
	myUUID := uuid.NewV1()
	log.Debugln("UUID Generated:", myUUID.String())
	return myUUID.Bytes()
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
	}

	log.Debugln("GetOsType =", osType)
	log.Debugln("GetOsType LEAVE")
	return osType
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
