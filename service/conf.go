package service

import (
	"encoding/json"
	"fmt"
	"github.com/Ericwyn/Andex/api"
	"github.com/Ericwyn/Andex/modal"
	"github.com/Ericwyn/Andex/storage"
	"github.com/Ericwyn/GoTools/file"
	"strings"
	"time"
)

const AndexServerVersion = "V1.1"

const UserConfigFilePath = "config.json"

const UserReadmeFilePath = "README.md"

type ConfigKey string

// Andex 系统配置， 存储在 .conf 文件夹内
type AndexSysConf struct {
	RefreshToken     string
	Authorization    string
	DriveID          string
	RootPath         string
	LastGetTokenTime time.Time
}

// Andex 用户配置, 存储在 ./
type AndexUserConf struct {
	RefreshToken  string `json:"refresh_token"`
	RootPath      string `json:"root_path"`
	Port          string `json:"port"`
	AdminPassword string `json:"admin_password"`
	SiteName      string `json:"site_name"`
}

var UserConfNow *AndexUserConf

// 默认的用户配置模板
var defaultUserConf = AndexUserConf{
	RefreshToken:  "NULL",
	RootPath:      "/",
	Port:          "8080",
	AdminPassword: "",
	SiteName:      "Andex云盘",
}

var SysConfigNow = AndexSysConf{
	RefreshToken:  "NULL",
	Authorization: "NULL",
	DriveID:       "NULL",
	RootPath:      "/",
}

var ReadmeText = ""

func readUserConfigFromFile() (*AndexUserConf, error) {
	// 载入用户配置
	logFile, err := storage.ReadFileAsString(UserConfigFilePath)

	if err != nil {
		panic(err)
	}

	var userConfNow *AndexUserConf

	err = json.Unmarshal([]byte(logFile), &userConfNow)
	if err != nil {
		fmt.Println("读取用户配置文件错误", err)
		return nil, err
	}

	return userConfNow, nil
}

// 载入配置, 程序启动时候调用
func LoadConfFromFile() {
	var err error
	UserConfNow, err = readUserConfigFromFile()
	if err != nil {
		//fmt.Println("")
		return
	}

	SysConfigNow.RefreshToken = modal.GetConfig(modal.TypeRefreshToken, "NULL").(string)
	SysConfigNow.Authorization = modal.GetConfig(modal.TypeAuthorization, "NULL").(string)
	SysConfigNow.DriveID = modal.GetConfig(modal.TypeDriveID, "NULL").(string)
	SysConfigNow.RootPath = modal.GetConfig(modal.TypeRootPath, "/").(string)
	SysConfigNow.LastGetTokenTime = modal.GetConfig(modal.TypeLastGetTokenTime, time.Now()).(time.Time)

	if SysConfigNow.RefreshToken == "NULL" {
		fmt.Println("首次创建 sys conf")
		// 首次创建 sys conf，需要将 user config 保存到 db 里面

		SysConfigNow.RefreshToken = UserConfNow.RefreshToken
		SysConfigNow.RootPath = UserConfNow.RootPath

		modal.SaveConf(modal.TypeRefreshToken, UserConfNow.RefreshToken)
		modal.SaveConf(modal.TypeRootPath, UserConfNow.RootPath)

	} else {
		var err error
		if SysConfigNow.RefreshToken != UserConfNow.RefreshToken {
			SysConfigNow.RefreshToken = UserConfNow.RefreshToken
			fmt.Println("同步 user 配置中 RefreshToken 至 sys 配置")
			err = modal.SaveConf(modal.TypeRefreshToken, SysConfigNow.RefreshToken)
		}
		if SysConfigNow.RootPath != UserConfNow.RootPath {
			SysConfigNow.RootPath = UserConfNow.RootPath
			fmt.Println("同步 user 配置中 RootPath 至 sys 配置")
			err = modal.SaveConf(modal.TypeRootPath, UserConfNow.RootPath)
		}
		if err != nil {
			fmt.Println("同步 user 配置至 sys 时发生错误", err)
		}
	}
}

func LoadReadmeFile() {
	readmeFile := file.OpenFile(UserReadmeFilePath)
	if readmeFile.Exits() {
		var err error
		ReadmeText, err = storage.ReadFileAsString(UserReadmeFilePath)
		if err != nil {
			fmt.Println("载入 README 文件失败", err)
		}
	}
}

// 保存 authorization/refreshToken/driveId 到 sysConf 和 userConf
func SaveTokenConf(authorization string, refreshToken string, driveId string) {
	SysConfigNow.Authorization = authorization
	modal.SaveConf(modal.TypeAuthorization, authorization)

	SysConfigNow.RefreshToken = refreshToken
	modal.SaveConf(modal.TypeRefreshToken, refreshToken)

	SysConfigNow.DriveID = driveId
	modal.SaveConf(modal.TypeDriveID, driveId)

	SysConfigNow.LastGetTokenTime = time.Now() // 保存上一次 token 获取的时间
	modal.SaveConf(modal.TypeLastGetTokenTime, SysConfigNow.LastGetTokenTime)

	updateRefreshTokenInUserConf(refreshToken)
}

func RefreshToken() {
	api.RefreshToken(SysConfigNow.RefreshToken, func(authorization string, refreshToken string, driveId string) {
		// 配置刷新成功
		SaveTokenConf(authorization, refreshToken, driveId)
	})
}

// 刷新用户配置里面的 token
func updateRefreshTokenInUserConf(newToken string) {
	if UserConfNow == nil {
		return
	}
	// 读取 user config 配置文本
	userConfStr, err := storage.ReadFileAsString(UserConfigFilePath)
	if err != nil {
		fmt.Println("用户配置文件更新失败", err)
	} else {
		fmt.Println("用户配置文件已更新")
	}
	userConfStr = strings.Replace(userConfStr, UserConfNow.RefreshToken, newToken, 1)
	storage.WriteStringToFile(UserConfigFilePath, userConfStr, false)
}

// 创建空白的用户配置模板
func CreateUserConfFile() {
	bytes, err := json.MarshalIndent(defaultUserConf, "", "  ")
	if err != nil {
		fmt.Println("序列化 user 配置发生错误", err)
	} else {

		err := storage.WriteStringToFile(UserConfigFilePath, string(bytes), false)
		if err != nil {
			fmt.Println("user 配置文件更新失败", err)
		} else {
			fmt.Println("user 配置文件已更新")
		}
	}
}
