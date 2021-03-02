package controller

import (
	"github.com/Ericwyn/Andex/service"
	"github.com/gin-gonic/gin"
)

func pages(ctx *gin.Context) {
	query, hasPathQuery := ctx.GetQuery("p")
	var pathDetail []service.PathDetailBean
	var hasDetail bool
	if hasPathQuery {
		pathDetail, hasDetail = service.GetPathDetail(query)
	} else {
		pathDetail, hasDetail = service.GetPathDetail("/")
	}

	if hasDetail {
		ctx.HTML(200, "index.html", gin.H{
			"pathDetail": pathDetail,
		})
	} else {
		ctx.String(200, "没有找到该路径")
	}
}
