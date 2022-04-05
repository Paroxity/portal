![Banner](https://raw.githubusercontent.com/Paroxity/portal/master/banner.png)

A lightweight transfer proxy written in Go for Minecraft: Bedrock Edition.

# Installation
1. Download the latest release for your platform from the [GitHub releases page](https://github.com/Paroxity/portal/releases/)
2. Move it to a directory of your choice, and run from the command line. 

*Note for Linux/macOS users: run `chmod +x` on the binary to make it executable.*
# Configuration
After running portal for the first time, a default configuration file called `config.json` will be created in the same directory as the program.
### Overview of the configuration file
- **query**
    - **maxPlayers:** The maximum number of players allowed to connect to the proxy simultaneously. 0 means no limit.
    - **motd:** A message to display to players on the server connect screen. Supports minecraft formatting codes.
- **whitelist**
    - **enabled:** Whether or to enable the whitelist.
    - **players:** An array of player names to whitelist.
- **proxy**
    - **bindAddress:** The address to bind the proxy to.
    - **groups:** A list of server groups to use. Each group contains a list of servers within their respective group.
    - **defaultGroup:** The default group to route players to when they connect.
    - **authentication:** Whether to authenticate players when they connect. Disabling this will allow users to use any name when connecting
    - **resourceDir:** Directory that resource packs are loaded from.
    - **forceTextures:** Force players to accept the resource pack before joining.
- **playerLatency**
    - **report:** Whether to report player latency to the server.
    - **updateInterval:** The interval at which to update the player latency, in seconds.
- **socket**
    - **bindAddress:** The address to bind the socket server to.
    - **secret:** The secret key for clients to authenticate with. Leaving this blank could allow anyone to register their server to the proxy.
- **logger**
    - **file:** File to output logs to.
    - **debug:** Whether to output debug logs.
