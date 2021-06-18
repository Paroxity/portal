package config

import (
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"github.com/sirupsen/logrus"
	"time"
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

	reportPlayerLatency         bool
	playerLatencyUpdateInterval time.Duration

	socketAddress string
	socketSecret  string

	logFile  string
	logLevel logrus.Level
)

// MaxPlayers returns the maximum amount of players allowed on the proxy. If this value is 0 then there will
// be no limit.
func MaxPlayers() int {
	return maxPlayers
}

// MOTD returns the "Message Of The Day" shown to clients in the server list.
func MOTD() string {
	return motd
}

// Whitelist returns if the proxy is whitelisted or not. If true, only the Whitelisted() players can join.
func Whitelist() bool {
	return whitelist
}

// Whitelisted returns a slice of usernames that are whitelisted on the proxy. They are the only ones who can
// join when Whitelist() is true.
func Whitelisted() []string {
	return whitelisted
}

// Authentication returns if the proxy should require Xbox Live authentication or not.
func Authentication() bool {
	return authentication
}

// BindAddress returns the address that the proxy should bind to. This should contain the IP address and the
// port, separated by a colon. E.g. "0.0.0.0:19132".
func BindAddress() string {
	return bindAddress
}

// ResourcePacks returns a slice of resource packs to send to a client when they connect.
func ResourcePacks() []*resource.Pack {
	return resourcePacks
}

// ForceTexturePacks returns if texture packs should be required to be downloaded before the player can join
// the server.
func ForceTexturePacks() bool {
	return forceTextures
}

// ReportPlayerLatency returns if the proxy should send a player's latency to their connected server. This
// can be disabled if not needed to save bandwidth over the network.
func ReportPlayerLatency() bool {
	return reportPlayerLatency
}

// PlayerLatencyUpdateInterval is how often to update the player's latency to their connected server.
func PlayerLatencyUpdateInterval() time.Duration {
	return playerLatencyUpdateInterval
}

// SocketAddress returns the address that the socket server should bind to. This should contain the IP
// address and the port, separated by a colon. E.g. "0.0.0.0:19131".
func SocketAddress() string {
	return socketAddress
}

// SocketSecret returns the secret key used for authenticating clients over the socket connection. The
// connecting clients must match this key completely to be able to authenticate.
func SocketSecret() string {
	return socketSecret
}

// LogFile returns the path of the file to log to. If the file does not exist, it will be created when the
// proxy starts.
func LogFile() string {
	return logFile
}

// LogLevel returns the log level to be used for Logrus. This is only to decide if debugging logs should be
// printed to console and the LogFile() or not.
func LogLevel() logrus.Level {
	return logLevel
}
