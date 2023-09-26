package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/chinaboard/brewing/dispatcher"
	"github.com/chinaboard/brewing/model"
	"github.com/chinaboard/brewing/pkg/bininfo"
	"github.com/chinaboard/brewing/pkg/cfg"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	videoUrl := ""
	version := false
	env := ""
	flag.StringVar(&cfg.WhisperEndpoint, "whisperEndpoint", "", "e.g whisper:9000")
	flag.StringVar(&cfg.WhisperEndpointSchema, "whisperEndpointSchema", "", "e.g http or https")
	flag.StringVar(&cfg.OpenAiBaseURL, "openApiUrl", "", "e.g https://api.openai.com/v1")
	flag.StringVar(&cfg.OpenAiToken, "openAiToken", "", "e.g sk-5Cpzm6j3lB8LxVqKG1UWs5FkN8HrSCF6x3WJq1ECsGmklx")
	flag.StringVar(&cfg.OpenAiProxy, "openAiProxy", "", "e.g http://proxy:59438")
	flag.StringVar(&videoUrl, "videoUrl", "", "e.g https://www.bilibili.com/video/Asdf45678")
	flag.StringVar(&env, "env", "", "e.g [\"env1=foo\", \"env2=bar\"]")
	flag.BoolVar(&version, "v", false, "show version")
	flag.Parse()

	if version {
		fmt.Println(bininfo.StringifySingleLine())
		os.Exit(0)
	}

	if cfg.WhisperEndpointSchema == "" {
		fmt.Println("Error: WhisperEndpointSchema must be set")
		os.Exit(1)
	}
	if cfg.WhisperEndpoint == "" {
		fmt.Println("Error: WhisperEndpoint must be set")
		os.Exit(1)
	}
	if cfg.OpenAiToken == "" {
		fmt.Println("Error: OpenAiToken must be set")
		os.Exit(1)
	}
	environ := []string{}
	if env != "" {
		json.Unmarshal([]byte(env), &environ)
		logrus.Debugln("env", environ)
	}
	brewingTask := &model.Task{
		ForcePull: false,
		Env:       environ,
		Command: []string{
			"brewing-runner",
			"-whisperEndpoint",
			fmt.Sprintf("%s://%s", cfg.WhisperEndpointSchema, cfg.WhisperEndpoint),
			"-videoUrl",
			videoUrl,
		},
	}

	d, err := dispatcher.NewTaskDispatcher()
	if err != nil {
		logrus.Fatalln(err)
	}

	if err = d.Add(brewingTask); err != nil {
		logrus.Fatalln(err)
	}
	if err = d.Run(brewingTask); err != nil {
		logrus.Fatalln(err)
	}

	logrus.Println("Done")
}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		ForceQuote:      false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
