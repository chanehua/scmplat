package controllers

import (
	"github.com/astaxie/beego"
	"github.com/chanehua/scmplat/models"
)

const (
	timeForm = "2006-01-02 15:04:05"
)

type ProcController struct {
	beego.Controller
}

var procPInfo = new(PageInfo)

func (c *ProcController) Get() {
	c.Data["IsProc"] = true
	c.TplName = "process.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)

	var procMgs []*models.ProcManerge
	var err error
	procPInfo.Page = c.Input().Get("p")
	procPInfo.PageSize = 5
	procPInfo.TableName = "proc_manerge"

	// 定义条件查询表达式
	procPInfo.Fields = []string{"project_name__icontains", "version__icontains",
		"deploy_type__icontains", "create_time__gte", "create_time__lte",
		"deploy_shell_name__icontains"}

	var tp int64
	// 获取查询条件
	procPInfo.SearchT = []string{c.Input().Get("sprojectName"),
		c.Input().Get("sversion"), c.Input().Get("sType"),
		c.Input().Get("startTime"), c.Input().Get("endTime"),
		c.Input().Get("sShellName"),
	}
	// 判断查询条件是否为空，为空则不进行条件查询，否则进行条件查询
	if procPInfo.SearchT[3] == "" {
		procPInfo.SearchT[3] = "1900-01-01 00:00:00"
	}
	if procPInfo.SearchT[4] == "" {
		procPInfo.SearchT[4] = "2388-12-31 23:59:59"
	}
	// 获取总页数以及根据总页数与当前页大小比对来确定当前页值
	tp, procPInfo.Page = models.GetTotalCurrentPage(procPInfo.PageSize,
		procPInfo.TableName, procPInfo.Page, procPInfo.SearchT, procPInfo.Fields)
	// 获取分页数据
	procMgs, err = models.GetProcList(procPInfo.Page, procPInfo.PageSize,
		procPInfo.SearchT, procPInfo.Fields)
	if err != nil {
		beego.Error(err)
	}

	c.Data["Procs"] = procMgs
	c.Data["PageNo"] = procPInfo.Page
	c.Data["TotalPage"] = tp
	c.Data["SprojectName"] = procPInfo.SearchT[0]
	c.Data["Sversion"] = procPInfo.SearchT[1]
	c.Data["Stype"] = procPInfo.SearchT[2]
	c.Data["StartTime"] = procPInfo.SearchT[3]
	c.Data["EndTime"] = procPInfo.SearchT[4]
	c.Data["SshellName"] = procPInfo.SearchT[5]

	// 权限校验
	uname := c.Ctx.GetCookie("uname")
	if len(uname) == 0 {
		c.Data["SecMg"] = "none"
		c.Data["DockerMg"] = "none"
	} else {
		operId := map[string]string{"SecMg": "none",
			"ProcAddButton":    "none",
			"ProcModifyButton": "none",
			"ProcDelButton":    "none",
			"ProcCreateButton": "none",
			"DockerMg":         "none",
		}
		err := models.CheckSec(operId, uname)
		if err != nil {
			beego.Error(err)
		}
		c.Data["SecMg"] = operId["SecMg"]
		c.Data["ProcAddButton"] = operId["ProcAddButton"]
		c.Data["ProcModifyButton"] = operId["ProcModifyButton"]
		c.Data["ProcDelButton"] = operId["ProcDelButton"]
		c.Data["ProcCreateButton"] = operId["ProcCreateButton"]
		c.Data["DockerMg"] = operId["DockerMg"]
	}
}

func (c *ProcController) Post() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	// 解析表单
	procMgs := map[string]string{
		"ProcId":         c.Input().Get("pid"),
		"projectName":    c.Input().Get("projectName"),
		"version":        c.Input().Get("version"),
		"oper":           c.Input().Get("operator"),
		"name":           c.Input().Get("deploy_name"),
		"host":           c.Input().Get("deploy_host"),
		"type":           c.Input().Get("deploy_type"),
		"shellName":      c.Input().Get("shell_name"),
		"memARGS":        c.Input().Get("mem_args"),
		"appframeKey":    c.Input().Get("appframe_key"),
		"appframeValue":  c.Input().Get("appframe_value"),
		"shellParams":    c.Input().Get("shell_params"),
		"parameterValue": c.Input().Get("parameter_value"),
		"startupClass":   c.Input().Get("startup_class"),
		"exeRelat":       c.Input().Get("exerelat"),
		"exeDataSource":  c.Input().Get("exedatasource"),
		"cTime":          c.Input().Get("create_time"),
		"remark":         c.Input().Get("remark"),
	}
	var err error
	if len(procMgs["ProcId"]) == 0 {
		err = models.AddProcManerge(procMgs, c.Ctx.GetCookie("uname"))
	} else {
		err = models.ModifyProc(procMgs, c.Ctx.GetCookie("uname"))
	}
	if err != nil {
		beego.Error(err)
	}
	// 检查当前页是否大于总页数
	_, procPInfo.Page = models.GetTotalCurrentPage(procPInfo.PageSize,
		procPInfo.TableName, procPInfo.Page, procPInfo.SearchT, procPInfo.Fields)
	c.Redirect("/proc?p="+procPInfo.Page+"&sprojectName="+
		procPInfo.SearchT[0]+"&sversion="+procPInfo.SearchT[1]+
		"&sType="+procPInfo.SearchT[2]+"&startTime="+procPInfo.SearchT[3]+
		"&endTime="+procPInfo.SearchT[4]+"&sShellName="+procPInfo.SearchT[5], 302)
}

func (c *ProcController) Modify() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	c.TplName = "proc_modify.html"
	pid := c.Input().Get("pid")
	proc, err := models.GetProc(pid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Proc"] = proc
	c.Data["ProcId"] = pid
	c.Data["IsLogin"] = true
}

func (c *ProcController) Delete() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	pidStr := c.Input().Get("pid")
	err := models.DeleteProcs(pidStr, c.Ctx.GetCookie("uname"))
	if err != nil {
		beego.Error(err)
	}
	// 检查当前页是否大于总页数
	_, procPInfo.Page = models.GetTotalCurrentPage(procPInfo.PageSize,
		procPInfo.TableName, procPInfo.Page, procPInfo.SearchT, procPInfo.Fields)
	c.Redirect("/proc?p="+procPInfo.Page+"&sprojectName="+
		procPInfo.SearchT[0]+"&sversion="+procPInfo.SearchT[1]+
		"&sType="+procPInfo.SearchT[2]+"&startTime="+procPInfo.SearchT[3]+
		"&endTime="+procPInfo.SearchT[4]+"&sShellName="+procPInfo.SearchT[5], 302)
}

func (c *ProcController) Create() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	pidStr := c.Input().Get("pid")
	err := models.CreateProcs(pidStr, c.Ctx.GetCookie("uname"))
	if err != nil {
		beego.Error(err)
	}
	// 检查当前页是否大于总页数
	_, procPInfo.Page = models.GetTotalCurrentPage(procPInfo.PageSize,
		procPInfo.TableName, procPInfo.Page, procPInfo.SearchT, procPInfo.Fields)
	c.Redirect("/proc?p="+procPInfo.Page+"&sprojectName="+
		procPInfo.SearchT[0]+"&sversion="+procPInfo.SearchT[1]+
		"&sType="+procPInfo.SearchT[2]+"&startTime="+procPInfo.SearchT[3]+
		"&endTime="+procPInfo.SearchT[4]+"&sShellName="+procPInfo.SearchT[5], 302)

}
