package controller

import (
	"github.com/Ericwyn/Andex/service"
	"github.com/gin-gonic/gin"
)

func pages(ctx *gin.Context) {
	path, hasPathQuery := ctx.GetQuery("p")

	if !hasPathQuery {
		path = "/"
	}

	// 格式化 query 参数
	path = service.FormatPathQuery(path)

	// 构造路径下的文件/文件夹列表
	pathDetail, hasDetail := service.GetPathDetail(path)

	// 获取面包屑参数
	navPathList := service.GetNavPathList(path)

	if hasDetail {
		ctx.HTML(200, "index.html", gin.H{
			"pathDetail":    pathDetail,
			"navPathList":   navPathList,
			"navPathLength": len(navPathList),
		})
	} else {
		ctx.String(200, "没有找到该路径")
	}
}
