package main

import (
	"fmt"
	"github.com/Ericwyn/Andex/api"
	"github.com/Ericwyn/Andex/conf"
	"github.com/Ericwyn/GoTools/file"
	"os"
)

func main() {

	//loadConfig()
	conf.LoadConfFromFile()

	info := api.UserInfo()
	fmt.Println("用户昵称:", info.NickName)


	// 获取文件夹文件列表
	list := api.FolderList("root")
	for _,item := range list.Items {
		fmt.Println("name:", item.Name, ", type:", item.Type, ", id:", item.FileID)
	}

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
