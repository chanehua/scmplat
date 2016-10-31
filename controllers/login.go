package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/chanehua/scmplat/auth"
	"github.com/chanehua/scmplat/models"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	// 判断是否为退出操作
	if c.Input().Get("exit") == "true" {
		c.Ctx.SetCookie("uname", "", -1, "/")
		c.Ctx.SetCookie("pwd", "", -1, "/")
		c.Redirect("/", 302)
		return
	}
	c.TplName = "login.html"
}

func (c *LoginController) Post() {
	// 获取表单信息
	uname := c.Input().Get("uname")
	pwd := c.Input().Get("pwd")
	autoLogin := c.Input().Get("autoLgin") == "on"

	//进行LDAP认证
	err := auth.LdapAuth(uname, pwd)
	if err != nil {
		beego.Error(err)
		c.Data["loginErrorInfo"] = "user name or password error.please enter again!"
		c.TplName = "login.html"
		return
	}
	maxAge := 0
	if autoLogin {
		maxAge = 1<<31 - 1
	}

	// 增加用户到用户信息表
	err = models.AddUserInfo(uname)
	if err != nil {
		beego.Error(err)
	}
	c.Ctx.SetCookie("uname", uname, maxAge, "/")
	c.Ctx.SetCookie("pwd", pwd, maxAge, "/")

	c.Redirect("/", 302)
	return
}

func checkAccount(ctx *context.Context) bool {
	ck, err := ctx.Request.Cookie("uname")
	if err != nil {
		return false
	}
	uname := ck.Value

	ck, err = ctx.Request.Cookie("pwd")
	if err != nil {
		return false
	}

	pwd := ck.Value

	err = auth.LdapAuth(uname, pwd)
	if err != nil {
		beego.Error(err)
		return false
	}
	// 增加用户到用户信息表
	err = models.AddUserInfo(uname)
	if err != nil {
		beego.Error(err)
	}
	return true
}
