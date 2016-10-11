package goxplatform

import (
	log "github.com/Sirupsen/logrus"

	fs "github.com/dvonthenen/goxplatform/fs"
	inst "github.com/dvonthenen/goxplatform/inst"
	run "github.com/dvonthenen/goxplatform/run"
	sys "github.com/dvonthenen/goxplatform/sys"
)

func init() {
	log.SetLevel(log.InfoLevel)
	log.Infoln("Initializing goxplatform...")
}

//NewFs generates a new Fs object
func NewFs() *fs.Fs {
	fs := &fs.Fs{}
	return fs
}

//NewRun generates a new Run object
func NewRun() *run.Run {
	run := &run.Run{}
	return run
}

//NewSys generates a new Sys object
func NewSys() *sys.Sys {
	return sys.NewSys()
}

//NewInst generates a new Run object
func NewInst() *inst.Inst {
	inst := &inst.Inst{}
	return inst
}
