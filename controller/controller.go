package controller

import (
	"fmt"
	"github.com/Ericwyn/Andex/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

type loginBody struct {
	Password string `json:"password"`
}



// 文件/文件夹页面获取接口
func pages(ctx *gin.Context) {
	path, hasPathQuery := ctx.GetQuery("p")

	session := sessions.Default(ctx)

	templateName := "folder.html"
	hadLogin := false
	if session.Get("hadLogin") != nil && session.Get("hadLogin").(bool) {
		//fmt.Println("用户已登录")
		hadLogin = true
		templateName = "adminFolder.html"
	}

	if !hasPathQuery {
		path = "/"
	}

	// 格式化 query 参数
	path = service.FormatPathQuery(path)

	// 判断访问路径是否正确
	if !service.IsPathTrue(path) {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "你访问的页面不存在, 或者路径未缓存",
			"andexVersion": service.AndexServerVersion,
		})
		return
	}

	// 获取面包屑参数
	navPathList := service.GetNavPathList(path)

	startTime := time.Now()

	if service.IsPathIsFile(path) {
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
				//"hadLogin": hadLogin,
			})
			return
		} else {
			ctx.HTML(200, "error.html", gin.H{
				"errorNote":    "获取文件详情失败了",
				"andexVersion": service.AndexServerVersion,
			})
			return
		}
	} else {
		// 文件夹路径访问
		// 构造路径下的文件/文件夹列表
		pathDetail, hasDetail := service.GetPathDetail(path)
		if hasDetail {

			var readmeText string = ""
			var hasReadme bool = false
			if path == "/" || path == "root" {
				readmeText, hasReadme = service.GetReadmeText()
			}

			ctx.HTML(200, templateName, gin.H{
				"pathDetail":     pathDetail,
				"navPathList":    navPathList,
				"navPathLength":  len(navPathList),
				"apiRequestTime": fmt.Sprint(1.0*(time.Now().UnixNano()-startTime.UnixNano())/1000000, "ms"),
				"readme":         readmeText,
				"hasReadme":      hasReadme,
				"andexVersion":   service.AndexServerVersion,
				"hadLogin":       hadLogin,
			})

			return
		} else {
			ctx.HTML(200, "error.html", gin.H{
				"errorNote":    "获取文件夹详情失败",
				"andexVersion": service.AndexServerVersion,
			})
			return
		}
	}
}

// 文件下载接口, /download?p=/a/v/c.pdf
func download(ctx *gin.Context) {
	path, hasPathQuery := ctx.GetQuery("p")

	if !hasPathQuery {
		path = "/"
	}

	// 格式化 query 参数
	path = service.FormatPathQuery(path)

	// 判断是否是文件
	if !service.IsPathTrue(path) {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "你访问的文件路径不存在, 或路径未缓存",
			"andexVersion": service.AndexServerVersion,
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
			})
			return
		}
	} else {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote":    "该路径不是可下载文件",
			"andexVersion": service.AndexServerVersion,
		})
		return
	}
}
