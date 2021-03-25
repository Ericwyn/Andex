package main

import (
	"github.com/Ericwyn/Andex/modal"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 需要移到数据库存储的类
	// 1. PathMap， 存储路径缓存
	// 2. config 存储配置, k-v 对
	// 3. pathPw 路径访问密码
	modal.InitDb(true)

	//pathsArr := make([]modal.AndexPath, 0)
	//for i:=0;i<52;i++ {
	//	pathsArr = append(pathsArr, modal.AndexPath{
	//		Path: "path" + strconv.Itoa(i),
	//	})
	//}
	//modal.SavePathMapList(pathsArr)

	//list, err := modal.ListAndexPath()
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//	for _,path := range list {
	//		fmt.Println(path.Path)
	//	}
	//}
	//modal.DeleteAllPath()

	riverRootPath1 := modal.AndexPath{
		Path:   "/",
		Name:   "/1",
		FileId: "600111c4bb3567a8f9a74162a24d0596fb8ebc07",
		IsDir:  true,
	}
	riverRootPath2 := modal.AndexPath{
		Path:   "/",
		Name:   "/2",
		FileId: "6050944156a91dc849af456e9bc13b146c39957d",
		IsDir:  true,
	}

	modal.DeleteOldAndexPath([]modal.AndexPath{
		riverRootPath1, riverRootPath2,
	})

	//err := modal.SaveAndexPath(riverRootPath1)
	//if err != nil {
	//	fmt.Println(err)
	//}
}
