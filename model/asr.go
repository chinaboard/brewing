package model

import (
	"encoding/json"
	"github.com/chinaboard/brewing/pkg/whisper"
	"time"
)

type AsrResponse struct {
	UniqueId string `json:"uniqueId" bson:"uniqueId"`
	whisper.AsrResp
	Content   []string  `json:"content" bson:"content"`
	Errors    []string  `json:"errors" bson:"errors"`
	Pretty    []string  `json:"pretty" bson:"pretty"`
	BarkToken string    `json:"barkToken" bson:"barkToken"`
	Name      string    `json:"name" bson:"name"`
	CreateAt  time.Time `json:"createAt" bson:"createAt"`
	UpdateAt  time.Time `bson:"updateAt" bson:"updateAt"`
}

func ConvertToAsrResponse(data string) (*AsrResponse, error) {
	var v AsrResponse
	err := json.Unmarshal([]byte(data), &v)
	if err != nil {
		return nil, err
	}
	v.CreateAt = time.Now()
	v.UpdateAt = v.CreateAt
	return &v, err
}
