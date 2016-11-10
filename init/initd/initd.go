package initd

import (
	"errors"
	"os"

	log "github.com/Sirupsen/logrus"

	run "github.com/dvonthenen/goxplatform/run"
)

var (
	//ErrExecEmptyOutput failed to generate any output
	ErrExecEmptyOutput = errors.New("Failed to generate any output")

	//ErrAddDependencyFailed Failed to add the dependency
	ErrAddDependencyFailed = errors.New("Failed to add the dependency to the service")

	//ErrDeleteDependencyFailed Failed to remove the dependency from the service
	ErrDeleteDependencyFailed = errors.New("Failed to remove the dependency from the service")
)

//InitD implementation for InitD
type InitD struct {
	run *run.Run
}

//NewInitD generates a InitD object
func NewInitD() *InitD {
	myRun := run.NewRun()
	myInitD := &InitD{
		run: myRun,
	}
	return myInitD
}

//Start the service
func (id *InitD) Start(serviceName string) error {
	return id.StartEx(serviceName, "running|process [0-9]+|PID [0-9]+")
}

//StartEx the service
func (id *InitD) StartEx(serviceName string, successRegex string) error {
	log.Debugln("InitD::StartEx ENTER")
	log.Debugln("serviceName:", serviceName)

	cmdLine := "service " + serviceName + " start"
	err := id.run.Command(cmdLine, successRegex, "")
	if err != nil {
		log.Debugln("StartEx: Failed to Start")
		log.Debugln("InitD::StartEx LEAVE")
		return err
	}

	log.Debugln("InitD::StartEx LEAVE")
	return nil
}

//Restart the service
func (id *InitD) Restart(serviceName string) error {
	err := id.Stop(serviceName)
	if err != nil {
		return err
	}
	err = id.Start(serviceName)
	if err != nil {
		return err
	}

	return nil
}

//RestartEx the service
func (id *InitD) RestartEx(serviceName string, successStopRegex string, successStartRegex string) error {
	err := id.StopEx(serviceName, successStopRegex)
	if err != nil {
		return err
	}
	err = id.StartEx(serviceName, successStartRegex)
	if err != nil {
		return err
	}

	return nil
}

//Status of the service
func (id *InitD) Status(serviceName string) (bool, error) {
	return id.StatusEx(serviceName, "running|process [0-9]+|PID [0-9]+")
}

//StatusEx of the service
func (id *InitD) StatusEx(serviceName string, successRegex string) (bool, error) {
	log.Debugln("InitD::StatusEx ENTER")
	log.Debugln("serviceName:", serviceName)

	cmdLine := "service " + serviceName + " status"
	err := id.run.Command(cmdLine, successRegex, "")
	if err == run.ErrExecuteFailed {
		log.Debugln("Status: Stopped")
		log.Debugln("InitD::StartEx LEAVE")
		return false, nil
	} else if err != nil {
		log.Debugln("Start Failed:", err)
		log.Debugln("InitD::StartEx LEAVE")
		return false, err
	}

	log.Debugln("Status: Running")
	log.Debugln("InitD::StatusEx LEAVE")
	return true, nil
}

//Stop the service
func (id *InitD) Stop(serviceName string) error {
	return id.StopEx(serviceName, "stop|Shutting down.*SUCCESS|process already finished")
}

//StopEx the service
func (id *InitD) StopEx(serviceName string, successRegex string) error {
	log.Debugln("InitD::StopEx ENTER")
	log.Debugln("serviceName:", serviceName)

	cmdLine := "service " + serviceName + " stop"
	err := id.run.Command(cmdLine, successRegex, "")
	if err != nil {
		log.Debugln("StopEx: Failed to Stop")
		log.Debugln("InitD::StartEx LEAVE")
		return err
	}

	log.Debugln("StopEx: Stopped")
	log.Debugln("InitD::StopEx LEAVE")
	return nil
}

//Enable the service
func (id *InitD) Enable(serviceName string) error {
	log.Debugln("InitD::Enable ENTER")
	log.Debugln("serviceName:", serviceName)

	cmdLine := "echo manual | sudo tee /etc/init/" + serviceName + ".override"
	err := id.run.Command(cmdLine, "manual", "")
	if err != nil {
		log.Debugln("StopEx: Failed to Stop")
		log.Debugln("InitD::StartEx LEAVE")
		return err
	}

	log.Debugln("InitD::Enable LEAVE")
	return nil
}

//Disable the service
func (id *InitD) Disable(serviceName string) error {
	log.Debugln("InitD::Disable ENTER")
	log.Debugln("serviceName:", serviceName)

	fullPath := "/etc/init/" + serviceName + ".override"
	err := os.Remove(fullPath)
	if err != nil {
		log.Debugln("Disable Failed:", err)
		log.Debugln("InitD::Disable LEAVE")
		return err		
	}

	log.Debugln("Disable Succeeded")
	log.Debugln("InitD::Disable LEAVE")
	return nil
}

//AddDependentService to the service
func (id *InitD) AddDependentService(serviceName string, depName string) error {
	log.Debugln("InitD::AddDependentService ENTER")
	log.Debugln("serviceName:", serviceName)

	cmdLine1 := "grep -e Required-Start.*" + depName + " /etc/init.d/" + serviceName
	output1, err1 := id.run.CommandOutput(cmdLine1)
	if err1 != nil {
		log.Debugln("AddDependentService Failed:", err1)
		log.Debugln("InitD::AddDependentService LEAVE")
		return err1
	}
	if len(output1) > 0 {
		log.Debugln("AddDependentService Succeeded")
		log.Debugln("InitD::AddDependentService LEAVE")
		return nil	
	}

	cmdLine2 := "sed -i 's/# Required-Start: /# Required-Start: " + depName + "/' /etc/init.d/" + serviceName
	output2, err2 := id.run.CommandOutput(cmdLine2)
	if err2 != nil {
		log.Errorln("AddDependentService Failed:", err2)
		log.Debugln("InitD::AddDependentService LEAVE")
		return err2
	}
	if len(output2) > 0 {
		log.Debugln("AddDependentService Failed")
		log.Debugln("InitD::AddDependentService LEAVE")
		return ErrAddDependencyFailed
	}

	log.Debugln("AddDependentService Succeeded")
	log.Debugln("InitD::AddDependentService LEAVE")
	return nil
}

//RemoveDependentService to the service
func (id *InitD) RemoveDependentService(serviceName string, depName string) error {
	log.Debugln("InitD::RemoveDependentService ENTER")
	log.Debugln("serviceName:", serviceName)

	cmdLine1 := "grep -e Required-Start.*" + depName + " /etc/init.d/" + serviceName
	output1, err1 := id.run.CommandOutput(cmdLine1)
	if err1 != nil {
		log.Debugln("RemoveDependentService Failed:", err1)
		log.Debugln("InitD::RemoveDependentService LEAVE")
		return err1
	}
	if len(output1) == 0 {
		log.Debugln("RemoveDependentService Succeeded")
		log.Debugln("InitD::RemoveDependentService LEAVE")
		return nil	
	}

	cmdLine2 := "sed -i 's/ " + depName + "//' /etc/init.d/" + serviceName
	output2, err2 := id.run.CommandOutput(cmdLine2)
	if err2 != nil {
		log.Errorln("RemoveDependentService Failed:", err2)
		log.Debugln("InitD::RemoveDependentService LEAVE")
		return err2
	}
	if len(output2) > 0 {
		log.Debugln("RemoveDependentService Failed")
		log.Debugln("InitD::RemoveDependentService LEAVE")
		return ErrDeleteDependencyFailed
	}

	log.Debugln("RemoveDependentService Succeeded")
	log.Debugln("InitD::RemoveDependentService LEAVE")
	return nil
}
