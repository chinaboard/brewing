package main

import (
	"fmt"
	"github.com/chinaboard/brewing/controller"
	"github.com/chinaboard/brewing/pkg/bininfo"
	"github.com/chinaboard/brewing/pkg/cfg"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	routersInit := controller.InitRouter(logrus.StandardLogger())

	server := &http.Server{
		Addr:    ":" + cfg.HttpPort,
		Handler: routersInit,
	}

	logrus.Printf("Starting server on :%v", cfg.HttpPort)

	logrus.Fatalln(server.ListenAndServe())
}

func init() {
	fmt.Println(bininfo.StringifyMultiLine())
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		ForceQuote:      false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}
