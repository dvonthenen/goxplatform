package init

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	common "github.com/dvonthenen/goxplatform/init/common"
	initd "github.com/dvonthenen/goxplatform/init/initd"
	systemd "github.com/dvonthenen/goxplatform/init/systemd"
	run "github.com/dvonthenen/goxplatform/run"
	sys "github.com/dvonthenen/goxplatform/sys"
)

const (
	//InitUnknown type is Unknown
	InitUnknown = 0

	//InitSystemD is SystemD
	InitSystemD = 1

	//InitUpdateRcD is UpdateRC
	InitUpdateRcD = 2

	//InitChkConfig is InitD
	InitChkConfig = 3
)

var (
	//ErrInvalidInitSystem the Init System is not valid
	ErrInvalidInitSystem = errors.New("Invalid Init System")
)

//Init is a static class that captures install package rules
type Init struct {
	sys  *sys.Sys
	run  *run.Run
	init common.IInit
}

//NewInit generates a Init object
func NewInit() *Init {
	mySys := sys.NewSys()
	myRun := run.NewRun()

	myInitSystem := &Init{
		sys: mySys,
		run: myRun,
	}

	var myInit common.IInit
	switch myInitSystem.GetInitSystemType() {
	case InitSystemD:
		myInit = systemd.NewSystemD()
	case InitUpdateRcD:
		myInit = initd.NewInitD()
	case InitChkConfig:
		myInit = initd.NewInitD()
	}
	myInitSystem.init = myInit

	return myInitSystem
}


//GetInitSystemType returns the Init type on the Operating System
func (init *Init) GetInitSystemType() int {
	log.Debugln("getInitSystemType ENTER")

	if init.run.ExecExistsInPath("systemctl") {
		log.Debugln("getInitSystemType = initSystemD")
		log.Debugln("getInitSystemType LEAVE")
		return InitSystemD
	}

	if init.run.ExecExistsInPath("update-rc.d") {
		log.Debugln("getInitSystemType = initUpdateRcD")
		log.Debugln("getInitSystemType LEAVE")
		return InitUpdateRcD
	}

	if init.run.ExecExistsInPath("chkconfig") {
		log.Debugln("getInitSystemType = initChkConfig")
		log.Debugln("getInitSystemType LEAVE")
		return InitChkConfig
	}

	log.Debugln("getInitSystemType = initUnknown")
	log.Debugln("getInitSystemType LEAVE")
	return InitUnknown
}

//Start is the package installed
func (init *Init) Start(serviceName string) error {
	if init.init == nil {
		return ErrInvalidInitSystem
	}

	return init.init.Start(serviceName)
}

//StartEx is the package installed
func (init *Init) StartEx(serviceName string, successRegex string) error {
	if init.init == nil {
		return ErrInvalidInitSystem
	}

	return init.init.StartEx(serviceName, successRegex)
}

//Restart the service
func (init *Init) Restart(serviceName string) error {
	if init.init == nil {
		return ErrInvalidInitSystem
	}

	return init.init.Restart(serviceName)
}

//RestartEx the service
func (init *Init) RestartEx(serviceName string, successStopRegex string, successStartRegex string) error {
	if init.init == nil {
		return ErrInvalidInitSystem
	}

	return init.init.RestartEx(serviceName, successStopRegex, successStartRegex)
}

//Status of the service
func (init *Init) Status(serviceName string) (bool, error) {
	if init.init == nil {
		return ErrInvalidInitSystem
	}

	return init.init.Status(serviceName)
}

//StatusEx of the service
func (init *Init) StatusEx(serviceName string, successRegex string) (bool, error) {
	if init.init == nil {
		return ErrInvalidInitSystem
	}

	return init.init.StatusEx(serviceName, successRegex)
}

//Stop the service
func (init *Init) Stop(serviceName string) error {
	if init.init == nil {
		return ErrInvalidInitSystem
	}

	return init.init.Stop(serviceName)
}

//StopEx the service
func (init *Init) StopEx(serviceName string, successRegex string) error {
	if init.init == nil {
		return ErrInvalidInitSystem
	}

	return init.init.StopEx(serviceName, successRegex)
}

//Enable the service
func (init *Init) Enable(serviceName string) error {
	if init.init == nil {
		return ErrInvalidInitSystem
	}

	return init.init.Enable(serviceName)
}

//Disable the service
func (init *Init) Disable(serviceName string) error {
	if init.init == nil {
		return ErrInvalidInitSystem
	}

	return init.init.Disable(serviceName)
}

//AddDependentService to the service
func (init *Init) AddDependentService(serviceName string, depName string) error {
	if init.init == nil {
		return ErrInvalidInitSystem
	}

	return init.init.AddDependentService(serviceName, depName)
}

//RemoveDependentService to the service
func (init *Init) RemoveDependentService(serviceName string, depName string) error {
	if init.init == nil {
		return ErrInvalidInitSystem
	}

	return init.init.RemoveDependentService(serviceName, depName)
}
