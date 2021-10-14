package session

import (
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Whitelist handles the players joining the proxy to decide which are allowed join.
type Whitelist interface {
	// Authorize returns whether a player with the given connection is allowed to join the proxy and a message
	// to display to the player on their disconnection screen.
	Authorize(conn *minecraft.Conn) (bool, string)
}

// SimpleWhitelist is a whitelist that, if enabled, only allows a set list of players to join.
type SimpleWhitelist struct {
	enabled bool
	players []string
}

// NewSimpleWhitelist returns a simple whitelist from the enabled status and a player list passed.
func NewSimpleWhitelist(enabled bool, players []string) *SimpleWhitelist {
	return &SimpleWhitelist{enabled, players}
}

// Authorize ...
func (s *SimpleWhitelist) Authorize(conn *minecraft.Conn) (bool, string) {
	if !s.enabled {
		return true, ""
	}
	u := conn.IdentityData().DisplayName
	for _, p := range s.players {
		if u == p {
			return true, ""
		}
	}
	return false, text.Colourf("<red>Server is whitelisted</red>")
}
