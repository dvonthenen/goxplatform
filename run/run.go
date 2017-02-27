package run

import (
	"bufio"
	"bytes"
	"errors"
	"math"
	"os/exec"
	"regexp"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	rootUID = 0
	rootGID = 0

	cmdReties = 4
)

var (
	//ErrBufferCreateFailed creating the buffer failed
	ErrBufferCreateFailed = errors.New("Unable to create the buffer object")

	//ErrScannerCreateFailed creating the scanner failed
	ErrScannerCreateFailed = errors.New("Unable to create the scanner object")

	//ErrReaderCreateFailed creating the reader failed
	ErrReaderCreateFailed = errors.New("Unable to create the reader object")

	//ErrCommandCreateFailed creating the command failed
	ErrCommandCreateFailed = errors.New("Unable to create the command object")

	//ErrExecuteFailed installation package failed
	ErrExecuteFailed = errors.New("The command line failed to execute correctly")
)

//Run is a static class that enables running and capturing command output
type Run struct{}

//NewRun generates a Run object
func NewRun() *Run {
	myRun := &Run{}
	return myRun
}

//ExecExistsInPath returns ture if exec exists in the given path
func (run *Run) ExecExistsInPath(exe string) bool {
	_, err := exec.LookPath(exe)
	return err == nil
}

func command(cmdLine string, successRegex string, failureRegex string) error {
	log.Debugln("command ENTER")
	log.Debugln("Cmdline:", cmdLine)
	log.Debugln("SuccessRegex:", successRegex)
	log.Debugln("FailureRegex:", failureRegex)

	cmd := exec.Command("bash", "-c", cmdLine)
	if cmd == nil {
		log.Errorln("Error creating cmd")
		log.Debugln("command LEAVE")
		return ErrCommandCreateFailed
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorln("Error starting Cmd:", err)
		log.Debugln("command LEAVE")
		return err
	}

	readbuffer := bytes.NewBuffer(out)
	if readbuffer == nil {
		log.Errorln("Error creating buffer")
		log.Debugln("command LEAVE")
		return ErrBufferCreateFailed
	}

	reader := bufio.NewScanner(readbuffer)
	if reader == nil {
		log.Errorln("Error creating reader")
		log.Debugln("command LEAVE")
		return ErrReaderCreateFailed
	}

	failure := false
	succeeded := false
	for reader.Scan() {
		line := reader.Text()
		log.Debugln("Line:", line)
		if failure {
			continue
		}
		if len(failureRegex) > 0 {
			myfail, _ := regexp.MatchString(failureRegex, line)
			if myfail {
				log.Debugln("Line Matched - FAILURE!")
				failure = true
			}
		}
		if succeeded {
			continue
		}
		if len(successRegex) > 0 {
			mysucceed, _ := regexp.MatchString(successRegex, line)
			if mysucceed {
				log.Debugln("Line Matched - SUCCEEDED!")
				succeeded = true
			}
		}
	}

	if failure {
		log.Debugln("Cmdline explicitly failed to execute correctly")
		log.Debugln("command LEAVE")
		return ErrExecuteFailed
	}
	if succeeded {
		log.Debugln("Cmdline executed successful")
		log.Debugln("command LEAVE")
		return nil
	}

	log.Debugln("Cmdline implicitly failed to execute correctly")
	log.Debugln("command LEAVE")
	return ErrExecuteFailed
}

//Command executes a command that monitors output for success or failure
func (run *Run) Command(cmdLine string, successRegex string, failureRegex string) error {
	log.Debugln("Command ENTER")
	log.Debugln("Cmdline:", cmdLine)
	log.Debugln("SuccessRegex:", successRegex)
	log.Debugln("FailureRegex:", failureRegex)

	var err error
	for i := 0; i < cmdReties; i++ {
		log.Debugln("Command attempt #", i+1)

		err = command(cmdLine, successRegex, failureRegex)
		if err == nil {
			log.Debugln("Command Succeeded")
			break
		}

		expDelay := math.Pow(2, float64(i+1))
		log.Debugln("Waiting", expDelay, "before retry.")
		time.Sleep(time.Duration(expDelay) * time.Second)
	}

	log.Debugln("Command LEAVE")
	return err
}

func commandEx(cmdLine string, successRegex string, failureRegex string, waitInSec int) error {
	log.Debugln("commandEx ENTER")
	log.Debugln("Cmdline:", cmdLine)
	log.Debugln("SuccessRegex:", successRegex)
	log.Debugln("FailureRegex:", failureRegex)

	cmd := exec.Command("bash", "-c", cmdLine)
	if cmd == nil {
		log.Errorln("Error creating cmd")
		log.Debugln("commandEx LEAVE")
		return ErrCommandCreateFailed
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorln("Error getting StdoutPipe:", err)
		log.Debugln("commandEx LEAVE")
		return err
	}

	err = cmd.Start()
	if err != nil {
		log.Errorln("Error on cmd start:", err)
		log.Debugln("commandEx LEAVE")
		return err
	}

	stdoutScanner := bufio.NewScanner(stdout)
	if cmd == nil {
		log.Errorln("Error creating scanner")
		log.Debugln("commandEx LEAVE")
		return ErrScannerCreateFailed
	}

	output := ""
	go func() {
		for stdoutScanner.Scan() {
			line := stdoutScanner.Text()
			log.Infoln(line)
			output += line
		}
	}()

	err = cmd.Wait()
	if err != nil {
		log.Warnln("Error on cmd wait:", err)
	}

	cmd.Process.Wait() //this should wait until all child processes are gone

	time.Sleep(time.Duration(waitInSec) * time.Second)

	outputBuffer := bytes.NewBuffer([]byte(output))
	if outputBuffer == nil {
		log.Errorln("Error creating buffer")
		log.Debugln("commandEx LEAVE")
		return ErrBufferCreateFailed
	}

	outputScanner := bufio.NewScanner(outputBuffer)
	if outputScanner == nil {
		log.Errorln("Error creating reader")
		log.Debugln("commandEx LEAVE")
		return ErrScannerCreateFailed
	}

	failure := false
	succeeded := false
	for outputScanner.Scan() {
		line := outputScanner.Text()
		log.Debugln("Line:", line)
		if failure {
			continue
		}
		if len(failureRegex) > 0 {
			myfail, _ := regexp.MatchString(failureRegex, line)
			if myfail {
				log.Debugln("Line Matched - FAILURE!")
				failure = true
			}
		}
		if succeeded {
			continue
		}
		if len(successRegex) > 0 {
			mysucceed, _ := regexp.MatchString(successRegex, line)
			if mysucceed {
				log.Debugln("Line Matched - SUCCEEDED!")
				succeeded = true
			}
		}
	}

	if failure {
		log.Debugln("Cmdline explicitly failed to execute correctly")
		log.Debugln("commandEx LEAVE")
		return ErrExecuteFailed
	}
	if succeeded {
		log.Debugln("Cmdline executed successful")
		log.Debugln("commandEx LEAVE")
		return nil
	}

	log.Debugln("Cmdline implicitly failed to execute correctly")
	log.Debugln("commandEx LEAVE")
	return ErrExecuteFailed
}

//CommandEx executes a command that monitors output for success or failure with a timeout
func (run *Run) CommandEx(cmdLine string, successRegex string, failureRegex string, waitInSec int) error {
	log.Debugln("CommandEx ENTER")
	log.Debugln("Cmdline:", cmdLine)
	log.Debugln("SuccessRegex:", successRegex)
	log.Debugln("FailureRegex:", failureRegex)

	var err error
	for i := 0; i < cmdReties; i++ {
		log.Debugln("CommandEx attempt #", i+1)

		err = commandEx(cmdLine, successRegex, failureRegex, waitInSec)
		if err == nil {
			log.Debugln("CommandEx Succeeded")
			break
		}

		expDelay := math.Pow(2, float64(i+1))
		log.Debugln("Waiting", expDelay, "before retry.")
		time.Sleep(time.Duration(expDelay) * time.Second)
	}

	log.Debugln("CommandEx LEAVE")
	return err
}

func commandOutput(cmdLine string) (string, error) {
	log.Debugln("commandOutput ENTER")
	log.Debugln("Cmdline:", cmdLine)

	cmd := exec.Command("bash", "-c", cmdLine)
	if cmd == nil {
		log.Errorln("Error creating cmd")
		log.Debugln("commandOutput LEAVE")
		return "", ErrCommandCreateFailed
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorln("Error getting output:", err)
		log.Debugln("commandOutput LEAVE")
		return "", err
	}

	output := strings.TrimSpace(string(out))

	log.Debugln("commandOutput Succeeded")
	log.Debugln(output)
	log.Debugln("commandOutput LEAVE")
	return output, nil
}

//CommandOutput executes a command that returns the output
func (run *Run) CommandOutput(cmdLine string) (string, error) {
	log.Debugln("CommandOutput ENTER")
	log.Debugln("Cmdline:", cmdLine)

	var output string
	var err error
	for i := 0; i < cmdReties; i++ {
		log.Debugln("CommandOutput attempt #", i+1)

		output, err = commandOutput(cmdLine)
		if err == nil {
			log.Debugln("CommandOutput Succeeded")
			break
		}

		expDelay := math.Pow(2, float64(i+1))
		log.Debugln("Waiting", expDelay, "before retry.")
		time.Sleep(time.Duration(expDelay) * time.Second)
	}

	log.Debugln("CommandOutput LEAVE")
	return output, err
}
