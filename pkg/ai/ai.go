package ai

import (
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/chinaboard/brewing/pkg/cfg"
	"github.com/pkoukk/tiktoken-go"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	MAX_TOKEN = 1000
	MODEL     = openai.GPT3Dot5Turbo
)

type GPTService struct {
	tkm   *tiktoken.Tiktoken
	model string
}

func init() {
	retry.Attempts(3)
	tkm, _ := tiktoken.EncodingForModel(MODEL)
	Service = &GPTService{model: MODEL, tkm: tkm}
}

var Service *GPTService

func newClient() *openai.Client {
	config := openai.DefaultConfig(cfg.OpenAiToken)
	config.BaseURL = "https://api.openai.com/v1"
	if cfg.OpenAiBaseURL != "" {
		config.BaseURL = cfg.OpenAiBaseURL
	}
	if cfg.OpenAiProxy != "" {
		proxyUrl, err := url.Parse(cfg.OpenAiProxy)
		if err != nil {
			panic(err)
		}
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
		config.HTTPClient = &http.Client{
			Transport: transport,
		}
	}
	return openai.NewClientWithConfig(config)
}

func (g *GPTService) Text(i int, content string) (string, error) {
	msg := renderRequest(content)
	logrus.Debugln("token", i, "request:", g.numTokensFromMessages(&msg))
	resp, err := newClient().CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: msg,
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func (g *GPTService) TextStream(content string) {
	msg := renderRequest(content)
	logrus.Debugln("token request:", g.numTokensFromMessages(&msg))
	stream, err := newClient().CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: msg,
			Stream:   true,
		},
	)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			return
		}

		if err != nil {
			fmt.Printf("\nSummary error: %v\n", err)
			return
		}

		fmt.Print(response.Choices[0].Delta.Content)
	}
}

func (g *GPTService) Summary(content []string) (string, error) {
	parts := g.SplitContent(content)
	var resultParts []string
	logrus.Debugln("req content", len(parts), "parts")
	for i, c := range parts {
		dt := time.Now()
		resultPart, err := g.Text(i, c)
		if err != nil {
			return "", err
		}
		logrus.Debugln("content", i, "usage", time.Since(dt))
		resultParts = append(resultParts, resultPart)
	}
	return strings.Join(resultParts, ""), nil
}

func (g *GPTService) SummaryParallel(parts []string) ([]string, []error) {
	var resultParts []string
	m := sync.Map{}
	var wg sync.WaitGroup
	var e []error
	st := time.Now()
	ch := make(chan struct{}, runtime.NumCPU())
	logrus.Debugln("req content", len(parts), "parts")
	for i, c := range parts {
		ch <- struct{}{}
		wg.Add(1)
		i, c := i, c
		go func() {
			defer wg.Done()
			dt := time.Now()

			resultPart, err := retry.DoWithData(func() (string, error) {
				retry.RetryIf(func(err error) bool {
					return strings.Contains(err.Error(), "503") ||
						strings.Contains(err.Error(), "401") ||
						strings.Contains(err.Error(), "429")
				})
				r, e := g.Text(i, c)
				if e != nil {
					if strings.Contains(e.Error(), "You exceeded your current quota, please check your plan and billing details") {
						e = retry.Unrecoverable(e)
					}
				}
				return r, e
			})

			e = append(e, err)
			if err != nil {
				logrus.Debugln("content", i, "err", err)
			}
			logrus.Debugln("content", i, "usage", time.Since(dt))
			m.Store(i, resultPart)
			<-ch
		}()
	}
	wg.Wait()

	logrus.Debugln("content", "process", "done", time.Since(st))
	for i, _ := range parts {
		v, _ := m.Load(i)
		resultParts = append(resultParts, v.(string))
	}
	return resultParts, e
}

func (g *GPTService) SummaryStream(content []string) {
	parts := g.SplitContent(content)
	for _, part := range parts {
		g.TextStream(part)
	}
}

func (g *GPTService) SplitContent(content []string) []string {
	logrus.Debugln("raw content", len(content), "parts")
	dt := time.Now()
	var parts []string
	var str string
	var strToken int
	for _, c := range content {
		msg := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: c,
			}}
		token := g.numTokensFromMessages(&msg)
		if strToken+token < MAX_TOKEN {
			strToken += token
			str += c
		} else {
			parts = append(parts, str)
			strToken = token
			str = c
		}
	}
	parts = append(parts, str)
	logrus.Debugln("split content usage", time.Since(dt))
	return parts
}

func (g *GPTService) numTokensFromMessages(messages *[]openai.ChatCompletionMessage) (numTokens int) {
	model, tkm := g.model, g.tkm

	var tokens_per_message int
	var tokens_per_name int
	if model == "gpt-3.5-turbo-0301" || model == "gpt-3.5-turbo" {
		tokens_per_message = 4
		tokens_per_name = -1
	} else if model == "gpt-4-0314" || model == "gpt-4" {
		tokens_per_message = 3
		tokens_per_name = 1
	} else {
		fmt.Println("Warning: collection not found. Using cl100k_base encoding.")
		tokens_per_message = 3
		tokens_per_name = 1
	}

	for _, message := range *messages {
		numTokens += tokens_per_message
		numTokens += len(tkm.Encode(message.Content, nil, nil))
		numTokens += len(tkm.Encode(message.Role, nil, nil))
		numTokens += len(tkm.Encode(message.Name, nil, nil))
		if message.Name != "" {
			numTokens += tokens_per_name
		}
	}

	numTokens += 3
	return numTokens
}

func renderRequest(content string) []openai.ChatCompletionMessage {
	return []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "请你修改错字错词整理成文章格式，如果无法处理，请原样返回",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: content,
		}}
}
