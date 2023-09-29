package controller

import (
	"errors"
	"github.com/chinaboard/brewing/dispatcher"
	"github.com/chinaboard/brewing/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type taskRepo struct {
	td dispatcher.Dispatcher
}

func newTaskRepo() (*taskRepo, error) {
	result, err := dispatcher.NewTaskDispatcher()
	if err != nil {
		return nil, err
	}
	return &taskRepo{td: result}, nil
}

func (tr *taskRepo) Add(c *gin.Context) {
	var task model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	if task.UniqueId == "" {
		task.UniqueId = task.Hash()
		task.Status = "Init"
	}

	if strings.TrimSpace(task.ImageName) == "" || !strings.HasPrefix(task.ImageName, "brewing") {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": errors.New("invalid image name").Error(),
		})
		return
	}

	err := tr.td.Add(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "success", "data": task})
}

func (tr *taskRepo) Run(c *gin.Context) {
	taskId := c.Param("id")
	if taskId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	r, err := tr.td.Get(taskId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"msg": err.Error()})
	}

	go tr.td.Run(r)

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "success"})
}

func (tr *taskRepo) Get(c *gin.Context) {
	taskId := c.Param("id")
	if taskId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r, err := tr.td.Get(taskId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"msg": err.Error()})
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "success", "data": r})
}
