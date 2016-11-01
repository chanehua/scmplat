package models

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

// 获取总页数以及根据总页数与当前页大小比对来确定当前页值,给权限管理模块调用
func TotalCurrentPage(pageSize int64, tableName, page string,
	searchT, fields []string) (int64, string) {
	o := orm.NewOrm()

	// 获取数据总数
	qs := o.QueryTable(tableName).Filter(fields[0], searchT[0]).Filter(fields[1],
		searchT[1]).Filter(fields[2], searchT[2])
	count, err := qs.Count()
	if err != nil {
		return 1, "1"
	}

	// 若总数据量都为0，总页数，当前页就为都为1
	if count == 0 {
		return 1, "1"
	}

	// 获取总页数
	tp := count / pageSize
	if count%pageSize > 0 {
		tp = count/pageSize + 1
	}

	// 根据总页数与当前页大小比对来确定当前页值
	p, _ := strconv.ParseInt(page, 10, 64)
	if p > tp {
		page = strconv.FormatInt(tp, 10)
	}
	return tp, page
}

func GetTotalCurrentPage(pageSize int64, tableName, page string,
	searchT, fields []string) (int64, string) {
	o := orm.NewOrm()

	// 获取数据总数
	qs := o.QueryTable(tableName).Filter(fields[0], searchT[0]).Filter(fields[1],
		searchT[1]).Filter(fields[2], searchT[2]).Filter(fields[3],
		searchT[3]).Filter(fields[4], searchT[4]).Filter(fields[5], searchT[5])
	count, err := qs.Count()
	if err != nil {
		return 1, "1"
	}
	// 若总数据量都为0，总页数，当前页就为都为1
	if count == 0 {
		return 1, "1"
	}

	// 获取总页数
	tp := count / pageSize
	if count%pageSize > 0 {
		tp = count/pageSize + 1
	}

	// 根据总页数与当前页大小比对来确定当前页值
	p, _ := strconv.ParseInt(page, 10, 64)
	if p > tp {
		page = strconv.FormatInt(tp, 10)
	}
	return tp, page
}
