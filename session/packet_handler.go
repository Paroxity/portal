package session

import (
	"github.com/paroxity/portal/event"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"log"
)

// handlePackets handles the packets sent between the client and the server. Processes such as runtime
// translations are also handled here.
func handlePackets(s *Session) {
	go func() {
		defer s.Close()
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
					s.serverMu.Lock()
					gameData := s.tempServerConn.GameData()
					_ = s.conn.WritePacket(&packet.ChangeDimension{
						Dimension: packet.DimensionOverworld,
						Position:  gameData.PlayerPosition,
					})

					_ = s.serverConn.Close()

					s.serverConn = s.tempServerConn
					s.tempServerConn = nil
					s.serverMu.Unlock()

					s.updateTranslatorData(gameData)

					// TODO: Set gamemode and stuff
					continue
				}
			case *packet.Text:
				pk.XUID = ""
			case *packet.BookEdit:
				pk.XUID = ""
			}

			ctx := event.C()
			s.handler().HandleServerBoundPacket(ctx, pk)

			ctx.Continue(func() {
				_ = s.ServerConn().WritePacket(pk)
			})
		}
	}()

	go func() {
		defer s.Close()
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

			ctx := event.C()
			s.handler().HandleClientBoundPacket(ctx, pk)

			ctx.Continue(func() {
				_ = s.Conn().WritePacket(pk)
			})
		}
	}()
}
