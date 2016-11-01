package controllers

import (
	"github.com/astaxie/beego"
	"github.com/chanehua/scmplat/models"
)

type HomeController struct {
	beego.Controller
}

func (c *HomeController) Get() {
	c.Data["IsHome"] = true
	c.TplName = "home.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)

	// 权限校验
	uname := c.Ctx.GetCookie("uname")
	if len(uname) == 0 {
		c.Data["SecMg"] = "none"
	} else {
		operId := map[string]string{"SecMg": "none"}
		err := models.CheckSec(operId, uname)
		if err != nil {
			beego.Error(err)
		}
		c.Data["SecMg"] = operId["SecMg"]
	}

}
