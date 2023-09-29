package whisper

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func Asr(whisperEndpoint, filePath, language string) (*AsrResp, error) {
	client := resty.New()
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetFiles(map[string]string{"audio_file": filePath}).
		SetQueryString(fmt.Sprintf("task=transcribe&language=%s&encode=false&output=json", language)).
		Post(whisperEndpoint + "/asr")

	if err != nil {
		return nil, err
	}

	var result AsrResp
	err = json.Unmarshal(resp.Body(), &result)
	result.Decode()

	return &result, err
}
