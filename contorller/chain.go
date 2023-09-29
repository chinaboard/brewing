package contorller

import (
	"errors"
	"fmt"
	"github.com/chinaboard/brewing/model"
	"github.com/chinaboard/brewing/pkg/cfg"
	"github.com/chinaboard/brewing/pkg/notify"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (ch *Chain) chain(c *gin.Context) {
	var task model.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	if !strings.HasPrefix(task.Command[len(task.Command)-1], "http") {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": errors.New("invalid video url").Error(),
		})
		return
	}

	if strings.TrimSpace(task.ImageName) == "" || !strings.Contains(task.ImageName, "brewing") {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": errors.New("invalid image name").Error(),
		})
		return
	}

	go func() {
		err := runner(ch.repo, &task)

		if err != nil {
			notify.Send(task.Name, err.Error(), "brewing", task.BarkToken, "")
		} else {
			msg := fmt.Sprintf("%s/share/%s", strings.TrimSuffix(cfg.ShareDomain, "/"), task.UniqueId)
			notify.Send(task.Name, "success", "brewing", task.BarkToken, msg)
		}
	}()

	taskId := task.UniqueId
	if taskId == "" {
		taskId = task.Hash()
	}
	c.AbortWithStatusJSON(http.StatusOK, gin.H{"uniqueId": taskId})
}

type Chain struct {
	repo *asrRepo
}

func runner(repo *asrRepo, task *model.Task) error {
	if task.UniqueId == "" {
		task.UniqueId = task.Hash()
		task.Status = "Init"
	}

	if task.BarkToken == "" {
		task.BarkToken = cfg.BarkNotifyToken
	}

	v, err := repo.taskDispatcher.Get(task.UniqueId)
	if err == nil && v.(*model.Task).ExitCode == 0 {
		goto jump
	}

	err = repo.taskDispatcher.Add(task)
	if err != nil {
		return err
	}

	err = repo.taskDispatcher.Run(task)
	if err != nil {
		return err
	}

	v, err = repo.taskDispatcher.Get(task.UniqueId)
	if err != nil {
		return err
	}

jump:
	task = v.(*model.Task)
	asr, err := model.ConvertToAsrReponse(task.Stdout)
	if err != nil {
		return err
	}

	asr.Name = task.Name
	asr.UniqueId = task.UniqueId
	asr.BarkToken = task.BarkToken
	asr.Decode()

	err = repo.openaiDispatcher.Add(asr)
	if err != nil {
		return err
	}

	err = repo.openaiDispatcher.Run(asr)
	return err
}
