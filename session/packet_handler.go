package session

import (
	"fmt"
	"github.com/paroxity/wormhole/server"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"log"
)

// handlePackets handles the packets sent between the client and the server. Processes such as runtime
// translations are also handled here.
func handlePackets(s *Session) {
	go func() {
		// TODO: Defer and close connections
		for {
			pk, err := s.Conn().ReadPacket()
			if err != nil {
				log.Println(err)
				return
			}

			switch pk := pk.(type) {
			case *packet.PlayerAction:
				if pk.ActionType == packet.PlayerActionDimensionChangeDone && s.Transferring() {
					s.SetTransferring(false)

					old := s.serverConn
					conn := s.tempServerConn

					pos := conn.GameData().PlayerPosition
					_ = s.conn.WritePacket(&packet.ChangeDimension{
						Dimension: packet.DimensionOverworld,
						Position:  pos,
					})

					_ = old.Close()

					s.serverConn = conn
					s.tempServerConn = nil

					// TODO: Set gamemode and stuff
					continue
				}
			case *packet.Text:
				switch pk.Message {
				case "hub1":
					if s.server.Name() == "Hub1" {
						_ = s.conn.WritePacket(&packet.Text{
							Message: "You are already on Hub1",
						})
						continue
					}
					srv, ok := server.DefaultGroup().Server("hub1")
					if !ok {
						_ = s.conn.WritePacket(&packet.Text{
							Message: text.Colourf("<red>Server not found</red>"),
						})
						continue
					}
					if err := s.Transfer(srv); err != nil {
						_ = s.conn.WritePacket(&packet.Text{
							Message: text.Colourf("<red>%s</red>", err.Error()),
						})
					}
					continue
				case "hub2":
					if s.server.Name() == "Hub2" {
						_ = s.conn.WritePacket(&packet.Text{
							Message: "You are already on Hub1",
						})
						continue
					}
					srv, ok := server.DefaultGroup().Server("hub2")
					if !ok {
						_ = s.conn.WritePacket(&packet.Text{
							Message: text.Colourf("<red>Server not found</red>"),
						})
						continue
					}
					if err := s.Transfer(srv); err != nil {
						_ = s.conn.WritePacket(&packet.Text{
							Message: text.Colourf("<red>%s</red>", err.Error()),
						})
					}
					continue
				case "server":
					_ = s.conn.WritePacket(&packet.Text{
						Message: fmt.Sprintf("You are on %s", s.server.Name()),
					})
					continue
				}
			}

			_ = s.ServerConn().WritePacket(pk)
		}
	}()

	go func() {
		// TODO: Defer and close connections
		for {
			conn := s.ServerConn()
			pk, err := conn.ReadPacket()
			if err != nil {
				log.Println(err)
				if conn != s.ServerConn() {
					continue
				}
				return
			}

			// TODO: Runtime ID translations

			_ = s.Conn().WritePacket(pk)
		}
	}()
}
