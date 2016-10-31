package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

//执行发布任务
func PublishSer(Pubs []*PubMg) {
	// 保证gorouting进入到发布阶段，再进入下一个主机的发布操作
	ch := make(chan int)

	var srcFileInit, srcFile, dstFile, dstAddr string

	for i, _ := range Pubs {

		dstAddr = Pubs[i].TargetServer + ":" + Pubs[i].SshPort
		// 判断发布类型
		switch Pubs[i].PubType {
		//将后台进程打包以便发布
		case "prc":
			srcDir := Pubs[i].PubSrcDir
			tarDirPath := "shell/" + srcDir + "/app"
			// 判断是否已经有tar包，若有则不再进行打包操作
			files, _ := GetPubFileInfo(tarDirPath)
			for fi, fn := range files {
				if strings.Contains(fn.Name(), ".tar.gz") {
					srcFileInit = fn.Name()
					break
				} else {
					if fi == (len(files) - 1) {
						srcFileInit = "aiprcShell_" + time.Now().Format(timeStamp) + ".tar.gz"
						tarCmdSub := `tar -zcf ` + tarDirPath + "/" + srcFileInit +
							" --exclude=\"*.tar.gz\" -C " + tarDirPath + " ./"
						tarCmd := exec.Command("/bin/sh", "-c", tarCmdSub)
						err := tarCmd.Run()
						if err != nil {
							beego.Error(err)
						}
					}
				}
			}

			srcFile = tarDirPath + "/" + srcFileInit
			if len(Pubs[i].PubDstDir) == 0 {
				dstFile = "app/" + srcFileInit
			} else {
				dstFile = Pubs[i].PubDstDir + "/" + srcFileInit
			}
		case "patch", "plat":
			srcDir := Pubs[i].PubSrcDir + "/" + Pubs[i].Opertor
			tarDirPath := "upload/" + srcDir
			files, fileNum := GetPubFileInfo(tarDirPath)
			if fileNum == 1 {
				srcFileInit = files[0].Name()
			} else {
				// 判断是否已经有tar包，若有则不再进行打包操作
				for fi, fn := range files {
					if strings.Contains(fn.Name(), ".tar.gz") {
						srcFileInit = fn.Name()
						break
					} else {
						if fi == (len(files) - 1) {
							srcFileInit = Pubs[i].PubType + "_" +
								time.Now().Format(timeStamp) + ".tar.gz"
							tarCmdSub := `tar -zcf ` + tarDirPath + "/" + srcFileInit +
								" --exclude=\"*.tar.gz\" -C " + tarDirPath + " ./"
							tarCmd := exec.Command("/bin/sh", "-c", tarCmdSub)
							err := tarCmd.Run()
							if err != nil {
								beego.Error(err)
							}
						}
					}
				}

			}

			srcFile = tarDirPath + "/" + srcFileInit
			if len(Pubs[i].PubDstDir) == 0 {
				dstFile = Pubs[i].PubType + "/" + srcFileInit
			} else {
				dstFile = Pubs[i].PubDstDir + "/" + srcFileInit
			}
		default:
			err := fmt.Errorf("publish type incorrect,please rechoice!")
			beego.Error(err)

		}

		// 并发处理发布以及包的解压
		go func(i int) {
			//连接到目标服务器
			server := NewServerConn(dstAddr, Pubs[i].SshUser, Pubs[i].SshPwd, Pubs[i].SshKey)
			defer server.Close()
			_, err := server.getSshConnect()
			if err != nil {
				ch <- 1
				beego.Error(err)
			} else {
				ch <- 1
				//将源文件传输到目标主机
				err = server.CopyFile(srcFile, dstFile)
				if err != nil {
					beego.Error(err)
				}

			}

		}(i)
		<-ch
	}

}

func GetPubFileInfo(pubDir string) (files []os.FileInfo, fileNum int) {
	files, _ = ioutil.ReadDir(pubDir)
	fileNum = len(files)
	return files, fileNum
}
