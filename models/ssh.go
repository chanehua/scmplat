package models

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type ServerConn struct {
	addr       string
	user       string
	pwd        string
	key        string
	conn       *ssh.Client
	sftpClient *sftp.Client
}

func NewServerConn(addr, user, pwd, key string) *ServerConn {
	return &ServerConn{
		addr: addr,
		user: user,
		pwd:  pwd,
		key:  key,
	}
}

// 连接ssh服务器
func (s *ServerConn) getSshConnect() (*ssh.Client, error) {
	if s.conn != nil {
		return s.conn, nil
	}
	config := ssh.ClientConfig{
		User: s.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.pwd),
		},
		Timeout: 15 * time.Second,
	}

	keys := []ssh.Signer{}
	if pk, err := readPrivateKey(s.key); err == nil {
		keys = append(keys, pk)
	}
	config.Auth = append(config.Auth, ssh.PublicKeys(keys...))

	conn, err := ssh.Dial("tcp", s.addr, &config)
	if err != nil {
		return nil, fmt.Errorf("无法连接到服务器 [%s]: %v", s.addr, err)
	}
	s.conn = conn
	return s.conn, nil
}

func readPrivateKey(path string) (ssh.Signer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return ssh.ParsePrivateKey(b)
}

// 返回sftp连接
func (s *ServerConn) getSftpConnect() (*sftp.Client, error) {
	if s.sftpClient != nil {
		return s.sftpClient, nil
	}

	conn, err := s.getSshConnect()
	if err != nil {
		return nil, err
	}

	s.sftpClient, err = sftp.NewClient(conn, sftp.MaxPacket(1<<15))
	return s.sftpClient, err
}

// 关闭连接
func (s *ServerConn) Close() {
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
	if s.sftpClient != nil {
		s.sftpClient.Close()
		s.sftpClient = nil
	}
}

// 在远程服务器执行命令
func (s *ServerConn) RunCmd(cmd string) (string, error) {

	conn, err := s.getSshConnect()
	if err != nil {
		return "", err
	}

	session, err := conn.NewSession()
	if err != nil {
		return "", fmt.Errorf("创建会话失败: %v", err)
	}
	defer session.Close()

	var buf bytes.Buffer

	session.Stdout = &buf
	session.Stdin = &buf

	if err := session.Run(cmd); err != nil {
		return "", fmt.Errorf("执行命令失败: %v", err)
	}

	return buf.String(), nil

}

// 拷贝本机文件到远程服务器
func (s *ServerConn) CopyFile(srcFile, dstFile string) error {
	client, err := s.getSftpConnect()
	if err != nil {
		return err
	}

	toPath := path.Dir(dstFile)
	if _, err := s.RunCmd("mkdir -p " + toPath); err != nil {
		return fmt.Errorf("创建目录失败：%v", err)
	}

	f, err := os.Open(srcFile)
	if err != nil {
		return fmt.Errorf("打开本地文件失败: %v", err)
	}
	defer f.Close()

	w, err := client.Create(dstFile)
	if err != nil {
		return fmt.Errorf("创建文件失败 [%s]: %v", dstFile, err)
	}
	defer w.Close()

	n, err := io.Copy(w, f)
	if err != nil {
		return fmt.Errorf("拷贝文件失败: %v", err)
	}

	fstat, _ := f.Stat()
	if fstat.Size() != n {
		return fmt.Errorf("写入文件大小错误，源文件大小：%d, 写入大小：%d", fstat.Size(), n)
	}
	err = s.ExecRemoteCmd(dstFile)
	if err != nil {
		return err
	}
	return nil
}

// 在远程主机执行命令
func (s *ServerConn) ExecRemoteCmd(dstFile string) error {
	dstDir, dstFileName := path.Split(dstFile)

	//远程主机上进行备份
	bakName := strings.Replace(dstDir, "/", "_", -1) +
		time.Now().Format(timeStamp) + "_bak.tar.gz"
	bakCmd := "cd " + dstDir + " && " +
		"tar -zcf " + bakName + " --exclude=\"*.log\" --exclude=\"*.tar.gz\" *"
	_, err := s.RunCmd(bakCmd)
	if err != nil {
		return fmt.Errorf("执行备份操作错误：%v", err)
	}

	//远程主机上进行解压
	fileExt := path.Ext(dstFileName)
	switch fileExt {
	case ".gz":
		tarCmd := "cd " + dstDir + " && " + "tar -zxf " + dstFileName
		_, err = s.RunCmd(tarCmd)
		if err != nil {
			return fmt.Errorf("解压远程服务器tar文件错误：%v", err)
		}
	case ".zip":
		unZip := "cd " + dstDir + " && " + "unzip " + dstFileName
		_, err = s.RunCmd(unZip)
		if err != nil {
			return fmt.Errorf("解压远程服务器zip文件错误：%v", err)
		}
	case ".ear", ".war":
		jarEAR := "cd " + dstDir + " && " + "jar -xf " + dstFileName
		_, err = s.RunCmd(jarEAR)
		if err != nil {
			return fmt.Errorf("解压远程服务器ear or war 文件错误：%v", err)
		}
	default:
		break
	}

	return nil
}

// 执行启停脚本
func (s *ServerConn) ExecRemoteScript(execScript string) error {
	dstPath, dstScript := path.Split(execScript)
	execCmd := "cd " + dstPath + " && " + "./" + dstScript + " > executeScript.log"
	_, err := s.RunCmd(execCmd)
	if err != nil {
		return fmt.Errorf("执行%s脚本发生错误: %v", execScript, err)
	}
	return nil
}
