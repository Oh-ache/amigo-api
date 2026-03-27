/*
Package pkg 提供通用的 SQL 处理函数

包含 SQL 查询条件处理和分页查询等功能
*/
package utils

import (
	"fmt"
	"strings"
)

// DealSQL 处理 SQL 查询条件
// key: 字段名
// value: 查询值
// op: 操作符，支持 =, !=, >, >=, <, <=, like, notlike, in, notin
func DealSQL(key, value, op string) string {
	if value == "" || value == "0" {
		return ""
	}
	// 处理时间字段
	if key == "create_time" || key == "update_time" {
		// 验证值是否为逗号分隔的两个数字
		parts := strings.Split(value, ",")
		if len(parts) != 2 {
			return ""
		}
		start := strings.TrimSpace(parts[0])
		end := strings.TrimSpace(parts[1])
		if start == "" || end == "" {
			return ""
		}
		// 生成时间区间查询，左闭右开
		return fmt.Sprintf(" %s >= %s and %s < %s ", key, start, key, end)
	}

	// 处理其他字段
	switch strings.ToLower(op) {
	case "=":
		return fmt.Sprintf(" %s = '%s' ", key, value)
	case "!=":
		return fmt.Sprintf(" %s != '%s' ", key, value)
	case ">":
		return fmt.Sprintf(" %s > '%s' ", key, value)
	case ">=":
		return fmt.Sprintf(" %s >= '%s' ", key, value)
	case "<":
		return fmt.Sprintf(" %s < '%s' ", key, value)
	case "<=":
		return fmt.Sprintf(" %s <= '%s' ", key, value)
	case "like":
		return fmt.Sprintf(" %s like '%%%s%%' ", key, value)
	case "notlike":
		return fmt.Sprintf(" %s not like '%%%s%%' ", key, value)
	case "in":
		return fmt.Sprintf(" %s in (%s) ", key, value)
	case "notin":
		return fmt.Sprintf(" %s not in (%s) ", key, value)
	default:
		return ""
	}
}

// DelSQLPage 处理分页查询
func DelSQLPage(page, pageSize int64) string {
	if page < 0 {
		return ""
	}

	if page == 0 {
		page = 10
	}

	if pageSize == 0 {
		pageSize = 10
	}
	return fmt.Sprintf(" limit %d offset %d", pageSize, (page-1)*pageSize)
}
