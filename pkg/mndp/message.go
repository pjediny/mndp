package mndp

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sort"
)

// Message ...
type Message struct {
	Src     net.Addr
	TypeTag uint16
	SeqNo   uint16
	Fields  map[TLVTag]TLV
}

func (msg Message) fieldsString() string {
	count := len(msg.Fields)
	if count < 1 {
		return ""
	}

	keys := make([]TLVTag, 0, count)
	for k := range msg.Fields {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	result := msg.Fields[keys[0]].String()
	if count == 1 {
		return result
	}
	separator := "; "
	for _, key := range keys[1:] {
		result = result + separator + msg.Fields[key].String()
	}
	return result
}

func (msg Message) IsResponse() bool {
	return !(msg.TypeTag == TagMNDP && msg.SeqNo == 0)
}

func (msg Message) IsRefreshRequest() bool {
	return msg.TypeTag == TagMNDP && msg.SeqNo == 0
}

func (msg Message) String() string {
	if msg.IsRefreshRequest() {
		return fmt.Sprintf("MNDP[%d/%d] (%s): <%s>", msg.TypeTag, msg.SeqNo, msg.Src, "Refresh")
	}
	return fmt.Sprintf("MNDP[%d/%d] (%s): <%s>\n  %s", msg.TypeTag, msg.SeqNo, msg.Src, "Response", msg.fieldsString())
}

func ReadMsg(r io.Reader) *Message {
	var record Message
	record.Fields = make(map[TLVTag]TLV, TagMAX+1)
	err := binary.Read(r, binary.BigEndian, &record.TypeTag)
	if err == io.EOF {
		return nil
	} else if err != nil {
		return nil
	}

	err = binary.Read(r, binary.BigEndian, &record.SeqNo)
	if err == io.EOF {
		return nil
	} else if err != nil {
		return nil
	}

	for {
		tlv := ReadTLV(r)
		if tlv == nil {
			break
		}
		record.Fields[tlv.Tag] = *tlv
	}

	return &record
}
