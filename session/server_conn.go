package session

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"io"
	"time"
)

// ServerConn represents a connection that can be used between the proxy and a server to communicate on behalf of a client.
type ServerConn interface {
	io.Closer
	// GameData returns specific game data set to the connection for the player to be initialised with. This data is
	// obtained from the server during the login process.
	GameData() minecraft.GameData
	// DoSpawnTimeout starts the game for the client in the server with a timeout after which an error is returned if the
	// client has not yet spawned by that time. DoSpawnTimeout will start the spawning sequence using the game data found
	// in conn.GameData(), which was sent earlier by the server.
	DoSpawnTimeout(timeout time.Duration) error
	// ReadPacket reads a packet from the Conn, depending on the packet ID that is found in front of the packet data. If
	// a read deadline is set, an error is returned if the deadline is reached before any packet is received. ReadPacket
	// must not be called on multiple goroutines simultaneously. If the packet read was not implemented, a *packet.Unknown
	// is returned, containing the raw payload of the packet read.
	ReadPacket() (packet.Packet, error)
	// WritePacket encodes the packet passed and writes it to the Conn. The encoded data is buffered until the next 20th
	// of a second, after which the data is flushed and sent over the connection.
	WritePacket(packet.Packet) error
}
