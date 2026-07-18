package system

import (
	"errors"
	"strconv"

	"thirteen-admin/app/models"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DictHandler 字典 HTTP 处理器。
type DictHandler struct {
	handler.Api
}

// TypeList GET /system/dict/type/list
func (e DictHandler) TypeList(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	total, rows, err := s.TypeList(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Table(c, total, rows)
}

// TypeGet GET /system/dict/type/:dictId
func (e DictHandler) TypeGet(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("dictId"), 10, 64)
	data, err := s.TypeGet(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.NotFound(c, "字典类型不存在")
			return
		}
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

type dictTypeBody struct {
	DictID   uint   `json:"dictId"`
	DictName string `json:"dictName"`
	DictType string `json:"dictType"`
	Status   string `json:"status"`
	Remark   string `json:"remark"`
}

// TypeAdd POST /system/dict/type
func (e DictHandler) TypeAdd(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req dictTypeBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	item := &models.DictType{
		DictName: req.DictName,
		DictType: req.DictType,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	if item.Status == "" {
		item.Status = "0"
	}
	if err := s.TypeAdd(item, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "新增成功")
}

// TypeEdit PUT /system/dict/type
func (e DictHandler) TypeEdit(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req dictTypeBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	item := &models.DictType{
		DictName: req.DictName,
		DictType: req.DictType,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	item.ID = req.DictID
	if err := s.TypeEdit(item, services.ContextUserName(c)); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// TypeDelete DELETE /system/dict/type/:dictId
func (e DictHandler) TypeDelete(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	ids := services.ParseUintIDs(c.Param("dictId"))
	if len(ids) == 0 {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.TypeDelete(ids); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "删除成功")
}

// TypeRefreshCache DELETE /system/dict/type/refreshCache
func (e DictHandler) TypeRefreshCache(c *gin.Context) {
	e.SuccessMsg(c, "刷新缓存成功")
}

// TypeOptionSelect GET /system/dict/type/optionselect
func (e DictHandler) TypeOptionSelect(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	data, err := s.TypeOptionSelect()
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

// TypeExport POST /system/dict/type/export
func (e DictHandler) TypeExport(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	rows, err := s.TypeExport(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	services.WriteCSV(c, "dict_type.csv", []string{"字典编号", "字典名称", "字典类型", "状态", "创建时间"}, rows)
}

// DataList GET /system/dict/data/list
func (e DictHandler) DataList(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	total, rows, err := s.DataList(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Table(c, total, rows)
}

// DataGet GET /system/dict/data/:dictCode
func (e DictHandler) DataGet(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("dictCode"), 10, 64)
	data, err := s.DataGet(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.NotFound(c, "字典数据不存在")
			return
		}
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

// DataByType GET /system/dict/data/type/:dictType
func (e DictHandler) DataByType(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	data, err := s.DataByType(c.Param("dictType"))
	if err != nil {
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

type dictDataBody struct {
	DictCode  uint   `json:"dictCode"`
	DictSort  int    `json:"dictSort"`
	DictLabel string `json:"dictLabel"`
	DictValue string `json:"dictValue"`
	DictType  string `json:"dictType"`
	CSSClass  string `json:"cssClass"`
	ListClass string `json:"listClass"`
	IsDefault string `json:"isDefault"`
	Status    string `json:"status"`
	Remark    string `json:"remark"`
}

// DataAdd POST /system/dict/data
func (e DictHandler) DataAdd(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req dictDataBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	item := &models.DictData{
		DictSort:  req.DictSort,
		DictLabel: req.DictLabel,
		DictValue: req.DictValue,
		DictType:  req.DictType,
		CSSClass:  req.CSSClass,
		ListClass: req.ListClass,
		IsDefault: req.IsDefault,
		Status:    req.Status,
		Remark:    req.Remark,
	}
	if item.IsDefault == "" {
		item.IsDefault = "N"
	}
	if item.Status == "" {
		item.Status = "0"
	}
	if err := s.DataAdd(item, services.ContextUserName(c)); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "新增成功")
}

// DataEdit PUT /system/dict/data
func (e DictHandler) DataEdit(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req dictDataBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
		return
	}
	item := &models.DictData{
		DictSort:  req.DictSort,
		DictLabel: req.DictLabel,
		DictValue: req.DictValue,
		DictType:  req.DictType,
		CSSClass:  req.CSSClass,
		ListClass: req.ListClass,
		IsDefault: req.IsDefault,
		Status:    req.Status,
		Remark:    req.Remark,
	}
	item.ID = req.DictCode
	if err := s.DataEdit(item, services.ContextUserName(c)); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// DataDelete DELETE /system/dict/data/:dictCode
func (e DictHandler) DataDelete(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	ids := services.ParseUintIDs(c.Param("dictCode"))
	if len(ids) == 0 {
		e.BadRequest(c, "参数错误")
		return
	}
	if err := s.DataDelete(ids); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "删除成功")
}

// DataExport POST /system/dict/data/export
func (e DictHandler) DataExport(c *gin.Context) {
	s := services.DictService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	rows, err := s.DataExport(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	services.WriteCSV(c, "dict_data.csv", []string{"字典编码", "字典标签", "字典键值", "字典类型", "字典排序", "状态"}, rows)
}
