package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	//"用户名:密码@[连接方式](主机名:端口号)/数据库名"
	db, err := sql.Open("mysql", "root:root@(127.0.0.1:3306)/go_learn") // 设置连接数据库的参数
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	// 查看数据库版本
	var version string
	err2 := db.QueryRow("SELECT VERSION()").Scan(&version)
	if err2 != nil {
		log.Fatal(err2)
	}

	fmt.Println(version)

}
