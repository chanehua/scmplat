package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

const (
	timeForm   = "2006-01-02 15:04:05"
	timeStamp  = "200601021504"
	timeStamp1 = "20060102"
)

type UserInfo struct {
	Id     int64
	Uname  string
	OpRole string
}

type SecMg struct {
	SecId     int64  `orm:"pk;auto"`
	RoleName  string `orm:"size(50)"`
	OperID    string `orm:size(50)`   // 前端定义菜单或按钮的ID
	DplStatus string `orm:"size(50)"` // 控制菜单显示或隐藏
}

type ProcManerge struct {
	ProcId                  int64     `orm:"pk;auto"`
	PROJECT_NAME            string    `orm:"size(50)"`
	VERSION                 string    `orm:"size(50)"`
	OPERATOR                string    `orm:"size(50)"`
	DEPLOY_NAME             string    `orm:"size(50)"`
	DEPLOY_HOST             string    `orm:"size(50)"`
	DEPLOY_TYPE             string    `orm:"size(10)"`
	DEPLOY_SHELL_NAME       string    `orm:"size(50)"`
	DEPLOY_MEM_ARGS         string    `orm:"size(50)"`
	DEPLOY_APPFRAME_KEY     string    `orm:"size(100)"`
	DEPLOY_APPFRAME_VALUE   string    `orm:"size(50)"`
	DEPLOY_SHELL_PARAMS     string    `orm:"size(100);null"`
	DEPLOY_PARAMETER_VALUE  string    `orm:"size(100)"`
	DEPLOY_CONNECTION_VALUE string    `orm:"size(100);null"`
	DEPLOY_STARTUP_CLASS    string    `orm:"size(150)"`
	EXERELAT                string    `orm:"size(150);null"`
	EXEDATASOURCE           string    `orm:"size(150);null"`
	CREATE_TIME             time.Time `orm:"index"`
	REMARK                  string    `orm:"null"`
	STANDARD                string    `orm:"size(10)"`
}

type PubMg struct {
	PubId        int64     `orm:"pk;auto"`
	ProjectName  string    `orm:"size(50)"`
	Version      string    `orm:"size(50)"`
	PubType      string    `orm:"size(20)"`
	TargetServer string    `orm:"size(50)"`
	SshUser      string    `orm:"size(20)"` // 发布帐号
	SshPort      string    `orm:"size(10)"` // SSH端口
	SshPwd       string    `orm:"size(50)"`
	SshKey       string    `orm:"size(200)"` // SSH KEY路径
	PubSrcDir    string    `orm:"size(200)"`
	PubDstDir    string    `orm:"size(200)"`
	PubTime      time.Time `orm:"index"`
	Opertor      string    `orm:"size(50)"`
	Remark       string    `orm:"null"`
}

type ExecMg struct {
	ExecId       int64     `orm:"pk;auto"`
	ProjectName  string    `orm:"size(50)"`
	Version      string    `orm:"size(50)"`
	TargetServer string    `orm:"size(50)"`
	OperType     string    `orm:"size(50)"` // 启停类型
	SshUser      string    `orm:"size(20)"` // 用户名
	SshPort      string    `orm:"size(10)"` // SSH端口
	SshPwd       string    `orm:"size(50)"`
	SshKey       string    `orm:"size(200)"` // SSH KEY路径
	ExecScript   string    `orm:"size(200)"`
	ExecTime     time.Time `orm:"index"`
	Opertor      string    `orm:"size(50)"`
	Remark       string    `orm:"null"`
}

type DockerMg struct {
	DockerId     int64     `orm:"pk;auto;index"`
	ProjectName  string    `orm:"size(50)"`
	Version      string    `orm:"size(50)"`
	TargetServer string    `orm:"size(50)"`
	OperType     string    `orm:"size(50)"` // 操作类型
	SshUser      string    `orm:"size(20)"` // 用户名
	SshPort      string    `orm:"size(10)"` // SSH端口
	SshPwd       string    `orm:"size(50)"`
	SshKey       string    `orm:"size(200)"` // SSH KEY路径
	OsType       string    `orm:"size(20)"`  // 系统类型
	Ctime        time.Time `orm:"index"`
	Opertor      string    `orm:"size(50)"`
	Remark       string    `orm:"null"`
}

func RegisterDB() {
	//get connect db infomation
	dbSection, err := beego.AppConfig.GetSection("db")
	if err != nil {
		fmt.Printf("get db connet db infomation error:%s\n", err)
		return
	}
	dbhost := dbSection["dbhost"]
	dbport := dbSection["dbport"]
	dbuser := dbSection["dbuser"]
	passwd := dbSection["passwd"]
	dbname := dbSection["dbname"]
	dbdriver := "mysql"
	dsn := dbuser + ":" + passwd + "@tcp(" + dbhost +
		":" + dbport + ")/" + dbname + "?charset=utf8"

	//register db models
	orm.RegisterModel(new(UserInfo), new(ProcManerge),
		new(PubMg), new(ExecMg), new(SecMg), new(DockerMg))
	//register db driver("mysql" is default driver,can ignor)
	orm.RegisterDriver(dbdriver, orm.DRMySQL)
	//register database
	orm.RegisterDataBase("default", dbdriver, dsn)
}

func AddUserInfo(uname string) error {
	o := orm.NewOrm()

	uinfo := &UserInfo{Uname: uname, OpRole: "operator"}

	// 查询数据
	qs := o.QueryTable("user_info")
	err := qs.Filter("uname", uname).One(uinfo)
	if err == nil {
		return err
	}

	// 插入数据
	_, err = o.Insert(uinfo)
	if err != nil {
		return err
	}
	return nil
}

func AddProcManerge(procMg map[string]string, uname string) error {
	o := orm.NewOrm()
	createTime, _ := time.Parse(timeForm, procMg["cTime"])
	procMgs := &ProcManerge{
		PROJECT_NAME:           procMg["projectName"],
		VERSION:                procMg["version"],
		OPERATOR:               procMg["oper"],
		DEPLOY_NAME:            procMg["name"],
		DEPLOY_HOST:            procMg["host"],
		DEPLOY_TYPE:            procMg["type"],
		DEPLOY_SHELL_NAME:      procMg["shellName"],
		DEPLOY_MEM_ARGS:        procMg["memARGS"],
		DEPLOY_APPFRAME_KEY:    procMg["appframeKey"],
		DEPLOY_APPFRAME_VALUE:  procMg["appframeValue"],
		DEPLOY_SHELL_PARAMS:    procMg["shellParams"],
		DEPLOY_PARAMETER_VALUE: procMg["parameterValue"],
		DEPLOY_STARTUP_CLASS:   procMg["startupClass"],
		EXERELAT:               procMg["exeRelat"],
		EXEDATASOURCE:          procMg["exeDataSource"],
		CREATE_TIME:            createTime,
		REMARK:                 procMg["remark"],
	}

	_, err := o.Insert(procMgs)
	if err != nil {
		return err
	}
	// 调用操作记录函数
	err = LogRecord(uname, "ProcessMg", "add", procMgs)
	if err != nil {
		return err
	}
	return nil
}

func GetProcList(page string, pageSize int64, searchT,
	fields []string) (procMgs []*ProcManerge, err error) {
	o := orm.NewOrm()

	p, _ := strconv.ParseInt(page, 10, 64)
	procMgs = make([]*ProcManerge, 0)

	qs := o.QueryTable("proc_manerge").Filter(fields[0], searchT[0]).Filter(fields[1],
		searchT[1]).Filter(fields[2], searchT[2]).Filter(fields[3],
		searchT[3]).Filter(fields[4], searchT[4]).Filter(fields[5],
		searchT[5]).OrderBy("proc_id")
	if pageSize > 0 {
		qs = qs.Limit(pageSize, (p-1)*pageSize)
	}
	_, err = qs.All(&procMgs)
	return procMgs, err
}

func GetProc(pid string) (*ProcManerge, error) {
	pidNum, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	proc := new(ProcManerge)
	qs := o.QueryTable("proc_manerge")
	err = qs.Filter("proc_id", pidNum).One(proc)
	if err != nil {
		return nil, err
	}
	return proc, err
}

func ModifyProc(procMg map[string]string, uname string) error {
	pidNum, err := strconv.ParseInt(procMg["ProcId"], 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	createTime, _ := time.Parse(timeForm, procMg["cTime"])
	proc := &ProcManerge{ProcId: pidNum}
	if o.Read(proc) == nil {
		proc.PROJECT_NAME = procMg["projectName"]
		proc.VERSION = procMg["version"]
		proc.OPERATOR = procMg["oper"]
		proc.DEPLOY_NAME = procMg["name"]
		proc.DEPLOY_HOST = procMg["host"]
		proc.DEPLOY_TYPE = procMg["type"]
		proc.DEPLOY_SHELL_NAME = procMg["shellName"]
		proc.DEPLOY_MEM_ARGS = procMg["memARGS"]
		proc.DEPLOY_APPFRAME_KEY = procMg["appframeKey"]
		proc.DEPLOY_APPFRAME_VALUE = procMg["appframeValue"]
		proc.DEPLOY_SHELL_PARAMS = procMg["shellParams"]
		proc.DEPLOY_PARAMETER_VALUE = procMg["parameterValue"]
		proc.DEPLOY_STARTUP_CLASS = procMg["startupClass"]
		proc.EXERELAT = procMg["exeRelat"]
		proc.EXEDATASOURCE = procMg["exeDataSource"]
		proc.CREATE_TIME = createTime
		proc.REMARK = procMg["remark"]

		_, err = o.Update(proc)
		if err != nil {
			return err
		}
	}
	// 调用操作记录函数
	err = LogRecord(uname, "ProcessMg", "modify", proc)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProcs(pidstr, uname string) error {
	var pidNum int64
	var err error
	var proc *ProcManerge
	pids := strings.Split(pidstr, ",")
	o := orm.NewOrm()
	for i, _ := range pids {
		pidNum, err = strconv.ParseInt(pids[i], 10, 64)
		if err != nil {
			return err
		}
		proc = &ProcManerge{ProcId: pidNum}
		if o.Read(proc) == nil {
			_, err = o.Delete(proc)
			if err != nil {
				return err
			}
		}
		// 调用操作记录函数
		err = LogRecord(uname, "ProcessMg", "delete", proc)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateProcs(pidStr, uname string) error {
	o := orm.NewOrm()

	Procs := make([]*ProcManerge, 0)
	// 获取 QueryBuilder 对象. 需要指定数据库驱动参数。
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return err
	}
	// 构建查询对象
	qb.Select("*").From("SCMPLAT.proc_manerge").Where("proc_id").In(pidStr)
	// 导出SQL语句
	sql := qb.String()

	// 执行SQL语句
	o.Raw(sql).QueryRows(&Procs)
	err = CreateShell(Procs)
	if err != nil {
		return err
	}
	// 调用操作记录函数
	err = LogRecord(uname, "ProcessMg", "modify", Procs)
	if err != nil {
		return err
	}
	return nil
}

func GetPubList(page string, pageSize int64, searchT,
	fields []string) (pubs []*PubMg, err error) {
	o := orm.NewOrm()

	p, _ := strconv.ParseInt(page, 10, 64)
	pubs = make([]*PubMg, 0)

	qs := o.QueryTable("pub_mg").Filter(fields[0], searchT[0]).Filter(fields[1],
		searchT[1]).Filter(fields[2], searchT[2]).Filter(fields[3],
		searchT[3]).Filter(fields[4], searchT[4]).Filter(fields[5],
		searchT[5]).OrderBy("pub_id")
	if pageSize > 0 {
		qs = qs.Limit(pageSize, (p-1)*pageSize)
	}
	_, err = qs.All(&pubs)
	return pubs, err
}

func AddPubTask(pubTask map[string]string, uname string) error {
	o := orm.NewOrm()
	pubTime, _ := time.Parse(timeForm, pubTask["pubTime"])
	pub := &PubMg{
		ProjectName:  pubTask["projectName"],
		Version:      pubTask["version"],
		PubType:      pubTask["pubType"],
		TargetServer: pubTask["targetSer"],
		SshUser:      pubTask["sshUser"],
		SshPwd:       pubTask["sshPwd"],
		SshPort:      pubTask["sshPort"],
		SshKey:       pubTask["sshKey"],
		PubSrcDir:    pubTask["pubSrcDir"],
		PubDstDir:    pubTask["pubDstDir"],
		PubTime:      pubTime,
		Opertor:      pubTask["operator"],
		Remark:       pubTask["remark"],
	}

	_, err := o.Insert(pub)
	if err != nil {
		return err
	}
	// 调用操作记录函数
	err = LogRecord(uname, "PublishMg", "add", pub)
	if err != nil {
		return err
	}
	return nil
}

func GetPub(pid string) (*PubMg, error) {
	pidNum, err := strconv.ParseInt(pid, 10, 64)
	if err != nil {
		return nil, err
	}

	o := orm.NewOrm()
	pub := new(PubMg)
	qs := o.QueryTable("pub_mg")
	err = qs.Filter("pub_id", pidNum).One(pub)
	if err != nil {
		return nil, err
	}
	return pub, nil

}

func ModifyPub(pubTask map[string]string, uname string) error {
	pidNum, err := strconv.ParseInt(pubTask["pid"], 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	pubTime, _ := time.Parse(timeForm, pubTask["pubTime"])
	pub := &PubMg{PubId: pidNum}
	if o.Read(pub) == nil {
		pub.ProjectName = pubTask["projectName"]
		pub.Version = pubTask["version"]
		pub.PubType = pubTask["pubType"]
		pub.TargetServer = pubTask["targetSer"]
		pub.SshUser = pubTask["sshUser"]
		pub.SshPwd = pubTask["sshPwd"]
		pub.SshPort = pubTask["sshPort"]
		pub.PubSrcDir = pubTask["pubSrcDir"]
		pub.PubDstDir = pubTask["pubDstDir"]
		pub.PubTime = pubTime
		pub.Opertor = pubTask["operator"]
		pub.Remark = pubTask["remark"]
		// 判断密钥是否更新
		if pubTask["sshKey"] != "" {
			pub.SshKey = pubTask["sshKey"]
		}
		_, err := o.Update(pub)
		if err != nil {
			return err
		}
	}
	// 调用操作记录函数
	err = LogRecord(uname, "PublishMg", "modify", pub)
	if err != nil {
		return err
	}
	return nil
}

func DeletePubs(pidstr, uname string) error {
	var pidNum int64
	var err error
	var pub *PubMg
	pids := strings.Split(pidstr, ",")
	o := orm.NewOrm()
	for i, _ := range pids {
		pidNum, err = strconv.ParseInt(pids[i], 10, 64)
		if err != nil {
			return err
		}
		pub = &PubMg{PubId: pidNum}
		if o.Read(pub) == nil {
			_, err = o.Delete(pub)
			if err != nil {
				return err
			}
		}
		// 调用操作记录函数
		err = LogRecord(uname, "PublishMg", "delete", pub)
		if err != nil {
			return err
		}
	}
	return nil
}

func PublishTask(pidStr, uname string) error {
	o := orm.NewOrm()

	Pubs := make([]*PubMg, 0)
	// 获取 QueryBuilder 对象. 需要指定数据库驱动参数。
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return err
	}
	// 构建查询对象
	qb.Select("*").From("SCMPLAT.pub_mg").Where("pub_id").In(pidStr)
	// 导出SQL语句
	sql := qb.String()

	// 执行SQL语句
	o.Raw(sql).QueryRows(&Pubs)
	// 调用执行发布函数
	PublishSer(Pubs)

	// 调用操作记录函数
	err = LogRecord(uname, "PublishMg", "publish", Pubs)
	if err != nil {
		return err
	}
	return nil
}

func GetExecList(page string, pageSize int64, searchT,
	fields []string) (execs []*ExecMg, err error) {
	o := orm.NewOrm()

	p, _ := strconv.ParseInt(page, 10, 64)
	execs = make([]*ExecMg, 0)

	qs := o.QueryTable("exec_mg").Filter(fields[0], searchT[0]).Filter(fields[1],
		searchT[1]).Filter(fields[2], searchT[2]).Filter(fields[3],
		searchT[3]).Filter(fields[4], searchT[4]).Filter(fields[5],
		searchT[5]).OrderBy("exec_id")
	if pageSize > 0 {
		qs = qs.Limit(pageSize, (p-1)*pageSize)
	}
	_, err = qs.All(&execs)
	return execs, err
}

func AddExecTask(execTask map[string]string, uname string) error {
	o := orm.NewOrm()
	execTime, _ := time.Parse(timeForm, execTask["execTime"])

	exec := &ExecMg{
		ProjectName:  execTask["projectName"],
		Version:      execTask["version"],
		OperType:     execTask["operType"],
		TargetServer: execTask["targetSer"],
		SshUser:      execTask["sshUser"],
		SshPwd:       execTask["sshPwd"],
		SshPort:      execTask["sshPort"],
		SshKey:       execTask["sshKey"],
		ExecScript:   execTask["execScript"],
		ExecTime:     execTime,
		Opertor:      execTask["operator"],
		Remark:       execTask["remark"],
	}

	_, err := o.Insert(exec)
	if err != nil {
		return err
	}
	// 调用操作记录函数
	err = LogRecord(uname, "StartStopMg", "add", exec)
	if err != nil {
		return err
	}
	return nil
}

func GetExec(eid string) (*ExecMg, error) {
	eidNum, err := strconv.ParseInt(eid, 10, 64)
	if err != nil {
		return nil, err
	}

	exec := new(ExecMg)
	o := orm.NewOrm()
	qs := o.QueryTable("exec_mg")
	err = qs.Filter("exec_id", eidNum).One(exec)
	if err != nil {
		return nil, err
	}

	return exec, err
}

// 修改执行任务
func ModifyExec(execTask map[string]string, uname string) error {
	eidNum, err := strconv.ParseInt(execTask["eid"], 10, 64)
	if err != nil {
		return err
	}
	execTime, _ := time.Parse(timeForm, execTask["execTime"])

	o := orm.NewOrm()
	exec := &ExecMg{ExecId: eidNum}
	if o.Read(exec) == nil {
		exec.ProjectName = execTask["projectName"]
		exec.Version = execTask["version"]
		exec.OperType = execTask["operType"]
		exec.TargetServer = execTask["targetSer"]
		exec.SshUser = execTask["sshUser"]
		exec.SshPwd = execTask["sshPwd"]
		exec.SshPort = execTask["sshPort"]
		exec.ExecScript = execTask["execScript"]
		exec.ExecTime = execTime
		exec.Opertor = execTask["operator"]
		exec.Remark = execTask["remark"]
		// 判断密钥是否更新
		if execTask["sshKey"] != "" {
			exec.SshKey = execTask["sshKey"]
		}

		_, err = o.Update(exec)
		if err != nil {
			return err
		}
	}
	// 调用操作记录函数
	err = LogRecord(uname, "StartStopMg", "modify", exec)
	if err != nil {
		return err
	}
	return nil
}

// 删除执行任务
func DeleteExecs(eidStr, uname string) error {
	var eidNum int64
	var err error
	var exec *ExecMg
	eids := strings.Split(eidStr, ",")
	o := orm.NewOrm()

	for i, _ := range eids {
		eidNum, err = strconv.ParseInt(eids[i], 10, 64)
		if err != nil {
			return err
		}
		exec = &ExecMg{ExecId: eidNum}
		if o.Read(exec) == nil {
			_, err = o.Delete(exec)
			if err != nil {
				return err
			}
		}
		// 调用操作记录函数
		err = LogRecord(uname, "StartStopMg", "delete", exec)
		if err != nil {
			return err
		}
	}
	return nil
}

func ExecuteTask(eidStr, uname string) error {
	o := orm.NewOrm()

	Execs := make([]*ExecMg, 0)
	// 获取 QueryBuilder 对象. 需要指定数据库驱动参数。
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return err
	}
	// 构建查询对象
	qb.Select("*").From("SCMPLAT.exec_mg").Where("exec_id").In(eidStr)
	// 导出SQL语句
	sql := qb.String()

	// 执行SQL语句
	o.Raw(sql).QueryRows(&Execs)
	// 调用执行远程主机启停脚本函数
	ExecuteCmd(Execs)

	// 调用操作记录函数
	err = LogRecord(uname, "StartStopMg", "execute", Execs)
	if err != nil {
		return err
	}
	return nil

}

func GetSecList(page string, pageSize int64, searchT,
	fields []string) (secs []*SecMg, err error) {
	o := orm.NewOrm()

	p, _ := strconv.ParseInt(page, 10, 64)
	secs = make([]*SecMg, 0)

	qs := o.QueryTable("sec_mg").Filter(fields[0], searchT[0]).Filter(fields[1],
		searchT[1]).Filter(fields[2], searchT[2]).OrderBy("sec_id")
	if pageSize > 0 {
		qs = qs.Limit(pageSize, (p-1)*pageSize)
	}
	_, err = qs.All(&secs)
	return secs, err
}

// 增加权限
func AddSec(secItem map[string]string) error {
	o := orm.NewOrm()

	sec := &SecMg{
		RoleName:  secItem["roleName"],
		OperID:    secItem["operID"],
		DplStatus: secItem["dplStatus"],
	}

	_, err := o.Insert(sec)
	if err != nil {
		return err
	}
	return nil
}

func GetSec(sid string) (*SecMg, error) {
	sidNum, err := strconv.ParseInt(sid, 10, 64)
	if err != nil {
		return nil, err
	}

	sec := new(SecMg)
	o := orm.NewOrm()
	qs := o.QueryTable("sec_mg")
	err = qs.Filter("sec_id", sidNum).One(sec)
	if err != nil {
		return nil, err
	}

	return sec, err
}

// 修改权限
func ModifySec(secItem map[string]string) error {
	sidNum, err := strconv.ParseInt(secItem["sid"], 10, 64)
	if err != nil {
		return err
	}

	o := orm.NewOrm()
	sec := &SecMg{SecId: sidNum}
	if o.Read(sec) == nil {
		sec.RoleName = secItem["roleName"]
		sec.OperID = secItem["operID"]
		sec.DplStatus = secItem["dplStatus"]

		_, err = o.Update(sec)
		if err != nil {
			return err
		}
	}
	return nil
}

// 删除权限
func DeleteSecs(sidStr string) error {
	var sidNum int64
	var err error
	var sec *SecMg
	sids := strings.Split(sidStr, ",")
	o := orm.NewOrm()

	for i, _ := range sids {
		sidNum, err = strconv.ParseInt(sids[i], 10, 64)
		if err != nil {
			return err
		}
		sec = &SecMg{SecId: sidNum}
		if o.Read(sec) == nil {
			_, err = o.Delete(sec)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 权限校验
func CheckSec(operId map[string]string, uname string) error {
	o := orm.NewOrm()

	// 根据操作员获取操作员角色
	role := &UserInfo{Uname: uname}
	qs := o.QueryTable("user_info")
	err := qs.Filter("uname", uname).One(role)
	if err != nil {
		return err
	}

	// 根据角色获取该角色权限
	secInfo := new(SecMg)
	qs = o.QueryTable("sec_mg")
	for k, _ := range operId {
		qs.Filter("role_name", role.OpRole).Filter("oper_id", k).One(secInfo)
		operId[k] = secInfo.DplStatus

	}
	return err
}

// 获取docker管理任务列表
func GetDockerList(page string, pageSize int64, searchT,
	fields []string) (dockers []*DockerMg, err error) {
	o := orm.NewOrm()

	p, _ := strconv.ParseInt(page, 10, 64)
	dockers = make([]*DockerMg, 0)

	qs := o.QueryTable("docker_mg").Filter(fields[0], searchT[0]).Filter(fields[1],
		searchT[1]).Filter(fields[2], searchT[2]).Filter(fields[3],
		searchT[3]).Filter(fields[4], searchT[4]).Filter(fields[5],
		searchT[5]).OrderBy("docker_id")
	if pageSize > 0 {
		qs = qs.Limit(pageSize, (p-1)*pageSize)
	}
	_, err = qs.All(&dockers)
	return dockers, err
}

// 增加docker管理任务
func AddDocker(dockerList map[string]string, uname string) error {
	o := orm.NewOrm()
	cTime, _ := time.Parse(timeForm, dockerList["cTime"])

	docker := &DockerMg{
		ProjectName:  dockerList["projectName"],
		Version:      dockerList["version"],
		OperType:     dockerList["operType"],
		TargetServer: dockerList["targetSer"],
		SshUser:      dockerList["sshUser"],
		SshPwd:       dockerList["sshPwd"],
		SshPort:      dockerList["sshPort"],
		SshKey:       dockerList["sshKey"],
		OsType:       dockerList["osType"],
		Ctime:        cTime,
		Opertor:      dockerList["operator"],
		Remark:       dockerList["remark"],
	}

	_, err := o.Insert(docker)
	if err != nil {
		return err
	}
	// 调用操作记录函数
	err = LogRecord(uname, "DockerMg", "add", docker)
	if err != nil {
		return err
	}
	return nil
}

func GetDocker(did string) (*DockerMg, error) {
	didNum, err := strconv.ParseInt(did, 10, 64)
	if err != nil {
		return nil, err
	}

	docker := new(DockerMg)
	o := orm.NewOrm()
	qs := o.QueryTable("docker_mg")
	err = qs.Filter("docker_id", didNum).One(docker)
	if err != nil {
		return nil, err
	}

	return docker, err
}

// 修改docker管理任务
func ModifyDocker(dockerList map[string]string, uname string) error {
	didNum, err := strconv.ParseInt(dockerList["did"], 10, 64)
	if err != nil {
		return err
	}
	cTime, _ := time.Parse(timeForm, dockerList["cTime"])

	o := orm.NewOrm()
	docker := &DockerMg{DockerId: didNum}
	if o.Read(docker) == nil {
		docker.ProjectName = dockerList["projectName"]
		docker.Version = dockerList["version"]
		docker.OperType = dockerList["operType"]
		docker.TargetServer = dockerList["targetSer"]
		docker.SshUser = dockerList["sshUser"]
		docker.SshPwd = dockerList["sshPwd"]
		docker.SshPort = dockerList["sshPort"]
		docker.OsType = dockerList["osType"]
		docker.Ctime = cTime
		docker.Opertor = dockerList["operator"]
		docker.Remark = dockerList["remark"]
		// 判断密钥是否更新
		if dockerList["sshKey"] != "" {
			docker.SshKey = dockerList["sshKey"]
		}

		_, err = o.Update(docker)
		if err != nil {
			return err
		}
	}
	// 调用操作记录函数
	err = LogRecord(uname, "DockerMg", "modify", docker)
	if err != nil {
		return err
	}
	return nil
}

// 删除docker管理任务
func DeleteDockers(didStr, uname string) error {
	var didNum int64
	var err error
	var docker *DockerMg
	dids := strings.Split(didStr, ",")
	o := orm.NewOrm()

	for i, _ := range dids {
		didNum, err = strconv.ParseInt(dids[i], 10, 64)
		if err != nil {
			return err
		}
		docker = &DockerMg{DockerId: didNum}
		if o.Read(docker) == nil {
			_, err = o.Delete(docker)
			if err != nil {
				return err
			}
		}
		// 调用操作记录函数
		err = LogRecord(uname, "DockerMg", "delete", docker)
		if err != nil {
			return err
		}
	}
	return nil
}

func StartTask(didStr, uname string) error {
	o := orm.NewOrm()

	Dockers := make([]*DockerMg, 0)
	// 获取 QueryBuilder 对象. 需要指定数据库驱动参数。
	qb, err := orm.NewQueryBuilder("mysql")
	if err != nil {
		return err
	}
	// 构建查询对象
	qb.Select("*").From("SCMPLAT.docker_mg").Where("docker_id").In(didStr)
	// 导出SQL语句
	sql := qb.String()

	// 执行SQL语句
	o.Raw(sql).QueryRows(&Dockers)
	// 到远程主机进行docker管理操作
	ExecDocker(Dockers)

	// 调用操作记录函数
	err = LogRecord(uname, "DockerMg", "start", Dockers)
	if err != nil {
		return err
	}
	return nil

}
