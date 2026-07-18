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
	handler.Handler
	svc *services.DictService
}

// NewDictHandler 创建字典 Handler。
func NewDictHandler(svc *services.DictService) *DictHandler {
	return &DictHandler{svc: svc}
}

// TypeList GET /system/dict/type/list
func (h *DictHandler) TypeList(ctx *gin.Context) {
	total, rows, err := h.svc.TypeList(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// TypeGet GET /system/dict/type/:dictId
func (h *DictHandler) TypeGet(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("dictId"), 10, 64)
	data, err := h.svc.TypeGet(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.NotFound(ctx, "字典类型不存在")
			return
		}
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

type dictTypeBody struct {
	DictID   uint   `json:"dictId"`
	DictName string `json:"dictName"`
	DictType string `json:"dictType"`
	Status   string `json:"status"`
	Remark   string `json:"remark"`
}

// TypeAdd POST /system/dict/type
func (h *DictHandler) TypeAdd(ctx *gin.Context) {
	var req dictTypeBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
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
	if err := h.svc.TypeAdd(item, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "新增成功")
}

// TypeEdit PUT /system/dict/type
func (h *DictHandler) TypeEdit(ctx *gin.Context) {
	var req dictTypeBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	item := &models.DictType{
		DictName: req.DictName,
		DictType: req.DictType,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	item.ID = req.DictID
	if err := h.svc.TypeEdit(item, services.ContextUserName(ctx)); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// TypeDelete DELETE /system/dict/type/:dictId
func (h *DictHandler) TypeDelete(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Param("dictId"))
	if len(ids) == 0 {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.TypeDelete(ids); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "删除成功")
}

// TypeRefreshCache DELETE /system/dict/type/refreshCache
func (h *DictHandler) TypeRefreshCache(ctx *gin.Context) {
	h.SuccessMsg(ctx, "刷新缓存成功")
}

// TypeOptionSelect GET /system/dict/type/optionselect
func (h *DictHandler) TypeOptionSelect(ctx *gin.Context) {
	data, err := h.svc.TypeOptionSelect()
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// TypeExport POST /system/dict/type/export
func (h *DictHandler) TypeExport(ctx *gin.Context) {
	rows, err := h.svc.TypeExport(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	services.WriteCSV(ctx, "dict_type.csv", []string{"字典编号", "字典名称", "字典类型", "状态", "创建时间"}, rows)
}

// DataList GET /system/dict/data/list
func (h *DictHandler) DataList(ctx *gin.Context) {
	total, rows, err := h.svc.DataList(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// DataGet GET /system/dict/data/:dictCode
func (h *DictHandler) DataGet(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("dictCode"), 10, 64)
	data, err := h.svc.DataGet(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.NotFound(ctx, "字典数据不存在")
			return
		}
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// DataByType GET /system/dict/data/type/:dictType
func (h *DictHandler) DataByType(ctx *gin.Context) {
	data, err := h.svc.DataByType(ctx.Param("dictType"))
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
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
func (h *DictHandler) DataAdd(ctx *gin.Context) {
	var req dictDataBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
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
	if err := h.svc.DataAdd(item, services.ContextUserName(ctx)); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "新增成功")
}

// DataEdit PUT /system/dict/data
func (h *DictHandler) DataEdit(ctx *gin.Context) {
	var req dictDataBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
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
	if err := h.svc.DataEdit(item, services.ContextUserName(ctx)); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// DataDelete DELETE /system/dict/data/:dictCode
func (h *DictHandler) DataDelete(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Param("dictCode"))
	if len(ids) == 0 {
		h.BadRequest(ctx, "参数错误")
		return
	}
	if err := h.svc.DataDelete(ids); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "删除成功")
}

// DataExport POST /system/dict/data/export
func (h *DictHandler) DataExport(ctx *gin.Context) {
	rows, err := h.svc.DataExport(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	services.WriteCSV(ctx, "dict_data.csv", []string{"字典编码", "字典标签", "字典键值", "字典类型", "字典排序", "状态"}, rows)
}
