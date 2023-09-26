package model

import (
	"crypto/sha1"
	"fmt"
)

type Task struct {
	Name string `json:"name" bson:"name" binding:"required"`

	UniqueId string `json:"uniqueId" bson:"uniqueId"`

	Command    []string `json:"command" bson:"command"`
	AutoRemove bool     `json:"autoRemove" bson:"autoRemove"`

	ForcePull bool     `json:"forcePull" bson:"forcePull"`
	ImageName string   `json:"imageName" bson:"imageName" binding:"required"`
	Env       []string `json:"env" bson:"env"`

	Stdout string `json:"stdout" bson:"stdout"`
	Stderr string `json:"stderr" bson:"stderr"`

	Status   string `json:"status" bson:"status"`
	ExitCode int    `json:"exitCode" bson:"exitCode"`

	ContainerId string `json:"containerId" bson:"containerId"`

	Language  string `json:"language" bson:"language" binding:"required"`
	BarkToken string `json:"barkToken" bson:"barkToken"`
}

func (b *Task) Hash() string {
	h := sha1.New()
	str := fmt.Sprint(b.Name, b.Command, b.AutoRemove, b.ForcePull, b.Env, b.ImageName)
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
