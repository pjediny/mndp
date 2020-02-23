package mndp

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

type TLVTag uint16

const (
	TagMNDP          = 0
	TagMACAddr       = 1
	TagIdentity      = 5
	TagVersion       = 7
	TagPlatform      = 8
	TagUptime        = 10
	TagSoftwareID    = 11
	TagBoard         = 12
	TagUnpack        = 14
	TagIPv6Addr      = 15
	TagInterfaceName = 16
	TagIPv4Addr      = 17

	TagMAX = TagIPv4Addr
)

func (t TLVTag) String() string {
	switch t {
	case TagMACAddr:
		return "MACAddr"
	case TagIdentity:
		return "Identity"
	case TagVersion:
		return "Version"
	case TagPlatform:
		return "Platform"
	case TagUptime:
		return "Uptime"
	case TagSoftwareID:
		return "Software-ID"
	case TagBoard:
		return "Board"
	case TagUnpack:
		return "Unpack"
	case TagIPv6Addr:
		return "IPv6Addr"
	case TagInterfaceName:
		return "Interface"
	case TagIPv4Addr:
		return "IPv4Addr"
	default:
		return fmt.Sprintf("UnknownTag%d", t)
	}
}

// TLV handles message field
type TLV struct {
	Tag    TLVTag
	Length uint16
	Value  []byte
}

func ReadTLV(r io.Reader) *TLV {
	var record TLV
	err := binary.Read(r, binary.BigEndian, &record.Tag)
	if err == io.EOF {
		return nil
	} else if err != nil {
		return nil
	}

	err = binary.Read(r, binary.BigEndian, &record.Length)
	if err == io.EOF {
		return nil
	} else if err != nil {
		return nil
	}

	var buff [maxPacketSize]byte

	c, err := r.Read(buff[:record.Length])
	if err != nil {
		return nil
	}
	if uint16(c) != record.Length {
		return nil
	}

	record.Value = buff[:record.Length]

	return &record
}

func (tlv TLV) String() string {
	value := ""
	switch tlv.Tag {
	case TagMACAddr:
		value = tlv.ValAsHardwareAddr().String()
	case TagIdentity:
		value = tlv.ValAsString()
	case TagVersion:
		value = tlv.ValAsString()
	case TagPlatform:
		value = tlv.ValAsString()
	case TagUptime:
		value = tlv.ValAsDuration().String()
	case TagSoftwareID:
		value = tlv.ValAsString()
	case TagBoard:
		value = tlv.ValAsString()
	case TagUnpack:
		value = tlv.ValAsHexString()
	case TagIPv6Addr:
		value = tlv.ValAsIP().String()
	case TagInterfaceName:
		value = tlv.ValAsString()
	case TagIPv4Addr:
		value = tlv.ValAsIP().String()
	default:
		value = tlv.ValAsHexString()
	}
	return fmt.Sprintf("%s: %s", tlv.Tag, value)
}

func (tlv TLV) ValAsString() string {
	return string(tlv.Value)
}

func (tlv TLV) ValAsHexString() string {
	return fmt.Sprintf("%x", tlv.Value)
}

func (tlv TLV) ValAsHardwareAddr() net.HardwareAddr {
	return tlv.Value
}

func (tlv TLV) ValAsDuration() time.Duration {
	return time.Duration(binary.LittleEndian.Uint32(tlv.Value)) * time.Second
}

func (tlv TLV) ValAsIP() net.IP {
	return tlv.Value
}
