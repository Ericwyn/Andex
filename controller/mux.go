package controller

import (
	"crypto/rand"
	"fmt"
	"github.com/Ericwyn/Andex/modal"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"math/big"

	"github.com/Ericwyn/GoTools/file"
	"github.com/gin-gonic/gin"
)

// 设置 API 路由
func initAPI(router *gin.Engine) {
	router.GET("/", apiPages)
	router.GET("/download", apiDownload)

	router.POST("/adminLogin", apiAdminLogin)
	router.POST("/adminLogout", apiAdminLogout)
	router.POST("/setPassword", apiSetPassword)
	router.POST("/removePassword", apiRemovePassword)

	// 申请某个路径的访问权限
	router.POST("/pathPermRequest", apiPathPermRequest)
}

// 返回全局路由, 包括静态资源
func NewMux() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	loadStaticPath(router)

	store := cookie.NewStore(GetCookieKey())
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

func GetCookieKey() []byte {
	key := modal.GetConfig(modal.TypeCookieKey, "NULL").(string)
	if key != "NULL" {
		return []byte(key)
	} else {
		key = GeneralRandomStr(keyParisLen)
		err := modal.SaveConf(modal.TypeCookieKey, key)
		if err != nil {
			fmt.Println("save cookie key error", err)
		}
		return []byte(key)
	}
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
