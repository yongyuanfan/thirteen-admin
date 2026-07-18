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

// ConfigHandler 参数 HTTP 处理器。
type ConfigHandler struct {
	handler.Api
}

// List GET /system/config/list
func (e ConfigHandler) List(c *gin.Context) {
	s := services.ConfigService{}
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

// Get GET /system/config/:configId
func (e ConfigHandler) Get(c *gin.Context) {
	s := services.ConfigService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("configId"), 10, 64)
	data, err := s.Get(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.NotFound(c, "参数不存在")
			return
		}
		e.InternalError(c, err)
		return
	}
	e.Success(c, data)
}

// GetByKey GET /system/config/configKey/:configKey
func (e ConfigHandler) GetByKey(c *gin.Context) {
	s := services.ConfigService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	value, err := s.GetByKey(c.Param("configKey"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.NotFound(c, "参数不存在")
			return
		}
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, value)
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
func (e ConfigHandler) Add(c *gin.Context) {
	s := services.ConfigService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req configBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
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
	if err := s.Add(cfg, services.ContextUserName(c)); err != nil {
		e.Fail(c, handler.CodeError, err.Error())
		return
	}
	e.SuccessMsg(c, "新增成功")
}

// Edit PUT /system/config
func (e ConfigHandler) Edit(c *gin.Context) {
	s := services.ConfigService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	var req configBody
	if err := c.ShouldBindJSON(&req); err != nil {
		e.BadRequest(c, "参数错误")
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
	if err := s.Edit(cfg, services.ContextUserName(c)); err != nil {
		e.InternalError(c, err)
		return
	}
	e.SuccessMsg(c, "修改成功")
}

// Delete DELETE /system/config/:configId
func (e ConfigHandler) Delete(c *gin.Context) {
	s := services.ConfigService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	ids := services.ParseUintIDs(c.Param("configId"))
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

// RefreshCache DELETE /system/config/refreshCache
func (e ConfigHandler) RefreshCache(c *gin.Context) {
	e.SuccessMsg(c, "刷新缓存成功")
}

// Export POST /system/config/export
func (e ConfigHandler) Export(c *gin.Context) {
	s := services.ConfigService{}
	if err := e.MakeContext(c).MakeOrm().MakeService(&s.Service).Errors; err != nil {
		e.InternalError(c, err)
		return
	}
	rows, err := s.Export(c)
	if err != nil {
		e.InternalError(c, err)
		return
	}
	services.WriteCSV(c, "config.csv", []string{"参数编号", "参数名称", "参数键名", "参数键值", "系统内置", "创建时间"}, rows)
}
