package controller

import (
	"github.com/chinaboard/brewing/pkg/cfg"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html",
		gin.H{
			"WhisperEndpointSchema": cfg.WhisperEndpointSchema,
			"WhisperEndpoint":       cfg.WhisperEndpoint,
		},
	)
}
