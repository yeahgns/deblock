<img width="2172" height="378" alt="deblock_logo" src="https://github.com/user-attachments/assets/e448a63d-51b3-4e84-93de-326d28382a9f" />

<div align="center">

![Static Badge](https://img.shields.io/badge/Release-1.1.0-green)

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
→ optionally expose the server to the internet via playit.gg (no account? Deblock walks you through it)
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
- **Optional internet exposure via playit.gg** — no port forwarding, no router changes. Deblock automates linking your playit.gg account (opens your browser once, no token to copy/paste); the last step, creating the actual tunnel, is currently manual on playit.gg's end (see [Known limitations](#known-limitations))

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

You'll be asked for a server name, a Minecraft version (defaults to the latest release), and the basics of `server.properties` (MOTD, max players, port, difficulty, gamemode, whitelist, online-mode, memory). Right after you accept Mojang's EULA, it asks if you want to expose the server to the internet — and, at the end, if you want to start the server right away.

### Exposing your server to the internet (optional)

Right after accepting the EULA, Deblock offers to set up a tunnel through [playit.gg](https://playit.gg), a free tunneling service built for game servers — no port forwarding or router configuration needed. This requires a free playit.gg account.

If you say yes:

1. **Deblock checks whether the [playit agent](https://github.com/playit-cloud/playit-agent) is installed.** If it isn't, it prints install instructions for your OS and stops there — just run Deblock again after installing it.
2. **Deblock links the agent to your playit.gg account automatically.** It opens your browser to a one-time approval page — no token to copy or paste, just click "Continue" → "Add Agent".
3. **Deblock prints the couple of steps still left to do by hand** (pointing the agent at the secret it just got, and starting it), then opens the playit.gg dashboard for you to create the tunnel itself (pointing it at your server's port, protocol TCP, type "Minecraft Java").

That last step — creating the tunnel — isn't automatable yet because playit.gg doesn't currently expose a public API for it (see [Known limitations](#known-limitations)). Everything before it is fully automatic.

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
- The optional playit.gg tunnel isn't fully automatic yet: linking the agent to your account is (no token copy/pasting needed), but creating the actual tunnel still has to be done once, by hand, on the playit.gg dashboard (which Deblock opens for you). This is a limitation of playit.gg's current public API, not a design choice on Deblock's side — we'll automate this last step too as soon as they expose a way to do it.
- Deblock does not install or manage the playit agent itself (no bundled binary, no background service) — it only automates the account-linking step and tells you what to run.
- Each loader talks to a third-party public API (Mojang, PaperMC, Fabric). If one of them changes its response format, that specific loader may need a fix — they're isolated in their own files under `internal/loaders` specifically so that's a small patch, not a rewrite.

## Contributing

Want to work on Deblock itself (project structure, running from source, releasing)?
See [CONTRIBUTING.md](CONTRIBUTING.md).
