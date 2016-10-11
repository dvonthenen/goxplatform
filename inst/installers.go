package inst

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"

	fs "github.com/dvonthenen/goxplatform/fs"
	common "github.com/dvonthenen/goxplatform/inst/common"
	deb "github.com/dvonthenen/goxplatform/inst/deb"
	ipm "github.com/dvonthenen/goxplatform/inst/ipackagemgr"
	sys "github.com/dvonthenen/goxplatform/sys"
)

var (
	//ErrInvalidOsType the OS is not valid
	ErrInvalidOsType = errors.New("Invalid OS Type")
)

//Inst is a static class that captures install package rules
type Inst struct {
	fs  *fs.Fs
	sys *sys.Sys
	ipm ipm.IPackageMgr
}

//NewInst generates a Inst object
func NewInst() *Inst {
	myFs := fs.NewFs()
	mySys := sys.NewSys()

	var myIpm ipm.IPackageMgr
	switch mySys.GetOsType() {
	case sys.OsUbuntu:
		myIpm = deb.NewDeb()
	}

	myInst := &Inst{
		fs:  myFs,
		sys: mySys,
		ipm: myIpm,
	}

	return myInst
}

//IsInstalled is the package installed
func (inst *Inst) IsInstalled(packageName string) error {
	if inst.ipm == nil {
		return ErrInvalidOsType
	}

	return inst.ipm.IsInstalled(packageName)
}

//GetInstalledVersion get the installed version of the package
func (inst *Inst) GetInstalledVersion(packageName string, parseVersion bool) (string, error) {
	if inst.ipm == nil {
		return "", ErrInvalidOsType
	}

	return inst.ipm.GetInstalledVersion(packageName, parseVersion)
}

//DownloadPackage downloads a payload specified by the URI and
//returns the local path for where the bits land
func (inst *Inst) DownloadPackage(installPackageURI string) (string, error) {
	log.Infoln("downloadPackage ENTER")
	log.Infoln("installPackageURI=", installPackageURI)

	path, err := inst.fs.GetFullPath()
	if err != nil {
		log.Errorln("GetFullPath Failed:", err)
		log.Infoln("downloadPackage LEAVE")
		return "", err
	}

	filename := inst.fs.GetFilenameFromURIOrFullPath(installPackageURI)
	log.Infoln("Filename:", filename)

	fullpath := inst.fs.AppendSlash(path) + filename
	log.Infoln("Fullpath:", fullpath)

	//create a downloaded file
	output, err := os.Create(fullpath)
	if err != nil {
		log.Errorln("Create File Failed:", err)
		log.Infoln("downloadPackage LEAVE")
		return "", err
	}

	//get the "executor" file
	resp, err := http.Get(installPackageURI)
	if err != nil {
		log.Errorln("HTTP GET Failed:", err)
		log.Infoln("downloadPackage LEAVE")
		return "", err
	}
	defer resp.Body.Close()

	_, err = io.Copy(output, resp.Body)
	if err != nil {
		log.Errorln("IO Copy Failed:", err)
		log.Infoln("downloadPackage LEAVE")
		return "", err
	}
	output.Close()

	log.Infoln("downloadPackage Succeeded:", fullpath)
	log.Infoln("downloadPackage LEAVE")
	return fullpath, nil
}

//ParseVersionFromFilename this parses the version string out of the filename
func (inst *Inst) ParseVersionFromFilename(filename string) (string, error) {
	return common.ParseVersionFromFilename(filename)
}

//IsVersionStringHigher checks to see if one version is higher than the current
func (inst *Inst) IsVersionStringHigher(existing string, comparing string) bool {
	log.Debugln("IsVersionStringHigher ENTER")
	log.Debugln("existing:", existing)
	log.Debugln("comparing:", comparing)

	arr1 := strings.Split(existing, ".")
	arr2 := strings.Split(comparing, ".")

	for i := 0; i < len(arr1); i++ {
		tok2, err2 := strconv.Atoi(arr2[i])
		if err2 != nil {
			continue
		}
		tok1, err1 := strconv.Atoi(arr1[i])
		if err1 != nil {
			continue
		}
		if tok2 > tok1 {
			log.Debugln("New Higher:", comparing, ">", existing)
			log.Debugln("IsVersionStringHigher LEAVE")
			return true
		}
	}

	if len(arr2) > len(arr1) {
		log.Debugln("New Higher:", comparing, ">", existing)
		log.Debugln("IsVersionStringHigher LEAVE")
		return true
	}

	log.Debugln("Is Lower")
	log.Debugln("IsVersionStringHigher LEAVE")
	return false
}
