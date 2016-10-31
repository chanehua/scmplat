package controllers

type PageInfo struct {
	Page      string
	PageSize  int64
	TableName string
	Fields    []string
	SearchT   []string
}
