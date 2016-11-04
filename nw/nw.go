package util

import (
	"net"
	"strings"

	log "github.com/Sirupsen/logrus"
)

//Nw is a static class that provides Network related functions
type Nw struct{}

//NewNw generates a Network object
func NewNw() *Nw {
	myNet := &Nw{}
	return myNet
}

//AutoDiscoverIP attempt to discover the IP for this host
func (nw *Nw) AutoDiscoverIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Warnln("Failed to get Interfaces", err)
		return "", err
	}

	var ip string
	for _, i := range ifaces {
		if strings.Contains(i.Name, "lo") || strings.Contains(i.Name, "docker") {
			log.Debugln("Skipping interface:", i.Name)
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			log.Infoln("Failed to get IPs on Interface", err)
			continue
		}
		// handle err
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP.String()
				log.Debugln("IPNet:", ip)
			case *net.IPAddr:
				ip = v.IP.String()
				log.Debugln("IPAddr:", ip)
			}
			if len(ip) > 0 {
				break
			}
		}

		log.Infoln("IP Discovered:", ip)
		break
	}

	return ip, nil
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
