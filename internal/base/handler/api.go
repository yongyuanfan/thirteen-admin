package handler

import (
	"errors"
	"fmt"

	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ContextDBKey gin.Context 中请求级 DB 的键。
const ContextDBKey = "ormDB"

// Api 请求级处理器基类（对齐 go-admin Make* 链式装配）。
type Api struct {
	Handler
	Context *gin.Context
	Orm     *gorm.DB
	Errors  error
}

// MakeContext 绑定 gin.Context。
func (e *Api) MakeContext(c *gin.Context) *Api {
	e.Context = c
	return e
}

// MakeOrm 从 Context 取出请求级 DB。
func (e *Api) MakeOrm() *Api {
	if e.Errors != nil {
		return e
	}
	if e.Context == nil {
		e.AddError(errors.New("MakeOrm: context is nil"))
		return e
	}
	v, ok := e.Context.Get(ContextDBKey)
	if !ok {
		e.AddError(errors.New("MakeOrm: db not found in context"))
		return e
	}
	db, ok := v.(*gorm.DB)
	if !ok || db == nil {
		e.AddError(errors.New("MakeOrm: invalid db in context"))
		return e
	}
	e.Orm = db
	return e
}

// MakeService 将请求级 Orm 注入 Service。
func (e *Api) MakeService(s *service.Service) *Api {
	if e.Errors != nil {
		return e
	}
	if e.Orm == nil {
		e.AddError(errors.New("MakeService: orm is nil"))
		return e
	}
	s.DB = e.Orm
	return e
}

// AddError 累积错误。
func (e *Api) AddError(err error) *Api {
	if err == nil {
		return e
	}
	if e.Errors == nil {
		e.Errors = err
		return e
	}
	e.Errors = fmt.Errorf("%v; %w", e.Errors, err)
	return e
}
