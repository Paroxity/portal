package portal

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"go.uber.org/atomic"
)

// MOTDStatusProvider represents a status provider that shows a custom MOTD which can be changed at any time.
type MOTDStatusProvider struct {
	motd atomic.String
}

// NewMOTDStatusProvider creates a new server status provider which shows a custom message in the server list.
func NewMOTDStatusProvider(motd string) *MOTDStatusProvider {
	p := &MOTDStatusProvider{}
	p.motd.Store(motd)
	return p
}

// MOTD sets the MOTD for the current provider.
func (p *MOTDStatusProvider) MOTD(v string) {
	p.motd.Store(v)
}

// ServerStatus ...
func (p *MOTDStatusProvider) ServerStatus(playerCount, maxPlayers int) minecraft.ServerStatus {
	return minecraft.ServerStatus{
		ServerName:  p.motd.Load(),
		PlayerCount: playerCount,
		MaxPlayers:  maxPlayers,
	}
}
