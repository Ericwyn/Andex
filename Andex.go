package main

import (
	"fmt"
	"github.com/Ericwyn/Andex/api"
	"github.com/Ericwyn/Andex/conf"
	"github.com/Ericwyn/Andex/controller"
	"github.com/Ericwyn/GoTools/file"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

func main() {

	loadConfig()
	//conf.LoadConfFromFile()
	//
	//info := api.UserInfo()
	//fmt.Println("用户昵称:", info.NickName)

	////// 获取文件夹文件列表
	//list := api.FolderList("root")
	//for _,item := range list.Items {
	//	fmt.Println("name:", item.Name, ", type:", item.Type, ", id:", item.FileID)
	//}

	startServer()

}

func loadConfig() {
	configFile := file.OpenFile(conf.ConfigFilePath)
	if !configFile.Exits() {
		fmt.Println("未检测到配置文件, 创建空白配置文件")
		conf.SaveConf()
		os.Exit(0)
	}
	// 载入配置
	conf.LoadConfFromFile()
	if conf.ConfigNow.RefreshToken == "NULL" {
		fmt.Println("RefreshToken 未配置")
		os.Exit(0)
	}

	// 刷新配置
	api.RefreshToken()

	if conf.ConfigNow.Authorization == "NULL" {
		fmt.Println("RefreshToken 错误或已过期")
		os.Exit(0)
	}
}

func startServer() {
	gin.SetMode(gin.ReleaseMode)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        controller.NewMux(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	_ = s.ListenAndServe()
}
