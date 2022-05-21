![Banner](https://raw.githubusercontent.com/Paroxity/portal/master/banner.png)

A lightweight transfer proxy written in Go for Minecraft: Bedrock Edition.

# Installation

1. Download the latest release for your platform from
   the [GitHub releases page](https://github.com/Paroxity/portal/releases/)
2. Move it to a directory of your choice, and run from the command line.

*Note for Linux/macOS users: run `chmod +x` on the binary to make it executable.*

# Configuration

After running portal for the first time, a default configuration file called `config.json` will be created in the same
directory as the program.

### Overview of the configuration file

- **network**
    - **address**: The address on which the proxy should listen. Players may connect to this address in order to join.
      It should be in the format of "ip:port"
    - **communication**
        - **address**: Address is the address on which the communication service should listen. External connections can
          use this address in order to communicate with the proxy. It should be in the format of "ip:port"
        - **secret**: Secret is the authentication secret required by external connections in order to authenticate to
          the proxy and start communicating
- **logger**
    - **file**: File is the path to the file in which logs should be stored. If the path is empty then logs will not be
      written to a file
    - **level**: Level is the required level logs should have to be shown in console or in the file above
- **player_latency**
    - **report**: Determines if the proxy should send the proxy of a player to their server at a regular interval
    - **update_interval**: The interval to report a player's ping if report is true
- **whitelist**
    - **enabled**: Determines if the whitelist is enabled
    - **players**: A list of whitelisted players' usernames
- **resource_packs**
    - **required**: Determines if players are required to download the resource packs before connecting
    - **directory**: The directory to load resource packs from. They can be directories, .zip files or .mcpack files
    - **encryption_keys**: A map of resource pack UUIDs to their encryption key
