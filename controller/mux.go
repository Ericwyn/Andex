package controller

import (
	"crypto/rand"
	"github.com/Ericwyn/Andex/modal"
	"github.com/Ericwyn/Andex/util/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"math/big"
	"strings"

	"github.com/gin-gonic/gin"
)

// 设置 API 路由
func initAPI(router *gin.Engine) {
	//router.GET("/", apiPages)
	//router.GET("/download", apiDownload)

	// 使用自定义路
	router.GET("/*requestPath", func(ctx *gin.Context) {
		requestPath := ctx.Param("requestPath")
		//fmt.Println("request: " + requestPath)
		downloadFlag := ctx.Query("dl") == "1"
		if isStaticPath(requestPath) {
			apiStatic(requestPath, ctx)
		} else if downloadFlag {
			apiDownload(requestPath, ctx)
		} else {
			apiPages(requestPath, ctx)
		}
	})

	router.POST("/adminLogin", apiAdminLogin)
	router.POST("/adminLogout", apiAdminLogout)
	router.POST("/setPassword", apiSetPassword)
	router.POST("/removePassword", apiRemovePassword)
	router.POST("/getRedirectLink", apiGetRedirectLink)

	// 申请某个路径的访问权限
	router.POST("/pathPermRequest", apiPathPermRequest)
}

var staticDirAndPath = []string{"/css/", "/fonts/", "/icons/", "/js/", "/assets/", "/favicon.ico", "favicon.ico"}

func isStaticPath(requestPath string) bool {
	for _, staicPath := range staticDirAndPath {
		if strings.Index(requestPath, staicPath) == 0 {
			return true
		}
	}
	return false
}

// 返回全局路由, 包括静态资源
func NewMux() *gin.Engine {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	store := cookie.NewStore(getCookieKey())
	router.Use(sessions.Sessions("andex-session", store))

	router.Use(gin.Logger())
	router.LoadHTMLGlob("static/*.html")

	initAPI(router)
	return router
}

var keyParisLen = 64

// 获取 cookie 加密钥匙，首次启动的时候生成随机 cookie key
func getCookieKey() []byte {
	key := modal.GetConfig(modal.TypeCookieKey, "NULL").(string)
	if key != "NULL" {
		return []byte(key)
	} else {
		key = GeneralRandomStr(keyParisLen)
		err := modal.SaveConf(modal.TypeCookieKey, key)
		if err != nil {
			log.E("save cookie key error", err)
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
