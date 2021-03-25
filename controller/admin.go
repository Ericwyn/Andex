package controller

import (
	"fmt"
	"github.com/Ericwyn/Andex/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"strings"
)

const RestApiParamError = "4000"
const RestApiAuthorizationError = "4001"
const RestApiServerError = "4003"
const RestApiSuccess = "1000"

func adminLogin(ctx *gin.Context) {
	var loginBody loginBody
	err := ctx.BindJSON(&loginBody)

	if err != nil {
		ctx.JSON(200, gin.H{
			"code": RestApiParamError,
			"msg":  "参数错误",
		})
		fmt.Println("登录参数错误", err)
		return
	} else {
		if service.UserConfNow.AdminPassword == "" {
			ctx.JSON(200, gin.H{
				"code": RestApiParamError,
				"msg":  "参数错误",
			})
			fmt.Println("管理员密码未设置", err)
			return
		}
		if loginBody.Password == service.UserConfNow.AdminPassword {
			session := sessions.Default(ctx)

			// 设置session数据
			session.Set("hadLogin", true)
			// 保存session数据
			err := session.Save()

			if err != nil {
				ctx.JSON(200, gin.H{
					"code": RestApiServerError,
					"msg":  "服务器错误",
				})
				fmt.Println("服务器错误", err)
				return
			} else {
				ctx.JSON(200, gin.H{
					"code": RestApiSuccess,
					"msg":  "登录成功",
				})
				return
			}
		} else {
			ctx.JSON(200, gin.H{
				"code": RestApiAuthorizationError,
				"msg":  "密码错误",
			})
			return
		}
	}
}

func adminLogout(ctx *gin.Context) {
	session := sessions.Default(ctx)

	if session.Get("hadLogin") != nil && session.Get("hadLogin").(bool) {
		//fmt.Println("用户已登录")
		session.Delete("hadLogin")
		err := session.Save()

		if err != nil {
			ctx.JSON(200, gin.H{
				"code": RestApiServerError,
				"msg":  "服务器错误",
			})
			fmt.Println("服务器错误", err)
			return
		} else {
			ctx.JSON(200, gin.H{
				"code": RestApiSuccess,
				"msg":  "操作成功",
			})
			return
		}
	}
}

type passwordBody struct {
	Path string `json:"path"`
	Password string `json:"password"`
}

// 加密某个文件夹，会将该文件夹及其子文件夹全部加密
func setPassword(ctx *gin.Context) {
	if checkLogin(ctx){
		var body passwordBody
		err := ctx.BindJSON(&body)
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": RestApiServerError,
				"msg":  "服务器错误",
			})
			fmt.Println("服务器错误", err)
			return
		}

		body.Path = service.FormatPathQuery(body.Path)
		if !service.IsPathTrue(body.Path) || strings.Trim(body.Password," ") == "" {
			ctx.JSON(200, gin.H{
				"code": RestApiParamError,
				"msg":  "参数错误",
			})
			fmt.Println("服务器错误", err)
			return
		}

		service.SetPathPassword(body.Path, body.Password)

		ctx.JSON(200, gin.H{
			"code": "1000",
			"msg":  "操作成功",
		})
		return
	}
}

// 解密文件夹，会将该文件夹及其所有子文件夹全部解密
func removePassword(ctx *gin.Context) {
	if checkLogin(ctx){
		var body passwordBody
		err := ctx.BindJSON(&body)
		if err != nil {
			ctx.JSON(200, gin.H{
				"code": RestApiServerError,
				"msg":  "服务器错误",
			})
			fmt.Println("服务器错误", err)
			return
		}

		body.Path = service.FormatPathQuery(body.Path)
		if !service.IsPathTrue(body.Path) {
			ctx.JSON(200, gin.H{
				"code": RestApiParamError,
				"msg":  "参数错误",
			})
			fmt.Println("服务器错误", err)
			return
		}

		// 将密码设置为 "" 就是去除解密了
		service.SetPathPassword(body.Path, "")

		ctx.JSON(200, gin.H{
			"code": "1000",
			"msg":  "操作成功",
		})
		return
	}
}


func checkLogin(ctx *gin.Context) bool {
	session := sessions.Default(ctx)

	if session.Get("hadLogin") != nil && session.Get("hadLogin").(bool) {
		return true
	} else {
		ctx.JSON(200, gin.H{
			"code": RestApiAuthorizationError,
			"msg":  "用户未登录",
		})
		return false
	}

}