package database

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB
var err error

func init() {

	mysql_username := beego.AppConfig.String("kg_mysql_name")
	mysql_password := beego.AppConfig.String("kg_mysql_password")
	mysql_host := beego.AppConfig.String("kg_mysql_host")
	mysql_dbname := beego.AppConfig.String("kg_mysql_dbname")

	Db, err = gorm.Open("mysql", mysql_username+":"+mysql_password+"@tcp("+mysql_host+":3306)/"+mysql_dbname+"?charset=utf8")

	if err != nil {
		fmt.Println(err)
		panic("连接数据库失败")
	}

	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)

	// 全局禁用表名复数
	Db.SingularTable(true)
}
