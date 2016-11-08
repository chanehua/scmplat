package controllers

import (
	"os"
	"path"

	"github.com/astaxie/beego"

	"github.com/chanehua/scmplat/models"
)

type DockerController struct {
	beego.Controller
}

var dockerPInfo = new(PageInfo)

func (c *DockerController) Get() {
	c.Data["IsDocker"] = true
	c.TplName = "dockerMg.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)

	var dockers []*models.DockerMg
	var err error
	dockerPInfo.Page = c.Input().Get("p")
	dockerPInfo.PageSize = 5
	dockerPInfo.TableName = "docker_mg"
	dockerPInfo.Fields = []string{"project_name__icontains", "version__icontains",
		"oper_type__icontains", "ctime__gte", "ctime__lte",
		"target_server__icontains"}
	var tp int64

	// 获取查询条件
	dockerPInfo.SearchT = []string{c.Input().Get("sprojectName"),
		c.Input().Get("sversion"), c.Input().Get("sType"),
		c.Input().Get("startTime"), c.Input().Get("endTime"),
		c.Input().Get("stargetSer"),
	}

	// 判断查询条件是否为空，为空则不进行条件查询，否则进行条件查询
	if dockerPInfo.SearchT[3] == "" {
		dockerPInfo.SearchT[3] = "1900-01-01 00:00:00"
	}
	if dockerPInfo.SearchT[4] == "" {
		dockerPInfo.SearchT[4] = "2388-12-31 23:59:59"
	}
	// 获取总页数以及根据总页数与当前页大小比对来确定当前页值
	tp, dockerPInfo.Page = models.GetTotalCurrentPage(dockerPInfo.PageSize,
		dockerPInfo.TableName, dockerPInfo.Page, dockerPInfo.SearchT, dockerPInfo.Fields)
	// 获取分页数据
	dockers, err = models.GetDockerList(dockerPInfo.Page, dockerPInfo.PageSize,
		dockerPInfo.SearchT, dockerPInfo.Fields)
	if err != nil {
		beego.Error(err)
	}
	c.Data["dockers"] = dockers
	c.Data["PageNo"] = dockerPInfo.Page
	c.Data["TotalPage"] = tp
	c.Data["SprojectName"] = dockerPInfo.SearchT[0]
	c.Data["Sversion"] = dockerPInfo.SearchT[1]
	c.Data["Stype"] = dockerPInfo.SearchT[2]
	c.Data["StartTime"] = dockerPInfo.SearchT[3]
	c.Data["EndTime"] = dockerPInfo.SearchT[4]
	c.Data["StargetSer"] = dockerPInfo.SearchT[5]

}

func (c *DockerController) Post() {
	/*if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}*/

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
	dockerList := map[string]string{
		"did":         c.Input().Get("did"),
		"projectName": c.Input().Get("projectName"),
		"version":     c.Input().Get("version"),
		"operType":    c.Input().Get("operType"),
		"targetSer":   c.Input().Get("targetSer"),
		"sshUser":     c.Input().Get("sshUser"),
		"sshPwd":      c.Input().Get("sshPwd"),
		"sshPort":     c.Input().Get("sshPort"),
		"sshKey":      sshKey,
		"osType":      c.Input().Get("osType"),
		"cTime":       c.Input().Get("cTime"),
		"operator":    c.Input().Get("operator"),
		"remark":      c.Input().Get("remark"),
	}

	if len(dockerList["did"]) == 0 {
		err = models.AddDocker(dockerList, c.Ctx.GetCookie("uname"))
	} else {
		err = models.ModifyDocker(dockerList, c.Ctx.GetCookie("uname"))
	}
	if err != nil {
		beego.Error(err)
	}
	// 检查当前页是否大于总页数
	_, dockerPInfo.Page = models.GetTotalCurrentPage(dockerPInfo.PageSize,
		dockerPInfo.TableName, dockerPInfo.Page, dockerPInfo.SearchT, dockerPInfo.Fields)
	c.Redirect("/docker?p="+dockerPInfo.Page+"&sprojectName="+
		dockerPInfo.SearchT[0]+"&sversion="+dockerPInfo.SearchT[1]+
		"&sType="+dockerPInfo.SearchT[2]+"&startTime="+dockerPInfo.SearchT[3]+
		"&endTime="+dockerPInfo.SearchT[4]+"&stargetSer="+dockerPInfo.SearchT[5], 302)
}

func (c *DockerController) Modify() {
	/*if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}*/

	c.TplName = "dockerMg_modify.html"
	did := c.Input().Get("did")
	docker, err := models.GetDocker(did)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Docker"] = docker
	c.Data["DockerId"] = did
	c.Data["IsLogin"] = true
}

func (c *DockerController) Delete() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	didStr := c.Input().Get("did")
	err := models.DeleteDockers(didStr, c.Ctx.GetCookie("uname"))
	if err != nil {
		beego.Error(err)
	}
	// 检查当前页是否大于总页数
	_, dockerPInfo.Page = models.GetTotalCurrentPage(dockerPInfo.PageSize,
		dockerPInfo.TableName, dockerPInfo.Page, dockerPInfo.SearchT, dockerPInfo.Fields)
	c.Redirect("/docker?p="+dockerPInfo.Page+"&sprojectName="+
		dockerPInfo.SearchT[0]+"&sversion="+dockerPInfo.SearchT[1]+
		"&sType="+dockerPInfo.SearchT[2]+"&startTime="+dockerPInfo.SearchT[3]+
		"&endTime="+dockerPInfo.SearchT[4]+"&stargetSer="+dockerPInfo.SearchT[5], 302)
}

func (c *DockerController) Start() {
	/*if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}*/
	didStr := c.Input().Get("did")
	err := models.StartTask(didStr, c.Ctx.GetCookie("uname"))
	if err != nil {
		beego.Error(err)
	}
	// 检查当前页是否大于总页数
	_, dockerPInfo.Page = models.GetTotalCurrentPage(dockerPInfo.PageSize,
		dockerPInfo.TableName, dockerPInfo.Page, dockerPInfo.SearchT, dockerPInfo.Fields)
	c.Redirect("/docker?p="+dockerPInfo.Page+"&sprojectName="+
		dockerPInfo.SearchT[0]+"&sversion="+dockerPInfo.SearchT[1]+
		"&sType="+dockerPInfo.SearchT[2]+"&startTime="+dockerPInfo.SearchT[3]+
		"&endTime="+dockerPInfo.SearchT[4]+"&stargetSer="+dockerPInfo.SearchT[5], 302)

}
