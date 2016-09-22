package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func printResult(query *sql.Rows) {
	column, _ := query.Columns()              //读出查询出的列字段名
	values := make([][]byte, len(column))     //values是每个列的值，这里获取到byte里
	scans := make([]interface{}, len(column)) //因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
	for i := range values {                   //让每一行数据都填充到[][]byte里面
		scans[i] = &values[i]
	}
	results := make(map[int]map[string]string) //最后得到的map
	i := 0
	for query.Next() { //循环，让游标往下移动
		if err := query.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			fmt.Println(err)
			return
		}
		row := make(map[string]string) //每行数据
		for k, v := range values {     //每行数据是放在values里面，现在把它挪到row里
			key := column[k]
			row[key] = string(v)
		}
		results[i] = row //装入结果集中
		i++
	}
	for k, v := range results { //查询出来的数组
		fmt.Println(k, v)
	}
}
func main() {
	db, _ := sql.Open("mysql", "bpa:bpa123@tcp(localhost:3306)/bpa?charset=utf8")

	fmt.Println(db)

	rows, err := db.Query("select * from client_info")

	checkErr(err)
	colums, _ := rows.Columns()
	fmt.Println(colums)
	for rows.Next() {
		var id int
		var uuid, softType, platform, arch, ip, add_date, active_time string
		rows.Scan(&id, &uuid, &softType, &platform, &arch, &ip, &add_date, &active_time)
		// colums, _ := rows.Columns()
		// fmt.Println(colums)

		fmt.Printf("id = %d, uuid = %s, soft_type = %s, platform = %s, arch = %s, ip = %s, add_date = %s, active_time = %s\n",
			id, uuid, softType, platform, arch, ip, add_date, active_time)
	}

	printResult(rows)
}
