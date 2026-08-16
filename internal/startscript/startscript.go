package startscript

import (
	"fmt"
	"os"
	"path/filepath"
)

func Write(dir, memory string) error {
	sh := fmt.Sprintf("#!/bin/sh\ncd \"$(dirname \"$0\")\"\njava -Xmx%s -Xms%s -jar server.jar nogui\n", memory, memory)
	if err := os.WriteFile(filepath.Join(dir, "start.sh"), []byte(sh), 0o755); err != nil {
		return err
	}

	bat := fmt.Sprintf("@echo off\r\ncd /d %%~dp0\r\njava -Xmx%s -Xms%s -jar server.jar nogui\r\npause\r\n", memory, memory)
	return os.WriteFile(filepath.Join(dir, "start.bat"), []byte(bat), 0o644)
}