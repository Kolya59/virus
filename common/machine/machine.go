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

/*func (p *Port) ToProtobuf() *machine.Port {
	converted := &machine.Port{
		Port:    uint32(p.Port),
		Service: p.Service,
	}
	return converted
}*/

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

/*func (a *Address) ToProtobuf() *machine.Address {
	converted := &machine.Address{}
	converted.Ip = a.IP.String()
	converted.Ports = make([]*machine.Port, len(a.OpenedPorts))
	for i, port := range a.OpenedPorts {
		converted.Ports[i] = port.ToProtobuf()
	}
	return converted
}*/

type ExtendedIface struct {
	Iface     string
	Addresses []Address
	Err       error
}

/*func (i *ExtendedIface) ToProtobuf() *machine.Iface {
	converted := &machine.Iface{}
	if i.Err == nil {
		addresses := make([]*machine.Address, len(i.Addresses))
		for i, address := range i.Addresses {
			addresses[i] = address.ToProtobuf()
		}
	} else {
		converted.Error = i.Err.Error()
	}

	return converted
}*/

type Machine struct {
	ID         string
	ExternalIP string
	Ifaces     []ExtendedIface
	Err        error
}

func (m *Machine) GetIPS() {
	ifaces, err := net.Interfaces()
	if err != nil {
		m.Err = err
		return
	}
	var extenedIfaces []ExtendedIface
	for i, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			m.Ifaces[i].Err = err
			continue
		}
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
			extenedIfaces = append(extenedIfaces, ExtendedIface{
				Iface:     iface.Name,
				Addresses: extendedAddresses,
			})
		}
	}
	m.Ifaces = extenedIfaces
}

/*func (m *Machine) ToProtobuf() *machine.Machine {
	converted := &machine.Machine{}
	if m.Err != nil {
		converted.Error = m.Err.Error()
	} else {
		ifaces := make([]*machine.Iface, len(m.Ifaces))
		for i, iface := range m.Ifaces {
			ifaces[i] = iface.ToProtobuf()
		}
		converted.Ifaces = ifaces
	}

	return converted
}*/
