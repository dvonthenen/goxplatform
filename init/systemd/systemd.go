package systemd

import (
	"errors"
	"strings"

	log "github.com/Sirupsen/logrus"
	ini "github.com/go-ini/ini"

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
	return sd.StopEx(serviceName, "")
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

	iniFile := "/etc/systemd/system/" + serviceName + ".service"
	cfg, err := ini.Load(iniFile)
	if err != nil {
		log.Errorln("Load INI Failed. Err:", err)
		log.Debugln("SystemD::AddDependentService LEAVE")
		return err
	}

	key, err := cfg.Section("Unit").GetKey("After")
	if err != nil {
		log.Debugln("Key After does not exist. Create one!")
		_, err = cfg.Section("Unit").NewKey("After", depName)
		if err != nil {
			log.Errorln("Failed to create NewKey(After). Err:", err)
			log.Debugln("SystemD::AddDependentService LEAVE")
			return err
		}

		err = cfg.SaveTo(iniFile)
		if err != nil {
			log.Errorln("Failed to SaveTo File. Err:", err)
			log.Debugln("SystemD::AddDependentService LEAVE")
			return err
		}

		log.Debugln("AddDependentService Succeeded")
		log.Debugln("SystemD::AddDependentService LEAVE")
		return nil
	}

	value := key.Value()
	if strings.Contains(value, serviceName) {
		log.Debugln("Already contains dependency", serviceName, ". AddDependentService Succeeded")
		log.Debugln("SystemD::AddDependentService LEAVE")
		return nil
	}

	newValue := serviceName + " " + value
	key.SetValue(newValue)

	err = cfg.SaveTo(iniFile)
	if err != nil {
		log.Errorln("Failed to SaveTo File. Err:", err)
		log.Debugln("SystemD::AddDependentService LEAVE")
		return err
	}

	log.Debugln("AddDependentService Succeeded")
	log.Debugln("SystemD::AddDependentService LEAVE")
	return nil
}

//RemoveDependentService to the service
func (sd *SystemD) RemoveDependentService(serviceName string, depName string) error {
	log.Debugln("SystemD::RemoveDependentService ENTER")
	log.Debugln("serviceName:", serviceName)

	iniFile := "/etc/systemd/system/" + serviceName + ".service"
	cfg, err := ini.Load(iniFile)
	if err != nil {
		log.Errorln("Load INI Failed. Err:", err)
		log.Debugln("SystemD::RemoveDependentService LEAVE")
		return err
	}

	key, err := cfg.Section("Unit").GetKey("After")
	if err != nil {
		log.Debugln("RemoveDependentService Succeeded")
		log.Debugln("SystemD::RemoveDependentService LEAVE")
		return nil
	}

	value := key.Value()
	if !strings.Contains(value, serviceName) {
		log.Debugln("Doesnt contain dependency", serviceName, ". RemoveDependentService Succeeded")
		log.Debugln("SystemD::RemoveDependentService LEAVE")
		return nil
	}

	newValue := strings.Replace(value, " "+serviceName, "", -1)
	newValue = strings.Replace(newValue, serviceName+" ", "", -1)
	newValue = strings.Replace(newValue, serviceName, "", -1)
	key.SetValue(newValue)

	err = cfg.SaveTo(iniFile)
	if err != nil {
		log.Errorln("Failed to SaveTo File. Err:", err)
		log.Debugln("SystemD::RemoveDependentService LEAVE")
		return err
	}

	log.Debugln("RemoveDependentService Succeeded")
	log.Debugln("SystemD::RemoveDependentService LEAVE")
	return nil
}
