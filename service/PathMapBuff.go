package service

import (
	"encoding/json"
	"fmt"
	"github.com/Ericwyn/Andex/conf"
	"github.com/Ericwyn/Andex/storage"
	"github.com/Ericwyn/GoTools/file"
)

//================================== FileId 缓存 相关逻辑==================================

const localPathMapConf = conf.SysConfigDirPath + "/pathMap.json"

func LoadPathMapFromLocal() {
	fmt.Println("开始载入本地 path map 缓存文件")
	openFile := file.OpenFile(localPathMapConf)
	if !openFile.Exits() {
		fmt.Println("path map 缓存文件不存在")
		return
	} else {
		pathMapString, err := storage.ReadFileAsString(localPathMapConf)
		if err != nil {
			fmt.Println("读取 path map 缓存文件失败", err)
		} else {
			err := json.Unmarshal([]byte(pathMapString), &pathMap)
			if err != nil {
				fmt.Println("读取 path map 缓存文件失败, 序列化失败", err)
			}
		}
	}
	//if driverRootPath != nil {
	//	pathMap["/"] = *driverRootPath
	//	pathMap["/root"] = *driverRootPath
	//}
}

func savePathMapToLocal() {
	jsonRes, err := json.MarshalIndent(pathMap, "", "  ")
	if err != nil {
		fmt.Println("输出 path map 缓存到本地时候发生格式化错误", err)
	} else {
		err := storage.WriteStringToFile(localPathMapConf, string(jsonRes), false)
		if err != nil {
			fmt.Println("输出 path map 缓存到本地时候发生 io 错误", err)
		} else {
			fmt.Println("刷新 path map 缓存到本地")
		}
	}
}
