package ffmpeg

import (
	"encoding/json"
	"math"
	"os/exec"
)

func GetVideoAspectRatio(filepath string) (string, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v",
		"error",
		"-print_format",
		"json",
		"-show_streams",
		filepath,
	)

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var streams struct {
		Streams []struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(output, &streams); err != nil {
		return "", err
	}

	q := float64(streams.Streams[0].Width) / float64(streams.Streams[0].Height)
	tolerance := 5.0
	if withinPercent(q, 1.778, tolerance) {
		return "16:9", nil
	} else if withinPercent(q, 0.5625, tolerance) {
		return "9:16", nil
	}

	return "other", nil
}

func withinPercent(val, target, p float64) bool {
	tolerance := math.Abs(target) * p / 100.0
	return math.Abs(val-target) <= tolerance
}
