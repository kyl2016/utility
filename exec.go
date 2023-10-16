package utility

import "os/exec"

func ConvertAudio(src, dst string) error {
	return exec.Command("ffmpeg", []string{"-i", src, dst}...).Run()
}
