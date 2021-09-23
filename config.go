package portal

// Config represents the base configuration for portal. It holds settings that affect different aspects of the
// proxy.
type Config struct {
	// Network holds settings related to network aspects of the proxy.
	Network struct {
		// Address is the address on which the proxy should listen. Players may connect to this address in
		// order to join. It should be in the format of "ip:port".
		Address string
		// Communication holds settings related to the communication aspects of the proxy.
		Communication struct {
			// Address is the address on which the communication service should listen. External connections
			// can use this address in order to communicate with the proxy. It should be in the format of
			// "ip:port".
			Address string
			// Secret is the authentication secret required by external connections in order to authenticate
			// to the proxy and start communicating.
			Secret string
		}
	}
	// Logger holds settings related to the logging aspects of the proxy.
	Logger struct {
		// File is the path to the file in which logs should be stored. If the path is empty then logs will
		// not be written to a file.
		File string
		// Level is the required level logs should have to be shown in console or in the file above.
		Level string
	}
}

// DefaultConfig returns a configuration with the default values filled out.
func DefaultConfig() (c Config) {
	c.Network.Address = ":19132"
	c.Network.Communication.Address = ":19131"
	c.Logger.File = "proxy.log"
	c.Logger.Level = "debug"
	return
}
