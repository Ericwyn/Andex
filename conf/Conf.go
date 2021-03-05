package conf

import (
	"encoding/json"
	"fmt"
	"github.com/Ericwyn/Andex/storage"
	"github.com/Ericwyn/GoTools/file"
	"strings"
	"time"
)

const UserConfigFilePath = "config.json"

const UserReadmeFilePath = "README.md"

const SysConfigDirPath = ".conf"
const SysConfigFilePath = SysConfigDirPath + "/" + "config.json"

type ConfigKey string

// Andex 系统配置， 存储在 .conf 文件夹内
type SystemSysConf struct {
	RefreshToken     string    `json:"refresh_token"`
	Authorization    string    `json:"authorization"`
	DriveID          string    `json:"drive_id"`
	RootPath         string    `json:"root_path"`
	LastGetTokenTime time.Time `json:"last_get_token_time"`
}

// Andex 用户配置, 存储在 ./
type AndexUserConf struct {
	RefreshToken string `json:"refresh_token"`
	RootPath     string `json:"root_path"`
}

var userConfNow *AndexUserConf

var SysConfigNow = SystemSysConf{
	RefreshToken:  "NULL",
	Authorization: "NULL",
	DriveID:       "NULL",
	RootPath:      "/",
}

var ReadmeText = ""

// 载入配置, 程序启动时候调用
func LoadConfFromFile() {
	//fi, err := os.Open(UserConfigFilePath)

	// 载入用户配置
	logFile, err := storage.ReadFileAsString(UserConfigFilePath)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(logFile), &userConfNow)
	if err != nil {
		fmt.Println("读取用户配置文件错误", err)
	}

	// 载入 sys 配置
	systemConfFile := file.OpenFile(SysConfigFilePath)
	sysConfExit := false
	if systemConfFile.Exits() {
		sysConfExit = true
		// 如果文件存在的话就读取
		sysLogFile, err := storage.ReadFileAsString(SysConfigFilePath)
		if err != nil {
			panic(err)
		}

		err = json.Unmarshal([]byte(sysLogFile), &SysConfigNow)
		if err != nil {
			fmt.Println("读取 sys 配置文件错误", err)
		}
	}

	if !sysConfExit {
		fmt.Println("首次创建 sys conf")
		SaveSysConf()
	} else {
		saveFlag := false
		if SysConfigNow.RefreshToken != userConfNow.RefreshToken {
			SysConfigNow.RefreshToken = userConfNow.RefreshToken
			saveFlag = true
		}
		if SysConfigNow.RootPath != userConfNow.RootPath {
			SysConfigNow.RootPath = userConfNow.RootPath
			saveFlag = true
		}
		if saveFlag {
			fmt.Println("同步 user 配置中 RefreshToken 至 sys 配置")
			SaveSysConf()
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
	SysConfigNow.RefreshToken = refreshToken
	SysConfigNow.DriveID = driveId
	SysConfigNow.LastGetTokenTime = time.Now() // 保存上一次 token 获取的时间

	SaveSysConf()
	UpdateRefreshTokenInUserConf(refreshToken)
}

func SaveSysConf() {
	bytes, err := json.MarshalIndent(SysConfigNow, "", "  ")
	if err != nil {
		fmt.Println("序列化 sys 配置发生错误", err)
	} else {

		err := storage.WriteStringToFile(SysConfigFilePath, string(bytes), false)
		if err != nil {
			fmt.Println("sys 配置文件更新失败", err)
		} else {
			fmt.Println("sys 配置文件已更新")
		}
	}
}

// 刷新用户配置里面的 token
func UpdateRefreshTokenInUserConf(newToken string) {
	if userConfNow == nil {
		return
	}
	// 读取 user config 配置文本
	userConfStr, err := storage.ReadFileAsString(UserConfigFilePath)
	if err != nil {
		fmt.Println("用户配置文件更新失败", err)
	} else {
		fmt.Println("用户配置文件已更新")
	}
	userConfStr = strings.Replace(userConfStr, userConfNow.RefreshToken, newToken, 1)
	storage.WriteStringToFile(UserConfigFilePath, userConfStr, false)
}

// 创建空白的用户配置模板
func CreateUserConfFile() {
	bytes, err := json.MarshalIndent(AndexUserConf{
		RefreshToken: "NULL",
		RootPath:     "/",
	}, "", "  ")
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
