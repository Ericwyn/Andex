package modal

import (
	"fmt"
	"os"
	"xorm.io/xorm"
)

var sqlEngine *xorm.Engine

//var sqlBuilder = builder.SQLite()

func InitDb(showSql bool) {
	var err error

	sqlEngine, err = xorm.NewEngine("sqlite3", "./.conf/andex.db")
	if err != nil {
		fmt.Println(err)
		fmt.Println("\n\n SQL ENGINE INIT FAIL!!")
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
		fmt.Println("SYNC TABLE ERROR!!")
		os.Exit(-1)
	}
}
