# Contributing to Deblock

This file is only relevant if you want to work on Deblock's own source code.
If you just want to use Deblock to set up a Minecraft server, check the
[Installation instructions in the README](README.md#installation) instead.

## Project structure

```
deblock/
├── README.md
├── README.pt-br.md
├── README.es.md
├── LICENSE
├── go.mod
├── go.sum
├── main.go
├── main_test.go
└── internal/
    ├── loaders/       # talks to the Mojang/PaperMC/Fabric APIs
    ├── download/       # downloads the server.jar with a progress bar
    ├── props/          # reads and writes server.properties
    └── startscript/    # generates start.sh / start.bat
```

## Running from source

```bash
go run .          # run the wizard from source
go test ./...     # run the test suite
```

## Releasing

Releases are built and published automatically by GitHub Actions
(`.github/workflows/release.yml`) using [GoReleaser](https://goreleaser.com/),
triggered whenever a new `v*` tag is pushed:

```bash
git tag vX.Y.Z
git push origin vX.Y.Z
```
