package api

import (
	"encoding/json"
	"github.com/Ericwyn/Andex/util/log"
	"time"

	"github.com/Ericwyn/Andex/ajax"
)

const baseUrl = "https://api.aliyundrive.com"
const apiVersion = "v2"

const aliUrl = baseUrl + "/" + apiVersion

//=========================================
//   获取文件夹 list
//=========================================

type FileMsgBean struct {
	DriveID         string    `json:"drive_id"`
	DomainID        string    `json:"domain_id"`
	FileID          string    `json:"file_id"`
	Name            string    `json:"name"`
	Type            string    `json:"type"`
	ContentType     string    `json:"content_type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	FileExtension   string    `json:"file_extension"`
	Hidden          bool      `json:"hidden"`
	Size            int64     `json:"size"`
	Starred         bool      `json:"starred"`
	Status          string    `json:"status"`
	UploadID        string    `json:"upload_id"`
	ParentFileID    string    `json:"parent_file_id"`
	Crc64Hash       string    `json:"crc64_hash"`
	ContentHash     string    `json:"content_hash"`
	ContentHashName string    `json:"content_hash_name"`
	DownloadURL     string    `json:"download_url"`
	URL             string    `json:"url"`
	Category        string    `json:"category"`
	EncryptMode     string    `json:"encrypt_mode"`
	PunishFlag      int       `json:"punish_flag"`
}

type FolderListBean struct {
	Items      []FileMsgBean `json:"items"`
	NextMarker string        `json:"next_marker"`
}

// driveId 网盘id
// parentDirId 父文件夹的 id, 根目录为 root
func FolderList(authorization string, driveId string, parentDirId string) *FolderListBean {
	var result *FolderListBean = nil
	ajax.Send(ajax.Request{
		Url:    aliUrl + "/file/list",
		Method: ajax.POST,
		Json: map[string]interface{}{
			"limit":                   50,
			"marker":                  nil,
			"drive_id":                driveId,
			"parent_file_id":          parentDirId,
			"image_thumbnail_process": "image/resize,w_160/format,jpeg",
			"image_url_process":       "image/resize,w_1920/format,jpeg",
			"video_thumbnail_process": "video/snapshot,t_0,f_jpg,w_300",
			"fields":                  "*",
			"order_by":                "updated_at",
			"order_direction":         "DESC",
			"content-type":            "application/json;charset=UTF-8",
		},
		Header: buildHeader(authorization, true),
		Success: func(response *ajax.Response) {
			//fmt.Println("code:", response.Code)
			//fmt.Println("response:")
			//fmt.Println(response.Body)
			err := json.Unmarshal([]byte(response.Body), &result)
			if err != nil {
				log.E("JSON 解析发生错误", err)
			}
		},
		Fail: func(status int, errMsg string) {
			log.E("网络连接失败-FolderList")
			log.E("status:", status, ", errMsg:"+errMsg)

		},
		Always: nil,
	})
	return result
}

//=========================================
//   获取用户信息
//=========================================

type UserInfoBean struct {
	DomainID       string `json:"domain_id"`
	UserID         string `json:"user_id"`
	Avatar         string `json:"avatar"`
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt      int64  `json:"updated_at"`
	Email          string `json:"email"`
	NickName       string `json:"nick_name"`
	Phone          string `json:"phone"`
	Role           string `json:"role"`
	Status         string `json:"status"`
	UserName       string `json:"user_name"`
	Description    string `json:"description"`
	DefaultDriveID string `json:"default_drive_id"`
	UserData       struct {
		Share string `json:"share"`
	} `json:"user_data"`
}

func UserInfo(authorization string) *UserInfoBean {

	var result *UserInfoBean = nil

	ajax.Send(ajax.Request{
		Url:    aliUrl + "/user/get",
		Method: ajax.POST,
		Json:   map[string]interface{}{},
		Header: buildHeader(authorization, true),
		Success: func(response *ajax.Response) {
			//fmt.Println("code:", response.Code)
			//fmt.Println("response:")
			//fmt.Println(response.Body)
			err := json.Unmarshal([]byte(response.Body), &result)
			if err != nil {
				log.E("JSON 解析发生错误", err)
			}
		},
		Fail: func(status int, errMsg string) {
			log.E("网络连接失败-UserInfo")
			log.E("status:", status, ", errMsg:"+errMsg)

		},
		Always: nil,
	})
	return result
}

//=========================================
//   刷新 token
//=========================================
type RefreshTokenBean struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresIn    int       `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	UserID       string    `json:"user_id"`
	UserName     string    `json:"user_name"`
	Avatar       string    `json:"avatar"`
	NickName     string    `json:"nick_name"`
	ExpireTime   time.Time `json:"expire_time"`
	State        string    `json:"state"`

	DefaultDriveID     string `json:"default_drive_id"`
	DefaultSboxDriveID string `json:"default_sbox_drive_id"`
	//Role               string        `json:"role"`
	//Status             string        `json:"status"`

	//ExistLink          []interface{} `json:"exist_link"`
	//NeedLink           bool          `json:"need_link"`
	//PinSetup     bool   `json:"pin_setup"`
	//IsFirstLogin bool   `json:"is_first_login"`
	//NeedRpVerify bool   `json:"need_rp_verify"`
	DeviceID string `json:"device_id"`
}

type RefreshTokenCallback func(authorization string, refreshToken string, driveId string)

func RefreshToken(refreshToken string, callback RefreshTokenCallback) {
	//if conf.SysConfigNow.RefreshToken == "" {
	//	fmt.Println("config.json 中没有配置 refresh token")
	//	return
	//}

	var result *RefreshTokenBean = nil

	ajax.Send(ajax.Request{
		//Url:    "https://websv.aliyundrive.com/token/refresh", // 旧接口好像已经失效了？
		Url:    "https://auth.aliyundrive.com/v2/account/token",
		Method: ajax.POST,
		Json: map[string]interface{}{
			"refresh_token": refreshToken,
			"grant_type":    "refresh_token",
		},
		Header: map[string]string{
			"Content-Type":    "application/json; charset=utf-8",
			"origin":          "https://www.aliyundrive.com",
			"accept":          "*/*",
			"Accept-Language": "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
			"Connection":      "keep-alive",
			"user-agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36",
		},
		Success: func(response *ajax.Response) {
			//fmt.Println("code:", response.Code)
			//fmt.Println("response:")
			//fmt.Println(response.Body)
			err := json.Unmarshal([]byte(response.Body), &result)
			if err != nil {
				log.E("JSON 解析发生错误", err)
			}
			callback(result.AccessToken, result.RefreshToken, result.DefaultDriveID)
		},
		Fail: func(status int, errMsg string) {
			log.E("网络连接失败-RefreshToken")
			log.E("status:", status, ", errMsg:"+errMsg)

		},
		Always: nil,
	})

}

//=========================================
//   获取下载链接
//=========================================
type GetDownloadUrlBean struct {
	Method     string    `json:"method"`
	URL        string    `json:"url"`
	Expiration time.Time `json:"expiration"`
	Size       int       `json:"size"`
	Ratelimit  struct {
		PartSpeed int `json:"part_speed"`
		PartSize  int `json:"part_size"`
	} `json:"ratelimit"`
}

func GeoDownloadUrl(authorization string, driveId string, fileMsg FileMsgBean) string {
	return GetDownloadUrlByFileIdAndFileName(authorization, driveId, fileMsg.FileID, fileMsg.Name)
}

func GetDownloadUrlByFileIdAndFileName(authorization string, driveId string, fileId string, fileName string) string {
	var result *GetDownloadUrlBean = nil

	ajax.Send(ajax.Request{
		Url:    aliUrl + "/file/get_download_url",
		Method: ajax.POST,
		Json: map[string]interface{}{
			"drive_id":   driveId,
			"file_id":    fileId,
			"file_name":  fileName,
			"expire_sec": 7200,
		},
		Header: buildHeader(authorization, true),
		Success: func(response *ajax.Response) {
			err := json.Unmarshal([]byte(response.Body), &result)
			if err != nil {
				log.E("JSON 解析发生错误", err)
			}
		},
		Fail: func(status int, errMsg string) {
			log.E("网络连接失败-GetDownloadUrlByFileIdAndFileName")
			log.E("status:", status, ", errMsg:"+errMsg)
		},
		Always: nil,
	})
	if result != nil {
		return result.URL
	} else {
		return ""
	}
}

//=========================================
func buildHeader(authorization string, auth bool) map[string]string {
	res := map[string]string{
		"origin": "https://www.aliyundrive.com",
		//"referer": "https://www.aliyundrive.com/",
		"accept":          "*/*",
		"Accept-Language": "zh-CN,zh;q=0.8,en-US;q=0.5,en;q=0.3",
		"Connection":      "keep-alive",
		"user-agent":      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36",
	}
	if auth {
		res["authorization"] = "Bearer " + authorization
	}

	return res
}
