package service

import (
	"encoding/json"
	"fmt"
	"github.com/Ericwyn/Andex/api"
	"github.com/Ericwyn/Andex/fileIO"
	"github.com/Ericwyn/GoTools/file"
	"strings"
)

// 获取这个路径下面的文件列表
type PathDetailBean struct {
	Name       string
	ParentPath string
	Type       string
}

// 用来存储 path 与 fileId 的对应, 会同步到 pathLog.json 里面
var pathMap = map[string]string{
	"/": "root",
}

// 通过一个 path, 如果 /share/wx, 来获取这个 path 下面对应的 PathDetailBean
// 如果可以找到的话, 就返回 PathDetailBean, true, 如果不行的话返回 nil, false
// 如果 api 请求失败的话, 会返回 [], true
func GetPathDetail(path string) ([]PathDetailBean, bool) {

	// 载入本地配置文件
	if !loadPathFromLocalFlag {
		LoadPathMapFromLocal()
		loadPathFromLocalFlag = true
	}

	// 去除最后的 /
	path = strings.ReplaceAll(path, "//", "/")
	if path[0] != '/' {
		path = "/" + path
	}
	if path[len(path)-1] == '/' {
		path = path[0 : len(path)-1]
	}
	if path == "" {
		path = "/"
	}

	// 查找能否从 pathMap 里面找到这个 path 对应的 fileId
	fmt.Println("请求", path)
	if _, ok := pathMap[path]; !ok {
		fmt.Println("无法从该路径找到对应的 fileId: ", path)
		return nil, false
	}
	folderList := api.FolderList(pathMap[path])
	if folderList == nil || folderList.Items == nil {
		fmt.Println("路径:", path, "获取 folderList 失败")
		return make([]PathDetailBean, 0), true
	}
	result := make([]PathDetailBean, 0)

	pathMapUpdateFlag := false

	// 更新 pathMap 缓存, 刷新到本地
	for _, fileMsgBean := range folderList.Items {
		// 构造 PathDetailBean
		result = append(result, PathDetailBean{
			Name:       fileMsgBean.Name,
			Type:       fileMsgBean.Type,
			ParentPath: "",
		})

		// 更新缓存, 如果是文件夹的话, 拼接路径
		if fileMsgBean.Type == "folder" {
			pathMapKey := path + "/" + fileMsgBean.Name
			pathMapKey = strings.ReplaceAll(pathMapKey, "//", "/")
			// 如果已有缓存
			if _, ok := pathMap[pathMapKey]; ok {
				if pathMap[pathMapKey] != fileMsgBean.FileID {
					// 需要更新缓存
					pathMapUpdateFlag = true
					pathMap[pathMapKey] = fileMsgBean.FileID
				}
			} else {
				// 需要加入缓存
				pathMapUpdateFlag = true
				pathMap[pathMapKey] = fileMsgBean.FileID
			}
		}
	}

	if pathMapUpdateFlag {
		// 将缓存刷新到本地
		// TODO 如果 /a/b -> 映射为 xxx, 之后改名为 /a/bb, 那么 /a/bb 映射为 xxx, 但是 /a/b 将不会删除
		// TODO 清除旧的缓存
		savePathMapToLocal()
	}

	// 构造 PathDetailBean 返回
	return result, true
}

const localPathMapConf = "./pathMap.json"

var loadPathFromLocalFlag = false

func LoadPathMapFromLocal() {
	openFile := file.OpenFile(localPathMapConf)
	if !openFile.Exits() {
		fmt.Println("path map 缓存文件不存在")
		return
	} else {
		pathMapString, err := fileIO.ReadFileAsString(localPathMapConf)
		if err != nil {
			fmt.Println("读取 path map 缓存文件失败", err)
		} else {
			err := json.Unmarshal([]byte(pathMapString), &pathMap)
			if err != nil {
				fmt.Println("读取 path map 缓存文件失败, 序列化失败", err)
			}
		}
	}
	pathMap["/"] = "root"
	pathMap["/root"] = "root"
}

func savePathMapToLocal() {
	jsonRes, err := json.MarshalIndent(pathMap, "", "  ")
	if err != nil {
		fmt.Println("输出 path map 缓存到本地时候发生格式化错误", err)
	} else {
		err := fileIO.WriteStringToFile(localPathMapConf, string(jsonRes), false)
		if err != nil {
			fmt.Println("输出 path map 缓存到本地时候发生 io 错误", err)
		} else {
			fmt.Println("刷新 path map 缓存到本地")
		}
	}
}
