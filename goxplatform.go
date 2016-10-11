package goxplatform

import (
	log "github.com/Sirupsen/logrus"

	fs "github.com/dvonthenen/goxplatform/fs"
	inst "github.com/dvonthenen/goxplatform/inst"
	nw "github.com/dvonthenen/goxplatform/nw"
	run "github.com/dvonthenen/goxplatform/run"
	sys "github.com/dvonthenen/goxplatform/sys"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.Infoln("Initializing goxplatform...")
}

//NewFs generates a new Fs object
func NewFs() *fs.Fs {
	myFs := &fs.Fs{}
	return myFs
}

//NewNw generates a new Run object
func NewNw() *nw.Nw {
	myNw := &nw.Nw{}
	return myNw
}

//NewRun generates a new Run object
func NewRun() *run.Run {
	myRun := &run.Run{}
	return myRun
}

//NewSys generates a new Sys object
func NewSys() *sys.Sys {
	return sys.NewSys()
}

//NewInst generates a new Run object
func NewInst() *inst.Inst {
	myInst := &inst.Inst{}
	return myInst
}
