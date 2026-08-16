package props

import (
	"path/filepath"
	"testing"
)

func TestWriteThenRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "server.properties")

	original := Config{
		MOTD:       "A cool server",
		MaxPlayers: 12,
		Port:       25566,
		Difficulty: "hard",
		Gamemode:   "creative",
		Whitelist:  true,
		OnlineMode: false,
	}

	if err := Write(path, original); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	read, ok := Read(path)
	if !ok {
		t.Fatalf("Read did not find the newly written file")
	}

	if read != original {
		t.Fatalf("The values ​​read do not match the values ​​written:\n written: %+v\n read:    %+v", original, read)
	}
}

func TestReadMissingFileReturnsDefaults(t *testing.T) {
	cfg, ok := Read("/path/that/does/not/exist/server.properties")
	if ok {
		t.Fatalf("expected ok=false for non-existent file")
	}
	if cfg != Default() {
		t.Fatalf("expected default values ​​when the file does not exist")
	}
}
