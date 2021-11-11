package main

import (
	"database/sql"
	"log"

)

func main() {
	//修改成自己数据账号及密码
	db, err := sql.Open("oci8", "sccdb/scc@dbnms")
	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("select * from dual")
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var i string
		err = rows.Scan(&i)
		if err != nil {
			log.Fatal(err)
		}
		println(i)
	}
}
