package config

var (
	maxPlayers int
	motd       string

	whitelist   bool
	whitelisted []string
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
