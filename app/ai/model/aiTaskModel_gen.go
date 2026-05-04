package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"amigo-api/common/utils"

	"github.com/zeromicro/go-zero/core/stores/builder"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/core/stringx"
)

var (
	aiTaskFieldNames          = builder.RawFieldNames(&AiTask{})
	aiTaskRows                = strings.Join(aiTaskFieldNames, ",")
	aiTaskRowsExpectAutoSet   = strings.Join(stringx.Remove(aiTaskFieldNames, "`id`", "`created_at`", "`create_at`", "`create_time`", "`updated_at`", "`update_at`"), ",")
	aiTaskRowsWithPlaceHolder = strings.Join(stringx.Remove(aiTaskFieldNames, "`id`", "`created_at`", "`create_at`", "`create_time`", "`updated_at`", "`update_at`"), "=?,") + "=?"

	cacheAmigoAiTaskIdPrefix = "cache:amigo:aiTask:id:"
)

type (
	aiTaskModel interface {
		Insert(ctx context.Context, data *AiTask) (sql.Result, error)
		FindOne(ctx context.Context, id int64) (*AiTask, error)
		Update(ctx context.Context, data *AiTask) error
		Delete(ctx context.Context, id int64) error
		List(ctx context.Context, search *AiTaskSearch) ([]*AiTask, int64, error)
	}

	defaultAiTaskModel struct {
		sqlc.CachedConn
		table string
	}

	AiTask struct {
		Id           int64  `db:"id"`
		UserId       int64  `db:"user_id"`
		TaskId       string `db:"task_id"`
		TaskType     string `db:"task_type"`
		Prompt       string `db:"prompt"`
		RequestInfo  string `db:"request_info"`
		ResponseInfo string `db:"response_info"`
		ResultUrl    string `db:"result_url"`
		Status       int    `db:"status"`
		ErrorMsg     string `db:"error_msg"`
		CreatedAt    int64  `db:"created_at"`
		UpdatedAt    int64  `db:"updated_at"`
	}

	AiTaskSearch struct {
		TaskId   string
		TaskType string
		Status   int
		UserId   int64
		Page     int64
		PageSize int64
	}
)

func newAiTaskModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) *defaultAiTaskModel {
	return &defaultAiTaskModel{
		CachedConn: sqlc.NewConn(conn, c, opts...),
		table:      "`ai_task`",
	}
}

func (m *defaultAiTaskModel) Delete(ctx context.Context, id int64) error {
	data, err := m.FindOne(ctx, id)
	if err != nil && data.Id > 0 {
		return err
	}

	amigoAiTaskIdKey := fmt.Sprintf("%s%v", cacheAmigoAiTaskIdPrefix, id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("delete from %s where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, amigoAiTaskIdKey)
	return err
}

func (m *defaultAiTaskModel) FindOne(ctx context.Context, id int64) (*AiTask, error) {
	amigoAiTaskIdKey := fmt.Sprintf("%s%v", cacheAmigoAiTaskIdPrefix, id)
	var resp AiTask
	err := m.QueryRowCtx(ctx, &resp, amigoAiTaskIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", aiTaskRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *defaultAiTaskModel) Insert(ctx context.Context, data *AiTask) (sql.Result, error) {
	now := time.Now().Unix()
	data.CreatedAt = now
	data.UpdatedAt = now

	amigoAiTaskIdKey := fmt.Sprintf("%s%v", cacheAmigoAiTaskIdPrefix, data.Id)
	ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		insertFieldNames := stringx.Remove(aiTaskFieldNames, "`id`")
		insertFields := strings.Join(insertFieldNames, ",")
		placeholders := strings.Repeat("?,", len(insertFieldNames)-1) + "?"

		query := fmt.Sprintf("insert into %s (%s) values (%s)", m.table, insertFields, placeholders)
		params := make([]interface{}, 0, len(insertFieldNames))
		for _, field := range insertFieldNames {
			switch field {
			case "`user_id`":
				params = append(params, data.UserId)
			case "`task_id`":
				params = append(params, data.TaskId)
			case "`task_type`":
				params = append(params, data.TaskType)
			case "`prompt`":
				params = append(params, data.Prompt)
			case "`request_info`":
				params = append(params, data.RequestInfo)
			case "`response_info`":
				params = append(params, data.ResponseInfo)
			case "`result_url`":
				params = append(params, data.ResultUrl)
			case "`status`":
				params = append(params, data.Status)
			case "`error_msg`":
				params = append(params, data.ErrorMsg)
			case "`created_at`":
				params = append(params, data.CreatedAt)
			case "`updated_at`":
				params = append(params, data.UpdatedAt)
			}
		}
		return conn.ExecCtx(ctx, query, params...)
	}, amigoAiTaskIdKey)
	return ret, err
}

func (m *defaultAiTaskModel) Update(ctx context.Context, newData *AiTask) error {
	data, err := m.FindOne(ctx, newData.Id)
	if err != nil {
		return err
	}

	now := time.Now().Unix()
	newData.UpdatedAt = now

	amigoAiTaskIdKey := fmt.Sprintf("%s%v", cacheAmigoAiTaskIdPrefix, data.Id)
	_, err = m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		updateFieldNames := stringx.Remove(aiTaskFieldNames, "`id`", "`created_at`", "`create_at`", "`create_time`")
		updateFields := strings.Join(updateFieldNames, "=?,") + "=?"
		query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, updateFields)

		params := make([]interface{}, 0, len(updateFieldNames))
		for _, field := range updateFieldNames {
			switch field {
			case "`user_id`":
				params = append(params, newData.UserId)
			case "`task_id`":
				params = append(params, newData.TaskId)
			case "`task_type`":
				params = append(params, newData.TaskType)
			case "`prompt`":
				params = append(params, newData.Prompt)
			case "`request_info`":
				params = append(params, newData.RequestInfo)
			case "`response_info`":
				params = append(params, newData.ResponseInfo)
			case "`result_url`":
				params = append(params, newData.ResultUrl)
			case "`status`":
				params = append(params, newData.Status)
			case "`error_msg`":
				params = append(params, newData.ErrorMsg)
			case "`updated_at`":
				params = append(params, newData.UpdatedAt)
			}
		}
		params = append(params, newData.Id)

		return conn.ExecCtx(ctx, query, params...)
	}, amigoAiTaskIdKey)
	return err
}

func (m *defaultAiTaskModel) formatPrimary(primary any) string {
	return fmt.Sprintf("%s%v", cacheAmigoAiTaskIdPrefix, primary)
}

func (m *defaultAiTaskModel) queryPrimary(ctx context.Context, conn sqlx.SqlConn, v, primary any) error {
	query := fmt.Sprintf("select %s from %s where `id` = ? limit 1", aiTaskRows, m.table)
	return conn.QueryRowCtx(ctx, v, query, primary)
}

func (m *defaultAiTaskModel) tableName() string {
	return m.table
}

func (m *defaultAiTaskModel) List(ctx context.Context, search *AiTaskSearch) ([]*AiTask, int64, error) {
	var conditions []string

	if search.TaskId != "" {
		conditions = append(conditions, "`task_id` = '"+search.TaskId+"'")
	}
	if search.TaskType != "" {
		conditions = append(conditions, "`task_type` = '"+search.TaskType+"'")
	}
	if search.Status > 0 {
		conditions = append(conditions, "`status` = "+fmt.Sprintf("%d", search.Status))
	}
	if search.UserId > 0 {
		conditions = append(conditions, "`user_id` = "+fmt.Sprintf("%d", search.UserId))
	}

	queryWhere := ""
	if len(conditions) > 0 {
		queryWhere = " where " + strings.Join(conditions, " and ")
	}

	countQuery := fmt.Sprintf("select count(*) from %s %s", m.table, queryWhere)

	var total int64
	if m.QueryRowNoCacheCtx(ctx, &total, countQuery) != nil {
		return nil, 0, fmt.Errorf("failed to get total count")
	}

	pageSql := utils.DelSQLPage(search.Page, search.PageSize)
	query := fmt.Sprintf("select %s from %s %s order by `id` desc %s", aiTaskRows, m.table, queryWhere, pageSql)

	var list []*AiTask
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query); err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
