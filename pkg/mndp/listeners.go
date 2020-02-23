package mndp

import (
	"bytes"
	"encoding/binary"
	"net"
)

const maxPacketSize = 1500
const mndpPort = 5678

type udp6Listener struct {
	packetConnection net.PacketConn
}

type udp4Listener struct {
	packetConnection net.PacketConn
}

// Listener ...
type Listener struct {
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

	udp6l.packetConnection = connection

	return &udp6l
}

func newUDP4Listener() *udp4Listener {
	var udp4l udp4Listener

	lAddr := net.UDPAddr{IP: net.IPv4zero, Port: mndpPort}

	connection, err := net.ListenPacket("udp4", lAddr.String())
	if err != nil {
		return nil
	}

	udp4l.packetConnection = connection

	return &udp4l
}

// NewListener ...
func NewListener() *Listener {
	var mndpl Listener

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
			if _, err := net.ListenMulticastUDP("udp6", &inetInterface, multicastGroupUDP6); err != nil {
				continue
			}
		}
	}
}

func (listener *udp6Listener) listen(ch chan *Message) {
	var buff [maxPacketSize]byte

	listener.joinMNDPMulticastGroup()

	for {
		byteCount, addr, err := listener.packetConnection.ReadFrom(buff[:])
		if err != nil {
			return
		}

		byteReader := bytes.NewReader(buff[0:byteCount])
		msg := ReadMsg(byteReader)
		msg.Src = addr
		ch <- msg
	}
}

func (listener *udp4Listener) listen(ch chan *Message) {
	var buff [maxPacketSize]byte

	for {
		byteCount, addr, err := listener.packetConnection.ReadFrom(buff[:])
		if err != nil {
			return
		}

		byteReader := bytes.NewReader(buff[0:byteCount])
		msg := ReadMsg(byteReader)
		msg.Src = addr
		ch <- msg
	}
}

// Listen ...
func (listener *Listener) Listen(ch chan *Message) {
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
			_, _ = listener.packetConnection.WriteTo([]byte{0, 0, 0, 0}, &net.UDPAddr{IP: broadcast, Port: mndpPort})
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
			_, _ = listener.packetConnection.WriteTo([]byte{0, 0, 0, 0}, &net.UDPAddr{IP: net.IPv6linklocalallnodes, Port: mndpPort, Zone: iface.Name})
		}
	}

}

// RequestRefresh ...
func (listener *Listener) RequestRefresh() {
	listener.udp4listener.requestRefresh()
	listener.udp6listener.requestRefresh()
}
