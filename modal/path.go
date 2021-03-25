package modal

import (
	"fmt"
)

// Andex 索引缓存
type AndexPath struct {
	Path     string `xorm:"pk"` // Andex 的访问路径 /share/sw/ps.zip
	Name     string // 文件或文件夹名称 ps.zip
	FileId   string // 路径的 fileId
	IsDir    bool   // 类型
	Password string
}

func SaveAndexPath(path AndexPath) error {
	var sql = "INSERT OR REPLACE INTO `andex_path` (`path`,`name`,`file_id`,`is_dir`,`password`) VALUES (?,?,?,?,?)"

	_, err := sqlEngine.Exec(sql, path.Path, path.Name, path.FileId, path.IsDir, path.Password)
	if err != nil {
		return err
	}
	return nil
}

func SaveAndexPathList(paths []AndexPath) error {
	session := sqlEngine.NewSession()
	defer session.Close()

	var sql = "INSERT OR REPLACE INTO `andex_path` (`path`,`name`,`file_id`,`is_dir`,`password`) VALUES (?,?,?,?,?)"
	for _, path := range paths {
		_, err := session.Exec(sql, path.Path, path.Name, path.FileId, path.IsDir, path.Password)
		if err != nil {
			session.Rollback()
			return err
		}
	}
	err := session.Commit()
	if err != nil {
		return err
	}
	return nil
}

type void interface{}

// 通过 fileId 删除一些旧的索引链接关系
func DeleteOldAndexPath(paths []AndexPath) {
	if len(paths) == 0 {
		return
	}

	idSet := make(map[string]void)
	for _, p := range paths {
		if _, ok := idSet[p.FileId]; !ok {
			idSet[p.FileId] = new(void)
		}
	}

	ids := ""
	index := 0
	for fileId, _ := range idSet {
		ids += "'" + fileId + "'"
		if index < len(idSet)-1 {
			ids += ", "
		}

		index++
	}

	var sql = "DELETE FROM `andex_path` WHERE `file_id` in (" + ids + ")"

	_, err := sqlEngine.Exec(sql)
	if err != nil {
		fmt.Println("DeleteOldAndexPath:", err)
	}
}

func ListAndexPath() ([]AndexPath, error) {
	paths := make([]AndexPath, 0)
	err := sqlEngine.Find(&paths)
	if err != nil {
		return nil, err
	} else {
		return paths, nil
	}
}

func LoadPathMap() (map[string]AndexPath, error) {
	paths := make([]AndexPath, 0)
	err := sqlEngine.Find(&paths)
	if err != nil {
		return nil, err
	} else {
		resMap := make(map[string]AndexPath)
		for _, path := range paths {
			resMap[path.Path] = path
		}
		return resMap, nil
	}
}

func DeleteAllPath() {
	_, err := sqlEngine.Exec("DELETE FROM `andex_path`")
	if err != nil {
		fmt.Println("DeleteAllPath ", err)
		return
	}
}

func SavePathMap(pathMap map[string]AndexPath) {
	pathList := make([]AndexPath, 0)
	for _, value := range pathMap {
		pathList = append(pathList, value)
	}
	SaveAndexPathList(pathList)
}

// 获取加密/解密操作的所有路径，包含该 fileId 对应的路径，及其所有子路径
func GetAllSubPathsByFileId(fileId string) ([]AndexPath, error) {

	resPathList := make([]AndexPath, 0)

	andexPathsListOfFileId := make([]AndexPath, 0)
	// 获取路径
	err := sqlEngine.Where("`file_id` like ?", fileId).Find(&andexPathsListOfFileId)
	if err != nil {
		return nil, err
	}

	resPathList = append(resPathList, andexPathsListOfFileId...)
	// 获取 paths 的所有子路径
	for _, andexPathOfFileId := range andexPathsListOfFileId {
		subPathList := make([]AndexPath, 0)
		var subPathStart string
		if andexPathOfFileId.Path != "/" {
			subPathStart = andexPathOfFileId.Path + "/"
		} else {
			subPathStart = andexPathOfFileId.Path
		}

		sqlEngine.Where("`path` like '"+subPathStart+"%'").Find(&subPathList)
		resPathList = append(resPathList, subPathList...)
	}

	return resPathList, nil
}