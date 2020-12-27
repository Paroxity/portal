package query

import (
	"github.com/paroxity/portal/config"
	"github.com/sandertv/gophertunnel/minecraft"
	"go.uber.org/atomic"
)

var playerCount atomic.Int64

type StatusProvider struct{}

func (s StatusProvider) ServerStatus(_, _ int) minecraft.ServerStatus {
	return minecraft.ServerStatus{
		ServerName:  config.MOTD(),
		PlayerCount: PlayerCount(),
		MaxPlayers:  config.MaxPlayers(),
	}
}

func PlayerCount() int {
	return int(playerCount.Load())
}

func IncrementPlayerCount() {
	playerCount.Add(1)
}

func DecrementPlayerCount() {
	playerCount.Sub(1)
}
