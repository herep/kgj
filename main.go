package main

import (
	_ "finance/routers"
	"fmt"
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {

	mysql_username := beego.AppConfig.String("kg_mysql_name")
	mysql_password := beego.AppConfig.String("kg_mysql_password")
	mysql_host := beego.AppConfig.String("kg_mysql_host")
	mysql_dbname := beego.AppConfig.String("kg_mysql_dbname")

	//数据库连接 -- kuguanjia
	if err := orm.RegisterDataBase("default", "mysql", mysql_username+":"+mysql_password+"@tcp("+mysql_host+":3306)/"+mysql_dbname+"?charset=utf8"); err != nil {
		fmt.Println(err)
	}

}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
