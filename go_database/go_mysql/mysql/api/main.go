package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

const (
	UserName = "root"
	PassWord = "root"
	Host     = "127.0.0.1"
	Port     = "3306"
	Database = "go_learn"
	Charset  = "utf8"
)

type Model struct {
	line      *sql.DB
	tableName string
	field     string
	allFields []string
	where     string
	order     string
	limit     string
}

func inArray(need interface{}, needArr []string) bool {
	for _, v := range needArr {
		if need == v {
			return true
		}
	}
	return false
}

func returnRes(errCode int, res interface{}, msg interface{}) string {
	result := make(map[string]interface{})
	result["errCode"] = errCode
	result["res"] = res
	result["msg"] = msg
	data, _ := json.Marshal(result)
	return string(data)
}

func NewModel(table string) Model {
	var m Model
	m.field = "*"
	m.tableName = table
	m.getConnect()
	m.getFields()
	return m
}

func (m *Model) getConnect() {
	//dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s",UserName,PassWord, Host, Port, Database, Charset)
	//println(dbDSN)
	//dbDSN := fmt.Sprintf("#{UserName}:#{PassWord}@tcp(#{Host}:#{Port})/#{Database}?charset=#{Charset}")
	//db, err := mysql.Open("mysql", string(dbDSN))
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/go_learn?charset=utf8")
	if err != nil {
		fmt.Printf("Connect mysql fail ![%s]\n", err)
	} else {
		fmt.Println("Connect mysql successful")
	}
	m.line = db
}

func (m *Model) getFields() {
	sqlState := "DESC " + m.tableName
	fmt.Println(sqlState)
	result, err := m.line.Query(sqlState)
	if err != nil {
		fmt.Printf("sql fail ![%s]\n", err)
	}
	m.allFields = make([]string, 0)
	for result.Next() {
		var field string
		var Type interface{}
		var Null string
		var key string
		var Default interface{}
		var Extra string
		err := result.Scan(&field, &Type, &Null, &key, &Default, &Extra)
		if err != nil {
			fmt.Printf("Scan fail ![%s]\n", err)
		}
		m.allFields = append(m.allFields, field)
	}
}

func (m *Model) Field(field string) *Model {
	m.field = field
	return m
}

func (m *Model) Order(order string) *Model {
	m.order = ` order by ` + order
	return m
}

func (m *Model) Limit(limit int) *Model {
	m.limit = ` limit ` + strconv.Itoa(limit)
	return m
}

func (m *Model) Where(where string) *Model {
	m.where = ` where ` + where
	return m
}

func (m *Model) Count(count string) interface{} {
	sql := ` select count(*) as total from ` + m.tableName + ` limit 1`
	result := m.query(sql)
	return result
}

func (m *Model) query(sql string) interface{} {
	row2, err := m.line.Query(sql)
	if err != nil {
		return returnRes(0, "", err)
	}
	cols, err := row2.Columns()
	if err != nil {
		return returnRes(0, "", err)
	}
	value := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for k, _ := range value {
		scans[k] = &value[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for row2.Next() {
		_ = row2.Scan(scans...)
		row := make(map[string]string)
		for k, v := range value {
			key := cols[k]
			row[key] = string(v)
		}
		result[i] = row
		i++
	}
	return returnRes(1, result, "Success")
}

func (m *Model) exec(sql string) interface{} {
	res, err := m.line.Exec(sql)
	if err != nil {
		return returnRes(0, "", err)
	}
	result, err := res.LastInsertId()
	if err != nil {
		return returnRes(0, "", err)
	}
	return returnRes(1, result, "Success")
}

func (m *Model) selectAll() interface{} {
	// 查询所有
	sql := ` select ` + m.field + ` from ` + m.tableName + ` ` + m.where + ` ` + m.order + ` ` + m.limit
	result := m.query(sql)
	return result
}

func (m *Model) selectOnc(userId int) interface{} {
	where := ` where userId = ` + strconv.Itoa(userId)
	sql := ` select ` + m.field + ` from ` + m.tableName + where + `limit 1`
	result := m.query(sql)
	return result
}

func (m *Model) add(data map[string]interface{}) interface{} {
	var key string
	var value string
	for k, v := range data {
		if res := inArray(k, m.allFields); res != true {
			delete(data, k)
		} else {
			key += `,` + k
			value += `,` + ` · ` + v.(string) + ` · `
		}
	}
	key = strings.TrimLeft(key, `,`)
	value = strings.TrimLeft(value, `,`)
	sql := ` insert info ` + m.tableName + `(` + key + `) values( ` + value + `)`
	fmt.Println(sql)
	result := m.exec(sql)
	return result
}

func (m *Model) delete(userId int) interface{} {
	var conditions string
	if m.where == "" {
		conditions = ` where userId = ` + strconv.Itoa(userId)
	} else {
		conditions = m.where + ` and userId = ` + strconv.Itoa(userId)
	}
	sql := ` delete from ` + m.tableName + ` ` + conditions
	fmt.Println(sql)
	result := m.exec(sql)
	return result
}

func (m *Model) update(data map[string]interface{}) interface{} {
	var str string
	for k, v := range data {
		if res := inArray(k, m.allFields); res != true {
			delete(data, k)
		} else {
			str += k + ` = '` + v.(string) + `',`
		}
	}

	str = strings.TrimLeft(str, `,`)
	if m.where == "" {
		fmt.Println("No condition")
	}
	sql := ` update ` + m.tableName + ` set ` + str + ` ` + m.where
	fmt.Println(sql)
	result := m.exec(sql)
	return result
}

func main() {
	M := NewModel("users")

	data := make(map[string]interface{})
	data["user_name"] = "Make"
	data["email"] = "12345678@qq.com"
	data["add_time"] = "2021-06-25 21:25:20"
	res := M.Where("user_id = 1").update(data)
	fmt.Println("res", res)
}
