package controller

import (
	"github.com/Ericwyn/Andex/service"
	"github.com/gin-gonic/gin"
	"strings"
)

func pages(ctx *gin.Context) {
	path, hasPathQuery := ctx.GetQuery("p")

	if !hasPathQuery {
		path = "/"
	}
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

	// 构造路径下的文件/文件夹列表
	pathDetail, hasDetail := service.GetPathDetail(path)

	// 构造路径面包屑
	type navPath struct {
		Name string
		Path string
		Last bool
	}

	split := strings.Split(path, "/")
	navPathList := make([]navPath, 0)
	navPathList = append(navPathList, navPath{
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
		navPathList = append(navPathList, navPath{
			Name: navPathTemp,
			Path: navPathNow,
			Last: i >= (len(split) - 1),
		})
	}

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
