package portal

import (
	"github.com/sandertv/gophertunnel/minecraft/resource"
	"os"
	"path/filepath"
)

// Config represents the base configuration for portal. It holds settings that affect different aspects of the
// proxy.
type Config struct {
	// Network holds settings related to network aspects of the proxy.
	Network struct {
		// Address is the address on which the proxy should listen. Players may connect to this address in
		// order to join. It should be in the format of "ip:port".
		Address string `json:"address"`
		// Communication holds settings related to the communication aspects of the proxy.
		Communication struct {
			// Address is the address on which the communication service should listen. External connections
			// can use this address in order to communicate with the proxy. It should be in the format of
			// "ip:port".
			Address string `json:"address"`
			// Secret is the authentication secret required by external connections in order to authenticate
			// to the proxy and start communicating.
			Secret string `json:"secret"`
		} `json:"communication"`
	} `json:"network"`
	// Logger holds settings related to the logging aspects of the proxy.
	Logger struct {
		// File is the path to the file in which logs should be stored. If the path is empty then logs will
		// not be written to a file.
		File string `json:"file"`
		// Level is the required level logs should have to be shown in console or in the file above.
		Level string `json:"level"`
	} `json:"logger"`
	// PlayerLatency holds settings related to the latency reporting aspects of the proxy.
	PlayerLatency struct {
		// Report is if the proxy should send the proxy of a player to their server at a regular interval.
		Report bool `json:"report"`
		// UpdateInterval is the interval to report a player's ping if Report is true.
		UpdateInterval int `json:"update_interval"`
	} `json:"player_latency"`
	// Whitelist holds settings related to the proxy whitelist.
	Whitelist struct {
		// Enabled is if the whitelist is enabled.
		Enabled bool `json:"enabled"`
		// Players is a list of whitelisted players' usernames.
		Players []string `json:"players"`
	} `json:"whitelist"`
	// ResourcePacks holds settings related to sending resource packs to players.
	ResourcePacks struct {
		// Required is if players are required to download the resource packs before connecting.
		Required bool `json:"required"`
		// Directory is the directory to load resource packs from. They can be directories, .zip files or .mcpack files.
		Directory string
	}
}

// DefaultConfig returns a configuration with the default values filled out.
func DefaultConfig() (c Config) {
	c.Network.Address = ":19132"
	c.Network.Communication.Address = ":19131"
	c.Logger.File = "proxy.log"
	c.Logger.Level = "debug"
	c.PlayerLatency.Report = true
	c.PlayerLatency.UpdateInterval = 5
	c.ResourcePacks.Directory = "resource_packs"
	return
}

// LoadResourcePacks attempts to load all the resource packs in the provided directory. If the directory does not exist,
// it will be created. If any pack fails to compile, the error will be returned.
func LoadResourcePacks(dir string) ([]*resource.Pack, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	packs := make([]*resource.Pack, 0, len(files))
	for _, file := range files {
		pack, err := resource.Compile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}
		packs = append(packs, pack)
	}
	return packs, nil
}
