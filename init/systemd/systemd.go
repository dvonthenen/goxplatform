package systemd

import (
	"errors"

	log "github.com/Sirupsen/logrus"

	run "github.com/dvonthenen/goxplatform/run"
)

var (
	//ErrSystemCtlFailed Systemctl command failed
	ErrSystemCtlFailed = errors.New("Systemctl command failed")
)

//SystemD implementation for SystemD
type SystemD struct {
	run *run.Run
}

//NewSystemD generates a SystemD object
func NewSystemD() *SystemD {
	myRun := run.NewRun()
	mySystemD := &SystemD{
		run: myRun,
	}
	return mySystemD
}

//Start the service
func (sd *SystemD) Start(serviceName string) error {
	return sd.StartEx(serviceName, "")
}

//StartEx the service
func (sd *SystemD) StartEx(serviceName string, successRegex string) error {
	log.Debugln("SystemD::StartEx ENTER")
	log.Debugln("serviceName:", serviceName)
	//successRegex is ignored

	cmdLine := "systemctl start " + serviceName
	output, err := sd.run.CommandOutput(cmdLine)
	if err != nil {
		log.Debugln("StartEx Failed:", err)
		log.Debugln("SystemD::StartEx LEAVE")
		return err
	}
	if len(output) > 0 {
		log.Debugln("StartEx Failed")
		log.Debugln("SystemD::StartEx LEAVE")
		return ErrSystemCtlFailed	
	}

	log.Debugln("StartEx Succeeded")
	log.Debugln("SystemD::StartEx LEAVE")
	return nil
}

//Restart the service
func (sd *SystemD) Restart(serviceName string) error {
	err := sd.Stop(serviceName)
	if err != nil {
		return err
	}
	err = sd.Start(serviceName)
	if err != nil {
		return err
	}

	return nil
}

//RestartEx the service
func (sd *SystemD) RestartEx(serviceName string, successStopRegex string, successStartRegex string) error {
	err := sd.StopEx(serviceName, successStopRegex)
	if err != nil {
		return err
	}
	err = sd.StartEx(serviceName, successStartRegex)
	if err != nil {
		return err
	}

	return nil
}

//Status of the service
func (sd *SystemD) Status(serviceName string) (bool, error) {
	return sd.StatusEx(serviceName, "active")
}

//StatusEx of the service
func (sd *SystemD) StatusEx(serviceName string, successRegex string) (bool, error) {
	log.Debugln("SystemD::StatusEx ENTER")
	log.Debugln("serviceName:", serviceName)
	log.Debugln("successRegex:", successRegex)

	cmdLine := "systemctl is-active " + serviceName
	output, err := sd.run.CommandOutput(cmdLine)
	if err != nil {
		log.Debugln("StatusEx Failed:", err)
		log.Debugln("SystemD::StatusEx LEAVE")
		return false, err
	}
	if output == "active" {
		log.Debugln("StatusEx Succeeded")
		log.Debugln("SystemD::StatusEx LEAVE")
		return true, nil	
	}

	log.Debugln("StatusEx Failed:", output)
	log.Debugln("SystemD::StatusEx LEAVE")
	return false, nil
}

//Stop the service
func (sd *SystemD) Stop(serviceName string) error {
	return sd.StopEx(serviceName, "TODO")
}

//StopEx the service
func (sd *SystemD) StopEx(serviceName string, successRegex string) error {
	log.Debugln("SystemD::StopEx ENTER")
	log.Debugln("serviceName:", serviceName)
	//successRegex is ignored

	cmdLine := "systemctl stop " + serviceName
	output, err := sd.run.CommandOutput(cmdLine)
	if err != nil {
		log.Debugln("StopEx Failed:", err)
		log.Debugln("SystemD::StopEx LEAVE")
		return err
	}
	if len(output) > 0 {
		log.Debugln("StopEx Failed")
		log.Debugln("SystemD::StopEx LEAVE")
		return ErrSystemCtlFailed	
	}

	log.Debugln("StopEx Succeeded")
	log.Debugln("SystemD::StopEx LEAVE")
	return nil
}

//Enable the service
func (sd *SystemD) Enable(serviceName string) error {
	log.Debugln("SystemD::Enable ENTER")
	log.Debugln("serviceName:", serviceName)

	cmdLine1 := "systemctl is-enabled " + serviceName
	output1, err1 := sd.run.CommandOutput(cmdLine1)
	if err1 != nil {
		log.Debugln("Enable Failed:", err1)
		log.Debugln("SystemD::Enable LEAVE")
		return err1
	}
	if output1 == "enabled" {
		log.Debugln("Enable Succeeded")
		log.Debugln("SystemD::Enable LEAVE")
		return nil	
	}

	cmdLine2 := "systemctl enable " + serviceName
	output2, err2 := sd.run.CommandOutput(cmdLine2)
	if err2 != nil {
		log.Debugln("Enable Failed:", err2)
		log.Debugln("SystemD::Enable LEAVE")
		return err2
	}
	if len(output2) > 0 {
		log.Debugln("Enable Failed")
		log.Debugln("SystemD::Enable LEAVE")
		return ErrSystemCtlFailed	
	}

	log.Debugln("Enable Succeeded")
	log.Debugln("SystemD::Enable LEAVE")
	return nil
}

//Disable the service
func (sd *SystemD) Disable(serviceName string) error {
	log.Debugln("SystemD::Disable ENTER")
	log.Debugln("serviceName:", serviceName)

	cmdLine1 := "systemctl is-enabled " + serviceName
	output1, err1 := sd.run.CommandOutput(cmdLine1)
	if err1 != nil {
		log.Debugln("Disable Failed:", err1)
		log.Debugln("SystemD::Disable LEAVE")
		return err1
	}
	if output1 == "disabled" {
		log.Debugln("Disable Succeeded")
		log.Debugln("SystemD::Disable LEAVE")
		return nil	
	}

	cmdLine2 := "systemctl disable " + serviceName
	err2 := sd.run.Command(cmdLine2, "Removed symlink", "")
	if err2 != nil {
		log.Debugln("Disable Failed:", err2)
		log.Debugln("SystemD::Disable LEAVE")
		return err2
	}

	log.Debugln("Disable Succeeded")
	log.Debugln("SystemD::Disable LEAVE")
	return nil
}

func doesAfterExist(run *run.Run, serviceName string) bool {
	cmdLine := "grep -e After= /etc/systemd/system/" + serviceName + ".service"
	output, err := run.CommandOutput(cmdLine)
	if err != nil {
		return false
	}
	if len(output) > 0 {
		return true	
	}
	return false
}

//AddDependentService to the service
func (sd *SystemD) AddDependentService(serviceName string, depName string) error {
	log.Debugln("SystemD::AddDependentService ENTER")
	log.Debugln("serviceName:", serviceName)

	cmdLine1 := "grep -e After=.*" + depName + " /etc/systemd/system/" + serviceName + ".service"
	output1, err1 := sd.run.CommandOutput(cmdLine1)
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

	if !doesAfterExist(sd.run, serviceName) {
		cmdLine2 := "sed -i 's/Before=/After=" + depName + ".service\\nBefore=/' /etc/systemd/system/" +
			serviceName + ".service"
		output2, err2 := sd.run.CommandOutput(cmdLine2)
		if err2 != nil {
			log.Debugln("AddDependentService Failed:", err2)
			log.Debugln("InitD::AddDependentService LEAVE")
			return err2
		}
		if len(output2) > 0 {
			log.Errorln("AddDependentService Failed:", output2)
			log.Debugln("InitD::AddDependentService LEAVE")
			return ErrSystemCtlFailed
		}
	}

	cmdLine2 := "sed -i 's/After=/After=" + depName + ".service /' /etc/systemd/system/" +
		serviceName + ".service"
	output2, err2 := sd.run.CommandOutput(cmdLine2)
	if err2 != nil {
		log.Errorln("AddDependentService Failed:", err2)
		log.Debugln("InitD::AddDependentService LEAVE")
		return err2
	}
	if len(output2) > 0 {
		log.Debugln("AddDependentService Failed")
		log.Debugln("InitD::AddDependentService LEAVE")
		return ErrSystemCtlFailed
	}

	log.Debugln("AddDependentService Succeeded")
	log.Debugln("SystemD::AddDependentService LEAVE")
	return nil
}

//RemoveDependentService to the service
func (sd *SystemD) RemoveDependentService(serviceName string, depName string) error {
	log.Debugln("SystemD::RemoveDependentService ENTER")
	log.Debugln("serviceName:", serviceName)

	cmdLine1 := "grep -e After=.*" + depName + " /etc/systemd/system/" + serviceName + ".service"
	output1, err1 := sd.run.CommandOutput(cmdLine1)
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

	cmdLine2 := "sed -i 's/ " + depName + "/ /' /etc/systemd/system/" +
		serviceName + ".service"
	output2, err2 := sd.run.CommandOutput(cmdLine2)
	if err2 != nil {
		log.Errorln("RemoveDependentService Failed:", err2)
		log.Debugln("InitD::RemoveDependentService LEAVE")
		return err2
	}
	if len(output2) > 0 {
		log.Debugln("RemoveDependentService Failed")
		log.Debugln("InitD::RemoveDependentService LEAVE")
		return ErrSystemCtlFailed
	}

	log.Debugln("RemoveDependentService Succeeded")
	log.Debugln("SystemD::RemoveDependentService LEAVE")
	return nil
}
