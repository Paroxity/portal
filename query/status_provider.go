package query

import (
	"github.com/paroxity/portal/config"
	"github.com/sandertv/gophertunnel/minecraft"
	"go.uber.org/atomic"
)

var playerCount atomic.Int64

// StatusProvider represents a status provider that displays information based on the configuration file and
// an internal player count.
type StatusProvider struct{}

// ServerStatus ...
func (s StatusProvider) ServerStatus(_, _ int) minecraft.ServerStatus {
	return minecraft.ServerStatus{
		ServerName:  config.MOTD(),
		PlayerCount: PlayerCount(),
		MaxPlayers:  config.MaxPlayers(),
	}
}

// PlayerCount returns the global player count of the proxy.
func PlayerCount() int {
	return int(playerCount.Load())
}

// IncrementPlayerCount increases the global player count by 1.
func IncrementPlayerCount() {
	playerCount.Add(1)
}

// DecrementPlayerCount decreases the global player count by 1.
func DecrementPlayerCount() {
	playerCount.Sub(1)
}
