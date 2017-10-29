package api

import (
	"gopkg.in/gin-gonic/gin.v1"
)

func SetupRoutes(router *gin.Engine) {
	router.HEAD("/:uid", FilesHandler)
	router.GET("/:uid", FilesHandler)

}
