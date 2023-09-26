package whisper

import (
	"fmt"
	"strconv"
)

type AsrResp struct {
	Segments []struct {
		Start float64 `json:"start"`
		End   float64 `json:"end"`
		Text  string  `json:"text"`
	} `json:"segments"`
}

func (r *AsrResp) Decode() {
	for _, segment := range r.Segments {
		if temp, err := strconv.Unquote(segment.Text); err == nil {
			segment.Text = temp
		}
	}
}

func (r *AsrResp) MakeContent() []string {
	var content []string
	for _, segment := range r.Segments {
		content = append(content, segment.Text)
	}
	return content
}

func (r *AsrResp) MakeContentWithTime() []string {
	var content []string
	for _, segment := range r.Segments {
		content = append(content, fmt.Sprintf("<span>%s</span><span class=\"indent\">%s</span>", convertTime(segment.Start), segment.Text))
	}
	return content
}

func convertTime(num float64) string {
	seconds := int(num)
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	seconds = seconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}
