package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	"github.com/chanehua/scmplat/models"
	_ "github.com/chanehua/scmplat/routers"
)

//register db
func init() {
	models.RegisterDB()
}

func main() {
	//open orm debug
	orm.Debug = true
	//open auto create table
	orm.RunSyncdb("default", false, true)

	//run beego
	beego.Run()
}
