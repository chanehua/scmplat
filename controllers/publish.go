package controllers

import (
	"io"
	"os"
	"path"

	"github.com/astaxie/beego"
	"github.com/chanehua/scmplat/models"
)

type PubController struct {
	beego.Controller
}

var pubPInfo = new(PageInfo)

func (c *PubController) Get() {
	c.Data["IsPub"] = true
	c.TplName = "publish.html"
	c.Data["IsLogin"] = checkAccount(c.Ctx)

	var pubs []*models.PubMg
	var err error
	pubPInfo.Page = c.Input().Get("p")
	pubPInfo.PageSize = 5
	pubPInfo.TableName = "pub_mg"
	pubPInfo.Fields = []string{"project_name__icontains", "version__icontains",
		"pub_type__icontains", "pub_time__gte", "pub_time__lte",
		"target_server__icontains"}
	var tp int64

	// 获取查询条件
	pubPInfo.SearchT = []string{c.Input().Get("sprojectName"),
		c.Input().Get("sversion"), c.Input().Get("sType"),
		c.Input().Get("startTime"), c.Input().Get("endTime"),
		c.Input().Get("stargetSer"),
	}
	// 判断查询条件是否为空，为空则不进行条件查询，否则进行条件查询
	if pubPInfo.SearchT[3] == "" {
		pubPInfo.SearchT[3] = "1900-01-01 00:00:00"
	}
	if pubPInfo.SearchT[4] == "" {
		pubPInfo.SearchT[4] = "2388-12-31 23:59:59"
	}
	// 获取总页数以及根据总页数与当前页大小比对来确定当前页值
	tp, pubPInfo.Page = models.GetTotalCurrentPage(pubPInfo.PageSize,
		pubPInfo.TableName, pubPInfo.Page, pubPInfo.SearchT, pubPInfo.Fields)
	// 获取分页数据
	pubs, err = models.GetPubList(pubPInfo.Page, pubPInfo.PageSize,
		pubPInfo.SearchT, pubPInfo.Fields)
	if err != nil {
		beego.Error(err)
	}
	c.Data["pubs"] = pubs
	c.Data["PageNo"] = pubPInfo.Page
	c.Data["TotalPage"] = tp
	c.Data["SprojectName"] = pubPInfo.SearchT[0]
	c.Data["Sversion"] = pubPInfo.SearchT[1]
	c.Data["Stype"] = pubPInfo.SearchT[2]
	c.Data["StartTime"] = pubPInfo.SearchT[3]
	c.Data["EndTime"] = pubPInfo.SearchT[4]
	c.Data["StargetSer"] = pubPInfo.SearchT[5]
	// 权限校验
	uname := c.Ctx.GetCookie("uname")
	if len(uname) == 0 {
		c.Data["SecMg"] = "none"
	} else {
		operId := map[string]string{"SecMg": "none",
			"PubAddButton":     "none",
			"PubModifyButton":  "none",
			"PubDelButton":     "none",
			"PubPublishButton": "none",
		}
		err := models.CheckSec(operId, uname)
		if err != nil {
			beego.Error(err)
		}
		c.Data["SecMg"] = operId["SecMg"]
		c.Data["PubAddButton"] = operId["PubAddButton"]
		c.Data["PubModifyButton"] = operId["PubModifyButton"]
		c.Data["PubDelButton"] = operId["PubDelButton"]
		c.Data["PubPublishButton"] = operId["PubPublishButton"]
	}
}

func (c *PubController) Post() {
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
	pubTask := map[string]string{
		"pid":         c.Input().Get("pid"),
		"projectName": c.Input().Get("projectName"),
		"version":     c.Input().Get("version"),
		"pubType":     c.Input().Get("pubType"),
		"targetSer":   c.Input().Get("targetSer"),
		"sshUser":     c.Input().Get("sshUser"),
		"sshPwd":      c.Input().Get("sshPwd"),
		"sshPort":     c.Input().Get("sshPort"),
		"sshKey":      sshKey,
		"pubSrcDir":   c.Input().Get("pubSrcDir"),
		"pubDstDir":   c.Input().Get("pubDstDir"),
		"pubTime":     c.Input().Get("pubTime"),
		"operator":    c.Input().Get("operator"),
		"remark":      c.Input().Get("remark"),
	}

	// 获取附件
	fhs, err := c.GetFiles("upload")
	if err != nil {
		beego.Error(err)
	}
	// 保存附件
	var fName string
	if fhs != nil {
		// 创建保存附件目录
		saveDir := "upload/" + pubTask["pubSrcDir"] + "/" + pubTask["operator"]
		err = os.Mkdir("upload", os.ModePerm)
		if os.IsExist(err) {
			os.RemoveAll(saveDir)
			os.MkdirAll(saveDir, os.ModePerm)
		} else {
			os.MkdirAll(saveDir, os.ModePerm)
		}

		//循环对上传的每个文件进行处理
		for i, _ := range fhs {
			// 获取文件名
			fName = fhs[i].Filename

			// 结束文件
			file, err := fhs[i].Open()
			if err != nil {
				beego.Error(err)
			} else {
				//保存文件
				defer file.Close()
				f, err := os.Create(path.Join(saveDir, fName))
				if err != nil {
					beego.Error(err)
				} else {
					defer f.Close()
					io.Copy(f, file)
				}
			}
		}
	}
	if len(pubTask["pid"]) == 0 {
		err = models.AddPubTask(pubTask, c.Ctx.GetCookie("uname"))
	} else {
		err = models.ModifyPub(pubTask, c.Ctx.GetCookie("uname"))
	}
	if err != nil {
		beego.Error(err)
	}

	// 检查当前页是否大于总页数
	_, pubPInfo.Page = models.GetTotalCurrentPage(pubPInfo.PageSize,
		pubPInfo.TableName, pubPInfo.Page, pubPInfo.SearchT, pubPInfo.Fields)
	c.Redirect("/pub?p="+pubPInfo.Page+"&sprojectName="+
		pubPInfo.SearchT[0]+"&sversion="+pubPInfo.SearchT[1]+
		"&sType="+pubPInfo.SearchT[2]+"&startTime="+pubPInfo.SearchT[3]+
		"&endTime="+pubPInfo.SearchT[4]+"&stargetSer="+pubPInfo.SearchT[5], 302)
}

func (c *PubController) Modify() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}

	c.TplName = "pub_modify.html"
	pid := c.Input().Get("pid")
	pub, err := models.GetPub(pid)
	if err != nil {
		beego.Error(err)
		c.Redirect("/", 302)
		return
	}
	c.Data["Pub"] = pub
	c.Data["PubId"] = pid
	c.Data["IsLogin"] = true
}

func (c *PubController) Delete() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	pidStr := c.Input().Get("pid")
	err := models.DeletePubs(pidStr, c.Ctx.GetCookie("uname"))
	if err != nil {
		beego.Error(err)
	}
	// 检查当前页是否大于总页数
	_, pubPInfo.Page = models.GetTotalCurrentPage(pubPInfo.PageSize,
		pubPInfo.TableName, pubPInfo.Page, pubPInfo.SearchT, pubPInfo.Fields)
	c.Redirect("/pub?p="+pubPInfo.Page+"&sprojectName="+
		pubPInfo.SearchT[0]+"&sversion="+pubPInfo.SearchT[1]+
		"&sType="+pubPInfo.SearchT[2]+"&startTime="+pubPInfo.SearchT[3]+
		"&endTime="+pubPInfo.SearchT[4]+"&stargetSer="+pubPInfo.SearchT[5], 302)
}

func (c *PubController) Publish() {
	if !checkAccount(c.Ctx) {
		c.Redirect("/login", 302)
		return
	}
	pidStr := c.Input().Get("pid")
	err := models.PublishTask(pidStr, c.Ctx.GetCookie("uname"))
	if err != nil {
		beego.Error(err)
	}
	// 检查当前页是否大于总页数
	_, pubPInfo.Page = models.GetTotalCurrentPage(pubPInfo.PageSize,
		pubPInfo.TableName, pubPInfo.Page, pubPInfo.SearchT, pubPInfo.Fields)
	c.Redirect("/pub?p="+pubPInfo.Page+"&sprojectName="+
		pubPInfo.SearchT[0]+"&sversion="+pubPInfo.SearchT[1]+
		"&sType="+pubPInfo.SearchT[2]+"&startTime="+pubPInfo.SearchT[3]+
		"&endTime="+pubPInfo.SearchT[4]+"&stargetSer="+pubPInfo.SearchT[5], 302)
}
