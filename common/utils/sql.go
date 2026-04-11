/*
Package pkg 提供通用的 SQL 处理函数

包含 SQL 查询条件处理和分页查询等功能
*/
package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// SQLCondition 表示一个参数化的 SQL 查询条件
type SQLCondition struct {
	Clause string // 包含 ? 占位符的条件子句，如 "`username` = ?"
	Args   []any  // 对应占位符的参数值
}

// DealSQL 处理 SQL 查询条件（参数化查询，防止 SQL 注入）
// key: 字段名（必须是可信的列名，不能来自用户输入）
// value: 查询值
// op: 操作符，支持 =, !=, >, >=, <, <=, like, notlike, in, notin
func DealSQL(key, value, op string) SQLCondition {
	if value == "" || value == "0" {
		return SQLCondition{}
	}
	// 处理时间字段
	if key == "create_time" || key == "update_time" {
		// 验证值是否为逗号分隔的两个数字
		parts := strings.Split(value, ",")
		if len(parts) != 2 {
			return SQLCondition{}
		}
		start := strings.TrimSpace(parts[0])
		end := strings.TrimSpace(parts[1])
		if start == "" || end == "" {
			return SQLCondition{}
		}
		// 验证 start 和 end 是合法数字
		startInt, err1 := strconv.ParseInt(start, 10, 64)
		endInt, err2 := strconv.ParseInt(end, 10, 64)
		if err1 != nil || err2 != nil {
			return SQLCondition{}
		}
		// 生成时间区间查询，左闭右开（参数化）
		return SQLCondition{
			Clause: fmt.Sprintf(" %s >= ? and %s < ? ", key, key),
			Args:   []any{startInt, endInt},
		}
	}

	// 处理其他字段
	switch strings.ToLower(op) {
	case "=":
		return SQLCondition{Clause: fmt.Sprintf(" %s = ? ", key), Args: []any{value}}
	case "!=":
		return SQLCondition{Clause: fmt.Sprintf(" %s != ? ", key), Args: []any{value}}
	case ">":
		return SQLCondition{Clause: fmt.Sprintf(" %s > ? ", key), Args: []any{value}}
	case ">=":
		return SQLCondition{Clause: fmt.Sprintf(" %s >= ? ", key), Args: []any{value}}
	case "<":
		return SQLCondition{Clause: fmt.Sprintf(" %s < ? ", key), Args: []any{value}}
	case "<=":
		return SQLCondition{Clause: fmt.Sprintf(" %s <= ? ", key), Args: []any{value}}
	case "like":
		return SQLCondition{Clause: fmt.Sprintf(" %s like ? ", key), Args: []any{fmt.Sprintf("%%%s%%", value)}}
	case "notlike":
		return SQLCondition{Clause: fmt.Sprintf(" %s not like ? ", key), Args: []any{fmt.Sprintf("%%%s%%", value)}}
	case "in":
		// in 操作需要特殊处理：解析 value 为多个值
		values := strings.Split(value, ",")
		placeholders := strings.Repeat("?,", len(values))
		placeholders = placeholders[:len(placeholders)-1] // 去掉末尾逗号
		args := make([]any, len(values))
		for i, v := range values {
			args[i] = strings.TrimSpace(v)
		}
		return SQLCondition{Clause: fmt.Sprintf(" %s in (%s) ", key, placeholders), Args: args}
	case "notin":
		values := strings.Split(value, ",")
		placeholders := strings.Repeat("?,", len(values))
		placeholders = placeholders[:len(placeholders)-1]
		args := make([]any, len(values))
		for i, v := range values {
			args[i] = strings.TrimSpace(v)
		}
		return SQLCondition{Clause: fmt.Sprintf(" %s not in (%s) ", key, placeholders), Args: args}
	default:
		return SQLCondition{}
	}
}

// BuildWhereClause 将多个 SQLCondition 组合成 WHERE 子句和参数列表
func BuildWhereClause(conditions []SQLCondition) (string, []any) {
	var clauses []string
	var args []any
	for _, cond := range conditions {
		if cond.Clause != "" {
			clauses = append(clauses, cond.Clause)
			args = append(args, cond.Args...)
		}
	}
	where := ""
	if len(clauses) > 0 {
		where = " where " + strings.Join(clauses, " and ")
	}
	return where, args
}

// DelSQLPage 处理分页查询（参数化）
func DelSQLPage(page, pageSize int64) (string, []any) {
	if page < 0 {
		return "", nil
	}

	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 10
	}
	return " limit ? offset ?", []any{pageSize, (page - 1) * pageSize}
}
