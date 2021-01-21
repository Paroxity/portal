package session

import (
	"github.com/paroxity/portal/event"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sirupsen/logrus"
	"log"
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
				log.Println(err)
				return
			}
			s.translatePacket(pk)

			switch pk := pk.(type) {
			case *packet.BookEdit:
				pk.XUID = ""
			case *packet.PlayerAction:
				if pk.ActionType == packet.PlayerActionDimensionChangeDone && s.transferring.CAS(true, false) {
					s.serverMu.Lock()
					gameData := s.tempServerConn.GameData()
					_ = s.conn.WritePacket(&packet.ChangeDimension{
						Dimension: packet.DimensionOverworld,
						Position:  gameData.PlayerPosition,
					})

					var w sync.WaitGroup
					w.Add(5)
					go func() {
						s.clearEntities()
						w.Done()
					}()
					go func() {
						s.clearPlayerList()
						w.Done()
					}()
					go func() {
						s.clearEffects()
						w.Done()
					}()
					go func() {
						s.clearBossBars()
						w.Done()
					}()
					go func() {
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

					_ = s.conn.WritePacket(&packet.LevelEvent{EventType: packet.EventStopRain, EventData: 10000})
					_ = s.conn.WritePacket(&packet.LevelEvent{EventType: packet.EventStopThunder})
					_ = s.conn.WritePacket(&packet.SetDifficulty{Difficulty: uint32(gameData.Difficulty)})
					_ = s.conn.WritePacket(&packet.GameRulesChanged{GameRules: gameData.GameRules})
					_ = s.conn.WritePacket(&packet.SetPlayerGameType{GameType: gameData.PlayerGameMode})

					w.Wait()

					_ = s.serverConn.Close()

					s.serverConn = s.tempServerConn
					s.tempServerConn = nil
					s.serverMu.Unlock()

					s.updateTranslatorData(gameData)

					logrus.Infof("%s finished transferring\n", s.conn.IdentityData().DisplayName)
					continue
				}
			case *packet.Text:
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
				if conn != s.ServerConn() {
					continue
				}
				log.Println(err)
				return
			}
			s.translatePacket(pk)

			switch pk := pk.(type) {
			case *packet.AddActor:
				s.addEntity(pk.EntityUniqueID)
			case *packet.AddItemActor:
				s.addEntity(pk.EntityUniqueID)
			case *packet.AddPainting:
				s.addEntity(pk.EntityUniqueID)
			case *packet.AddPlayer:
				s.addEntity(pk.EntityUniqueID)
			case *packet.BossEvent:
				if pk.EventType == packet.BossEventShow {
					s.addBossBar(pk.BossEntityUniqueID)
				} else if pk.EventType == packet.BossEventHide {
					s.removeBossBar(pk.BossEntityUniqueID)
				}
			case *packet.MobEffect:
				if pk.Operation == packet.MobEffectAdd {
					s.addEffect(pk.EffectType)
				} else if pk.Operation == packet.MobEffectRemove {
					s.removeEffect(pk.EffectType)
				}
			case *packet.PlayerList:
				if pk.ActionType == packet.PlayerListActionAdd {
					for _, e := range pk.Entries {
						s.addToPlayerList(e.UUID)
					}
				} else {
					for _, e := range pk.Entries {
						s.removeFromPlayerList(e.UUID)
					}
				}
			case *packet.RemoveActor:
				s.removeEntity(pk.EntityUniqueID)
			case *packet.RemoveObjective:
				s.removeScoreboard(pk.ObjectiveName)
			case *packet.SetDisplayObjective:
				s.addScoreboard(pk.ObjectiveName)
			}

			ctx := event.C()
			s.handler().HandleClientBoundPacket(ctx, pk)

			ctx.Continue(func() {
				_ = s.Conn().WritePacket(pk)
			})
		}
	}()
}
