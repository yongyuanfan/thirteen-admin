package tool

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
	handler.Api
}

// List GET /tool/gen/list
func (e GenHandler) List(c *gin.Context) {
	s := services.GenService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	total, rows, err := s.List(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Table(c, total, rows)
}

// DbList GET /tool/gen/db/list
func (e GenHandler) DbList(c *gin.Context) {
	s := services.GenService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	total, rows, err := s.DbList(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Table(c, total, rows)
}

// Get GET /tool/gen/:tableId
func (e GenHandler) Get(c *gin.Context) {
	s := services.GenService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("tableId"), 10, 64)
	data, err := s.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.NotFound(c, "生成表不存在")
			return
		}
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
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
func (e GenHandler) Update(c *gin.Context) {
	s := services.GenService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req genUpdateBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
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
	for _, colReq := range req.Columns {
		tid := colReq.TableID
		if tid == 0 {
			tid = req.TableID
		}
		col := models.GenTableColumn{
			ColumnName:    colReq.ColumnName,
			ColumnComment: colReq.ColumnComment,
			ColumnType:    colReq.ColumnType,
			JavaType:      colReq.JavaType,
			JavaField:     colReq.JavaField,
			IsPk:          colReq.IsPk,
			IsIncrement:   colReq.IsIncrement,
			IsRequired:    colReq.IsRequired,
			IsInsert:      colReq.IsInsert,
			IsEdit:        colReq.IsEdit,
			IsList:        colReq.IsList,
			IsQuery:       colReq.IsQuery,
			QueryType:     colReq.QueryType,
			HTMLType:      colReq.HTMLType,
			DictType:      colReq.DictType,
			Sort:          colReq.Sort,
		}
		col.ID = colReq.ColumnID
		col.TableID = &tid
		columns = append(columns, col)
	}
	if err := s.Update(table, columns, services.ContextUserName(c)); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "保存成功")
}

// ImportTable POST /tool/gen/importTable
func (e GenHandler) ImportTable(c *gin.Context) {
	s := services.GenService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	tables := c.Query("tables")
	tplWebType := c.DefaultQuery("tplWebType", "element-plus")
	names := splitComma(tables)
	if err := s.ImportTable(names, tplWebType, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "导入成功")
}

// CreateTable POST /tool/gen/createTable
func (e GenHandler) CreateTable(c *gin.Context) {
	s := services.GenService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	sql := c.Query("sql")
	tplWebType := c.DefaultQuery("tplWebType", "element-plus")
	if err := s.CreateTable(sql, tplWebType, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "创建成功")
}

// Preview GET /tool/gen/preview/:tableId
func (e GenHandler) Preview(c *gin.Context) {
	s := services.GenService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("tableId"), 10, 64)
	data, err := s.Preview(uint(id))
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

// Delete DELETE /tool/gen/:tableId
func (e GenHandler) Delete(c *gin.Context) {
	s := services.GenService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	ids := services.ParseUintIDs(c.Param("tableId"))
	if len(ids) == 0 {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.Delete(ids); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "删除成功")
}

// GenCode GET /tool/gen/genCode/:tableName
func (e GenHandler) GenCode(c *gin.Context) {
	s := services.GenService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	if err := s.GenCode(c.Param("tableName")); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "生成成功")
}

// SynchDb GET /tool/gen/synchDb/:tableName
func (e GenHandler) SynchDb(c *gin.Context) {
	s := services.GenService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	if err := s.SynchDb(c.Param("tableName"), services.ContextUserName(c)); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "同步成功")
}

// BatchGenCode GET /tool/gen/batchGenCode
func (e GenHandler) BatchGenCode(c *gin.Context) {
	s := services.GenService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	names := splitComma(c.Query("tables"))
	data, filename, err := s.BatchGenCode(names)
	if err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", url.QueryEscape(filename)))
	c.Header("download-filename", url.QueryEscape(filename))
	c.Data(200, "application/octet-stream", data)
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
