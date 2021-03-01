package conf

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const ConfigFilePath = "./config.json"

type ConfigKey string

type AndexConf struct {
	RefreshToken  string `json:"refresh_token"`
	Authorization string `json:"authorization"`
	DriveID string `json:"drive_id"`
}

var ConfigNow = AndexConf{
	RefreshToken:  "NULL",
	Authorization: "NULL",
	DriveID: "NULL",
}

// 载入配置, 程序启动时候调用
func LoadConfFromFile() {

	logFile := ""
	fi, err := os.Open(ConfigFilePath)
	if err != nil {
		panic(err)
	}

	defer fi.Close()
	r := bufio.NewReader(fi)

	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
			//return
		}
		if 0 == n {
			break
		} else {
			// 将读取到的数据交给 callback 处理
			logFile += string(buf[:n])
		}
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
		fl, err := os.OpenFile(ConfigFilePath, os.O_CREATE|os.O_WRONLY, 0755)
		fl.Chdir()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer fl.Close()

		n, err := fl.WriteString(string(bytes))
		if err != nil {
			fmt.Println(err.Error())
			//return
		}
		if n < len(string(bytes)) {
			fmt.Println("write byte num error")
			return
		}
	}
	fmt.Println("配置文件已更新")
}
