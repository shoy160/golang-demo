package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "odbc/driver"
	"os"
)

func readData(path string, user string, pwd string) []map[string]interface{} {
	conn := fmt.Sprintf("driver={Microsoft Access Driver (*.mdb, *.accdb)};Dbq=%s;Uid=%s;Pwd=%s;", path, user, pwd)
	db, err := sql.Open("odbc", conn)
	if err != nil {
		fmt.Println("数据库打开错误:", err)
	}
	defer db.Close()
	stmt, err := db.Prepare("select * from lwmain")
	var list []map[string]interface{}
	if err != nil {
		fmt.Println("数据库查询错误:", err)
	} else {
		defer stmt.Close()
		rows, err := stmt.Query()
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		cols, _ := rows.Columns()
		for rows.Next() {
			columns := make([]interface{}, len(cols))
			columnPointers := make([]interface{}, len(cols))
			for i := range columns {
				columnPointers[i] = &columns[i]
			}

			if err := rows.Scan(columnPointers...); err != nil {
				log.Fatal(err)
			}
			m := make(map[string]interface{})
			for i, colName := range cols {
				val := columnPointers[i].(*interface{})
				m[colName] = *val
			}
			list = append(list, m)
		}
	}
	return list
}

func read(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	path := r.PostFormValue("path")
	_, err := os.Stat(path)
	if err != nil {
		data["status"] = false
		data["message"] = "数据文件不存在"
	} else {
		user := r.PostFormValue("user")
		pwd := r.PostFormValue("pwd")
		if pwd == "" {
			pwd = "LW"
		}
		data["status"] = true
		data["data"] = readData(path, user, pwd)
	}
	w.Header().Set("Content-Type", "application/json")
	buffer, _ := json.Marshal(data)
	w.Write(buffer)
}

func start() {
	http.HandleFunc("/xungen", read)
	log.Printf("server listen at 9988")
	if err := http.ListenAndServe(":9988", nil); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
