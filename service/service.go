package service

import (
	"encoding/json"
	"fmt"
	"github.com/Ericwyn/Andex/conf"
	"sort"
	"strings"

	"github.com/Ericwyn/Andex/api"
	"github.com/Ericwyn/Andex/storage"
	"github.com/Ericwyn/GoTools/date"
	"github.com/Ericwyn/GoTools/file"
)

// 获取这个路径下面的文件列表
type PathDetailBean struct {
	Name        string
	ParentPath  string
	Type        string
	CreateTime  string
	UpdateTime  string
	Path        string
	DownloadUrl string
	Size        string
}

// 获取这个路径下面的文件列表
type FileDetailBean struct {
	Name        string
	ParentPath  string
	Type        string
	DownloadUrl string
	CreateTime  string
	UpdateTime  string
	Size        string
	Path        string
}

type Path struct {
	Name   string // 路径地址
	FileId string // 路径的 fileId
	IsDir  bool   // 类型，
}

// 用来存储 path 与 fileId 的对应, 会同步到 pathLog.json 里面
// 默认是 root 文件夹
var driverRootPath *Path
var pathMap = map[string]Path{}

// 格式化 query 参数
func FormatPathQuery(path string) string {
	// 参数格式化
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
	return path
}

// 检查 config.json 当中的 rootPath 设置是否正确
func CheckRootPathSet() bool {

	conf.ConfigNow.RootPath = FormatPathQuery(conf.ConfigNow.RootPath)

	fmt.Println("检查 RootPath 设置:", conf.ConfigNow.RootPath)

	if conf.ConfigNow.RootPath == "" || conf.ConfigNow.RootPath == "/" {
		return true
	} else {
		pathSplit := strings.Split(conf.ConfigNow.RootPath, "/")
		fileIdNow := "root"
		fmt.Print("验证文件夹路径: ")
		for _, pathName := range pathSplit {
			if pathName == "" {
				fileIdNow = "root"
				continue
			}
			folderList := api.FolderList(fileIdNow)

			pathFlag := false

			for _, fileMsg := range folderList.Items {
				if fileMsg.Type == "folder" && fileMsg.Name == pathName {
					fmt.Print(" -> ", fileMsg.Name)
					fileIdNow = fileMsg.FileID
					pathFlag = true
					break
				}
			}
			if !pathFlag {
				fmt.Println("文件夹路径", conf.ConfigNow.RootPath, "错误, 无法找到文件夹:", pathName)
				return false
			}
		}
		fmt.Println()
		fmt.Println("文件夹路径验证成功")

		driverRootPath = &Path{
			Name:   "/",
			FileId: fileIdNow,
			IsDir:  true,
		}
		pathMap["/"] = *driverRootPath
		pathMap["/root"] = *driverRootPath
		// TODO pathMap 不需要每次都全部重写
		savePathMapToLocal()
		return true
	}
}

// 路径是否正确
func IsPathTrue(path string) bool {
	// 载入本地配置文件
	if !loadPathFromLocalFlag {
		LoadPathMapFromLocal()
		loadPathFromLocalFlag = true
	}

	// 查找能否从 pathMap 里面找到这个 path 对应的 fileId
	fmt.Println("请求", path)
	if _, ok := pathMap[path]; !ok {
		fmt.Println("无法从该路径找到对应的 fileId: ", path)
		return false
	}
	return true
}

// 是否为文件
// 返回参数1代表是否是文件， 参数2代表路径是否正确
func IsPathIsFile(path string) bool {
	// 载入本地配置文件
	if !loadPathFromLocalFlag {
		LoadPathMapFromLocal()
		loadPathFromLocalFlag = true
	}
	return !pathMap[path].IsDir
}

// 传入一个文件路径, 获取该文件的详情
// 先通过文件的 parent 路径, 获取其 parent 路径的 fileId
// 然后通过 fileList 的方式, 获取该文件的 fileId
func GetFileDetail(path string) (*FileDetailBean, bool) {
	// 载入本地配置文件
	if !loadPathFromLocalFlag {
		LoadPathMapFromLocal()
		loadPathFromLocalFlag = true
	}
	split := strings.Split(path, "/")
	parentPath := path[0 : len(path)-len(split[len(split)-1])]
	parentPath = FormatPathQuery(parentPath)

	dirDetails, haveDetail := GetPathDetail(parentPath)

	if !haveDetail {
		return nil, false
	}

	for _, fileDetail := range dirDetails {
		// 文件名判断
		if fileDetail.Name == split[len(split)-1] && fileDetail.Type == "file" {
			return &FileDetailBean{
				Name:       fileDetail.Name,
				ParentPath: fileDetail.ParentPath,
				Type:       "file",
				CreateTime: fileDetail.CreateTime,
				UpdateTime: fileDetail.UpdateTime,
				Path:       fileDetail.Path,
				Size:       fileDetail.Size,
			}, true

		}
	}
	return nil, false
}

type FileDownMsgBean struct {
	Name string
	Url  string
}

func GetFileDownloadUrl(path string) *FileDownMsgBean {

	// 载入本地配置文件
	if !loadPathFromLocalFlag {
		LoadPathMapFromLocal()
		loadPathFromLocalFlag = true
	}

	if _, ok := pathMap[path]; !ok {
		return nil
	}

	fileDetail := pathMap[path]
	if !fileDetail.IsDir {
		url := api.GetDownloadUrlByFileIdAndFileName(fileDetail.FileId, fileDetail.Name)
		if url != "" {
			return &FileDownMsgBean{
				Name: fileDetail.Name,
				Url:  url,
			}
		}
	}
	return nil
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

	folderList := api.FolderList(pathMap[path].FileId)
	if folderList == nil || folderList.Items == nil {
		fmt.Println("路径:", path, "获取 folderList 失败")
		return nil, false
	}

	result := make([]PathDetailBean, 0)

	pathMapUpdateFlag := false

	// 更新 pathMap 缓存, 刷新到本地
	for _, fileMsgBean := range folderList.Items {

		filePath := path + "/" + fileMsgBean.Name
		filePath = strings.ReplaceAll(filePath, "//", "/")

		// 构造 PathDetailBean
		result = append(result, PathDetailBean{
			Name:        fileMsgBean.Name,
			Type:        fileMsgBean.Type,
			Path:        filePath,
			CreateTime:  date.Format(fileMsgBean.CreatedAt, "yyyy-MM-dd HH:mm"),
			UpdateTime:  date.Format(fileMsgBean.UpdatedAt, "yyyy-MM-dd HH:mm"),
			DownloadUrl: fileMsgBean.DownloadURL,
			Size:        byteCountBinary(fileMsgBean.Size),
			ParentPath:  "",
		})

		// 更新缓存, 无论是文件还是文件夹都会缓存下来
		if _, ok := pathMap[filePath]; ok {
			if pathMap[filePath].FileId != fileMsgBean.FileID ||
				pathMap[filePath].Name != fileMsgBean.Name ||
				pathMap[filePath].IsDir != (fileMsgBean.Type == "folder") {
				// 需要更新缓存
				pathMapUpdateFlag = true
				pathMap[filePath] = Path{
					Name:   fileMsgBean.Name,
					FileId: fileMsgBean.FileID,
					IsDir:  fileMsgBean.Type == "folder",
				}
			}
		} else {
			// 需要加入缓存
			pathMapUpdateFlag = true
			pathMap[filePath] = Path{
				Name:   fileMsgBean.Name,
				FileId: fileMsgBean.FileID,
				IsDir:  fileMsgBean.Type == "folder",
			}
		}

	}

	if pathMapUpdateFlag {
		// 将缓存刷新到本地
		// TODO 如果 /a/b -> 映射为 xxx, 之后改名为 /a/bb, 那么 /a/bb 映射为 xxx, 但是 /a/b 将不会删除
		// TODO 清除旧的缓存
		savePathMapToLocal()
	}

	// result 排序
	sort.Slice(result, func(i, j int) bool {
		return strings.Compare(result[i].Name, result[j].Name) < 0
	})

	// 构造 PathDetailBean 返回
	return result, true
}

// 构造路径面包屑
type NavPath struct {
	Name string
	Path string
	Last bool
}

func GetNavPathList(path string) []NavPath {
	split := strings.Split(path, "/")
	navPathList := make([]NavPath, 0)
	navPathList = append(navPathList, NavPath{
		Name: "首页",
		Path: "/",
		// 判断是否是首页请求
		Last: path == "/" || path == "/root",
	})
	navPathNow := ""
	for i, navPathTemp := range split {
		if navPathTemp == "" {
			continue
		}
		navPathNow = navPathNow + "/" + navPathTemp
		navPathList = append(navPathList, NavPath{
			Name: navPathTemp,
			Path: navPathNow,
			Last: i >= (len(split) - 1),
		})
	}

	return navPathList
}

// 文件大小可读形式输出
func byteCountBinary(size int64) string {
	if size == 0 {
		return ""
	}
	const unit int64 = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := unit, 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(size)/float64(div), "KMGTPE"[exp])
}

//================================== FileId 缓存 相关逻辑==================================

const localPathMapConf = "./pathMap.json"

var loadPathFromLocalFlag = false

func LoadPathMapFromLocal() {
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

	if driverRootPath != nil {
		pathMap["/"] = *driverRootPath
		pathMap["/root"] = *driverRootPath
	}
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
