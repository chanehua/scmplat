package models

import (
	"github.com/astaxie/beego"
)

func ExecuteCmd(Execs []*ExecMg) {
	var dstAddr string
	var err error
	for i, _ := range Execs {
		// 并发执行远程主机启停脚本
		go func(i int) {
			dstAddr = Execs[i].TargetServer + ":" + Execs[i].SshPort
			//连接到目标服务器
			server := NewServerConn(dstAddr, Execs[i].SshUser, Execs[i].SshPwd, Execs[i].SshKey)
			defer server.Close()
			_, err = server.getSshConnect()
			if err != nil {
				beego.Error(err)
			}

			// 在远程主机执行启停脚本
			err = server.ExecRemoteScript(Execs[i].ExecScript)
			if err != nil {
				beego.Error(err)
			}
		}(i)

	}
}
