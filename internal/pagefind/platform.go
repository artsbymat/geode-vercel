package pagefind

import (
	"fmt"
	"runtime"
)

func platformKey() (string, error) {
	switch runtime.GOOS {
	case "linux":
		if runtime.GOARCH == "amd64" {
			return "x86_64-unknown-linux-musl", nil
		}
	case "darwin": // NOT TESTED
		if runtime.GOARCH == "amd64" {
			return "x86_64-apple-darwin", nil
		}
		if runtime.GOARCH == "arm64" {
			return "aarch64-apple-darwin", nil
		}
	case "windows": // NOT TESTED
		if runtime.GOARCH == "amd64" {
			return "x86_64-pc-windows-msvc", nil
		}
	}
	return "", fmt.Errorf("unsupported platform %s/%s", runtime.GOOS, runtime.GOARCH)
}
