package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginlogrus "github.com/toorop/gin-logrus"
	"html/template"
)

func InitRouter(logger *logrus.Logger) *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(ginlogrus.Logger(logger), gin.Recovery())

	tR, err := newTaskRepo()
	if err != nil {
		logrus.Fatalln(err)
	}
	aR, err := newAsrRepo()
	if err != nil {
		logrus.Fatalln(err)
	}

	r.SetFuncMap(template.FuncMap{
		"safeHTML": func(text string) template.HTML { return template.HTML(text) },
	})
	r.GET("/share/:id", aR.Summary)
	r.GET("/", Index)
	r.LoadHTMLGlob("templates/*")

	brewingApi := r.Group("v1/api")
	{
		d := brewingApi.Group("/docker")
		{
			d.POST("/job", tR.Add)
			d.POST("/job/:id", tR.Run)
			d.GET("/job/:id", tR.Get)
		}
		a := brewingApi.Group("/asr")
		{
			a.GET("/trigger/:id", aR.Add)
			a.POST("/job/:id", aR.Run)
			a.GET("/job/:id", aR.Get)
		}
		e := brewingApi.Group("/chain")
		{
			e.POST("/job", (&Chain{repo: aR}).chain)
		}
	}

	return r
}
