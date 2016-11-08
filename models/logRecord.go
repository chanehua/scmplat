package models

import (
	"fmt"
	"os"
	"path"
	"time"
)

func LogRecord(oper, operModel, operAction string, operData interface{}) error {
	// 定义日志输出目录以及文件
	logDir := "logs"
	logFile := "operate_" + time.Now().Format(timeStamp1) + ".log"
	os.MkdirAll(logDir, os.ModePerm)
	f, err := os.OpenFile(path.Join(logDir, logFile),
		os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	defer f.Close()
	if err != nil {
		return err
	}
	operTime := time.Now().Format("2006-01-02 15:04:05")
	var records string
	switch operData := operData.(type) {
	case []*PubMg:
		for i, _ := range operData {
			records += fmt.Sprintf("[%s] [%s] [%s] [%s] %v \n",
				operTime, operModel, operAction, oper, operData[i])
		}
	case []*ExecMg:
		for i, _ := range operData {
			records += fmt.Sprintf("[%s] [%s] [%s] [%s] %v \n",
				operTime, operModel, operAction, oper, operData[i])
		}
	case []*ProcManerge:
		for i, _ := range operData {
			records += fmt.Sprintf("[%s] [%s] [%s] [%s] %v \n",
				operTime, operModel, operAction, oper, operData[i])
		}
	case []*DockerMg:
		for i, _ := range operData {
			records += fmt.Sprintf("[%s] [%s] [%s] [%s] %v \n",
				operTime, operModel, operAction, oper, operData[i])
		}
	case *PubMg, *ExecMg, *ProcManerge, *DockerMg:
		records = fmt.Sprintf("[%s] [%s] [%s] [%s] %v \n",
			operTime, operModel, operAction, oper, operData)
	default:
		break
	}
	_, err = f.WriteString(records)
	if err != nil {
		return err
	}
	return nil
}
