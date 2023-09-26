package collection

import (
	"fmt"
	"github.com/chinaboard/brewing/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabaseRepo_Add(t *testing.T) {
	dd, _ := NewTaskCollection("brewing")
	dd.Add("59438", &model.Task{
		Name:     "hahahahha",
		UniqueId: "59438",
	})
}

func TestDatabaseRepo_Del(t *testing.T) {
	dd, _ := NewTaskCollection("brewing")
	err := dd.Del("59438")
	fmt.Println(err)
	assert.Nil(t, err)
	_, err = dd.Get("59438")
	assert.NotNil(t, err)
}

func TestDatabaseRepo_Get(t *testing.T) {
	dd, _ := NewTaskCollection("brewing")
	data, err := dd.Get("59438")
	fmt.Println(err)

	assert.Nil(t, err)
	assert.Equal(t, "59438", data.UniqueId)
}

func TestDatabaseRepo_List(t *testing.T) {
	dd, _ := NewTaskCollection("brewing")
	data, err := dd.List()

	assert.Nil(t, err)
	assert.Equal(t, "59438", data[0].UniqueId)
}

func TestDatabaseRepo_Update(t *testing.T) {
	dd, _ := NewTaskCollection("brewing")
	err := dd.Update("59438", &model.Task{
		Name:     "6666",
		UniqueId: "59438",
	})
	assert.Nil(t, err)
	data, _ := dd.Get("59438")
	assert.Equal(t, "6666", data.Name)
}
