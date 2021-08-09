package session

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"go.uber.org/atomic"
)

// translator represents a data structure which holds the data needed to transfer runtime IDs for a session.
type translator struct {
	originalRuntimeID uint64
	originalUniqueID  int64

	currentRuntimeID atomic.Uint64
	currentUniqueID  atomic.Int64
}

// newTranslator creates a new translator based off of the provided GameData from the initial server.
func newTranslator(data minecraft.GameData) *translator {
	return &translator{
		originalRuntimeID: data.EntityRuntimeID,
		originalUniqueID:  data.EntityUniqueID,

		currentRuntimeID: *atomic.NewUint64(data.EntityRuntimeID),
		currentUniqueID:  *atomic.NewInt64(data.EntityUniqueID),
	}
}

// updateTranslatorData updates the translator with the runtime IDs from a new server.
func (t *translator) updateTranslatorData(data minecraft.GameData) {
	t.currentRuntimeID.Store(data.EntityRuntimeID)
	t.currentUniqueID.Store(data.EntityUniqueID)
}

// translatePacket translates the runtime IDs in packets sent by the client and the connected server. If this
// process is not done, weird things would happen visually on the client.
func (t *translator) translatePacket(pk packet.Packet) {
	switch pk := pk.(type) {
	case *packet.ActorEvent:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.ActorPickRequest:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
	case *packet.AddActor:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.AddItemActor:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.AddPainting:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.AddPlayer:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.AdventureSettings:
		pk.PlayerUniqueID = t.translateUniqueID(pk.PlayerUniqueID)
	case *packet.Animate:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.AnimateEntity:
		for i := range pk.EntityRuntimeIDs {
			pk.EntityRuntimeIDs[i] = t.translateRuntimeID(pk.EntityRuntimeIDs[i])
		}
	case *packet.BossEvent:
		pk.BossEntityUniqueID = t.translateUniqueID(pk.BossEntityUniqueID)
		pk.PlayerUniqueID = t.translateUniqueID(pk.PlayerUniqueID)
	case *packet.Camera:
		pk.CameraEntityUniqueID = t.translateUniqueID(pk.CameraEntityUniqueID)
		pk.TargetPlayerUniqueID = t.translateUniqueID(pk.TargetPlayerUniqueID)
	case *packet.CommandOutput:
		pk.CommandOrigin.PlayerUniqueID = t.translateUniqueID(pk.CommandOrigin.PlayerUniqueID)
	case *packet.CommandRequest:
		pk.CommandOrigin.PlayerUniqueID = t.translateUniqueID(pk.CommandOrigin.PlayerUniqueID)
	case *packet.ContainerOpen:
		pk.ContainerEntityUniqueID = t.translateUniqueID(pk.ContainerEntityUniqueID)
	case *packet.DebugInfo:
		pk.PlayerUniqueID = t.translateUniqueID(pk.PlayerUniqueID)
	case *packet.Emote:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.EmoteList:
		pk.PlayerRuntimeID = t.translateRuntimeID(pk.PlayerRuntimeID)
	case *packet.Event:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.Interact:
		pk.TargetEntityRuntimeID = t.translateRuntimeID(pk.TargetEntityRuntimeID)
	case *packet.MobArmourEquipment:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.MobEffect:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.MobEquipment:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.MotionPredictionHints:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.MoveActorAbsolute:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.MoveActorDelta:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.MovePlayer:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
		pk.RiddenEntityRuntimeID = t.translateRuntimeID(pk.RiddenEntityRuntimeID)
	case *packet.NPCDialogue:
		pk.ActorUniqueID = uint64(t.translateUniqueID(int64(pk.ActorUniqueID)))
	case *packet.NPCRequest:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.PlayerAction:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.PlayerList:
		for i := range pk.Entries {
			pk.Entries[i].EntityUniqueID = t.translateUniqueID(pk.Entries[i].EntityUniqueID)
		}
	case *packet.RemoveActor:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
	case *packet.Respawn:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.SetActorData:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.SetActorLink:
		pk.EntityLink.RiddenEntityUniqueID = t.translateUniqueID(pk.EntityLink.RiddenEntityUniqueID)
		pk.EntityLink.RiderEntityUniqueID = t.translateUniqueID(pk.EntityLink.RiderEntityUniqueID)
	case *packet.SetActorMotion:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.SetLocalPlayerAsInitialised:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.SetScore:
		for i := range pk.Entries {
			if pk.Entries[i].IdentityType != protocol.ScoreboardIdentityFakePlayer {
				pk.Entries[i].EntityUniqueID = t.translateUniqueID(pk.Entries[i].EntityUniqueID)
			}
		}
	case *packet.SetScoreboardIdentity:
		if pk.ActionType != packet.ScoreboardIdentityActionClear {
			for i := range pk.Entries {
				pk.Entries[i].EntityUniqueID = t.translateUniqueID(pk.Entries[i].EntityUniqueID)
			}
		}
	case *packet.ShowCredits:
		pk.PlayerRuntimeID = t.translateRuntimeID(pk.PlayerRuntimeID)
	case *packet.SpawnParticleEffect:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
	case *packet.StartGame:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.StructureBlockUpdate:
		pk.Settings.LastEditingPlayerUniqueID = t.translateUniqueID(pk.Settings.LastEditingPlayerUniqueID)
	case *packet.StructureTemplateDataRequest:
		pk.Settings.LastEditingPlayerUniqueID = t.translateUniqueID(pk.Settings.LastEditingPlayerUniqueID)
	case *packet.TakeItemActor:
		pk.ItemEntityRuntimeID = t.translateRuntimeID(pk.ItemEntityRuntimeID)
		pk.TakerEntityRuntimeID = t.translateRuntimeID(pk.TakerEntityRuntimeID)
	case *packet.UpdateAttributes:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.UpdateEquip:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
	case *packet.UpdatePlayerGameType:
		pk.PlayerUniqueID = t.translateUniqueID(pk.PlayerUniqueID)
	case *packet.UpdateTrade:
		pk.VillagerUniqueID = t.translateUniqueID(pk.VillagerUniqueID)
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
	}
}

// translateRuntimeID returns the correct entity runtime ID for the client to function properly.
func (t *translator) translateRuntimeID(id uint64) uint64 {
	original := t.originalRuntimeID
	current := t.currentRuntimeID.Load()

	if original == id {
		return current
	} else if current == id {
		return original
	}
	return id
}

// translateUniqueID returns the correct entity unique ID for the client to function properly.
func (t *translator) translateUniqueID(id int64) int64 {
	original := t.originalUniqueID
	current := t.currentUniqueID.Load()

	if original == id {
		return current
	} else if current == id {
		return original
	}
	return id
}
