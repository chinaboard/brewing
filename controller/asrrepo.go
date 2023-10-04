package controller

import (
	"github.com/chinaboard/brewing/dispatcher"
	"github.com/chinaboard/brewing/model"
	"github.com/chinaboard/brewing/pkg/notify"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type asrRepo struct {
	taskDispatcher   dispatcher.Dispatcher
	openaiDispatcher dispatcher.Dispatcher
}

func newAsrRepo() (*asrRepo, error) {
	td, err := dispatcher.NewTaskDispatcher()
	if err != nil {
		return nil, err
	}
	od, err := dispatcher.NewOpenaiDispatcher()
	if err != nil {
		return nil, err
	}
	return &asrRepo{taskDispatcher: td, openaiDispatcher: od}, nil
}

func (j *asrRepo) Add(c *gin.Context) {
	taskId := c.Param("id")
	if taskId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	t, err := j.taskDispatcher.Get(taskId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": err.Error(),
		})
		return
	}
	asr, err := model.ConvertToAsrResponse(t.(*model.Task).Stdout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	asr.UniqueId = taskId
	asr.Decode()
	err = j.openaiDispatcher.Add(asr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "success", "data": asr})
}

func (j *asrRepo) Run(c *gin.Context) {
	taskId := c.Param("id")
	if taskId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	r, err := j.openaiDispatcher.Get(taskId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"msg": err.Error()})
	}

	go func() {
		msg := "success"
		e := j.openaiDispatcher.Run(r)
		if e != nil {
			msg = e.Error()
		}
		asr := r.(*model.AsrResponse)
		notify.Send(taskId, msg, "brewing", asr.BarkToken, "")
	}()

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "success"})
}

func (j *asrRepo) Get(c *gin.Context) {
	taskId := c.Param("id")
	if taskId == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	r, err := j.openaiDispatcher.Get(taskId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"msg": err.Error()})
	}

	c.AbortWithStatusJSON(http.StatusOK, gin.H{"msg": "success", "data": r})
}

func (j *asrRepo) Summary(c *gin.Context) {
	taskId := c.Param("id")
	if taskId == "" {
		c.Redirect(302, "https://mp.weixin.qq.com/mp/profile_ext?action=home&__biz=MzkyMDUwNDI5NQ==&scene=124#wechat_redirect")
		return
	}
	r, err := j.openaiDispatcher.Get(taskId)
	if err != nil {
		c.Redirect(302, "https://mp.weixin.qq.com/mp/profile_ext?action=home&__biz=MzkyMDUwNDI5NQ==&scene=124#wechat_redirect")
	}
	d, _ := j.taskDispatcher.Get(taskId)
	task := d.(*model.Task)
	url := ""
	for _, s := range task.Command {
		if strings.Contains(s, "http") {
			url = s
		}
	}
	asr := r.(*model.AsrResponse)
	content := asr.Pretty

	if c.Query("raw") != "" {
		content = asr.MakeContentWithTime()
	}
	c.HTML(http.StatusOK, "share.html",
		gin.H{
			"Name":    task.Name,
			"Content": content,
			"TaskId":  taskId,
			"Url":     url,
			"Raw":     c.Query("raw") != "",
		},
	)

}
