package session

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"log"
)

// handlePackets handles the packets sent between the client and the server. Processes such as runtime
// translations are also handled here.
func handlePackets(s *Session) {
	go func() {
		defer func() {
			s.Close()
		}()
		for {
			pk, err := s.Conn().ReadPacket()
			if err != nil {
				log.Println(err)
				return
			}
			s.translatePacket(pk)

			switch pk := pk.(type) {
			case *packet.PlayerAction:
				if pk.ActionType == packet.PlayerActionDimensionChangeDone && s.transferring.CAS(true, false) {
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

					s.updateTranslatorData(conn.GameData())

					// TODO: Set gamemode and stuff
					continue
				}
			case *packet.Text:
				pk.XUID = ""
			case *packet.BookEdit:
				pk.XUID = ""
			}

			if s.clientPacketFunc != nil {
				if s.clientPacketFunc(pk) {
					return
				}
			}

			_ = s.ServerConn().WritePacket(pk)
		}
	}()

	go func() {
		defer func() {
			s.Close()
		}()
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
			s.translatePacket(pk)

			if s.serverPacketFunc != nil {
				if s.serverPacketFunc(pk) {
					return
				}
			}

			_ = s.Conn().WritePacket(pk)
		}
	}()
}
