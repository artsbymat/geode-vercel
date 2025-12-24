package pagefind

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func Run(ctx context.Context, sourceDir string) error {
	binDir := "bin"
	binPath, err := Ensure(ctx, binDir)
	if err != nil {
		return fmt.Errorf("failed to ensure pagefind binary: %w", err)
	}

	cmd := exec.CommandContext(ctx, binPath, "--site", sourceDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run pagefind: %w", err)
	}

	return nil
}
