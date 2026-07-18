package handlers

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"thirteen-admin/app/models"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GenHandler 代码生成 HTTP 处理器。
type GenHandler struct {
	handler.Handler
	svc *services.GenService
}

// NewGenHandler 创建代码生成 Handler。
func NewGenHandler(svc *services.GenService) *GenHandler {
	return &GenHandler{svc: svc}
}

// List GET /tool/gen/list
func (h *GenHandler) List(ctx *gin.Context) {
	total, rows, err := h.svc.List(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// DbList GET /tool/gen/db/list
func (h *GenHandler) DbList(ctx *gin.Context) {
	total, rows, err := h.svc.DbList(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// Get GET /tool/gen/:tableId
func (h *GenHandler) Get(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("tableId"), 10, 64)
	data, err := h.svc.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.NotFound(ctx, "生成表不存在")
			return
		}
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

type genUpdateBody struct {
	TableID        uint            `json:"tableId"`
	TableName      string          `json:"tableName"`
	TableComment   string          `json:"tableComment"`
	SubTableName   string          `json:"subTableName"`
	SubTableFkName string          `json:"subTableFkName"`
	ClassName      string          `json:"className"`
	TplCategory    string          `json:"tplCategory"`
	TplWebType     string          `json:"tplWebType"`
	PackageName    string          `json:"packageName"`
	ModuleName     string          `json:"moduleName"`
	BusinessName   string          `json:"businessName"`
	FunctionName   string          `json:"functionName"`
	FunctionAuthor string          `json:"functionAuthor"`
	FormColNum     int             `json:"formColNum"`
	GenType        string          `json:"genType"`
	GenPath        string          `json:"genPath"`
	Options        string          `json:"options"`
	Remark         string          `json:"remark"`
	Columns        []genColumnBody `json:"columns"`
}

type genColumnBody struct {
	ColumnID      uint   `json:"columnId"`
	TableID       uint   `json:"tableId"`
	ColumnName    string `json:"columnName"`
	ColumnComment string `json:"columnComment"`
	ColumnType    string `json:"columnType"`
	JavaType      string `json:"javaType"`
	JavaField     string `json:"javaField"`
	IsPk          string `json:"isPk"`
	IsIncrement   string `json:"isIncrement"`
	IsRequired    string `json:"isRequired"`
	IsInsert      string `json:"isInsert"`
	IsEdit        string `json:"isEdit"`
	IsList        string `json:"isList"`
	IsQuery       string `json:"isQuery"`
	QueryType     string `json:"queryType"`
	HTMLType      string `json:"htmlType"`
	DictType      string `json:"dictType"`
	Sort          int    `json:"sort"`
}

// Update PUT /tool/gen
func (h *GenHandler) Update(ctx *gin.Context) {
	var req genUpdateBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	table := &models.GenTable{
		Name:           req.TableName,
		TableComment:   req.TableComment,
		SubTableName:   req.SubTableName,
		SubTableFkName: req.SubTableFkName,
		ClassName:      req.ClassName,
		TplCategory:    req.TplCategory,
		TplWebType:     req.TplWebType,
		PackageName:    req.PackageName,
		ModuleName:     req.ModuleName,
		BusinessName:   req.BusinessName,
		FunctionName:   req.FunctionName,
		FunctionAuthor: req.FunctionAuthor,
		FormColNum:     req.FormColNum,
		GenType:        req.GenType,
		GenPath:        req.GenPath,
		Options:        req.Options,
		Remark:         req.Remark,
	}
	table.ID = req.TableID
	columns := make([]models.GenTableColumn, 0, len(req.Columns))
	for _, c := range req.Columns {
		tid := c.TableID
		if tid == 0 {
			tid = req.TableID
		}
		col := models.GenTableColumn{
			ColumnName:    c.ColumnName,
			ColumnComment: c.ColumnComment,
			ColumnType:    c.ColumnType,
			JavaType:      c.JavaType,
			JavaField:     c.JavaField,
			IsPk:          c.IsPk,
			IsIncrement:   c.IsIncrement,
			IsRequired:    c.IsRequired,
			IsInsert:      c.IsInsert,
			IsEdit:        c.IsEdit,
			IsList:        c.IsList,
			IsQuery:       c.IsQuery,
			QueryType:     c.QueryType,
			HTMLType:      c.HTMLType,
			DictType:      c.DictType,
			Sort:          c.Sort,
		}
		col.ID = c.ColumnID
		col.TableID = &tid
		columns = append(columns, col)
	}
	if err := h.svc.Update(table, columns, services.ContextUserName(ctx)); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "保存成功")
}

// ImportTable POST /tool/gen/importTable
func (h *GenHandler) ImportTable(ctx *gin.Context) {
	tables := ctx.Query("tables")
	tplWebType := ctx.DefaultQuery("tplWebType", "element-plus")
	names := splitComma(tables)
	if err := h.svc.ImportTable(names, tplWebType, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "导入成功")
}

// CreateTable POST /tool/gen/createTable
func (h *GenHandler) CreateTable(ctx *gin.Context) {
	sql := ctx.Query("sql")
	tplWebType := ctx.DefaultQuery("tplWebType", "element-plus")
	if err := h.svc.CreateTable(sql, tplWebType, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "创建成功")
}

// Preview GET /tool/gen/preview/:tableId
func (h *GenHandler) Preview(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("tableId"), 10, 64)
	data, err := h.svc.Preview(uint(id))
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// Delete DELETE /tool/gen/:tableId
func (h *GenHandler) Delete(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Param("tableId"))
	if len(ids) == 0 {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.Delete(ids); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "删除成功")
}

// GenCode GET /tool/gen/genCode/:tableName
func (h *GenHandler) GenCode(ctx *gin.Context) {
	if err := h.svc.GenCode(ctx.Param("tableName")); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "生成成功")
}

// SynchDb GET /tool/gen/synchDb/:tableName
func (h *GenHandler) SynchDb(ctx *gin.Context) {
	if err := h.svc.SynchDb(ctx.Param("tableName"), services.ContextUserName(ctx)); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "同步成功")
}

// BatchGenCode GET /tool/gen/batchGenCode
func (h *GenHandler) BatchGenCode(ctx *gin.Context) {
	names := splitComma(ctx.Query("tables"))
	data, filename, err := h.svc.BatchGenCode(names)
	if err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.QueryEscape(filename)))
	ctx.Header("download-filename", url.QueryEscape(filename))
	ctx.Data(200, "application/octet-stream", data)
}

func splitComma(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
