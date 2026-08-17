<img width="2172" height="378" alt="deblock_logo" src="https://github.com/user-attachments/assets/e448a63d-51b3-4e84-93de-326d28382a9f" />

<div align="center">

![Static Badge](https://img.shields.io/badge/Release-1.0.2-green)

</div>

Deblock
=============

A terminal wizard that spins up a Minecraft Java Edition server in a handful of keystrokes. No jar-hunting, no hand-editing `server.properties`, no googling "how do I accept the Minecraft EULA" at 11pm.

## Why this exists

Setting up a Minecraft server the manual way means finding the right download link, remembering which `-Xmx`/`-Xms` flags to pass Java, editing `server.properties` in a text editor, and accepting an EULA file by hand before anything even boots.

None of that is hard, really. It's just tedious enough to get in the way of the part you actually care about: playing with your friends.

Deblock turns all of that into a short back-and-forth in your terminal.

## How it works

```
$ deblock
→ pick a name for the server (Deblock creates a folder for it right here)
→ fetch the latest Minecraft version straight from Mojang's API
   (or type a specific version yourself)
→ set MOTD, difficulty, gamemode, whitelist, port, memory
→ accept Mojang's EULA
→ download the official server.jar
→ generate start.sh / start.bat
→ optionally start the server right there, in the same terminal
```

## What it does

- **Interactive setup wizard** — answer a few prompts, skip the manual file editing
- **Talks to the official APIs** (Mojang) to always resolve the correct, current `server.jar`
- **Self-contained per server** — everything lives in one folder, named after your server
- **Detects existing installs** — re-run Deblock in the same folder to reconfigure, reinstall, or just start it
- **Cross-platform** — one Go binary, no runtime dependencies of its own (you still need Java to run the actual Minecraft server, just not to run Deblock)

## Compatibility

| OS | Running Deblock | Running the Minecraft server |
|---|---|---|
| Linux | ✅ | needs a JDK (21+ recommended) |
| macOS | ✅ | needs a JDK (21+ recommended) |
| Windows | ✅ | needs a JDK (21+ recommended) |

Older Minecraft versions may require an older Java version. If the server won't boot, that's the first thing to check.

Only **Vanilla** shows up in the menu right now. Paper and Fabric support is already implemented and tested under the hood (see [Known limitations](#known-limitations)).

## Installation

### Linux / macOS
```bash
curl -sSL https://raw.githubusercontent.com/yeahgns/deblock/main/install.sh | bash
```

If the `deblock` command isn't recognized after installing, add `~/.local/bin`
to your PATH:

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc   # or ~/.bashrc, depending on your shell
source ~/.zshrc
```

### Windows
Download the `.zip` from the [latest release](https://github.com/yeahgns/deblock/releases/latest),
extract it, and run `deblock.exe`.

## Usage

Once installed, run:

- **Linux / macOS:** `deblock`
- **Windows:** `deblock.exe`

You'll be asked for a server name, a Minecraft version (defaults to the latest release), and the basics of `server.properties` (MOTD, max players, port, difficulty, gamemode, whitelist, online-mode, memory). At the end it asks if you want to start the server right away.

### Managing a server you already set up

Run Deblock again inside that same folder and it'll notice the existing install:

- **Edit the settings** — change `server.properties` without touching the jar
- **Reinstall from scratch** — wipe the jar and download it again
- **Just start it** — skip configuration entirely

Or start it directly without going through the wizard:

```bash
cd my-server-folder
./start.sh       # Linux/Mac
start.bat        # Windows
```

## Known limitations

- Only Vanilla is exposed in the menu right now. Paper and Fabric are fully implemented (`internal/loaders`) but intentionally hidden until they get more real-world testing.
- One server per folder — there's no multi-server dashboard, just run Deblock again in a different folder for a second server.
- Deblock doesn't expose your server to the internet (port forwarding, tunnels, etc.) — that's on you, on purpose, since it depends heavily on your own network setup.
- Each loader talks to a third-party public API (Mojang, PaperMC, Fabric). If one of them changes its response format, that specific loader may need a fix — they're isolated in their own files under `internal/loaders` specifically so that's a small patch, not a rewrite.

## Contributing

Want to work on Deblock itself (project structure, running from source, releasing)?
See [CONTRIBUTING.md](CONTRIBUTING.md).
