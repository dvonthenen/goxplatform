package initd

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"

	fs "github.com/dvonthenen/goxplatform/fs"
	run "github.com/dvonthenen/goxplatform/run"
)

var (
	//ErrExecEmptyOutput failed to generate any output
	ErrExecEmptyOutput = errors.New("Failed to generate any output")

	//ErrAddDependencyFailed Failed to add the dependency
	ErrAddDependencyFailed = errors.New("Failed to add the dependency to the service")

	//ErrDeleteDependencyFailed Failed to remove the dependency from the service
	ErrDeleteDependencyFailed = errors.New("Failed to remove the dependency from the service")

	//ErrSrcNotExist src file doesnt exist
	ErrSrcNotExist = errors.New("Source file does not exist")

	//ErrSrcNotRegularFile src file is not a regular file
	ErrSrcNotRegularFile = errors.New("Source file is not a regular file")

	//ErrDstNotRegularFile dst file is not a regular file
	ErrDstNotRegularFile = errors.New("Destination file is not a regular file")
)

//InitD implementation for InitD
type InitD struct {
	run *run.Run
	fs  *fs.Fs
}

//NewInitD generates a InitD object
func NewInitD() *InitD {
	myRun := run.NewRun()
	myFs := fs.NewFs()
	myInitD := &InitD{
		run: myRun,
		fs:  myFs,
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

func doesDependencyExist(serviceName string, depName string) (bool, error) {
	log.Debugln("doesDependencyExist ENTER")
	log.Debugln("serviceName:", serviceName)
	log.Debugln("depName:", depName)

	fileName := "/etc/init.d/" + serviceName
	log.Debugln("fileName:", fileName)

	file, err := os.Open(fileName)
	if err != nil {
		log.Debugln("Failed on file Open:", err)
		log.Debugln("doesDependencyExist LEAVE")
		return false, err
	}
	defer file.Close()

	r, err := regexp.Compile(depName)
	if err != nil {
		log.Debugln("regexp is invalid")
		log.Debugln("doesDependencyExist LEAVE")
		return false, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		log.Debugln("Line:", line)
		if len(line) == 0 {
			continue
		}

		strings := r.FindStringSubmatch(line)
		if strings != nil || len(strings) == 1 {
			log.Debugln("Match found:", line)
			log.Debugln("doesDependencyExist LEAVE")
			return true, nil
		}
	}

	log.Debugln("Dependency was not found")
	log.Debugln("doesDependencyExist LEAVE")

	return false, nil
}

func makeTmpFileWithNewDep(serviceName string, depName string) error {
	log.Debugln("makeTmpFileWithNewDep ENTER")
	log.Debugln("serviceName:", serviceName)
	log.Debugln("depName:", depName)

	fileName := "/etc/init.d/" + serviceName
	log.Debugln("fileName:", fileName)

	fileNameTmp := "/tmp/" + depName + ".tmp"
	log.Debugln("fileNameTmp:", fileNameTmp)

	sfi, err := os.Stat(fileName)
	if err != nil {
		log.Debugln("Src Stat Failed:", err)
		log.Debugln("makeTmpFileWithNewDep LEAVE")
		return ErrSrcNotExist
	}
	if !sfi.Mode().IsRegular() {
		//cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Debugln("Src file is not regular")
		log.Debugln("makeTmpFileWithNewDep LEAVE")
		return ErrSrcNotRegularFile
	}
	dfi, err := os.Stat(fileNameTmp)
	if err == nil {
		if !(dfi.Mode().IsRegular()) {
			log.Debugln("Dst file is not regular")
			log.Debugln("makeTmpFileWithNewDep LEAVE")
			return ErrDstNotRegularFile
		}
	}

	//Copy the file
	in, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		log.Debugln("Failed to open SRC file:", err)
		log.Debugln("makeTmpFileWithNewDep LEAVE")
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(fileNameTmp, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Debugln("Failed to open DST file:", err)
		log.Debugln("makeTmpFileWithNewDep LEAVE")
		return err
	}
	defer out.Close()

	r, err := regexp.Compile("Required-Start:")
	if err != nil {
		log.Debugln("regexp is invalid")
		log.Debugln("makeTmpFileWithNewDep LEAVE")
		return err
	}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		log.Debugln("Line:", line)
		if len(line) == 0 {
			continue
		}

		str := r.FindStringSubmatch(line)
		if str != nil {
			log.Debugln("Match found:", line)
			newLine := line + " scini"
			out.WriteString(newLine + "\n")
		} else {
			out.WriteString(line + "\n")
		}
	}

	err = out.Sync()
	if err != nil {
		log.Debugln("Failed to flush file:", err)
		log.Debugln("makeTmpFileWithNewDep LEAVE")
		return err
	}

	log.Debugln("makeTmpFileWithNewDep Succeeded")
	log.Debugln("makeTmpFileWithNewDep LEAVE")

	return nil
}

//AddDependentService to the service
func (id *InitD) AddDependentService(serviceName string, depName string) error {
	log.Debugln("InitD::AddDependentService ENTER")
	log.Debugln("serviceName:", serviceName)
	log.Debugln("depName:", depName)

	found, err := doesDependencyExist(serviceName, depName)
	if err != nil {
		log.Debugln("doesDependencyExist Failed. Err:", err)
		log.Debugln("InitD::AddDependentService LEAVE")
		return err
	}
	if found {
		log.Debugln("Dependency already exists!")
		log.Debugln("InitD::AddDependentService LEAVE")
		return nil
	}

	err = makeTmpFileWithNewDep(serviceName, depName)
	if err != nil {
		log.Debugln("makeTmpFileWithNewDep Failed. Err:", err)
		log.Debugln("InitD::AddDependentService LEAVE")
		return err
	}

	err = id.fs.CopyFile("/tmp/"+serviceName+".tmp", "/etc/init.d/"+serviceName)
	if err != nil {
		log.Debugln("CopyFile Failed. Err:", err)
		log.Debugln("InitD::AddDependentService LEAVE")
		return err
	}

	log.Debugln("AddDependentService Succeeded")
	log.Debugln("InitD::AddDependentService LEAVE")
	return nil
}

func makeTmpFileWithoutNewDep(serviceName string, depName string) error {
	log.Debugln("makeTmpFileWithoutNewDep ENTER")
	log.Debugln("serviceName:", serviceName)
	log.Debugln("depName:", depName)

	fileName := "/etc/init.d/" + serviceName
	log.Debugln("fileName:", fileName)

	fileNameTmp := "/tmp/" + depName + ".tmp"
	log.Debugln("fileNameTmp:", fileNameTmp)

	sfi, err := os.Stat(fileName)
	if err != nil {
		log.Debugln("Src Stat Failed:", err)
		log.Debugln("makeTmpFileWithoutNewDep LEAVE")
		return ErrSrcNotExist
	}
	if !sfi.Mode().IsRegular() {
		//cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Debugln("Src file is not regular")
		log.Debugln("makeTmpFileWithoutNewDep LEAVE")
		return ErrSrcNotRegularFile
	}
	dfi, err := os.Stat(fileNameTmp)
	if err == nil {
		if !(dfi.Mode().IsRegular()) {
			log.Debugln("Dst file is not regular")
			log.Debugln("makeTmpFileWithoutNewDep LEAVE")
			return ErrDstNotRegularFile
		}
	}

	//Copy the file
	in, err := os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		log.Debugln("Failed to open SRC file:", err)
		log.Debugln("makeTmpFileWithoutNewDep LEAVE")
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(fileNameTmp, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Debugln("Failed to open DST file:", err)
		log.Debugln("makeTmpFileWithoutNewDep LEAVE")
		return err
	}
	defer out.Close()

	r, err := regexp.Compile("Required-Start:")
	if err != nil {
		log.Debugln("regexp is invalid")
		log.Debugln("makeTmpFileWithoutNewDep LEAVE")
		return err
	}

	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		log.Debugln("Line:", line)
		if len(line) == 0 {
			continue
		}

		str := r.FindStringSubmatch(line)
		if str != nil {
			log.Debugln("Match found:", line)
			newLine := strings.Replace(line, " scini", "", -1)
			newLine = strings.Replace(newLine, "scini ", "", -1)
			newLine = strings.Replace(newLine, "scini", "", -1)
			out.WriteString(newLine + "\n")
		} else {
			out.WriteString(line + "\n")
		}
	}

	err = out.Sync()
	if err != nil {
		log.Debugln("Failed to flush file:", err)
		log.Debugln("makeTmpFileWithoutNewDep LEAVE")
		return err
	}

	log.Debugln("makeTmpFileWithoutNewDep Succeeded")
	log.Debugln("makeTmpFileWithoutNewDep LEAVE")

	return nil
}

//RemoveDependentService to the service
func (id *InitD) RemoveDependentService(serviceName string, depName string) error {
	log.Debugln("InitD::RemoveDependentService ENTER")
	log.Debugln("serviceName:", serviceName)
	log.Debugln("depName:", depName)

	found, err := doesDependencyExist(serviceName, depName)
	if err != nil {
		log.Debugln("doesDependencyExist Failed. Err:", err)
		log.Debugln("InitD::RemoveDependentService LEAVE")
		return err
	}
	if !found {
		log.Debugln("Dependency doesnt exists!")
		log.Debugln("InitD::RemoveDependentService LEAVE")
		return nil
	}

	err = makeTmpFileWithoutNewDep(serviceName, depName)
	if err != nil {
		log.Debugln("makeTmpFileWithoutNewDep Failed. Err:", err)
		log.Debugln("InitD::RemoveDependentService LEAVE")
		return err
	}

	err = id.fs.CopyFile("/tmp/"+serviceName+".tmp", "/etc/init.d/"+serviceName)
	if err != nil {
		log.Debugln("CopyFile Failed. Err:", err)
		log.Debugln("InitD::RemoveDependentService LEAVE")
		return err
	}

	log.Debugln("RemoveDependentService Succeeded")
	log.Debugln("InitD::RemoveDependentService LEAVE")
	return nil
}
