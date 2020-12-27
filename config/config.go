package config

import "github.com/sandertv/gophertunnel/minecraft/resource"

var (
	maxPlayers int
	motd       string

	whitelist   bool
	whitelisted []string

	authentication bool
	resourcePacks  []*resource.Pack
	forceTextures  bool
)

func MaxPlayers() int {
	return maxPlayers
}

func MOTD() string {
	return motd
}

func Whitelist() bool {
	return whitelist
}

func Whitelisted() []string {
	return whitelisted
}

func Authentication() bool {
	return authentication
}

func ResourcePacks() []*resource.Pack {
	return resourcePacks
}

func ForceTexturePacks() bool {
	return forceTextures
}
