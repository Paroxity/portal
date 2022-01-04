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
		pk.EntityMetadata = t.translateEntityMetadata(pk.EntityMetadata)
		for i := range pk.EntityLinks {
			pk.EntityLinks[i] = t.translateEntityLink(pk.EntityLinks[i])
		}
	case *packet.AddItemActor:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
		pk.EntityMetadata = t.translateEntityMetadata(pk.EntityMetadata)
	case *packet.AddPainting:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.AddPlayer:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
		pk.EntityMetadata = t.translateEntityMetadata(pk.EntityMetadata)
		for i := range pk.EntityLinks {
			pk.EntityLinks[i] = t.translateEntityLink(pk.EntityLinks[i])
		}
	case *packet.AddVolumeEntity:
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
	case *packet.ClientBoundMapItemData:
		for i, x := range pk.TrackedObjects {
			if x.Type == protocol.MapObjectTypeEntity {
				x.EntityUniqueID = t.translateUniqueID(x.EntityUniqueID)
				pk.TrackedObjects[i] = x
			}
		}
	case *packet.CommandBlockUpdate:
		if !pk.Block {
			pk.MinecartEntityRuntimeID = t.translateRuntimeID(pk.MinecartEntityRuntimeID)
		}
	case *packet.CommandOutput:
		pk.CommandOrigin.PlayerUniqueID = t.translateUniqueID(pk.CommandOrigin.PlayerUniqueID)
	case *packet.CommandRequest:
		pk.CommandOrigin.PlayerUniqueID = t.translateUniqueID(pk.CommandOrigin.PlayerUniqueID)
	case *packet.ContainerOpen:
		pk.ContainerEntityUniqueID = t.translateUniqueID(pk.ContainerEntityUniqueID)
	case *packet.CreatePhoto:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
	case *packet.DebugInfo:
		pk.PlayerUniqueID = t.translateUniqueID(pk.PlayerUniqueID)
	case *packet.Emote:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.EmoteList:
		pk.PlayerRuntimeID = t.translateRuntimeID(pk.PlayerRuntimeID)
	case *packet.Event:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
		switch data := pk.EventData.(type) {
		case *protocol.MobKilledEventData:
			data.KillerEntityUniqueID = t.translateUniqueID(data.KillerEntityUniqueID)
			data.VictimEntityUniqueID = t.translateUniqueID(data.VictimEntityUniqueID)
		case *protocol.BossKilledEventData:
			data.BossEntityUniqueID = t.translateUniqueID(data.BossEntityUniqueID)
		case *protocol.PetDiedEventData:
			data.KillerEntityUniqueID = t.translateUniqueID(data.KillerEntityUniqueID)
			data.PetEntityUniqueID = t.translateUniqueID(data.PetEntityUniqueID)
		}
	case *packet.Interact:
		pk.TargetEntityRuntimeID = t.translateRuntimeID(pk.TargetEntityRuntimeID)
	case *packet.InventoryTransaction:
		switch data := pk.TransactionData.(type) {
		case *protocol.UseItemOnEntityTransactionData:
			data.TargetEntityRuntimeID = t.translateRuntimeID(data.TargetEntityRuntimeID)
		}
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
	case *packet.PhotoTransfer:
		pk.OwnerEntityUniqueID = t.translateUniqueID(pk.OwnerEntityUniqueID)
	case *packet.PlayerAction:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.PlayerList:
		for i := range pk.Entries {
			pk.Entries[i].EntityUniqueID = t.translateUniqueID(pk.Entries[i].EntityUniqueID)
		}
	case *packet.RemoveActor:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
	case *packet.RemoveVolumeEntity:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.Respawn:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
	case *packet.SetActorData:
		pk.EntityRuntimeID = t.translateRuntimeID(pk.EntityRuntimeID)
		pk.EntityMetadata = t.translateEntityMetadata(pk.EntityMetadata)
	case *packet.SetActorLink:
		pk.EntityLink = t.translateEntityLink(pk.EntityLink)
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
	case *packet.UpdateBlockSynced:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
	case *packet.UpdateEquip:
		pk.EntityUniqueID = t.translateUniqueID(pk.EntityUniqueID)
	case *packet.UpdatePlayerGameType:
		pk.PlayerUniqueID = t.translateUniqueID(pk.PlayerUniqueID)
	case *packet.UpdateSubChunkBlocks:
		for i, entry := range pk.Blocks {
			pk.Blocks[i].SyncedUpdateEntityUniqueID = uint64(t.translateUniqueID(int64(entry.SyncedUpdateEntityUniqueID)))
		}
		for i, entry := range pk.Extra {
			pk.Extra[i].SyncedUpdateEntityUniqueID = uint64(t.translateUniqueID(int64(entry.SyncedUpdateEntityUniqueID)))
		}
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

// translateEntityLink returns the correct entity link for the client to function properly.
func (t *translator) translateEntityLink(x protocol.EntityLink) protocol.EntityLink {
	x.RiddenEntityUniqueID = t.translateUniqueID(x.RiddenEntityUniqueID)
	x.RiderEntityUniqueID = t.translateUniqueID(x.RiderEntityUniqueID)
	return x
}

// translateEntityMetadata returns the correct entity metadata for the client to function properly. It translates the
// entity IDs to make sure there are no conflicts after transferring servers.
func (t *translator) translateEntityMetadata(x map[uint32]interface{}) map[uint32]interface{} {
	for k, v := range x {
		switch k {
		case 5: // Owner ID
			x[5] = t.translateUniqueID(v.(int64))
		case 6: // Target ID
			x[6] = t.translateUniqueID(v.(int64))
		case 17: // Shooter ID
			x[17] = t.translateUniqueID(v.(int64))
		case 37: // Leash holder ID
			x[37] = t.translateUniqueID(v.(int64))
		case 88: // Player agent ID
			x[88] = t.translateUniqueID(v.(int64))
		case 124: // Base Runtime ID
			x[124] = t.translateRuntimeID(v.(uint64))
		}
	}
	return x
}
