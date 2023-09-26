package notify

import (
	"github.com/chinaboard/brewing/pkg/cfg"
	"github.com/jzksnsjswkw/go-bark"
	"strings"
)

var (
	client *bark.Client
)

func init() {
	client = bark.New(cfg.BarkNotifyDomain)
}

func Send(title, msg, group, token, url string) error {
	tks := strings.Split(token, ",")
	var err error
	for _, tk := range tks {
		if e := client.Push(&bark.Options{
			URL:   url,
			Title: title,
			Msg:   msg,
			Group: group,
			Token: tk,
		}); e != nil {
			err = e
		}
	}
	return err
}
