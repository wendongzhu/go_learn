package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"time"
	_ "time"
)

const (
	UserName = "root"
	PassWord = "root"
	Host     = "127.0.0.1"
	Port     = "3306"
	Database = "go_learn"
	Charset  = "utf8"
)

var (
	key        string
	value      string
	conditions string
	str        string
)

type Model struct {
	link      *sql.DB  // Storage connection object
	tableName string   // Storage table name
	field     string   // Storage field
	allFields []string // Storage all fields of the current table
	where     string   // Storage where conditions
	order     string   // Storage order conditions
	limit     string   // Storage limit conditions
}

func NewModel(table string) Model {
	// Construction method
	// 1. Table name for storage operation
	// 2. Initialize the connection to the database
	// 3. Get all the fields of the current table

	var this Model
	this.field = "*"
	this.tableName = table
	this.getConnect()
	this.getFields()
	return this
}

func (m *Model) getConnect() {
	//Initialize the connection to the database operation
	dbUserData := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", UserName, PassWord, Host, Port, Database, Charset)
	db, err := sql.Open("mysql", dbUserData)
	if err != nil {
		fmt.Printf("connect mysql fail ! [%s]", err)
	}
	m.link = db
}

func (m *Model) getFields() {
	// Get all the fields of the current table
	// 1. View the table structure;
	// 2. Execute and send SQL (query)
	sql := "DESC " + m.tableName
	result, err := m.link.Query(sql)
	if err != nil {
		fmt.Printf("sql fail ! [%s]", err)
	}

	m.allFields = make([]string, 0)

	for result.Next() {
		var field string
		var Type interface{}
		var Null string
		var Key string
		var Default interface{}
		var Extra string
		err := result.Scan(&field, &Type, &Null, &Key, &Default, &Extra)
		if err != nil {
			fmt.Printf("scan fail ! [%s]", err)
		}
		m.allFields = append(m.allFields, field)
	}

}

func (m *Model) query(sql string) interface{} {
	// Execute and send SQL (query)
	// SQL statement to be queried
	// returns the queried two-dimensional array

	// Query data, get all fields
	rows2, err := m.link.Query(sql)
	if err != nil {
		return returnRes(0, ``, err)
	}

	// Return all columns
	cols, err := rows2.Columns()
	if err != nil {
		return returnRes(0, ``, err)
	}

	// Here represents the value of all columns in a row, represented by []byte
	value := make([][]byte, len(cols))

	// This represents a row of filled data
	scans := make([]interface{}, len(cols))

	// Here scans refers to value and fills the data into []byte
	for k, _ := range value {
		scans[k] = &value[k]
	}

	i := 0
	result := make(map[int]map[string]string)

	for rows2.Next() {
		rows2.Scan(scans...)
		row := make(map[string]string)
		for k, v := range value {
			key := cols[k]
			row[key] = string(v)
		}
		result[i] = row
		i++
	}
	return returnRes(1, result, "success")

}

func (m *Model) get() interface{} {
	// Query multiple data
	sql := `select ` + m.field + ` from ` + m.tableName + ` ` + m.where + ` ` + m.order + ` ` + m.limit
	result := m.query(sql)
	return result
}

func (m *Model) find(user_id int) interface{} {
	// Query a piece of data according to user_id and return the result
	where := ` where user_id = ` + strconv.Itoa(user_id)
	sql := `select ` + m.field + ` from ` + m.tableName + where + ` limit 1`
	result := m.query(sql)
	return result
}

func (m *Model) Field(field string) *Model {
	// Set the field data to be queried and return the query result
	m.field = field
	return m
}

func (m *Model) Order(order string) *Model {
	m.order = `order by ` + order
	return m
}

func (m *Model) Limit(limit int) *Model {
	m.limit = "limit " + strconv.Itoa(limit)
	return m
}

func (m *Model) Where(where string) *Model {
	m.where = `where ` + where
	return m
}

func (m *Model) count() interface{} {
	sql := `select count(*) as total from ` + m.tableName + ` limit 1`
	result := m.query(sql)
	return returnRes(1, result, "success")
}

func (m *Model) exec(sql string) interface{} {
	// Execute and send SQL statements (addition, deletion and modification);
	// If the addition is successful, the last operation id will be returned;
	// It will return true if delete modification operation, false if it fails;
	res, err := m.link.Exec(sql)
	if err != nil {
		return returnRes(0, ``, err)
	}
	result, err := res.LastInsertId()
	if err != nil {
		return returnRes(0, ``, err)
	}
	return returnRes(1, result, "success")
}

func (m *Model) insert(data map[string]interface{}) interface{} {
	// Add operation
	// If the addition is successful, it returns the id of the last operation, and if it fails, it returns false.
	// Filter illegal fields
	for k, v := range data {
		if res := inArray(k, m.allFields); res != true {
			delete(data, k)
		} else {
			key += `,` + k
			value += `,` + `'` + v.(string) + `'`
		}
	}

	// Convert the keys and values retrieved from the map to string concatenation
	key = strings.TrimLeft(key, ",")
	value = strings.TrimLeft(value, ",")

	sql := `insert into ` + m.tableName + ` (` + key + `) values (` + value + `)`
	result := m.exec(sql)
	return result
}

func (m *Model) delete(user_id int) interface{} {
	// Delete operation
	//  1. The id to be deleted
	//  2. Return true if the deletion is successful, false if it fails
	if m.where == "" {
		conditions = `where user_id = ` + strconv.Itoa(user_id)
	} else {
		conditions = m.where + ` and user_id = ` + strconv.Itoa(user_id)
	}

	sql := `delete from ` + m.tableName + ` ` + conditions
	result := m.exec(sql)
	return result
}

/**
 * 修改操作
 * @param  array $data  要修改的数组
 * @return bool 修改成功返回true，失败返回false
 */
func (m *Model) update(data map[string]interface{}) interface{} {
	// Update operation
	// Return true if the modification is successful, false if it fails
	// Filter illegal fields
	for k, v := range data {
		if res := inArray(k, m.allFields); res != true {
			delete(data, k)
		} else {
			str += k + ` = '` + v.(string) + `',`
		}
	}

	// Remove the rightmost comma
	str = strings.TrimRight(str, ",")

	// Determine whether there are query conditions
	if m.where == "" {
		fmt.Println("没有条件")
	}

	sql := `update ` + m.tableName + ` set ` + str + ` ` + m.where
	result := m.exec(sql)
	return result
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
	result["result"] = res
	result["msg"] = msg
	data, _ := json.Marshal(result)
	return string(data)
}

func main() {

	M := NewModel("users")
	//查询链式操作
	//res := M.Field("user_id,user_name,email").
	//	//Order("user_id desc").
	//	//Where("user_id = 2").
	//	Limit(2).
	//	get()

	//添加操作
	data := make(map[string]interface{})
	data["email"] = "118284901@qq.com"
	data["user_name"] = "alice"
	data["add_time"] = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
	res := M.insert(data)

	//删除操作
	//res := M.delete(2)

	//更新操作
	//data := make(map[string]interface{})
	//data["email"] = "118284901@qq.com"
	//data["user_name"] = "make"
	//data["add_time"] = time.Unix(time.Now().Unix(),0).Format("2006-01-02 15:04:05")
	//res:=M.Where("user_id = 1").update(data)
	fmt.Println(res)

}
