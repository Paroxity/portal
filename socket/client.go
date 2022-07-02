package socket

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/paroxity/portal/internal"
	"github.com/paroxity/portal/socket/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"go.uber.org/atomic"
	"net"
	"sync"
)

// Client represents a client connected over the TCP socket system.
type Client struct {
	log  internal.Logger
	conn net.Conn

	pool packet.Pool

	sendMu sync.Mutex
	hdr    *packet.Header
	buf    *bytes.Buffer

	name          string
	authenticated atomic.Bool
}

// NewClient creates a new socket Client with default allocations and required data. It pre-allocates 4096
// bytes to prevent allocations during runtime as much as possible.
func NewClient(conn net.Conn, log internal.Logger) *Client {
	return &Client{
		log:  log,
		conn: conn,

		pool: packet.NewPool(),
		buf:  bytes.NewBuffer(make([]byte, 0, 4096)),
		hdr:  &packet.Header{},
	}
}

// Name returns the name the client authenticated with.
func (c *Client) Name() string {
	return c.name
}

// Close closes the client and related connections.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Authenticate marks the client as authenticated and gives it the provided name.
func (c *Client) Authenticate(name string) {
	if c.authenticated.CAS(false, true) {
		c.name = name
	}
}

// Authenticated returns if the client has been authenticated or not.
func (c *Client) Authenticated() bool {
	return c.authenticated.Load()
}

// ReadPacket reads a packet from the connection and returns it. The client is expected to prefix the packet
// payload with 4 bytes for the length of the payload.
func (c *Client) ReadPacket() (pk packet.Packet, err error) {
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

	defer func() {
		if recoveredErr := recover(); recoveredErr != nil {
			err = fmt.Errorf("%T: %w", pk, recoveredErr.(error))
		}
	}()
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

	return nil
}
