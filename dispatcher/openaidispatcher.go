package dispatcher

import (
	"github.com/chinaboard/brewing/collection"
	"github.com/chinaboard/brewing/model"
	"github.com/chinaboard/brewing/pkg/ai"
	"strings"
)

type OpenaiDispatcher struct {
	tc collection.Collection
	ac collection.Collection
}

func NewOpenaiDispatcher() (Dispatcher, error) {
	tc, err := collection.NewTaskCollection("brewing")
	if err != nil {
		return nil, err
	}
	ac, err := collection.NewAsrCollection("brewing")
	if err != nil {
		return nil, err
	}
	return &OpenaiDispatcher{tc: tc, ac: ac}, nil
}

func (od *OpenaiDispatcher) Add(taskAny any) error {
	task := taskAny.(*model.AsrResponse)
	return od.ac.Add(task)
}

func (od *OpenaiDispatcher) Run(taskAny any) error {
	task := taskAny.(*model.AsrResponse)
	var errors []string
	parts := ai.Service.SplitContent(task.MakeContent())
	c, e := ai.Service.SummaryParallel(parts)
	for _, err := range e {
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			errors = append(errors, "")
		}

	}
	task.Content = parts
	task.Errors = errors
	var results []string
	for _, cc := range c {
		results = append(results, strings.Split(cc, "\n")...)
	}
	task.Pretty = results
	return od.ac.Update(task)
}

func (od OpenaiDispatcher) Del(id string) error {
	return od.ac.Del(id)
}

func (od OpenaiDispatcher) Get(id string) (any, error) {
	return od.ac.Get(id)
}
