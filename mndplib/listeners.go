package mndplib

import (
	"bytes"
	"encoding/binary"
	"net"

	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

const maxPacketSize = 1500
const mndpPort = 5678

type udp6Listener struct {
	packetConnection *ipv6.PacketConn
}

type udp4Listener struct {
	packetConnection *ipv4.PacketConn
}

// MNDPListener ...
type MNDPListener struct {
	udp4listener *udp4Listener
	udp6listener *udp6Listener
}

func newUDP6Listener() *udp6Listener {
	var udp6l udp6Listener

	lAddr := net.UDPAddr{IP: net.IPv6unspecified, Port: mndpPort}

	connection, err := net.ListenPacket("udp6", lAddr.String())
	if err != nil {
		return nil
	}

	ipv6connection := ipv6.NewPacketConn(connection)
	udp6l.packetConnection = ipv6connection

	return &udp6l
}

func newUDP4Listener() *udp4Listener {
	var udp4l udp4Listener

	lAddr := net.UDPAddr{IP: net.IPv4zero, Port: mndpPort}

	connection, err := net.ListenPacket("udp4", lAddr.String())
	if err != nil {
		return nil
	}

	ipv4connection := ipv4.NewPacketConn(connection)
	udp4l.packetConnection = ipv4connection

	return &udp4l
}

// NewMNDPListener ...
func NewMNDPListener() *MNDPListener {
	var mndpl MNDPListener

	mndpl.udp4listener = newUDP4Listener()
	mndpl.udp6listener = newUDP6Listener()

	return &mndpl
}

func (listener *udp6Listener) joinMNDPMulticastGroup() {
	multicastGroupUDP6 := &net.UDPAddr{IP: net.IPv6linklocalallnodes, Port: mndpPort}
	netInterfaces, errInterfaces := net.Interfaces()
	if errInterfaces != nil {
		return
	}

	for _, inetInterface := range netInterfaces {
		if (inetInterface.Flags&net.FlagUp != 0) && (inetInterface.Flags&net.FlagMulticast != 0) {
			if err := listener.packetConnection.JoinGroup(&inetInterface, multicastGroupUDP6); err != nil {
				continue
			}
		}
	}
}

func (listener *udp6Listener) listen(ch chan *MNDPMessage) {
	var buff [maxPacketSize]byte

	listener.joinMNDPMulticastGroup()

	for {
		byteCount, _, addr, err := listener.packetConnection.ReadFrom(buff[:])
		if err != nil {
			return
		}

		byteReader := bytes.NewReader(buff[0:byteCount])
		msg := readMndp(byteReader)
		msg.src = addr.String()
		ch <- msg
	}
}

func (listener *udp4Listener) listen(ch chan *MNDPMessage) {
	var buff [maxPacketSize]byte

	for {
		byteCount, _, addr, err := listener.packetConnection.ReadFrom(buff[:])
		if err != nil {
			return
		}

		byteReader := bytes.NewReader(buff[0:byteCount])
		msg := readMndp(byteReader)
		msg.src = addr.String()
		ch <- msg
	}
}

// Listen ...
func (listener *MNDPListener) Listen(ch chan *MNDPMessage) {
	listener.RequestRefresh()
	go listener.udp4listener.listen(ch)
	go listener.udp6listener.listen(ch)
}

func directedBroadcast(inet *net.IPNet) net.IP { // works when the n is a prefix, otherwise...
	if inet.IP.To4() == nil {
		return nil
	}
	ip := make(net.IP, len(inet.IP.To4()))
	binary.BigEndian.PutUint32(ip, binary.BigEndian.Uint32(inet.IP.To4())|^binary.BigEndian.Uint32(net.IP(inet.Mask).To4()))
	return ip
}

func (listener *udp4Listener) requestRefresh() {
	networks, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, network := range networks {
		_, inet, err := net.ParseCIDR(network.String())
		if err != nil {
			continue
		}

		broadcast := directedBroadcast(inet)
		if broadcast != nil {
			listener.packetConnection.WriteTo([]byte{0, 0, 0, 0}, nil, &net.UDPAddr{IP: broadcast, Port: mndpPort})
		}
	}

}

func (listener *udp6Listener) requestRefresh() {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}

	for _, iface := range interfaces {
		if (iface.Flags&net.FlagUp != 0) && (iface.Flags&net.FlagMulticast != 0) {
			listener.packetConnection.WriteTo([]byte{0, 0, 0, 0}, nil, &net.UDPAddr{IP: net.IPv6linklocalallnodes, Port: mndpPort, Zone: iface.Name})
		}
	}

}

// RequestRefresh ...
func (listener *MNDPListener) RequestRefresh() {
	listener.udp4listener.requestRefresh()
	listener.udp6listener.requestRefresh()
}
