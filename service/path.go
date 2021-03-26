package service

import (
	"fmt"
	"github.com/Ericwyn/Andex/modal"
	"github.com/Ericwyn/GoTools/date"
	"sort"
	"strings"

	"github.com/Ericwyn/Andex/api"
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
	HadPassword bool
	HadAccess   bool
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

type Path1 struct {
	Name   string // 路径地址
	FileId string // 路径的 fileId
	IsDir  bool   // 类型，
}

type FileDownMsgBean struct {
	Name string
	Url  string
}

// 构造路径面包屑
type NavPath struct {
	Name string
	Path string
	Last bool
}

// 用来存储 path 与 fileId 的对应, 会同步到 pathLog.json 里面
// 默认是 root 文件夹
var driverRootPath *modal.AndexPath
var pathMapBuff = map[string]modal.AndexPath{}

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

	SysConfigNow.RootPath = FormatPathQuery(SysConfigNow.RootPath)

	fmt.Println("检查 RootPath 设置:", SysConfigNow.RootPath)

	pathSplit := strings.Split(SysConfigNow.RootPath, "/")
	fileIdNow := "root"
	fmt.Print("验证文件夹路径: ")
	for _, pathName := range pathSplit {
		if pathName == "" {
			fileIdNow = "root"
			continue
		}
		folderList := api.FolderList(SysConfigNow.Authorization, SysConfigNow.DriveID, fileIdNow)

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
			fmt.Println("文件夹路径", SysConfigNow.RootPath, "错误, 无法找到文件夹:", pathName)
			return false
		}
	}
	fmt.Println()
	fmt.Println("文件夹路径验证成功, FileId:", fileIdNow)

	driverRootPath = &modal.AndexPath{
		Name:   "/",
		FileId: fileIdNow,
		IsDir:  true,
	}

	// 只有在 user config 的 rootPath 更新的情况下， 才去更新 pathMap.json， 否则沿用旧的
	var err error
	pathMapBuff, err = modal.LoadPathMap()
	if err != nil {
		fmt.Println("从数据库载入 path map 缓存失败:", err)
		return false
	}

	// 看看本地缓存文件里面的 / 和 /root 其 fileId 与最新的是否一致
	if _, ok := pathMapBuff["/"]; ok {
		fmt.Println("旧 Root Path -> FileId:", pathMapBuff["/"].FileId)
		fmt.Println("旧 Root Path -> FileId:", driverRootPath.FileId)
		// 不一致的情况需要刷新 pathMap
		if pathMapBuff["/"].FileId != driverRootPath.FileId {
			fmt.Println("检测到 rootPath 设置已更新， 清除 path map 缓存")
			modal.DeleteAllPath()

			pathMapBuff["/"] = modal.AndexPath{
				Path:     "/",
				Name:     driverRootPath.Name,
				FileId:   driverRootPath.FileId,
				IsDir:    driverRootPath.IsDir,
				Password: driverRootPath.Password,
			}
			pathMapBuff["/root"] = modal.AndexPath{
				Path:     "/root",
				Name:     driverRootPath.Name,
				FileId:   driverRootPath.FileId,
				IsDir:    driverRootPath.IsDir,
				Password: driverRootPath.Password,
			}

			modal.SavePathMap(pathMapBuff)
		} else {
			fmt.Println("rootPath 设置未变更")
		}
	} else {
		fmt.Println("无法从已有 pathMap 缓存当中找到根目录 FileId 设置")
		// 如果 pathMap 压根没有 / 的映射
		pathMapBuff["/"] = modal.AndexPath{
			Path:     "/",
			Name:     driverRootPath.Name,
			FileId:   driverRootPath.FileId,
			IsDir:    driverRootPath.IsDir,
			Password: driverRootPath.Password,
		}
		pathMapBuff["/root"] = modal.AndexPath{
			Path:     "/root",
			Name:     driverRootPath.Name,
			FileId:   driverRootPath.FileId,
			IsDir:    driverRootPath.IsDir,
			Password: driverRootPath.Password,
		}

		// 根目录 file id 未设置，证明是第一次
		modal.SavePathMap(pathMapBuff)
	}
	return true
}

// 路径是否正确
func IsPathTrue(path string) bool {

	// 查找能否从 pathMap 里面找到这个 path 对应的 fileId
	fmt.Println("请求", path)
	if _, ok := pathMapBuff[path]; !ok {
		fmt.Println("无法从该路径找到对应的 fileId: ", path)
		return false
	}
	return true
}

// 是否为文件
// 返回参数1代表是否是文件， 参数2代表路径是否正确
func IsPathIsFile(path string) bool {
	return !pathMapBuff[path].IsDir
}

// 传入一个文件路径, 获取该文件的详情
// 先通过文件的 parent 路径, 获取其 parent 路径的 fileId
// 然后通过 fileList 的方式, 获取该文件的 fileId
func GetFileDetail(path string) (*FileDetailBean, bool) {
	split := strings.Split(path, "/")
	parentPath := path[0 : len(path)-len(split[len(split)-1])]
	parentPath = FormatPathQuery(parentPath)

	dirDetails, haveDetail := GetPathDetailFromAli(parentPath)

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

func GetFileDownloadUrl(path string) *FileDownMsgBean {
	if _, ok := pathMapBuff[path]; !ok {
		return nil
	}

	fileDetail := pathMapBuff[path]
	if !fileDetail.IsDir {
		url := api.GetDownloadUrlByFileIdAndFileName(
			SysConfigNow.Authorization,
			SysConfigNow.DriveID,
			fileDetail.FileId,
			fileDetail.Name)
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
func GetPathDetailFromAli(path string) ([]PathDetailBean, bool) {
	folderList := api.FolderList(
		SysConfigNow.Authorization,
		SysConfigNow.DriveID,
		pathMapBuff[path].FileId)
	if folderList == nil || folderList.Items == nil {
		fmt.Println("路径:", path, "获取 folderList 失败")
		return nil, false
	}

	result := make([]PathDetailBean, 0)

	//pathMapUpdateFlag := false
	// 需要被更新的 path 的 list
	listNeedToSave := make([]modal.AndexPath, 0)

	// 更新 pathMap 缓存, 刷新到本地
	for _, fileMsgBean := range folderList.Items {

		filePath := path + "/" + fileMsgBean.Name
		filePath = strings.ReplaceAll(filePath, "//", "/")

		// 更新缓存, 无论是文件还是文件夹都会缓存下来
		if _, ok := pathMapBuff[filePath]; ok {
			if pathMapBuff[filePath].FileId != fileMsgBean.FileID ||
				pathMapBuff[filePath].Name != fileMsgBean.Name ||
				pathMapBuff[filePath].IsDir != (fileMsgBean.Type == "folder") {

				// path 更改
				pathTemp := modal.AndexPath{
					Path:     filePath,
					Name:     fileMsgBean.Name,
					FileId:   fileMsgBean.FileID,
					IsDir:    fileMsgBean.Type == "folder",
					Password: pathMapBuff[filePath].Password, // 新路径继承旧路径的密码
				}
				pathMapBuff[filePath] = pathTemp
				listNeedToSave = append(listNeedToSave, pathTemp)
			}
		} else {
			// path 添加
			// TODO 验证自己父文件夹是否有锁
			pathTemp := modal.AndexPath{
				Path:     filePath,
				Name:     fileMsgBean.Name,
				FileId:   fileMsgBean.FileID,
				IsDir:    fileMsgBean.Type == "folder",
				Password: "",
			}
			pathMapBuff[filePath] = pathTemp
			listNeedToSave = append(listNeedToSave, pathTemp)
		}

		// 构造 PathDetailBean， 此处会返回给页面前端
		result = append(result, PathDetailBean{
			Name:        fileMsgBean.Name,
			Type:        fileMsgBean.Type,
			Path:        filePath,
			CreateTime:  date.Format(fileMsgBean.CreatedAt, "yyyy-MM-dd HH:mm"),
			UpdateTime:  date.Format(fileMsgBean.UpdatedAt, "yyyy-MM-dd HH:mm"),
			DownloadUrl: fileMsgBean.DownloadURL,
			Size:        byteCountBinary(fileMsgBean.Size),
			ParentPath:  "",
			HadPassword: pathMapBuff[filePath].Password != "", // 路径是否需要密码才能访问
			HadAccess:   false,
		})

	}

	if len(listNeedToSave) > 0 {
		// 将缓存刷新到本地
		// TODO 如果 /a/b -> 映射为 xxx, 之后改名为 /a/bb, 那么 /a/bb 映射为 xxx, 但是 /a/b 将不会删除
		// TODO 清除旧的缓存
		modal.SaveAndexPathList(listNeedToSave)
	}

	// result 排序
	sort.Slice(result, func(i, j int) bool {
		return strings.Compare(result[i].Name, result[j].Name) < 0
	})

	// 构造 PathDetailBean 返回
	return result, true
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

func SetPathPassword(path string, password string) bool {
	if _, ok := pathMapBuff[path]; !ok {
		fmt.Println("密码设置路径不存在", path)
		return false
	}
	// 密码设置是根据 fileId 设置，同个 fileId 的多个路径
	paths, err := modal.GetAllSubPathsByFileId(pathMapBuff[path].FileId)
	if err != nil {
		fmt.Println("加密路径时候获取 AndexPath 失败", err)
		return false
	} else {
		newPaths := make([]modal.AndexPath, 0)
		for _, pathTemp := range paths {
			fmt.Println("密码设置路径:", pathTemp.Path, ", PW:", password)
			pathTemp.Password = password
			newPaths = append(newPaths, pathTemp)
		}
		err := modal.SaveAndexPathList(newPaths)
		if err != nil {
			fmt.Println("加密路径数据保存失败", err)
			return false
		} else {
			// 更新到 map 缓存
			for _, pathTemp := range newPaths {
				pathMapBuff[pathTemp.Path] = pathTemp
			}
			return true
		}
	}
}

func GetAndexPath(path string) *modal.AndexPath {
	if _, ok := pathMapBuff[path]; ok {
		path := pathMapBuff[path]
		return &path
	}
	return nil
}

func GetReadmeText() (string, bool) {
	if ReadmeText == "" {
		return "", false
	} else {
		return ReadmeText, true
	}
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
