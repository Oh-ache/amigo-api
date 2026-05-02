package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type (
	AiTaskModel interface {
		Insert(ctx context.Context, data *AiTask) (int64, error)
		FindOne(ctx context.Context, id int64) (*AiTask, error)
		Update(ctx context.Context, data *AiTask) error
		Delete(ctx context.Context, id int64) error
		List(ctx context.Context, search *AiTaskSearch) ([]*AiTask, int64, error)
	}

	aiTaskModel struct {
		conn sqlx.SqlConn
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

func NewAiTaskModel(conn sqlx.SqlConn) AiTaskModel {
	return &aiTaskModel{conn: conn}
}

func (m *aiTaskModel) Insert(ctx context.Context, data *AiTask) (int64, error) {
	query := fmt.Sprintf("insert into ai_task (user_id, task_id, task_type, prompt, request_info, response_info, result_url, status, error_msg, created_at, updated_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	result, err := m.conn.ExecCtx(ctx, query, data.UserId, data.TaskId, data.TaskType, data.Prompt, data.RequestInfo, data.ResponseInfo, data.ResultUrl, data.Status, data.ErrorMsg, data.CreatedAt, data.UpdatedAt)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (m *aiTaskModel) FindOne(ctx context.Context, id int64) (*AiTask, error) {
	query := fmt.Sprintf("select id, user_id, task_id, task_type, prompt, request_info, response_info, result_url, status, error_msg, created_at, updated_at from ai_task where id = ? limit 1")
	var resp AiTask
	err := m.conn.QueryRowCtx(ctx, &resp, query, id)
	switch err {
	case nil:
		return &resp, nil
	case sqlx.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

func (m *aiTaskModel) Update(ctx context.Context, data *AiTask) error {
	query := fmt.Sprintf("update ai_task set user_id=?, task_id=?, task_type=?, prompt=?, request_info=?, response_info=?, result_url=?, status=?, error_msg=?, updated_at=? where id=?")
	_, err := m.conn.ExecCtx(ctx, query, data.UserId, data.TaskId, data.TaskType, data.Prompt, data.RequestInfo, data.ResponseInfo, data.ResultUrl, data.Status, data.ErrorMsg, data.UpdatedAt, data.Id)
	return err
}

func (m *aiTaskModel) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf("delete from ai_task where id = ?")
	_, err := m.conn.ExecCtx(ctx, query, id)
	return err
}

func (m *aiTaskModel) List(ctx context.Context, search *AiTaskSearch) ([]*AiTask, int64, error) {
	where := "where 1=1"
	if search.TaskId != "" {
		where += fmt.Sprintf(" and task_id = '%s'", search.TaskId)
	}
	if search.TaskType != "" {
		where += fmt.Sprintf(" and task_type = '%s'", search.TaskType)
	}
	if search.Status > 0 {
		where += fmt.Sprintf(" and status = %d", search.Status)
	}
	if search.UserId > 0 {
		where += fmt.Sprintf(" and user_id = %d", search.UserId)
	}

	countQuery := fmt.Sprintf("select count(*) from ai_task %s", where)
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, countQuery)
	if err != nil {
		return nil, 0, err
	}

	offset := (search.Page - 1) * search.PageSize
	query := fmt.Sprintf("select id, user_id, task_id, task_type, prompt, request_info, response_info, result_url, status, error_msg, created_at, updated_at from ai_task %s limit %d offset %d", where, search.PageSize, offset)

	var rows []*AiTask
	err = m.conn.QueryRowsCtx(ctx, &rows, query)
	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}