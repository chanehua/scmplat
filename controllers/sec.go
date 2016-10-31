package controllers

import (
	"github.com/astaxie/beego"
	"github.com/chanehua/scmplat/models"
)

type SecController struct {
	beego.Controller
}

var secPInfo = new(PageInfo)

func (c *SecController) Get() {
	c.Data["IsSec"] = true
	c.TplName = "sec.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)

	var secs []*models.SecMg
	var err error
	secPInfo.Page = c.Input().Get("p")
	secPInfo.PageSize = 5
	secPInfo.TableName = "sec_mg"

	// 设置查询表达式
	secPInfo.Fields = []string{"role_name__icontains", "oper_id__icontains",
		"dpl_status__icontains",
	}
	var tp int64

	// 获取查询条件
	secPInfo.SearchT = []string{c.Input().Get("sRoleName"),
		c.Input().Get("sOperID"), c.Input().Get("sDplStatus"),
	}
	// 获取总页数以及根据总页数与当前页大小比对来确定当前页值
	tp, secPInfo.Page = models.TotalCurrentPage(secPInfo.PageSize,
		secPInfo.TableName, secPInfo.Page, secPInfo.SearchT, secPInfo.Fields)
	// 获取分页数据
	secs, err = models.GetSecList(secPInfo.Page, secPInfo.PageSize,
		secPInfo.SearchT, secPInfo.Fields)
	if err != nil {
		beego.Error(err)
	}

	c.Data["secs"] = secs
	c.Data["PageNo"] = secPInfo.Page
	c.Data["TotalPage"] = tp
	c.Data["SroleName"] = secPInfo.SearchT[0]
	c.Data["SoperID"] = secPInfo.SearchT[1]
	c.Data["SdplStatus"] = secPInfo.SearchT[2]
}

func (c *SecController) Post() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	// 解析表单
	secItem := map[string]string{
		"sid":       c.Input().Get("sid"),
		"roleName":  c.Input().Get("roleName"),
		"operID":    c.Input().Get("operID"),
		"dplStatus": c.Input().Get("dplStatus"),
	}
	var err error
	if len(secItem["sid"]) == 0 {
		err = models.AddSec(secItem)
	} else {
		err = models.ModifySec(secItem)
	}
	if err != nil {
		beego.Error(err)
	}

	// 检查当前页是否大于总页数
	_, secPInfo.Page = models.TotalCurrentPage(secPInfo.PageSize,
		secPInfo.TableName, secPInfo.Page, secPInfo.SearchT, secPInfo.Fields)
	c.Redirect("/sec?p="+secPInfo.Page+"&sRoleName="+
		secPInfo.SearchT[0]+"&sOperID="+secPInfo.SearchT[1]+
		"&sDplStatus="+secPInfo.SearchT[2], 302)
}

func (c *SecController) Modify() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	c.TplName = "sec_modify.html"
	sid := c.Input().Get("sid")
	sec, err := models.GetSec(sid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Sec"] = sec
	c.Data["SecId"] = sid
	c.Data["IsLogin"] = true
}

func (c *SecController) Delete() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	sidStr := c.Input().Get("sid")
	err := models.DeleteSecs(sidStr)
	if err != nil {
		beego.Error(err)
	}

	// 检查当前页是否大于总页数
	_, secPInfo.Page = models.TotalCurrentPage(secPInfo.PageSize,
		secPInfo.TableName, secPInfo.Page, secPInfo.SearchT, secPInfo.Fields)
	c.Redirect("/sec?p="+secPInfo.Page+"&sRoleName="+
		secPInfo.SearchT[0]+"&sOperID="+secPInfo.SearchT[1]+
		"&sDplStatus="+secPInfo.SearchT[2], 302)
}
