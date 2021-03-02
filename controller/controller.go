package controller

import (
	"github.com/Ericwyn/Andex/service"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

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

	if service.IsPathIsFile(path) {
		fileDetail := service.GetFileDetail(path)
		if fileDetail != nil {
			jsonString, _ := jsoniter.MarshalIndent(fileDetail, "", "  ")
			str := string(jsonString)
			ctx.String(200, "获取文件成功:\n"+str)
		} else {
			ctx.String(200, "获取文件详情失败")
		}
	} else {
		// 构造路径下的文件/文件夹列表
		pathDetail := service.GetPathDetail(path)

		// 获取面包屑参数
		navPathList := service.GetNavPathList(path)

		ctx.HTML(200, "folder.html", gin.H{
			"pathDetail":    pathDetail,
			"navPathList":   navPathList,
			"navPathLength": len(navPathList),
		})
	}
}
