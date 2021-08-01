package modal

import (
	"fmt"
	"github.com/Ericwyn/Andex/util/log"
	"github.com/Ericwyn/GoTools/file"
	"os"
	"xorm.io/xorm"
)

var sqlEngine *xorm.Engine

//var sqlBuilder = builder.SQLite()

func InitDb(showSql bool) {
	var err error

	dbDir := file.OpenFile(".conf")
	if !dbDir.Exits() || !dbDir.IsDir() {
		log.I("can't find .conf/ path for save db and config")
	} else {
		mkdirs := dbDir.Mkdirs()
		if !mkdirs {
			log.E("create .conf/ dir fail")
		}
	}

	sqlEngine, err = xorm.NewEngine("sqlite3", "./.conf/andex.db")
	if err != nil {
		log.E(err)
		log.E("\n\n SQL ENGINE INIT FAIL!!")
		os.Exit(-1)
	}

	// 开启 SQL 打印
	if showSql {
		sqlEngine.ShowSQL(true)
	}

	// 同步表结构
	err = sqlEngine.Sync2(new(AndexPath), new(AndexConfig))
	if err != nil {
		fmt.Println(err)
		log.E("SYNC TABLE ERROR!!")
		os.Exit(-1)
	}
}
