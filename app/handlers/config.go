package handlers

import (
	"errors"
	"strconv"

	"thirteen-admin/app/models"
	"thirteen-admin/app/services"
	"thirteen-admin/internal/base/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ConfigHandler 参数 HTTP 处理器。
type ConfigHandler struct {
	handler.Handler
	svc *services.ConfigService
}

// NewConfigHandler 创建参数 Handler。
func NewConfigHandler(svc *services.ConfigService) *ConfigHandler {
	return &ConfigHandler{svc: svc}
}

// List GET /system/config/list
func (h *ConfigHandler) List(ctx *gin.Context) {
	total, rows, err := h.svc.List(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.Table(ctx, total, rows)
}

// Get GET /system/config/:configId
func (h *ConfigHandler) Get(ctx *gin.Context) {
	id, _ := strconv.ParseUint(ctx.Param("configId"), 10, 64)
	data, err := h.svc.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.NotFound(ctx, "参数不存在")
			return
		}
		h.InternalError(ctx, err)
		return
	}
	h.Success(ctx, data)
}

// GetByKey GET /system/config/configKey/:configKey
func (h *ConfigHandler) GetByKey(ctx *gin.Context) {
	value, err := h.svc.GetByKey(ctx.Param("configKey"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.NotFound(ctx, "参数不存在")
			return
		}
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, value)
}

type configBody struct {
	ConfigID    uint   `json:"configId"`
	ConfigName  string `json:"configName"`
	ConfigKey   string `json:"configKey"`
	ConfigValue string `json:"configValue"`
	ConfigType  string `json:"configType"`
	Remark      string `json:"remark"`
}

// Add POST /system/config
func (h *ConfigHandler) Add(ctx *gin.Context) {
	var req configBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	cfg := &models.Config{
		ConfigName:  req.ConfigName,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		ConfigType:  req.ConfigType,
		Remark:      req.Remark,
	}
	if cfg.ConfigType == "" {
		cfg.ConfigType = "N"
	}
	if err := h.svc.Add(cfg, services.ContextUserName(ctx)); err != nil {
		h.Fail(ctx, handler.CodeError, err.Error())
		return
	}
	h.SuccessMsg(ctx, "新增成功")
}

// Edit PUT /system/config
func (h *ConfigHandler) Edit(ctx *gin.Context) {
	var req configBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.BadRequest(ctx, "参数错误")
		return
	}
	cfg := &models.Config{
		ConfigName:  req.ConfigName,
		ConfigKey:   req.ConfigKey,
		ConfigValue: req.ConfigValue,
		ConfigType:  req.ConfigType,
		Remark:      req.Remark,
	}
	cfg.ID = req.ConfigID
	if err := h.svc.Edit(cfg, services.ContextUserName(ctx)); err != nil {
		h.InternalError(ctx, err)
		return
	}
	h.SuccessMsg(ctx, "修改成功")
}

// Delete DELETE /system/config/:configId
func (h *ConfigHandler) Delete(ctx *gin.Context) {
	ids := services.ParseUintIDs(ctx.Param("configId"))
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

// RefreshCache DELETE /system/config/refreshCache
func (h *ConfigHandler) RefreshCache(ctx *gin.Context) {
	h.SuccessMsg(ctx, "刷新缓存成功")
}

// Export POST /system/config/export
func (h *ConfigHandler) Export(ctx *gin.Context) {
	rows, err := h.svc.Export(ctx)
	if err != nil {
		h.InternalError(ctx, err)
		return
	}
	services.WriteCSV(ctx, "config.csv", []string{"参数编号", "参数名称", "参数键名", "参数键值", "系统内置", "创建时间"}, rows)
}
