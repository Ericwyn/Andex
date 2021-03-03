package main

import (
	"fmt"
	"github.com/Ericwyn/Andex/api"
	"github.com/Ericwyn/Andex/conf"
	"github.com/Ericwyn/Andex/controller"
	"github.com/Ericwyn/Andex/service"
	"github.com/Ericwyn/GoTools/file"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"net/http"
	"os"
	"time"
)

func main() {

	////// 获取文件夹文件列表
	//list := api.FolderList("root")
	//for _,item := range list.Items {
	//	fmt.Println("name:", item.Name, ", type:", item.Type, ", id:", item.FileID)
	//}

	// 载入配置
	loadConfig()

	// 运行配置
	runCorn()

	// 启动 service
	startServer()
}

func runCorn() {
	s := gocron.NewScheduler(time.UTC)

	// 每 30 分钟刷新一次配置
	s.Every(30).Minutes().Do(func() {
		fmt.Println("corn 刷新 token 配置")
		api.RefreshToken()
	})
}

func loadConfig() {
	// 创建系统配置文件夹
	sysConfDir := file.OpenFile(conf.SysConfigDirPath)
	sysConfDir.Mkdirs()

	userConfigFile := file.OpenFile(conf.UserConfigFilePath)
	if !userConfigFile.Exits() {
		fmt.Println("未检测到配置文件, 创建空白配置文件")
		conf.CreateUserConfFile()
		os.Exit(0)
	}

	// 载入 System 和 User 配置
	conf.LoadConfFromFile()
	if conf.SysConfigNow.RefreshToken == "NULL" {
		fmt.Println("RefreshToken 未配置")
		os.Exit(0)
	}

	// api.RefreshToken() // 不需要每次启动都刷新 refreshToken 和 token

	if conf.SysConfigNow.Authorization == "NULL" {
		// 如果 Authorization 为 NULL, 先尝试刷新一遍 token
		api.RefreshToken()
		if conf.SysConfigNow.Authorization == "NULL" {
			fmt.Println("token 无法获取 RefreshToken 错误或已过期")
			os.Exit(0)
		}
	}

	// 载入配置后验证 token 是否已过期
	info := api.UserInfo()
	if info != nil {
		fmt.Println("token 未过期")
	} else {
		fmt.Println("token 已过期, 尝试重新获取")
		// 刷新配置
		api.RefreshToken()
	}

	// 验证根目录配置
	pathSetTrue := service.CheckRootPathSet()
	if !pathSetTrue {
		fmt.Println("Andex RootPath 参数设置错误, 请设置为正确的网盘文件夹路径")
		os.Exit(0)
	}
}

func startServer() {
	gin.SetMode(gin.DebugMode)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        controller.NewMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}
