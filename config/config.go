package config

import (
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"github.com/sirupsen/logrus"
)

var (
	maxPlayers int
	motd       string

	whitelist   bool
	whitelisted []string

	authentication bool
	bindAddress    string
	resourcePacks  []*resource.Pack
	forceTextures  bool

	socketAddress string
	socketSecret  string

	logFile  string
	logLevel logrus.Level
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

func BindAddress() string {
	return bindAddress
}

func ResourcePacks() []*resource.Pack {
	return resourcePacks
}

func ForceTexturePacks() bool {
	return forceTextures
}

func SocketAddress() string {
	return socketAddress
}

func SocketSecret() string {
	return socketSecret
}

func LogFile() string {
	return logFile
}

func LogLevel() logrus.Level {
	return logLevel
}
