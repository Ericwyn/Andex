package controller

import (
	"fmt"
	"github.com/Ericwyn/Andex/service"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// 文件/文件夹页面获取接口
func pages(ctx *gin.Context) {
	path, hasPathQuery := ctx.GetQuery("p")

	if !hasPathQuery {
		path = "/"
	}

	// 格式化 query 参数
	path = service.FormatPathQuery(path)

	// 判断访问路径是否正确
	if !service.IsPathTrue(path) {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote": "你访问的页面不存在, 或者路径未缓存",
		})
		return
	}

	// 获取面包屑参数
	navPathList := service.GetNavPathList(path)

	startTime := time.Now()

	if service.IsPathIsFile(path) {
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
			})
			return
		} else {
			ctx.HTML(200, "error.html", gin.H{
				"errorNote": "获取文件详情失败了",
			})
			return
		}
	} else {
		// 构造路径下的文件/文件夹列表
		pathDetail, hasDetail := service.GetPathDetail(path)
		if hasDetail {

			ctx.HTML(200, "folder.html", gin.H{
				"pathDetail":     pathDetail,
				"navPathList":    navPathList,
				"navPathLength":  len(navPathList),
				"isMobil":        isFromMobile(ctx.GetHeader("User-Agent")),
				"apiRequestTime": fmt.Sprint(1.0*(time.Now().UnixNano()-startTime.UnixNano())/1000000, "ms"),
			})
			return
		} else {
			ctx.HTML(200, "error.html", gin.H{
				"errorNote": "获取文件夹详情失败",
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
			"errorNote": "你访问的文件路径不存在, 或路径未缓存",
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
				"errorNote": "文件下载地址获取失败",
			})
			return
		}
	} else {
		ctx.HTML(200, "error.html", gin.H{
			"errorNote": "该路径不是可下载文件",
		})
		return
	}
}

func isFromMobile(userAgent string) bool {
	if len(userAgent) == 0 {
		return false
	}

	isMobile := false
	mobileKeywords := []string{"Mobile", "Android", "Silk/", "Kindle",
		"BlackBerry", "Opera Mini", "Opera Mobi"}

	for i := 0; i < len(mobileKeywords); i++ {
		if strings.Contains(userAgent, mobileKeywords[i]) {
			isMobile = true
			break
		}
	}
	return isMobile
}
