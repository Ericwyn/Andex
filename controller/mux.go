package controller

import (
	"fmt"
	"github.com/Ericwyn/GoTools/file"
	"github.com/gin-gonic/gin"
)

// 设置 API 路由
func initAPI(router *gin.Engine) {
	router.GET("/", pages)
}

// 返回全局路由, 包括静态资源
func NewMux() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	loadStaticPath(router)

	router.Use(gin.Logger())
	router.LoadHTMLGlob("static/*.html")

	initAPI(router)
	return router
}

func loadStaticPath(router *gin.Engine) {
	staticDirPath := "./static"
	staticDir := file.OpenFile(staticDirPath)
	children := staticDir.Children()
	for _, child := range children {
		if child.IsDir() {
			fmt.Println("load static router:", "/"+child.Name(), "->", staticDirPath+"/"+child.Name())
			router.Static(child.Name(), staticDirPath+"/"+child.Name())
		}
	}
}
