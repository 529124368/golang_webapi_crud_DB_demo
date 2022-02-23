package main

import (
	"fmt"
	"net/http"
)

var dbCon *DBModel

//跨域中间件
func cros(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
		w.Header().Set("content-type", "application/json")             //返回数据格式是json
		f(w, r)
	}
}

func registerMethod(w http.ResponseWriter, r *http.Request) {
	//POST 请求验证
	if r.Method == "POST" {
		name := r.PostFormValue("name")
		account := r.PostFormValue("account")
		password := r.PostFormValue("password")
		if name == "" || account == "" || password == "" {
			fmt.Println("账号信息输入不正确")
			return
		}
		fmt.Fprintf(w, dbCon.InsertUserAccount(name, account, password))
	} else {
		fmt.Fprintf(w, "请求失败")
		return
	}
}

func checkMethod(w http.ResponseWriter, r *http.Request) {
	//POST 请求验证
	if r.Method == "POST" {
		account := r.PostFormValue("account")
		password := r.PostFormValue("password")
		if account == "" || password == "" {
			fmt.Println("账号信息输入不正确")
			return
		}
		fmt.Fprintf(w, dbCon.SelectUserAccount(account, password))
	} else {
		fmt.Fprintf(w, "请求失败")
		return
	}
}

func main() {
	fmt.Println("服务器启动")
	dbCon = NewDBCon("root", "3306", "root00", "127.0.0.1")
	defer dbCon.close()
	dbCon.start()
	//web服务器
	http.HandleFunc("/register", cros(registerMethod))
	http.HandleFunc("/checkAccount", cros(checkMethod))

	// 启动web服务，监听8080端口
	http.ListenAndServe(":8081", nil)

}
