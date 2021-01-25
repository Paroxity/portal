package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"io"
)

type Packet interface {
	// ID returns the ID of the packet. All of these identifiers of packets may be found in id.go.
	ID() uint16
	// Marshal encodes the packet to its binary representation into buf.
	Marshal(w *protocol.Writer)
	// Unmarshal decodes a serialised packet in buf into the Packet instance. The serialised packet passed
	// into Unmarshal will not have a header in it.
	Unmarshal(r *protocol.Reader)
}

type Header struct {
	PacketID uint16
}

func (header *Header) Write(w io.ByteWriter) error {
	if err := w.WriteByte(byte(header.PacketID)); err != nil {
		return err
	}

	if err := w.WriteByte(byte(header.PacketID >> 8)); err != nil {
		return err
	}

	return nil
}

func (header *Header) Read(r io.ByteReader) error {
	b1, err := r.ReadByte()
	if err != nil {
		return err
	}
	b2, err := r.ReadByte()
	if err != nil {
		return err
	}

	header.PacketID = uint16(b1) | uint16(b2)<<8
	return nil
}
