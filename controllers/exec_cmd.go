package controllers

import (
	"os"
	"path"

	"github.com/astaxie/beego"

	"github.com/chanehua/scmplat/models"
)

type ExecController struct {
	beego.Controller
}

var execPInfo = new(PageInfo)

func (c *ExecController) Get() {
	c.Data["IsExec"] = true
	c.TplName = "exec_cmd.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)

	var execs []*models.ExecMg
	var err error
	execPInfo.Page = c.Input().Get("p")
	execPInfo.PageSize = 5
	execPInfo.TableName = "exec_mg"
	execPInfo.Fields = []string{"project_name__icontains", "version__icontains",
		"oper_type__icontains", "exec_time__gte", "exec_time__lte",
		"target_server__icontains"}
	var tp int64

	// 获取查询条件
	execPInfo.SearchT = []string{c.Input().Get("sprojectName"),
		c.Input().Get("sversion"), c.Input().Get("sType"),
		c.Input().Get("startTime"), c.Input().Get("endTime"),
		c.Input().Get("stargetSer"),
	}

	// 判断查询条件是否为空，为空则不进行条件查询，否则进行条件查询
	if execPInfo.SearchT[3] == "" {
		execPInfo.SearchT[3] = "1900-01-01 00:00:00"
	}
	if execPInfo.SearchT[4] == "" {
		execPInfo.SearchT[4] = "2388-12-31 23:59:59"
	}
	// 获取总页数以及根据总页数与当前页大小比对来确定当前页值
	tp, execPInfo.Page = models.GetTotalCurrentPage(execPInfo.PageSize,
		execPInfo.TableName, execPInfo.Page, execPInfo.SearchT, execPInfo.Fields)
	// 获取分页数据
	execs, err = models.GetExecList(execPInfo.Page, execPInfo.PageSize,
		execPInfo.SearchT, execPInfo.Fields)
	if err != nil {
		beego.Error(err)
	}
	c.Data["execs"] = execs
	c.Data["PageNo"] = execPInfo.Page
	c.Data["TotalPage"] = tp
	c.Data["SprojectName"] = execPInfo.SearchT[0]
	c.Data["Sversion"] = execPInfo.SearchT[1]
	c.Data["Stype"] = execPInfo.SearchT[2]
	c.Data["StartTime"] = execPInfo.SearchT[3]
	c.Data["EndTime"] = execPInfo.SearchT[4]
	c.Data["StargetSer"] = execPInfo.SearchT[5]

	// 权限校验
	uname := c.Ctx.GetCookie("uname")
	if len(uname) == 0 {
		c.Data["SecMg"] = "none"
		c.Data["DockerMg"] = "none"
	} else {
		operId := map[string]string{"SecMg": "none",
			"StAddButton":    "none",
			"StModifyButton": "none",
			"StDelButton":    "none",
			"StExecButton":   "none",
			"DockerMg":       "none",
		}
		err := models.CheckSec(operId, uname)
		if err != nil {
			beego.Error(err)
		}
		c.Data["SecMg"] = operId["SecMg"]
		c.Data["StAddButton"] = operId["StAddButton"]
		c.Data["StModifyButton"] = operId["StModifyButton"]
		c.Data["StDelButton"] = operId["StDelButton"]
		c.Data["StExecButton"] = operId["StExecButton"]
		c.Data["DockerMg"] = operId["DockerMg"]
	}
}

func (c *ExecController) Post() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	// 获取密钥
	_, fh, err := c.GetFile("sshKey")
	if err != nil {
		beego.Error(err)
	}

	var sshKey string
	if fh != nil {
		// 创建保存附件目录
		os.MkdirAll("sshKey", os.ModePerm)
		// 保存密钥
		sshKey = path.Join("sshKey", fh.Filename)
		beego.Info(sshKey)
		err = c.SaveToFile("sshKey", sshKey)
		os.Chmod(sshKey, 0400)
		if err != nil {
			beego.Error(err)
		}
	}
	// 解析表单
	execTask := map[string]string{
		"eid":         c.Input().Get("eid"),
		"projectName": c.Input().Get("projectName"),
		"version":     c.Input().Get("version"),
		"operType":    c.Input().Get("operType"),
		"targetSer":   c.Input().Get("targetSer"),
		"sshUser":     c.Input().Get("sshUser"),
		"sshPwd":      c.Input().Get("sshPwd"),
		"sshPort":     c.Input().Get("sshPort"),
		"sshKey":      sshKey,
		"execScript":  c.Input().Get("execScript"),
		"execTime":    c.Input().Get("execTime"),
		"operator":    c.Input().Get("operator"),
		"remark":      c.Input().Get("remark"),
	}

	if len(execTask["eid"]) == 0 {
		err = models.AddExecTask(execTask, c.Ctx.GetCookie("uname"))
	} else {
		err = models.ModifyExec(execTask, c.Ctx.GetCookie("uname"))
	}
	if err != nil {
		beego.Error(err)
	}
	// 检查当前页是否大于总页数
	_, execPInfo.Page = models.GetTotalCurrentPage(execPInfo.PageSize,
		execPInfo.TableName, execPInfo.Page, execPInfo.SearchT, execPInfo.Fields)
	c.Redirect("/exec?p="+execPInfo.Page+"&sprojectName="+
		execPInfo.SearchT[0]+"&sversion="+execPInfo.SearchT[1]+
		"&sType="+execPInfo.SearchT[2]+"&startTime="+execPInfo.SearchT[3]+
		"&endTime="+execPInfo.SearchT[4]+"&stargetSer="+execPInfo.SearchT[5], 302)
}

func (c *ExecController) Modify() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	c.TplName = "exec_cmd_modify.html"
	eid := c.Input().Get("eid")
	exec, err := models.GetExec(eid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Exec"] = exec
	c.Data["ExecId"] = eid
	c.Data["IsLogin"] = true
}

func (c *ExecController) Delete() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	eidStr := c.Input().Get("eid")
	err := models.DeleteExecs(eidStr, c.Ctx.GetCookie("uname"))
	if err != nil {
		beego.Error(err)
	}
	// 检查当前页是否大于总页数
	_, execPInfo.Page = models.GetTotalCurrentPage(execPInfo.PageSize,
		execPInfo.TableName, execPInfo.Page, execPInfo.SearchT, execPInfo.Fields)
	c.Redirect("/exec?p="+execPInfo.Page+"&sprojectName="+
		execPInfo.SearchT[0]+"&sversion="+execPInfo.SearchT[1]+
		"&sType="+execPInfo.SearchT[2]+"&startTime="+execPInfo.SearchT[3]+
		"&endTime="+execPInfo.SearchT[4]+"&stargetSer="+execPInfo.SearchT[5], 302)
}

func (c *ExecController) Execute() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	eidStr := c.Input().Get("eid")
	err := models.ExecuteTask(eidStr, c.Ctx.GetCookie("uname"))
	if err != nil {
		beego.Error(err)
	}
	// 检查当前页是否大于总页数
	_, execPInfo.Page = models.GetTotalCurrentPage(execPInfo.PageSize,
		execPInfo.TableName, execPInfo.Page, execPInfo.SearchT, execPInfo.Fields)
	c.Redirect("/exec?p="+execPInfo.Page+"&sprojectName="+
		execPInfo.SearchT[0]+"&sversion="+execPInfo.SearchT[1]+
		"&sType="+execPInfo.SearchT[2]+"&startTime="+execPInfo.SearchT[3]+
		"&endTime="+execPInfo.SearchT[4]+"&stargetSer="+execPInfo.SearchT[5], 302)

}
