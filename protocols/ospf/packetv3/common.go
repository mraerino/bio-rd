package packetv3

import (
	"bytes"
	"encoding/binary"

	"github.com/bio-routing/bio-rd/net"
	"github.com/bio-routing/tflow2/convert"
	"github.com/pkg/errors"
)

// Serializable represents any packet which can be serialized
// to bytes to be on the wire
type Serializable interface {
	Serialize(buf *bytes.Buffer)
}

// ID is a common type used for 32-bit IDs in OSPF
type ID uint32

func DeserializeID(buf *bytes.Buffer) (ID, int, error) {
	var id uint32
	if err := binary.Read(buf, binary.BigEndian, id); err != nil {
		return ID(id), 0, errors.Wrap(err, "unable to read ID from buffer")
	}
	return ID(id), 4, nil
}

func (i ID) Serialize(buf *bytes.Buffer) {
	buf.Write(convert.Uint32Byte(uint32(i)))
}

// bitmasks for flags in RouterOptions
const (
	RouterOptV6 uint16 = 1 << iota
	RouterOptE
	_
	RouterOptN
	RouterOptR
	RouterOptDC
	_
	_
	RouterOptAF
)

type RouterOptions struct {
	_     uint8
	Flags uint16
}

func (r *RouterOptions) Serialize(buf *bytes.Buffer) {
	buf.WriteByte(0)
	buf.Write(convert.Uint16Byte(uint16(r.Flags)))
}

type LSType uint16

func (t LSType) Serialize(buf *bytes.Buffer) {
	buf.Write(convert.Uint16Byte(uint16(t)))
}

type deserializableIP struct {
	higher uint64
	lower  uint64
}

func (ip deserializableIP) ToNetIP() net.IP {
	return net.IPv6(ip.higher, ip.lower)
}

func serializeIPv6(ip net.IP, buf *bytes.Buffer) {
	if ip.IsIPv4() {
		for i := 0; i < 16; i++ {
			buf.WriteByte(0)
		}
		return
	}

	buf.Write(ip.Bytes())
}
