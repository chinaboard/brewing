package videodl

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func Download(videoUrl, savePath string) error {
	cmd := exec.Command("youtube-dl", "-f", "ba", "-x", "--audio-format", "wav", videoUrl, "-o", savePath)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %s", err, strings.TrimSpace(stderr.String()))
	}
	return nil
}
