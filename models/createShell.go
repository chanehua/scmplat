package models

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	"github.com/astaxie/beego"
)

func CreateShell(Procs []*ProcManerge) error {
	//get proc temlate init info
	tplSection, err := beego.AppConfig.GetSection("proctpl")
	if err != nil {
		return err
	}
	tplDir := tplSection["create.ftl.dir"]
	stdard := tplSection["create.ftl.standard"]
	appDir := tplSection["create.shell.output"]

	err = os.Mkdir(appDir, 0755)
	if os.IsExist(err) { //创建后台进程生成目录
		os.RemoveAll(appDir)
		os.Mkdir(appDir, 0755)
	}

	//call create shell script functions
	for i, _ := range Procs {
		cTime := time.Now()
		Procs[i].STANDARD = stdard
		Procs[i].DEPLOY_TYPE = strings.ToLower(Procs[i].DEPLOY_TYPE)
		Procs[i].CREATE_TIME = cTime
		files, err := GetTplFiles(tplDir, ".tpl")
		if err != nil {
			return err
		}
		for _, f := range files {
			scriptContenx, err := Procs[i].GetTplInfo(f) //获取脚本内容
			if err != nil {
				return err
			}
			err = Procs[i].GenScript(appDir, f, scriptContenx)
			if err != nil {
				return err
			}
		}
		if Procs[i].EXERELAT != "" && Procs[i].EXEDATASOURCE != "" {
			err = GenExeConf(Procs[i].EXERELAT, Procs[i].EXEDATASOURCE,
				Procs[i].DEPLOY_APPFRAME_VALUE, appDir, Procs[i].DEPLOY_HOST)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

//获取模板文件
func GetTplFiles(tplDir, suffix string) ([]string, error) {
	files := make([]string, 0)
	dir, err := ioutil.ReadDir(tplDir)
	if err != nil {
		return files, err
	}
	pthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		if fi.IsDir() { //忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //文件匹配
			files = append(files, tplDir+pthSep+fi.Name())
		}
	}
	return files, nil
}

//将模板文件与数据对象整合，输出到string
func (tplData ProcManerge) GetTplInfo(fName string) ([]byte, error) {
	scriptContenx := make([]byte, 0)
	tpl, err := template.ParseFiles(fName)
	if err != nil {
		return scriptContenx, err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, tplData)
	if err != nil {
		return scriptContenx, err
	}
	scriptContenx = buf.Bytes()
	return scriptContenx, nil
}

//生成脚本
func (tplData ProcManerge) GenScript(appdir, fName string, scriptContenx []byte) error {
	separator := string(os.PathSeparator)
	host := tplData.DEPLOY_HOST
	shellType := tplData.DEPLOY_TYPE
	dirName := tplData.DEPLOY_NAME
	shellName := tplData.DEPLOY_SHELL_NAME
	dirBin := appdir + separator + host + separator + "app" + separator + strings.ToLower(shellType) +
		separator + dirName + separator + "bin" + separator
	dirLog := appdir + separator + host + separator + "app" + separator + strings.ToLower(shellType) +
		separator + dirName + separator + "log" + separator
	err := os.MkdirAll(dirBin, 0755)
	if err != nil {
		return err
	}
	os.MkdirAll(dirLog, 0755)
	switch {
	case strings.Contains(fName, "start"):
		startShellName := dirBin + "start_" + shellName + ".sh"
		ioutil.WriteFile(startShellName, scriptContenx, 0644)
	case strings.Contains(fName, "stop"):
		stopShellName := dirBin + "stop_" + shellName + ".sh"
		ioutil.WriteFile(stopShellName, scriptContenx, 0644)
	case strings.Contains(fName, "monitor"):
		monitorShellName := dirBin + "monitor_" + shellName + ".sh"
		ioutil.WriteFile(monitorShellName, scriptContenx, 0644)
	default:
		err := errors.New("Not find match shell file name!")
		return err
	}
	// 赋予脚本可执行权限
	cmd := "chmod +x `find ./" + appdir + " -name \"*.sh\"`"
	preCmd := exec.Command("/bin/sh", "-c", cmd)
	err = preCmd.Run()
	if err != nil {
		return err
	}
	return nil
}

//生成exe.properties
func GenExeConf(EXERELAT, EXEDATASOURCE, AppframeValue, appdir, host string) error {
	separator := string(os.PathSeparator)
	cfgDir := appdir + separator + host + separator + "config"
	os.MkdirAll(cfgDir, 0755)
	fName := cfgDir + separator + "exe.properties"
	f, err := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	context := AppframeValue + ".relat=" + EXERELAT + "\n" + AppframeValue + ".datasource=" + EXEDATASOURCE + "\n\n"
	f.WriteString(context)
	return nil
}
