package socket

import (
	"fmt"
	"github.com/paroxity/portal/config"
	portalpacket "github.com/paroxity/portal/socket/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"net"
	"strings"
)

var (
	handlers map[uint32]packetHandler
	pool     packet.Pool
)

func init() {
	handlers = map[uint32]packetHandler{
		portalpacket.IDAuthRequest: &AuthRequestHandler{},
	}

	packet.Register(portalpacket.IDAuthRequest, func() packet.Packet { return &portalpacket.AuthRequest{} })
	packet.Register(portalpacket.IDAuthResponse, func() packet.Packet { return &portalpacket.AuthResponse{} })
	pool = packet.NewPool()
}

// Listen starts a TCP listener on the configured address to handle incoming connections.
func Listen() error {
	listener, err := net.Listen("tcp", config.SocketAddress())
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClient(NewClient(conn))
	}
}

// handleClient reads packets from a connected client and handles them with the appropriate handler.
func handleClient(c *Client) {
	defer c.Close()

	for {
		pk, err := c.ReadPacket()
		if err != nil {
			if containsAny(err.Error(), "EOF", "closed") {
				return
			}
			fmt.Println(err)
			continue
		}

		h, ok := handlers[pk.ID()]
		if ok {
			if err := h.Handle(pk, c); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("Unhandled packet %T\n", pk)
		}
	}
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}

	return false
}
