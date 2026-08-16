package startscript

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestWrite(t *testing.T) {
	dir := t.TempDir()
	if err := Write(dir, "3G"); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	sh, err := os.ReadFile(filepath.Join(dir, "start.sh"))
	if err != nil {
		t.Fatalf("start.sh was not created: %v", err)
	}
	if !strings.Contains(string(sh), "-Xmx3G") || !strings.Contains(string(sh), "-Xms3G") {
		t.Fatalf("start.sh does not contain the expected memory parameters:\n%s", sh)
	}

	bat, err := os.ReadFile(filepath.Join(dir, "start.bat"))
	if err != nil {
		t.Fatalf("start.bat was not created: %v", err)
	}
	if !strings.Contains(string(bat), "-Xmx3G") {
		t.Fatalf("start.bat does not contain the expected memory parameters:\n%s", bat)
	}

	if runtime.GOOS != "windows" {
		info, err := os.Stat(filepath.Join(dir, "start.sh"))
		if err != nil {
			t.Fatalf("stat start.sh: %v", err)
		}
		if info.Mode().Perm()&0o111 == 0 {
			t.Fatalf("start.sh should be executable; check permissions: %v", info.Mode())
		}
	}
}
