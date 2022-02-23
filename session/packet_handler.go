package session

import (
	"errors"
	"github.com/paroxity/portal/event"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sirupsen/logrus"
	"sync"
)

// handlePackets handles the packets sent between the client and the server. Processes such as runtime
// translations are also handled here.
func handlePackets(s *Session) {
	go func() {
		defer s.Close()
		for {
			pk, err := s.Conn().ReadPacket()
			if err != nil {
				s.log.Errorf("failed to read packet from connection: %v", err)
				return
			}
			s.translatePacket(pk)

			switch pk := pk.(type) {
			case *packet.BookEdit:
				pk.XUID = ""
			case *packet.PlayerAction:
				if pk.ActionType == protocol.PlayerActionDimensionChangeDone {
					if s.transferring.Load() {
						s.serverMu.Lock()
						gameData := s.tempServerConn.GameData()
						_ = s.conn.WritePacket(&packet.ChangeDimension{
							Dimension: packet.DimensionOverworld,
							Position:  gameData.PlayerPosition,
						})
						_ = s.conn.WritePacket(&packet.StopSound{StopAll: true})

						var w sync.WaitGroup
						w.Add(2)
						go func() {
							s.clearEntities()
							s.clearEffects()
							w.Done()
						}()
						go func() {
							s.clearPlayerList()
							s.clearBossBars()
							s.clearScoreboard()
							w.Done()
						}()

						_ = s.conn.WritePacket(&packet.MovePlayer{
							EntityRuntimeID: s.originalRuntimeID,
							Position:        gameData.PlayerPosition,
							Pitch:           gameData.Pitch,
							Yaw:             gameData.Yaw,
							Mode:            packet.MoveModeReset,
						})

						_ = s.conn.WritePacket(&packet.LevelEvent{EventType: packet.LevelEventStopRaining, EventData: 10000})
						_ = s.conn.WritePacket(&packet.LevelEvent{EventType: packet.LevelEventStopThunderstorm})
						_ = s.conn.WritePacket(&packet.SetDifficulty{Difficulty: uint32(gameData.Difficulty)})
						_ = s.conn.WritePacket(&packet.GameRulesChanged{GameRules: gameData.GameRules})
						_ = s.conn.WritePacket(&packet.SetPlayerGameType{GameType: gameData.PlayerGameMode})

						w.Wait()

						_ = s.serverConn.Close()

						s.serverConn = s.tempServerConn
						s.tempServerConn = nil
						s.serverMu.Unlock()

						s.updateTranslatorData(gameData)

						s.transferring.Store(false)
						s.postTransfer.Store(true)

						s.log.Infof("%s finished transferring to %s", s.Conn().IdentityData().DisplayName, s.Server().Name())
						continue
					} else if s.postTransfer.CAS(true, false) {
						continue
					}
				}
			case *packet.Text:
				pk.XUID = ""
			}

			if s.Transferring() {
				continue
			}

			ctx := event.C()
			s.handler().HandleServerBoundPacket(ctx, pk)

			ctx.Continue(func() {
				_ = s.ServerConn().WritePacket(pk)
			})
		}
	}()

	go func() {
		for {
			conn := s.ServerConn()
			pk, err := conn.ReadPacket()
			if err != nil {
				if conn != s.ServerConn() {
					continue
				}
				ctx := event.C()
				s.handler().HandleServerDisconnect(ctx)

				c := false
				ctx.Continue(func() {
					c = true
					if disconnect, ok := errors.Unwrap(err).(minecraft.DisconnectError); ok {
						logrus.Debugln(disconnect.Error())
						_ = s.conn.WritePacket(&packet.Disconnect{Message: disconnect.Error()})
					}
					s.Close()
				})
				if c {
					return
				}
				continue
			}
			s.translatePacket(pk)

			switch pk := pk.(type) {
			case *packet.AddActor:
				s.entities.Add(pk.EntityUniqueID)
			case *packet.AddItemActor:
				s.entities.Add(pk.EntityUniqueID)
			case *packet.AddPainting:
				s.entities.Add(pk.EntityUniqueID)
			case *packet.AddPlayer:
				s.entities.Add(pk.EntityUniqueID)
			case *packet.BossEvent:
				if pk.EventType == packet.BossEventShow {
					s.bossBars.Add(pk.BossEntityUniqueID)
				} else if pk.EventType == packet.BossEventHide {
					s.bossBars.Remove(pk.BossEntityUniqueID)
				}
			case *packet.MobEffect:
				if pk.Operation == packet.MobEffectAdd {
					s.effects.Add(pk.EffectType)
				} else if pk.Operation == packet.MobEffectRemove {
					s.effects.Remove(pk.EffectType)
				}
			case *packet.PlayerList:
				if pk.ActionType == packet.PlayerListActionAdd {
					for _, e := range pk.Entries {
						s.playerList.Add(e.UUID)
					}
				} else {
					for _, e := range pk.Entries {
						s.playerList.Remove(e.UUID)
					}
				}
			case *packet.RemoveActor:
				s.entities.Remove(pk.EntityUniqueID)
			case *packet.RemoveObjective:
				s.scoreboards.Remove(pk.ObjectiveName)
			case *packet.SetDisplayObjective:
				s.scoreboards.Add(pk.ObjectiveName)
			}

			ctx := event.C()
			s.handler().HandleClientBoundPacket(ctx, pk)

			ctx.Continue(func() {
				_ = s.Conn().WritePacket(pk)
			})
		}
	}()
}
