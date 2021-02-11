package socket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/paroxity/portal/server"
	"github.com/paroxity/portal/socket/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
)

// Client represents a client connected over the TCP socket system.
type Client struct {
	conn net.Conn

	pool packet.Pool

	sendMu sync.Mutex
	hdr    *packet.Header
	buf    *bytes.Buffer

	name       string
	clientType uint8
	extraData  map[string]interface{}
}

// NewClient creates a new socket Client with default allocations and required data. It pre-allocates 4096
// bytes to prevent allocations during runtime as much as possible.
func NewClient(conn net.Conn) *Client {
	return &Client{
		conn: conn,

		pool: packet.NewPool(),
		buf:  bytes.NewBuffer(make([]byte, 0, 4096)),
		hdr:  &packet.Header{},

		extraData: make(map[string]interface{}),
	}
}

// Name returns the name the client authenticated with.
func (c *Client) Name() string {
	return c.name
}

// Close closes the client and related connections.
func (c *Client) Close() error {
	logrus.Debugf("Socket connection \"%s\" closed\n", c.name)

	switch c.clientType {
	case packet.ClientTypeServer:
		if name, ok := c.extraData["group"]; ok {
			g, _ := server.GroupFromName(name.(string))
			s, _ := g.Server(c.name)

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

	pk, ok := c.pool[header.PacketID]
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
	c.sendMu.Lock()
	c.hdr.PacketID = pk.ID()
	_ = c.hdr.Write(c.buf)

	pk.Marshal(protocol.NewWriter(c.buf, 0))

	data := c.buf.Bytes()
	c.buf.Reset()
	c.sendMu.Unlock()

	buf := bytes.NewBuffer(make([]byte, 0, 2+len(data)))

	if err := binary.Write(buf, binary.LittleEndian, int32(len(data))); err != nil {
		return err
	}
	if _, err := buf.Write(data); err != nil {
		return err
	}

	if _, err := c.conn.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}
