package database

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
)

var Olddb *gorm.DB
var Olderr error

func init() {

	mysql_username := beego.AppConfig.String("kg_mysql_name")
	mysql_password := beego.AppConfig.String("kg_mysql_password")
	mysql_host := beego.AppConfig.String("kg_mysql_host")
	mysql_dbname := beego.AppConfig.String("kg_old_mysql_dbname")

	Olddb, Olderr = gorm.Open("mysql", mysql_username+":"+mysql_password+"@tcp("+mysql_host+":3306)/"+mysql_dbname+"?charset=utf8&parseTime=True&loc=Local")

	if Olderr != nil {
		fmt.Println(Olderr)
		panic("连接数据库失败")
	}

	Olddb.DB().SetMaxIdleConns(10)
	Olddb.DB().SetMaxOpenConns(100)

	// 全局禁用表名复数
	Olddb.SingularTable(true)
}
