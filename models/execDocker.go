package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"sort"
	"strings"
)

func ExecDocker(Dockers []*DockerMg) {
	var dstAddr string
	var err error
	for i, _ := range Dockers {
		// 并发执行远程主机启停脚本
		go func(i int) {
			dstAddr = Dockers[i].TargetServer + ":" + Dockers[i].SshPort
			//连接到目标服务器
			server := NewServerConn(dstAddr, Dockers[i].SshUser, Dockers[i].SshPwd, Dockers[i].SshKey)
			defer server.Close()
			_, err = server.getSshConnect()
			if err != nil {
				beego.Error(err)
			} else {
				// 获取执行命令
				execCmd, err := GetExeCmd(Dockers[i].OsType, Dockers[i].OperType)
				if err != nil {
					beego.Error(err)
				} else {
					execRes, err := server.RunCmd(execCmd)
					if err != nil {
						err = fmt.Errorf("In %s execute cmd Error: %v", dstAddr, err)
						beego.Error(err)
					} else {
						switch Dockers[i].OperType {
						case "uninstall":
							if len(execRes) != 0 {
								uninstallCmd := "yum -y remove " +
									strings.Replace(execRes, "\n", " ", -1)
								_, err = server.RunCmd(uninstallCmd)
								if err != nil {
									err = fmt.Errorf("In %s uninstall docker Error: %v",
										dstAddr, err)
								}
							}

						case "install", "upgrade":
							execRes = fmt.Sprintf("In %s execute result is: %s",
								dstAddr, execRes)
							beego.Info(execRes)
						}
					}
				}
			}

		}(i)

	}
}

// 根据不同的操作系统获取不同的执行命令
func GetExeCmd(osType, operType string) (execCmd string, err error) {
	osOper := osType + "_" + operType
	var sectionInfo map[string]string
	execCmds := make([]string, 0)
	sectionInfo, err = beego.AppConfig.GetSection(osOper)
	if err != nil {
		return "", err
	}

	// 对map排序
	sortKey := make([]string, 0)
	for k, _ := range sectionInfo {
		sortKey = append(sortKey, k)
	}
	sort.Strings(sortKey)

	// 根据排好序的map拼接执行命令
	for _, k := range sortKey {
		execCmds = append(execCmds, sectionInfo[k])
	}
	execCmd = strings.Join(execCmds, " && ")
	return execCmd, err

}
