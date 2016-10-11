package goxplatform

import (
	log "github.com/Sirupsen/logrus"

	fs "github.com/dvonthenen/goxplatform/fs"
	inst "github.com/dvonthenen/goxplatform/inst"
	nw "github.com/dvonthenen/goxplatform/nw"
	run "github.com/dvonthenen/goxplatform/run"
	str "github.com/dvonthenen/goxplatform/str"
	sys "github.com/dvonthenen/goxplatform/sys"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.Infoln("Initializing goxplatform...")
}

//XPlatform is a static class that provides System related functions
type XPlatform struct {
	Sys  *sys.Sys
	Fs   *fs.Fs
	Str  *str.Str
	Nw   *nw.Nw
	Run  *run.Run
	Inst *inst.Inst
}

//New generates a Sys object
func New() *XPlatform {
	mySys := sys.NewSys()
	myFs := fs.NewFs()
	myStr := str.NewStr()
	myNw := nw.NewNw()
	myRun := run.NewRun()
	myInst := inst.NewInst()

	myXPlatform := &XPlatform{
		Sys:  mySys,
		Fs:   myFs,
		Str:  myStr,
		Nw:   myNw,
		Run:  myRun,
		Inst: myInst,
	}

	return myXPlatform
}
