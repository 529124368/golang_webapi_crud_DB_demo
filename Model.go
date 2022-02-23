package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Userinfo struct {
	Id            int64          `json:"id"`
	Name          string         `json:"name"`
	Account       string         `json:"account"`
	Password      string         `json:"password"`
	Register_time sql.NullString `json:"Register_time"`
}
type DBModel struct {
	ID       string
	Password string
	Port     string
	Url      string
	Con      *sqlx.DB
}

func (d *DBModel) start() {

	dbconfig := d.ID + ":" + d.Password + "@tcp(" + d.Url + ":" + d.Port + ")/mir3"
	db, err := sqlx.Connect("mysql", dbconfig)
	if err != nil {
		fmt.Println("数据库连接错误：", err)
		return
	}
	db.SetMaxOpenConns(100)
	d.Con = db
}
func (d *DBModel) close() {
	d.Con.Close()
}
func NewDBCon(id string, port string, password string, url string) *DBModel {
	db := &DBModel{
		ID:       id,
		Password: password,
		Port:     port,
		Url:      url,
	}
	return db
}

//插入数据库
func (d *DBModel) InsertUserAccount(name string, account string, password string) string {
	str_time := time.Now().Format("2006-01-02")

	sqlStrs := "insert into user_account(name,account,password,Register_time) values(?,?,?,?)"
	_, err := d.Con.Exec(sqlStrs, name, account, password, str_time)
	if err != nil {
		fmt.Println("数据库插入错误：", err)
		return string([]byte(`{"state": "error", "data": "","message":"注册失败"}`))
	} else {
		fmt.Println("用户注册成功")
		return string([]byte(`{"state": "ok", "data": "","message":"注册成功"}`))
	}
}

//查询数据库
func (d *DBModel) SelectUserAccount(account string, password string) string {
	var user Userinfo
	sqlStrs := "select * from user_account where  account =? and password = ?"
	err := d.Con.Get(&user, sqlStrs, account, password)
	if err != nil {
		fmt.Println("数据库查询错误：", err)
		return string([]byte(`{"state": "error", "data": "","message":"数据库查询错误"}`))
	} else {
		fmt.Println("用户查询成功 用户ID：", user.Id)
		b, err := json.Marshal(&user)
		if err != nil {
			fmt.Println(err)
		}
		return string([]byte(`{"state": "ok", "data": ` + string(b) + `,"message":"查询成功"}`))
	}
}
