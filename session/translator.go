package session

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"go.uber.org/atomic"
)

type translator struct {
	originalRuntimeID uint64
	originalUniqueID  int64

	currentRuntimeID atomic.Uint64
	currentUniqueID  atomic.Int64
}

func newTranslator(data minecraft.GameData) *translator {
	return &translator{
		originalRuntimeID: data.EntityRuntimeID,
		originalUniqueID:  data.EntityUniqueID,

		currentRuntimeID: *atomic.NewUint64(data.EntityRuntimeID),
		currentUniqueID:  *atomic.NewInt64(data.EntityUniqueID),
	}
}

func (s *Session) updateTranslatorData(data minecraft.GameData) {
	s.translator.currentRuntimeID.Store(data.EntityRuntimeID)
	s.translator.currentUniqueID.Store(data.EntityUniqueID)
}

func (s *Session) translatePacket(pk packet.Packet) {
	switch pk := pk.(type) {
	case *packet.ActorEvent:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.ActorPickRequest:
		pk.EntityUniqueID = s.translateUniqueID(pk.EntityUniqueID)
	case *packet.AddActor:
		pk.EntityUniqueID = s.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.AddItemActor:
		pk.EntityUniqueID = s.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.AddPainting:
		pk.EntityUniqueID = s.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.AddPlayer:
		pk.EntityUniqueID = s.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.AdventureSettings:
		pk.PlayerUniqueID = s.translateUniqueID(pk.PlayerUniqueID)
	case *packet.Animate:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.AnimateEntity:
		for i := range pk.EntityRuntimeIDs {
			pk.EntityRuntimeIDs[i] = s.translateRuntimeID(pk.EntityRuntimeIDs[i])
		}
	case *packet.BossEvent:
		pk.BossEntityUniqueID = s.translateUniqueID(pk.BossEntityUniqueID)
		pk.PlayerUniqueID = s.translateUniqueID(pk.PlayerUniqueID)
	case *packet.Camera:
		pk.CameraEntityUniqueID = s.translateUniqueID(pk.CameraEntityUniqueID)
		pk.TargetPlayerUniqueID = s.translateUniqueID(pk.TargetPlayerUniqueID)
	case *packet.CommandOutput:
		pk.CommandOrigin.PlayerUniqueID = s.translateUniqueID(pk.CommandOrigin.PlayerUniqueID)
	case *packet.CommandRequest:
		pk.CommandOrigin.PlayerUniqueID = s.translateUniqueID(pk.CommandOrigin.PlayerUniqueID)
	case *packet.ContainerOpen:
		pk.ContainerEntityUniqueID = s.translateUniqueID(pk.ContainerEntityUniqueID)
	case *packet.DebugInfo:
		pk.PlayerUniqueID = s.translateUniqueID(pk.PlayerUniqueID)
	case *packet.Emote:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.EmoteList:
		pk.PlayerRuntimeID = s.translateRuntimeID(pk.PlayerRuntimeID)
	case *packet.Event:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.Interact:
		pk.TargetEntityRuntimeID = s.translateRuntimeID(pk.TargetEntityRuntimeID)
	case *packet.MobArmourEquipment:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.MobEffect:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.MobEquipment:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.MotionPredictionHints:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.MoveActorAbsolute:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.MoveActorDelta:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.MovePlayer:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
		pk.RiddenEntityRuntimeID = s.translateRuntimeID(pk.RiddenEntityRuntimeID)
	case *packet.NPCRequest:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.PlayerAction:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.PlayerList:
		for i := range pk.Entries {
			pk.Entries[i].EntityUniqueID = s.translateUniqueID(pk.Entries[i].EntityUniqueID)
		}
	case *packet.RemoveActor:
		pk.EntityUniqueID = s.translateUniqueID(pk.EntityUniqueID)
	case *packet.Respawn:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.SetActorData:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.SetActorLink:
		pk.EntityLink.RiddenEntityUniqueID = s.translateUniqueID(pk.EntityLink.RiddenEntityUniqueID)
		pk.EntityLink.RiderEntityUniqueID = s.translateUniqueID(pk.EntityLink.RiderEntityUniqueID)
	case *packet.SetActorMotion:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.SetLocalPlayerAsInitialised:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.SetScore:
		for i := range pk.Entries {
			pk.Entries[i].EntityUniqueID = s.translateUniqueID(pk.Entries[i].EntityUniqueID)
		}
	case *packet.SetScoreboardIdentity:
		for i := range pk.Entries {
			pk.Entries[i].EntityUniqueID = s.translateUniqueID(pk.Entries[i].EntityUniqueID)
		}
	case *packet.ShowCredits:
		pk.PlayerRuntimeID = s.translateRuntimeID(pk.PlayerRuntimeID)
	case *packet.SpawnParticleEffect:
		pk.EntityUniqueID = s.translateUniqueID(pk.EntityUniqueID)
	case *packet.StartGame:
		pk.EntityUniqueID = s.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.StructureBlockUpdate:
		pk.Settings.LastEditingPlayerUniqueID = s.translateUniqueID(pk.Settings.LastEditingPlayerUniqueID)
	case *packet.StructureTemplateDataRequest:
		pk.Settings.LastEditingPlayerUniqueID = s.translateUniqueID(pk.Settings.LastEditingPlayerUniqueID)
	case *packet.TakeItemActor:
		pk.ItemEntityRuntimeID = s.translateRuntimeID(pk.ItemEntityRuntimeID)
		pk.TakerEntityRuntimeID = s.translateRuntimeID(pk.TakerEntityRuntimeID)
	case *packet.UpdateAttributes:
		pk.EntityRuntimeID = s.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.UpdateEquip:
		pk.EntityUniqueID = s.translateUniqueID(pk.EntityUniqueID)
	case *packet.UpdatePlayerGameType:
		pk.PlayerUniqueID = s.translateUniqueID(pk.PlayerUniqueID)
	case *packet.UpdateTrade:
		pk.VillagerUniqueID = s.translateUniqueID(pk.VillagerUniqueID)
		pk.EntityUniqueID = s.translateUniqueID(pk.EntityUniqueID)
	}
}

func (s *Session) translateRuntimeID(id uint64) uint64 {
	original := s.translator.originalRuntimeID
	current := s.translator.currentRuntimeID.Load()

	if original == id {
		return current
	} else if current == id {
		return original
	}
	return id
}

func (s *Session) translateUniqueID(id int64) int64 {
	original := s.translator.originalUniqueID
	current := s.translator.currentUniqueID.Load()

	if original == id {
		return current
	} else if current == id {
		return original
	}
	return id
}
