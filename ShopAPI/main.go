package main

import (
	_ "ShopAPI/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	// drop table 后再建表
	force := false
	// 打印执行过程
	verbose := true
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "root:medex@/medex?charset=utf8")
	orm.RunSyncdb("default", force, verbose)
	orm.Debug = true
	beego.Run()
}
