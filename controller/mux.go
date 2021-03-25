package controller

import (
	"crypto/rand"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"math/big"

	"github.com/Ericwyn/GoTools/file"
	"github.com/gin-gonic/gin"
)

// 设置 API 路由
func initAPI(router *gin.Engine) {
	router.GET("/", pages)
	router.GET("/download", download)

	router.POST("/adminLogin", adminLogin)
	router.POST("/adminLogout", adminLogout)
}

// 返回全局路由, 包括静态资源
func NewMux() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	loadStaticPath(router)

	store := cookie.NewStore(GeneralSessionKey())
	router.Use(sessions.Sessions("andex-session", store))

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
			router.Static("static/"+child.Name(), staticDirPath+"/"+child.Name())
		}
	}
}

var keyParisLen = 64

func GeneralSessionKey() []byte {
	return []byte(string(GeneralRandomStr(keyParisLen)))
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*("

func GeneralRandomStr(length int) string {
	str := ""
	for i := 0; i < length; i++ {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(length)))
		index64 := index.Int64()
		str += letterBytes[int(index64) : int(index64)+1]
	}
	return str
}
