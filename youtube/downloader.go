package youtube

import (
	"os/exec"
)

func DownloadVideo(videoURL, outputPath string) error {
	cmd := exec.Command("yt-dlp", "-o", outputPath, videoURL)
	return cmd.Run()
}
