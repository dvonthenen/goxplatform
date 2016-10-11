package util

import (
	"net"

	log "github.com/Sirupsen/logrus"
)

//Nw is a static class that provides Network related functions
type Nw struct{}

//NewNw generates a Network object
func NewNw() *Nw {
	myNet := &Nw{}
	return myNet
}

//ParseIP creates an IP object from a string
func (nw *Nw) ParseIP(address string) net.IP {
	addr, err := net.LookupIP(address)
	if err != nil {
		log.Errorln("LookupIP:", err)
	}
	if len(addr) < 1 {
		log.Errorln("failed to parse IP from address", address)
	}
	return addr[0]
}
