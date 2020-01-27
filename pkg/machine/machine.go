package machine

import (
	"net"
	"time"

	portscanner "github.com/anvie/port-scanner"
)

type Port struct {
	Port    int
	Service string
}

type Address struct {
	IP          net.IP
	OpenedPorts []Port
}

func (a *Address) ScanPorts() {
	ps := portscanner.NewPortScanner(a.IP.String(), time.Millisecond, 5)
	opened := ps.GetOpenedPort(1, 1000)
	if len(opened) != 0 {
		a.OpenedPorts = make([]Port, len(opened))

		for i := 0; i < len(opened); i++ {
			a.OpenedPorts[i] = Port{
				Port:    opened[i],
				Service: ps.DescribePort(opened[i]),
			}
		}
	}
}

type ExtenedIface struct {
	Iface     string
	Addresses []Address
	Err       error
}

type Machine struct {
	Ifaces []ExtenedIface
	Err    error
}

func (m *Machine) GetIPS() {
	ifaces, err := net.Interfaces()
	if err != nil {
		m.Err = err
		return
	}
	var extenedIfaces []ExtenedIface
	for i, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			m.Ifaces[i].Err = err
			continue
		}
		// TODO Refactor this fucking spaghetti trash
		var extendedAddresses []Address
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip != nil {
				extendedAddress := Address{IP: ip}
				extendedAddress.ScanPorts()
				if len(extendedAddress.OpenedPorts) != 0 {
					extendedAddresses = append(extendedAddresses, extendedAddress)
				}
			}
		}
		if extendedAddresses != nil {
			extenedIfaces = append(extenedIfaces, ExtenedIface{
				Iface:     iface.Name,
				Addresses: extendedAddresses,
			})
		}
	}
	m.Ifaces = extenedIfaces
}
