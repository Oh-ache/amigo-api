package model

import (
	"errors"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var (
	ErrNotFound      = sqlx.ErrNotFound
	ErrDuplicate     = errors.New("duplicate key error")
	ErrInvalidParams = errors.New("invalid parameters")
)