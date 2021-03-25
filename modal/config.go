package modal

import (
	"fmt"
	"strconv"
	"time"
)

type ConfigType string

const TypeRefreshToken ConfigType = "refreshToken"
const TypeAuthorization ConfigType = "authorization"
const TypeDriveID ConfigType = "driveID"
const TypeRootPath ConfigType = "rootPath"
const TypeLastGetTokenTime ConfigType = "lastGetTokenTime"

// Andex 配置
type AndexConfig struct {
	Key   ConfigType `xorm:"pk"`
	Value string
}

func LoadConfMap() (map[ConfigType]AndexConfig, error) {
	confList := make([]AndexConfig, 0)
	err := sqlEngine.Find(&confList)
	if err != nil {
		return nil, err
	}
	var resMap = make(map[ConfigType]AndexConfig)
	for _, conf := range confList {
		resMap[conf.Key] = conf
	}
	return resMap, nil
}

func SaveConf(configType ConfigType, value interface{}) error {
	var err error
	if inList(configType,
		[]ConfigType{TypeRefreshToken, TypeAuthorization, TypeDriveID, TypeRootPath}) {
		// 字符串配置保存
		err = saveAndexConf(configType, value.(string))
	} else if inList(configType, []ConfigType{TypeLastGetTokenTime}) {
		// 日期配置保存
		err = saveAndexConf(configType, strconv.FormatInt(value.(time.Time).Unix(), 10))
	}
	return err
}

func GetConfig(configType ConfigType, defValue interface{}) interface{} {
	confValue := getAndexConf(configType, "can_not_get_config_value")
	if confValue == "can_not_get_config_value" {
		return defValue
	}
	if inList(configType,
		[]ConfigType{TypeRefreshToken, TypeAuthorization, TypeDriveID, TypeRootPath}) {
		// 字符串配置保存
		return confValue
	} else if inList(configType, []ConfigType{TypeLastGetTokenTime}) {
		// 日期配置保存
		// 日期配置读取
		unixTime, err := strconv.ParseInt(confValue, 10, 64)
		if err != nil {
			fmt.Println("日期配置读取错误", err)
			return defValue
		}
		return time.Unix(unixTime, 0)
	}
	return defValue
}

func saveAndexConf(key ConfigType, value string) error {
	var sql = "INSERT OR REPLACE INTO `andex_config` (`key`, `value`) VALUES (?,?)"

	_, err := sqlEngine.Exec(sql, key, value)
	if err != nil {
		return err
	}
	return nil
}

func getAndexConf(key ConfigType, def string) string {
	conf := AndexConfig{
		Key: key,
	}

	has, err := sqlEngine.Get(&conf)

	if !has || err != nil {
		return def
	} else {
		return conf.Value
	}
}

func inList(configType ConfigType, configTypeList []ConfigType) bool {
	for _, typeTemp := range configTypeList {
		if configType == typeTemp {
			return true
		}
	}
	return false
}
