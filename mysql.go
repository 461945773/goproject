// mysql.go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"strconv"
)

var db *sql.DB

type texts struct {
	Id   int
	Name string
	Time string
	Text string
}

var count = 0

func mysqlInsert(db *sql.DB, id int, name string, text string) {
	var str = "insert into texts (ID,name,submitTime,submitText) values (" + strconv.Itoa(id) + ",'" + name + "',now(),'" + text + "')"
	_, err := db.Query(str)
	checkError(err)
}

func mysqlQuery(db *sql.DB, w http.ResponseWriter, tem *template.Template) {
	rows, err := db.Query("select * from texts")
	if err != nil {
		fmt.Println("dbQuery : ", err)
	}
	var ts []texts
	count = 0
	for rows.Next() {
		count++
		var t texts
		err := rows.Scan(&t.Id, &t.Name, &t.Time, &t.Text)
		checkError(err)
		ts = append(ts, t)
	}
	err = tem.Execute(w, ts)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func loadHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		content := r.FormValue("content")
		if name == "" || content == "" {
		} else {
			mysqlInsert(db, count+1, name, content)
		}
	}
	t, err := template.ParseFiles("mysql.html")
	if err != nil {
		fmt.Println(err)
	}
	mysqlQuery(db, w, t)
}

func main() {
	var errs error
	db, errs = sql.Open("mysql", "root:nothing@/myboard")
	checkError(errs)
	fmt.Println("connect successfully.")
	defer db.Close()
	http.HandleFunc("/", loadHandle)
	http.ListenAndServe(":8888", nil)
}

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
