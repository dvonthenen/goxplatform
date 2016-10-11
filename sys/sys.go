package sys

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	uuid "github.com/twinj/uuid"

	run "github.com/dvonthenen/goxplatform/run"
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
type Sys struct{}

//GetUUID generates a UUID
func (Sys) GetUUID() []byte {
	myUUID := uuid.NewV1()
	log.Debugln("UUID Generated:", myUUID.String())
	return myUUID.Bytes()
}

//GetRunningKernelVersion returns the running kernel version
func (Sys) GetRunningKernelVersion() (string, error) {
	log.Debugln("GetRunningKernelVersion ENTER")

	cmdline := "uname -r"
	output, err := run.CommandOutput(cmdline)
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
