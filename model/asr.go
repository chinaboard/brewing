package model

import (
	"encoding/json"
	"github.com/chinaboard/brewing/pkg/whisper"
)

type AsrReponse struct {
	UniqueId string `json:"uniqueId" bson:"uniqueId"`
	whisper.AsrResp
	Content   []string `json:"content" bson:"content"`
	Errors    []string `json:"errors" bson:"errors"`
	Pretty    []string `json:"pretty" bson:"pretty"`
	BarkToken string   `json:"barkToken" bson:"barkToken"`
	Name      string   `json:"name" bson:"name"`
}

func ConvertToAsrReponse(data string) (*AsrReponse, error) {
	var v AsrReponse
	err := json.Unmarshal([]byte(data), &v)
	return &v, err
}
