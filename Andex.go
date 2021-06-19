package main

import (
	"flag"
	"github.com/Ericwyn/Andex/controller"
	"github.com/Ericwyn/Andex/modal"
	"github.com/Ericwyn/Andex/service"
	"github.com/Ericwyn/Andex/util/log"
	"github.com/Ericwyn/GoTools/file"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var debugMode = flag.Bool("debug", false, "print sql log and gin debug log")

func main() {

	// TODO, 密码重置之后，用户的访问权限
	flag.Parse()

	if *debugMode {
		log.I("DEBUG 模式已打开")
	}

	modal.InitDb(*debugMode)

	// 载入配置
	loadConfig()

	// 运行配置
	runCorn()

	// 启动 service
	startServer()
}

func loadConfig() {
	log.I("==================== 启动配置 Start ==================== ")

	userConfigFile := file.OpenFile(service.UserConfigFilePath)
	if !userConfigFile.Exits() {
		log.E("未检测到配置文件, 创建空白配置文件, 请再 ./config.json 中进行参数配置")
		service.CreateUserConfFile()
		os.Exit(0)
	}

	// 载入 System 和 User 配置
	service.LoadConfFromFile()
	if service.SysConfigNow.RefreshToken == "NULL" {
		log.E("RefreshToken 未配置")
		os.Exit(0)
	}

	// api.RefreshToken() // 不需要每次启动都刷新 refreshToken 和 token
	log.I("服务运行端口:", service.UserConfNow.Port)

	if service.SysConfigNow.Authorization == "NULL" {
		// 如果 Authorization 为 NULL, 先尝试刷新一遍 token
		service.RefreshToken()
		if service.SysConfigNow.Authorization == "NULL" {
			log.E("token 无法获取 RefreshToken 错误或已过期")
			os.Exit(0)
		}
	} else {
		log.I("从数据库中载入 Authorization")
	}

	// 通过配置的时间来确认是否过期, 而不是执行一次请求
	// 在距离上次刷新超过最大时间间隔
	// refresh token 接口里面 expires 时间是 7200 这里取其 3/4 长度作为最大过期间隔
	log.I("运行配置载入完毕")

	var maxTokenExpireTime int64 = 60 * 90
	log.I("token 过期时间设置:", maxTokenExpireTime)
	if service.SysConfigNow.LastGetTokenTime.Unix() == (time.Time{}).Unix() ||
		time.Now().Unix()-service.SysConfigNow.LastGetTokenTime.Unix() > maxTokenExpireTime {
		service.RefreshToken()
	} else {
		log.I("距离上次 token 获取时间已过去:",
			time.Now().Unix()-service.SysConfigNow.LastGetTokenTime.Unix())
	}

	log.I()
	// 验证根目录配置
	pathSetTrue := service.CheckRootPathSet()
	if !pathSetTrue {
		log.E("Andex RootPath 参数设置错误, 请设置为正确的网盘文件夹路径")
		os.Exit(0)
	}

	// 载入 README 文件
	service.LoadReadmeFile()

	log.I("===================== 启动配置 End =====================")
}

var cornFirstFlag = true

func runCorn() {
	s := gocron.NewScheduler(time.UTC)

	// 每 28 分钟刷新一次配置
	s.Every(28).Minutes().Do(func() {
		if cornFirstFlag {
			cornFirstFlag = false
		} else {
			log.I("corn 执行刷新 token 配置")
			service.RefreshToken()
		}
	})

	s.StartAsync()
}

func startServer() {
	if *debugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	s := &http.Server{
		Addr:           ":" + service.UserConfNow.Port,
		Handler:        controller.NewMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}
