![Banner](https://raw.githubusercontent.com/Paroxity/portal/master/banner.png)

A lightweight transfer proxy written in Go for Minecraft: Bedrock Edition.

# Installation
1. Download the latest release for your platform from the [GitHub releases page](https://github.com/Paroxity/portal/releases/)
2. Move it to a directory of your choice, and run from the command line. 

*Note for Linux/macOS users: run `chmod +x` on the binary to make it executable.*
# Configuration
After running portal for the first time, a default configuration file called `config.json` will be created in the same directory as the program.
### Default config.json with comments
```json
{
    "query": {
        "maxPlayers": 0, // The maximum number of players allowed to connect to the proxy simoultaneously. 0 means no limit.
        "motd": "Portal" // A message to display to players on the server connect screen. Supports minecraft formatting codes.
    },
    "whitelist": {
        "enabled": false, // Whether or not to enable the whitelist.
        "players": null // An array of player names to whitelist.
    },
    "proxy": {
        "bindAddress": "0.0.0.0:19132", // The address to bind the proxy to.
        "groups": { // A list of server groups. When a player joins a group, portal automatically distributes them between servers within that group.
            "Hub": { // Each group can have any name.
                "Hub1": "127.0.0.1:19133" // A list of servers within the group.
            }
        },
        "defaultGroup": "Hub", // The default group to route players to when they connect.
        "authentication": true, // Whether or not to require players to use Xbox Live to authenticate.
        "resourcesDir": "", // Directory for resource packs.
        "forceTextures": false // Force players to accept the resource pack before joining.
    },
    "playerLatency": {
        "report": true, // Whether or not to report player latency to the server.
        "updateInterval": 5 // The interval in seconds to update the player latency.
    },
    "socket": {
        "bindAddress": "127.0.0.1:19131", // The address to bind the socket server to.
        "secret": "" // The secret key for clients to authenticate with. Leaving this blank could allow 
                     // anyone to register their server with the proxy.
    },
    "logger": {
        "file": "proxy.log", // File to output logs to.
        "debug": false // Whether or not to output debug logs.
    }
}
```