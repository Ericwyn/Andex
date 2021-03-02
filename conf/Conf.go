package conf

import (
	"encoding/json"
	"fmt"
	"github.com/Ericwyn/Andex/fileIO"
)

const ConfigFilePath = "./config.json"

type ConfigKey string

type AndexConf struct {
	RefreshToken  string `json:"refresh_token"`
	Authorization string `json:"authorization"`
	DriveID       string `json:"drive_id"`
}

var ConfigNow = AndexConf{
	RefreshToken:  "NULL",
	Authorization: "NULL",
	DriveID:       "NULL",
}

// 载入配置, 程序启动时候调用
func LoadConfFromFile() {
	//fi, err := os.Open(ConfigFilePath)

	logFile, err := fileIO.ReadFileAsString(ConfigFilePath)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(logFile), &ConfigNow)
	if err != nil {
		fmt.Println("读取配置文件错误", err)
	}
}

func SaveTokenConf(authorization string, refreshToken string, driveId string) {
	ConfigNow.Authorization = authorization
	ConfigNow.RefreshToken = refreshToken
	ConfigNow.DriveID = driveId
	SaveConf()
}

func SaveConf() {
	bytes, err := json.MarshalIndent(ConfigNow, "", "  ")
	if err != nil {
		fmt.Println("序列化配置发生错误", err)
	} else {

		err := fileIO.WriteStringToFile(ConfigFilePath, string(bytes), false)
		if err != nil {
			fmt.Println("配置文件更新失败", err)
		} else {
			fmt.Println("配置文件已更新")
		}
	}
}
