package dispatcher

import (
	"bytes"
	"context"
	"github.com/chinaboard/brewing/collection"
	"github.com/chinaboard/brewing/model"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

type Dispatcher interface {
	Add(any) error
	Get(string) (any, error)
	Del(string) error
	Run(any) error
}

type TaskDispatcher struct {
	tc  collection.Collection
	cli *client.Client
}

func NewTaskDispatcher() (Dispatcher, error) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	tc, err := collection.NewTaskCollection("brewing")
	if err != nil {
		return nil, err
	}
	return &TaskDispatcher{cli: client, tc: tc}, nil
}

func (dd *TaskDispatcher) Add(taskAny any) error {
	task := taskAny.(*model.Task)
	task.ExitCode = -1
	return dd.tc.Add(task)
}

func (dd *TaskDispatcher) Run(taskAny any) error {
	task := taskAny.(*model.Task)
	ctx := context.Background()
	imageName := task.ImageName
	_, _, err := dd.cli.ImageInspectWithRaw(ctx, imageName)
	if task.ForcePull || client.IsErrNotFound(err) {
		task.Status = "ImagePull"
		if err = dd.tc.Update(task); err != nil {
			return err
		}

		logrus.Debugln("Image start pull", imageName)
		reader, err := dd.cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
		if err != nil {
			logrus.Errorln(err)
			return err
		}

		io.Copy(os.Stdout, reader)
		defer reader.Close()
		logrus.Debugln("Image pulled", imageName)
	} else if err != nil {
		logrus.Errorln(err)
		return err
	}

	task.Status = "ContainerCreate"
	logrus.Debugln(task.UniqueId, "ContainerCreate")
	if err = dd.tc.Update(task); err != nil {
		return err
	}
	resp, err := dd.cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Cmd:   task.Command,
		Env:   task.Env,
	}, nil, nil, nil, task.UniqueId)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	task.ContainerId = resp.ID
	if err = dd.tc.Update(task); err != nil {
		return err
	}

	task.Status = "ContainerStart"
	logrus.Debugln(task.UniqueId, resp.ID, "ContainerStart")
	if err = dd.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	task.Status = "ContainerWait"
	logrus.Debugln(task.UniqueId, resp.ID, "ContainerWait")
	if err = dd.tc.Update(task); err != nil {
		return err
	}

	statusCh, errCh := dd.cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}

	out, err := dd.cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return err
	}

	var stdout, stderr bytes.Buffer
	_, err = stdcopy.StdCopy(&stdout, &stderr, out)
	if err != nil {
		logrus.Errorln(err)
		return err
	}

	task.Stdout = stdout.String()
	task.Stderr = stderr.String()

	inspection, err := dd.cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return err
	}

	task.Status = inspection.State.Status
	task.ExitCode = inspection.State.ExitCode

	logrus.Debugln(task.UniqueId, resp.ID, "ContainerInspect")
	if err = dd.tc.Update(task); err != nil {
		return err
	}

	logrus.Debugln("Container", inspection.State.Status, "code", inspection.State.ExitCode)
	if task.AutoRemove && inspection.State.ExitCode == 0 {
		defer func() {
			if err = dd.cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			}); err != nil {
				logrus.Fatalln(err)
			}
			logrus.Debugln(task.UniqueId, resp.ID, "Container removed")
		}()
	}
	return nil
}

func (dd TaskDispatcher) Del(id string) error {
	return dd.tc.Del(id)
}

func (dd TaskDispatcher) Get(id string) (any, error) {
	return dd.tc.Get(id)
}
