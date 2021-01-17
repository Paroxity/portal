package socket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/paroxity/portal/server"
	portalpacket "github.com/paroxity/portal/socket/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"net"
)

// Client represents a client connected over the TCP socket system.
type Client struct {
	buf  *bytes.Buffer
	conn net.Conn

	name       string
	clientType uint8
	extraData  map[string]interface{}
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		buf:  bytes.NewBuffer(make([]byte, 0, 4096)),
		conn: conn,

		extraData: make(map[string]interface{}),
	}
}

// Name returns the name the client authenticated with.
func (c *Client) Name() string {
	return c.name
}

// Close closes the client and related connections.
func (c *Client) Close() error {
	fmt.Println("Closing connection")

	c.buf.Reset()

	switch c.clientType {
	case portalpacket.ClientTypeServer:
		if name, ok := c.extraData["group"]; ok {
			g, _ := server.GroupFromName(name.(string))
			s, ok := g.Server(c.name)
			if !ok {
				// wtf how
				break
			}

			server_setConn(s, nil)
		}
	}

	return c.conn.Close()
}

// ReadPacket reads a packet from the connection and returns it. The client is expected to prefix the packet
// payload with 4 bytes for the length of the payload.
func (c *Client) ReadPacket() (packet.Packet, error) {
	var l uint32
	if err := binary.Read(c.conn, binary.LittleEndian, &l); err != nil {
		return nil, err
	}

	data := make([]byte, l)
	read, err := c.conn.Read(data)
	if err != nil {
		return nil, err
	}
	if read != int(l) {
		return nil, fmt.Errorf("expected %v bytes, got %v", l, read)
	}

	buf := bytes.NewBuffer(data)

	header := &packet.Header{}
	if err := header.Read(buf); err != nil {
		return nil, err
	}

	pk, ok := pool[header.PacketID]
	if !ok {
		return nil, fmt.Errorf("unknown packet %v", header.PacketID)
	}

	pk.Unmarshal(protocol.NewReader(buf, 0))
	if buf.Len() > 0 {
		return nil, fmt.Errorf("still have %v bytes unread", buf.Len())
	}

	return pk, nil
}

// WritePacket writes a packet to the client. Since it's a TCP connection, the payload is prefixed with a
// length so the client can read the exact length of the packet.
func (c *Client) WritePacket(pk packet.Packet) error {
	header := &packet.Header{PacketID: pk.ID()}
	if err := header.Write(c.buf); err != nil {
		return err
	}

	pk.Marshal(protocol.NewWriter(c.buf, 0))

	data := c.buf.Bytes()
	buf := bytes.NewBuffer(make([]byte, 0, 4+len(data)))

	if err := binary.Write(buf, binary.LittleEndian, int32(len(data))); err != nil {
		return err
	}
	if _, err := buf.Write(data); err != nil {
		return err
	}

	if _, err := c.conn.Write(buf.Bytes()); err != nil {
		return err
	}
	c.buf.Reset()

	return nil
}
