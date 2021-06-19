package controller

import (
	"fmt"
	"github.com/Ericwyn/Andex/service"
	"github.com/Ericwyn/Andex/util/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func apiStatic(path string, ctx *gin.Context) {
	ctx.File("./static" + path)
}

type loginBody struct {
	Password string `json:"password"`
}

// 文件/文件夹页面获取接口
func apiPages(path string, ctx *gin.Context) {
	// 格式化 query 参数
	path = service.FormatPathQuery(path)

	// 判断访问路径是否正确
	if !service.IsPathTrue(path) {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "你访问的页面不存在, 或者路径未缓存",
			"andexVersion": service.AndexServerVersion,
			"siteName":     service.UserConfNow.SiteName,
		})
		return
	}

	session := sessions.Default(ctx)
	// TODO 验证用户对路径的访问权限
	permFlag, err := checkUserPathPerm(path, session)
	if err != nil {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "你访问的页面不存在, 或者路径未缓存",
			"andexVersion": service.AndexServerVersion,
			"siteName":     service.UserConfNow.SiteName,
		})
		return
	}
	if !permFlag {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "您无权访问该页面",
			"andexVersion": service.AndexServerVersion,
			"siteName":     service.UserConfNow.SiteName,
		})
		return
	}

	// 获取面包屑参数
	navPathList := service.GetNavPathList(path)

	if service.IsPathIsFile(path) {
		// 文件页面打开
		openFilePage(navPathList, path, ctx)
	} else {
		// 文件夹页面打开
		openDirPage(navPathList, session, path, ctx)
	}
}

// 检查用户对当前路径是否有权限打开
func checkUserPathPerm(path string, session sessions.Session) (bool, error) {
	// 检查该路径是否加密
	andexPath := service.GetAndexPath(path)
	if andexPath == nil {
		// 路径不存在，不给访问
		return false, fmt.Errorf("路径不存在")
	}

	if andexPath.Password == "" {
		return true, nil
	}
	// 管理员登录
	if session.Get("hadLogin") != nil && session.Get("hadLogin").(bool) {
		return true, nil
	}

	// 不然的话检查
	pathPerm := session.Get("pathPerm")
	if pathPerm != nil {
		pathPermString := pathPerm.(string)
		// 获取该用户有权限访问的路径列表
		pathPermArr := strings.Split(pathPermString, permCookieStrSplit)
		// 如果其中一个文件夹是当前路径的父路径，且其密码与当前文件夹密码一致的话，就可以访问
		for _, pathHasPerm := range pathPermArr {
			if pathHasPerm == "" {
				continue
			}
			if pathHasPerm == path {
				return true, nil
			}

			subStr := pathHasPerm
			if pathHasPerm != "/" {
				subStr += "/"
			}
			if strings.Index(path, subStr) == 0 {
				andexPathHasPerm := service.GetAndexPath(pathHasPerm)
				if andexPathHasPerm.Password == andexPath.Password {
					return true, nil
				}
			}
		}
		return false, nil
	}
	return false, nil
}

func openFilePage(navPathList []service.NavPath, path string, ctx *gin.Context) {
	startTime := time.Now()

	// 文件路径访问
	fileDetail, _ := service.GetFileDetail(path)

	var navPath service.NavPath
	if len(navPathList) >= 2 {
		navPath = navPathList[len(navPathList)-2]
	} else {
		navPath = service.NavPath{
			Name: "首页",
			Path: "/",
			// 判断是否是首页请求
			Last: false,
		}
	}
	if fileDetail != nil {
		ctx.HTML(200, "file.html", gin.H{
			"fileDetail":     fileDetail,
			"navPath":        navPath, // 父路径
			"apiRequestTime": fmt.Sprint(1.0*(time.Now().UnixNano()-startTime.UnixNano())/1000000, "ms"),
			"andexVersion":   service.AndexServerVersion,
			"siteName":       service.UserConfNow.SiteName,
			//"hadLogin": hadLogin,
		})
		return
	} else {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "获取文件详情失败了",
			"andexVersion": service.AndexServerVersion,
			"siteName":     service.UserConfNow.SiteName,
		})
		return
	}
}

func openDirPage(navPathList []service.NavPath, session sessions.Session, path string, ctx *gin.Context) {
	startTime := time.Now()

	templateName := "folder.html"
	hadLogin := false
	if session.Get("hadLogin") != nil && session.Get("hadLogin").(bool) {
		//fmt.Println("用户已登录")
		hadLogin = true
		templateName = "adminFolder.html"
	}

	// 文件夹路径访问
	// 构造路径下的文件/文件夹列表
	pathDetailList, hasDetail := service.GetPathDetailFromAli(path)
	if hasDetail {
		newPathDetailList := make([]service.PathDetailBean, 0)

		// 权限过滤
		// 获取该用户可访问的文件夹路径
		var pathPermList = make([]string, 0)
		pathPerm := session.Get("pathPerm")
		if pathPerm != nil {
			pathPermList = strings.Split(pathPerm.(string), permCookieStrSplit)
		}
		for _, pathDetail := range pathDetailList {
			// 如果这个 pathDetail 需要密码访问的话
			if pathDetail.HadPassword && len(pathPermList) > 0 {
				for _, pathPermTemp := range pathPermList {
					if pathPermTemp == pathDetail.Path {
						pathDetail.HadAccess = true
					}
					// 如果 pathPermTemp 是 pathDetail.Path 的父文件夹且密码一致的话，也是可以访问的
					parentPath := pathPermTemp
					if parentPath != "/" {
						parentPath += "/"
					}
					if strings.Index(pathDetail.Path, parentPath) == 0 {
						andexPath1 := service.GetAndexPath(pathDetail.Path)
						andexPath2 := service.GetAndexPath(pathPermTemp)
						if andexPath1 != nil && andexPath2 != nil && andexPath1.Password == andexPath2.Password {
							pathDetail.HadAccess = true
						}
					}
				}
			}

			newPathDetailList = append(newPathDetailList, pathDetail)
		}

		var readmeText string = ""
		var hasReadme bool = false
		if path == "/" || path == "root" {
			readmeText, hasReadme = service.GetReadmeText()
		}

		ctx.HTML(200, templateName, gin.H{
			"pathDetail":     newPathDetailList,
			"navPathList":    navPathList,
			"navPathLength":  len(navPathList),
			"apiRequestTime": fmt.Sprint(1.0*(time.Now().UnixNano()-startTime.UnixNano())/1000000, "ms"),
			"readme":         readmeText,
			"hasReadme":      hasReadme,
			"andexVersion":   service.AndexServerVersion,
			"siteName":       service.UserConfNow.SiteName,
			"hadLogin":       hadLogin,
		})

		return
	} else {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "获取文件夹详情失败",
			"andexVersion": service.AndexServerVersion,
			"siteName":     service.UserConfNow.SiteName,
		})
		return
	}

}

func apiPathPermRequest(ctx *gin.Context) {
	var body passwordBody
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": RestApiParamError,
			"msg":  "参数错误",
		})
		log.E("登录参数错误", err)
		return
	}

	body.Path = service.FormatPathQuery(body.Path)

	path := service.GetAndexPath(body.Path)
	if path == nil {
		ctx.JSON(200, gin.H{
			"code": RestApiParamError,
			"msg":  "访问路径错误",
		})
		log.E("访问路径错误", body.Path)
		return
	}

	if path.Password == body.Password {
		err := addPathPerm(body.Path, ctx)

		if err != nil {
			ctx.JSON(200, gin.H{
				"code": RestApiServerError,
				"msg":  "访问授权保存失败",
			})
			return
		}

		ctx.JSON(200, gin.H{
			"code": RestApiSuccess,
			"msg":  "路径访问授权成功: " + body.Path,
		})
		return
	} else {
		ctx.JSON(200, gin.H{
			"code": RestApiAuthorizationError,
			"msg":  "访问密码错误",
		})
		return
	}
}

const permCookieStrSplit = "${AND}$"

func addPathPerm(path string, ctx *gin.Context) error {
	session := sessions.Default(ctx)
	pathPermOld := session.Get("pathPerm")
	var err error = nil
	if pathPermOld == nil {
		session.Set("pathPerm", path)
		session.Options(sessions.Options{MaxAge: 60 * 60})
		err = session.Save()
	} else {
		pathPermOldStr := pathPermOld.(string)
		splitList := strings.Split(pathPermOldStr, "")

		alreadyAddFlag := false
		for _, pathOld := range splitList {
			if pathOld == path {
				alreadyAddFlag = true
			}
		}
		if !alreadyAddFlag {
			session.Set("pathPerm", pathPermOld.(string)+permCookieStrSplit+path)
			session.Options(sessions.Options{MaxAge: 60 * 60})
			err = session.Save()
		}
	}
	if err != nil {
		return err
	} else {
		return nil
	}
}

// 文件下载接口, /download?p=/a/v/c.pdf
func apiDownload(path string, ctx *gin.Context) {

	// 格式化 query 参数
	path = service.FormatPathQuery(path)

	// 权限校验
	session := sessions.Default(ctx)

	permFlag, err := checkUserPathPerm(path, session)
	if err != nil {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "你访问的页面不存在, 或者路径未缓存",
			"andexVersion": service.AndexServerVersion,
			"siteName":     service.UserConfNow.SiteName,
		})
		return
	}
	if !permFlag {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "您无权限访问该页面",
			"andexVersion": service.AndexServerVersion,
			"siteName":     service.UserConfNow.SiteName,
		})
		return
	}

	// 判断是否是文件
	if !service.IsPathTrue(path) {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "你访问的文件路径不存在, 或路径未缓存",
			"andexVersion": service.AndexServerVersion,
			"siteName":     service.UserConfNow.SiteName,
		})
		return
	}

	if service.IsPathIsFile(path) {

		fileDownMsgBean := service.GetFileDownloadUrl(path)

		if fileDownMsgBean != nil {
			ctx.Redirect(302, fileDownMsgBean.Url)
			return
		} else {
			ctx.HTML(200, "error.html", gin.H{
				"errorNote":    "文件下载地址获取失败",
				"andexVersion": service.AndexServerVersion,
				"siteName":     service.UserConfNow.SiteName,
			})
			return
		}
	} else {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "该路径不是可下载文件",
			"andexVersion": service.AndexServerVersion,
			"siteName":     service.UserConfNow.SiteName,
		})
		return
	}
}

// 获取加密路径的直接下载链接
func apiGetRedirectLink(ctx *gin.Context) {
	var body passwordBody
	err := ctx.BindJSON(&body)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": RestApiServerError,
			"msg":  "服务器错误",
		})
		log.E("服务器错误", err)
		return
	}

	// 格式化 query 参数
	body.Path = service.FormatPathQuery(body.Path)

	// 判断访问路径是否正确
	if !service.IsPathTrue(body.Path) {
		ctx.JSON(200, gin.H{
			"code": RestApiParamError,
			"msg":  "路径错误",
		})
		return
	}

	// 权限校验
	session := sessions.Default(ctx)

	permFlag, err := checkUserPathPerm(body.Path, session)
	if err != nil {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "你访问的页面不存在, 或者路径未缓存",
			"andexVersion": service.AndexServerVersion,
			"siteName":     service.UserConfNow.SiteName,
		})
		return
	}
	if !permFlag {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "您无权限访问该页面",
			"andexVersion": service.AndexServerVersion,
			"siteName":     service.UserConfNow.SiteName,
		})
		return
	}

	// 直链获取
	url := service.GetFileDownloadUrl(body.Path)
	if url == nil {
		ctx.JSON(200, gin.H{
			"code": RestApiServerError,
			"msg":  "无法获取下载直链",
		})
		return
	} else {
		ctx.JSON(200, gin.H{
			"code": RestApiSuccess,
			"msg":  url,
		})
		return
	}

}
