package routers

import (
	"github.com/astaxie/beego"
	"github.com/chanehua/scmplat/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/proc", &controllers.ProcController{})
	beego.Router("/pub", &controllers.PubController{})
	beego.Router("/exec", &controllers.ExecController{})
	beego.Router("/sec", &controllers.SecController{})
	beego.Router("/docker", &controllers.DockerController{})
	beego.AutoRouter(&controllers.ProcController{})
	beego.AutoRouter(&controllers.PubController{})
	beego.AutoRouter(&controllers.ExecController{})
	beego.AutoRouter(&controllers.SecController{})
	beego.AutoRouter(&controllers.DockerController{})
}
