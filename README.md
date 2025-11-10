# BedrockWormhole

UDP proxy for Minecraft Bedrock Edition. Forwards traffic from your local network to a remote server, allowing consoles (Xbox, PlayStation, Switch) and mobile devices to connect to remote servers as LAN games.

**The computer you run this on needs to be on the same network as your Xbox/etc. to function!**

May pop positive for some antivirus due to this establishing a UDP proxy, which is a similar behavior to some reverse shells.

## How Does It Work?

This is a transparent UDP proxy, using the computer you run it on to pose as the host of a Minecraft Bedrock server, when in fact it's just forwarding all traffic to the remote host. You can run the same thing yourself with something like `socat UDP4-LISTEN:19132,fork UDP4-SENDTO:example.com:19132`.

```
┌───────────────────────────────────────────────────────────────────┐
│                       Your Home Network                           │
│                                                                   │
│                                                                   │
│              "Yeah this is totally a    ┌───────────────────────┐ │
│               local Bedrock server;     │    Your Computer      │ │
│ ┌──────────┐  connect to me!"           │ ┌───────────────────┐ │ │
│ │ Xbox   ◄─┼────────────────────────────┼─┼─►                 │ │ │
│ └──────────┘                            │ │ Bedrock Wormhole  │ │ │
│ ┌──────────┐                            │ │                   │ │ │
│ │ Switch ◄─┼────────────────────────────┼─┼─►       ▲         │ │ │
│ └──────────┘                            │ └─────────┼─────────┘ │ │
│                                         └───────────┼───────────┘ │
└─────────────────────────────────────────────────────┼─────────────┘
                                                      │
                                                      ▼
                                              ┌─────────────────┐
                                              │ Some Internet   │
                                              │ Bedrock Server  │
                                              └─────────────────┘
```

## Usage

### Interactive Mode

Download from the __[Releases](https://github.com/jkingsman/BedrockWormhole/releases)__ tab on the right and download the appropriate architecture. Unpack and run the executable (interactive or via CLI).

Follow the prompts to enter remote host and port.

### Command-Line Mode

```bash
./BedrockWormhole -host <hostname> -port <port>
```

**Options:**
- `-host` - Remote server hostname or IP (required)
- `-port` - Remote server port (default: 19132)

**Example:**

```bash
./BedrockWormhole -host play.example.com -port 19132
```

## Development

Pretty vanilla Go. Build for your arch with

```bash
go build -o BedrockWormhole
```

or for all platforms with

```bash
./build.sh
```

Outputs to `build/` directory.
