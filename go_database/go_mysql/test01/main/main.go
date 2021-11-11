package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	USER_NAME = "root"
	PASS_WORD = "root"
	HOST      = "127.0.0.1"
	PORT      = "3306"
	DATABASE  = "go_learn"
	CHARSET   = "utf8"
)

//Message 是调用了api 约定返回的格式 结构体
type Message struct {
	Code     string
	Msg      string
	Time     int64
	UserInfo UserInfo
}

//UserInfo 是查询了数据库之后 用于接受数据的结构体
type UserInfo struct {
	Id   int
	Name string
	Age  int
}

func init() {
	//初始化mysql 驱动 获得db
	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", USER_NAME, PASS_WORD, HOST, PORT, DATABASE, CHARSET)

	db, err := sql.Open("mysql", dbDSN)

	//设置最大连接数
	db.SetMaxOpenConns(20)
	//设置空闲数
	db.SetMaxIdleConns(10)
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatalf("Error on opening database connection: %s", pingErr.Error())
	}
	checkErr(err)
}

//func main() {
//	//设置监听路径
//	http.HandleFunc("/userinfo", handler)
//	http.ListenAndServe("localhost:8899", nil)
//}

//func queryUserInfo(name string)([]byte, error) {
//	//查询数据库  name 是传进来的参数
//	rows, err := db.Query("select id, `name`, age from shop_user where name = ?", name)
//	checkErr(err)
//	var uInfo UserInfo
//	for rows.Next() {
//		var id int
//		var name string
//		var age int
//		if err:= rows.Scan(&id, &name, &age); err != nil {
//			log.Fatal(err);
//		}
//		uInfo = UserInfo{id, name, age};
//	}
//	//拼接数据格式 并且是json 返回
//	m := Message{"0","获取成功", time.Now().Unix(),uInfo}
//	return json.MarshalIndent(m,"","");
//}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
