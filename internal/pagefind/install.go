package pagefind

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Ensure(ctx context.Context, dir string) (string, error) {
	bin := filepath.Join(dir, exeName())
	if _, err := os.Stat(bin); err == nil {
		return bin, nil
	}

	url, isZip, err := downloadURL()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}

	if isZip {
		return extractZip(resp.Body, bin)
	}
	return extractTarGz(resp.Body, bin)
}

func extractTarGz(r io.Reader, bin string) (string, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if filepath.Base(h.Name) == "pagefind" {
			return writeBinary(tr, bin)
		}
	}
	return "", fmt.Errorf("pagefind binary not found")
}

func extractZip(r io.Reader, bin string) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	zr, err := zip.NewReader(bytesReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	for _, f := range zr.File {
		if filepath.Base(f.Name) == "pagefind.exe" {
			rc, _ := f.Open()
			defer rc.Close()
			return writeBinary(rc, bin)
		}
	}
	return "", fmt.Errorf("pagefind.exe not found")
}

func writeBinary(r io.Reader, bin string) (string, error) {
	out, err := os.Create(bin)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, r); err != nil {
		return "", err
	}

	_ = os.Chmod(bin, 0o755)
	return bin, nil
}

func exeName() string {
	if os.PathSeparator == '\\' {
		return "pagefind.exe"
	}
	return "pagefind"
}

type bytesReader []byte

func (b bytesReader) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(b)) {
		return 0, io.EOF
	}
	n := copy(p, b[off:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}
