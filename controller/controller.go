package controller

import (
	"github.com/Ericwyn/Andex/service"
	"github.com/gin-gonic/gin"
)

// 文件/文件夹页面获取接口
func pages(ctx *gin.Context) {
	path, hasPathQuery := ctx.GetQuery("p")

	if !hasPathQuery {
		path = "/"
	}

	// 格式化 query 参数
	path = service.FormatPathQuery(path)

	// 判断是否是文件
	if !service.IsPathTrue(path) {
		ctx.String(200, "没有找到该路径")
		return
	}

	// 获取面包屑参数
	navPathList := service.GetNavPathList(path)

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
				"fileDetail": fileDetail,
				"navPath":    navPath, // 父路径
			})
			return
		} else {
			ctx.String(200, "获取文件详情失败")
			return
		}
	} else {
		// 构造路径下的文件/文件夹列表
		pathDetail, hasDetail := service.GetPathDetail(path)
		if hasDetail {

			ctx.HTML(200, "folder.html", gin.H{
				"pathDetail":    pathDetail,
				"navPathList":   navPathList,
				"navPathLength": len(navPathList),
			})
			return
		} else {
			ctx.String(200, "获取文件夹详情失败")
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
		ctx.String(200, "没有找到该路径")
		return
	}

	if service.IsPathIsFile(path) {

		fileDownMsgBean := service.GetFileDownloadUrl(path)

		if fileDownMsgBean != nil {
			ctx.Redirect(302, fileDownMsgBean.Url)
			return
		} else {
			ctx.String(200, "文件下载地址获取失败")
			return
		}

	} else {
		ctx.String(200, "非文件路径")
		return
	}

}
