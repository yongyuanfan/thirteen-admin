package services

import (
	"errors"
	"strconv"
	"strings"

	"thirteen-admin/app/models"
)

// DeptToMap 部门 VO。
func DeptToMap(d *models.Dept) map[string]any {
	if d == nil {
		return nil
	}
	return map[string]any{
		"deptId":     d.ID,
		"parentId":   d.ParentID,
		"ancestors":  d.Ancestors,
		"deptName":   d.DeptName,
		"orderNum":   d.OrderNum,
		"leader":     d.Leader,
		"phone":      d.Phone,
		"email":      d.Email,
		"status":     d.Status,
		"createdBy":   d.CreatedBy,
		"createdAt": FormatTimeVal(d.CreatedAt),
		"updatedBy":   d.UpdatedBy,
		"updatedAt": FormatTimeVal(d.UpdatedAt),
	}
}

// MenuToMap 菜单 VO。
func MenuToMap(m *models.Menu) map[string]any {
	if m == nil {
		return nil
	}
	return map[string]any{
		"menuId":     m.ID,
		"menuName":   m.MenuName,
		"parentId":   m.ParentID,
		"orderNum":   m.OrderNum,
		"path":       m.Path,
		"component":  m.Component,
		"query":      m.Query,
		"routeName":  m.RouteName,
		"isFrame":    strconv.Itoa(m.IsFrame),
		"isCache":    strconv.Itoa(m.IsCache),
		"menuType":   m.MenuType,
		"visible":    m.Visible,
		"status":     m.Status,
		"perms":      m.Perms,
		"icon":       m.Icon,
		"createdBy":   m.CreatedBy,
		"createdAt": FormatTimeVal(m.CreatedAt),
		"updatedBy":   m.UpdatedBy,
		"updatedAt": FormatTimeVal(m.UpdatedAt),
		"remark":     m.Remark,
	}
}

// PostToMap 岗位 VO。
func PostToMap(p *models.Post) map[string]any {
	if p == nil {
		return nil
	}
	return map[string]any{
		"postId":     p.ID,
		"postCode":   p.PostCode,
		"postName":   p.PostName,
		"postSort":   p.PostSort,
		"status":     p.Status,
		"createdBy":   p.CreatedBy,
		"createdAt": FormatTimeVal(p.CreatedAt),
		"updatedBy":   p.UpdatedBy,
		"updatedAt": FormatTimeVal(p.UpdatedAt),
		"remark":     p.Remark,
	}
}

// ConfigToMap 参数 VO。
func ConfigToMap(c *models.Config) map[string]any {
	if c == nil {
		return nil
	}
	return map[string]any{
		"configId":    c.ID,
		"configName":  c.ConfigName,
		"configKey":   c.ConfigKey,
		"configValue": c.ConfigValue,
		"configType":  c.ConfigType,
		"createdBy":    c.CreatedBy,
		"createdAt":  FormatTimeVal(c.CreatedAt),
		"updatedBy":    c.UpdatedBy,
		"updatedAt":  FormatTimeVal(c.UpdatedAt),
		"remark":      c.Remark,
	}
}

// DictTypeToMap 字典类型 VO。
func DictTypeToMap(d *models.DictType) map[string]any {
	if d == nil {
		return nil
	}
	return map[string]any{
		"dictId":     d.ID,
		"dictName":   d.DictName,
		"dictType":   d.DictType,
		"status":     d.Status,
		"createdBy":   d.CreatedBy,
		"createdAt": FormatTimeVal(d.CreatedAt),
		"updatedBy":   d.UpdatedBy,
		"updatedAt": FormatTimeVal(d.UpdatedAt),
		"remark":     d.Remark,
	}
}

// DictDataToMap 字典数据 VO。
func DictDataToMap(d *models.DictData) map[string]any {
	if d == nil {
		return nil
	}
	return map[string]any{
		"dictCode":   d.ID,
		"dictSort":   d.DictSort,
		"dictLabel":  d.DictLabel,
		"dictValue":  d.DictValue,
		"dictType":   d.DictType,
		"cssClass":   d.CSSClass,
		"listClass":  d.ListClass,
		"isDefault":  d.IsDefault,
		"status":     d.Status,
		"createdBy":   d.CreatedBy,
		"createdAt": FormatTimeVal(d.CreatedAt),
		"updatedBy":   d.UpdatedBy,
		"updatedAt": FormatTimeVal(d.UpdatedAt),
		"remark":     d.Remark,
	}
}

// NoticeToMap 公告 VO。
func NoticeToMap(n *models.Notice) map[string]any {
	if n == nil {
		return nil
	}
	content := ""
	if len(n.NoticeContent) > 0 {
		content = string(n.NoticeContent)
	}
	return map[string]any{
		"noticeId":      n.ID,
		"noticeTitle":   n.NoticeTitle,
		"noticeType":    n.NoticeType,
		"noticeContent": content,
		"status":        n.Status,
		"createdBy":      n.CreatedBy,
		"createdAt":    FormatTimeVal(n.CreatedAt),
		"updatedBy":      n.UpdatedBy,
		"updatedAt":    FormatTimeVal(n.UpdatedAt),
		"remark":        n.Remark,
	}
}

// RoleToMap 角色 VO。
func RoleToMap(r *models.Role) map[string]any {
	if r == nil {
		return nil
	}
	return map[string]any{
		"roleId":            r.ID,
		"roleName":          r.RoleName,
		"roleKey":           r.RoleKey,
		"roleSort":          r.RoleSort,
		"dataScope":         r.DataScope,
		"menuCheckStrictly": r.MenuCheckStrictly,
		"deptCheckStrictly": r.DeptCheckStrictly,
		"status":            r.Status,
		"createdBy":          r.CreatedBy,
		"createdAt":        FormatTimeVal(r.CreatedAt),
		"updatedBy":          r.UpdatedBy,
		"updatedAt":        FormatTimeVal(r.UpdatedAt),
		"remark":            r.Remark,
	}
}

// BuildDeptTreeSelect 构建部门 TreeSelect。
func BuildDeptTreeSelect(depts []models.Dept, parentID uint) []map[string]any {
	out := make([]map[string]any, 0)
	for _, d := range depts {
		if d.ParentID != parentID {
			continue
		}
		children := BuildDeptTreeSelect(depts, d.ID)
		node := map[string]any{
			"id":       d.ID,
			"label":    d.DeptName,
			"disabled": d.Status == "1",
		}
		if len(children) > 0 {
			node["children"] = children
		} else {
			node["children"] = []map[string]any{}
		}
		out = append(out, node)
	}
	return out
}

// BuildMenuTreeSelect 构建菜单 TreeSelect。
func BuildMenuTreeSelect(menus []models.Menu, parentID uint) []map[string]any {
	out := make([]map[string]any, 0)
	for _, m := range menus {
		if m.ParentID != parentID {
			continue
		}
		children := BuildMenuTreeSelect(menus, m.ID)
		node := map[string]any{
			"id":       m.ID,
			"label":    m.MenuName,
			"disabled": m.Status == "1",
		}
		if len(children) > 0 {
			node["children"] = children
		} else {
			node["children"] = []map[string]any{}
		}
		out = append(out, node)
	}
	return out
}

// ParseSortUpdate 解析排序更新参数。
func ParseSortUpdate(idsStr, numsStr string) ([]uint, []int, error) {
	idParts := strings.Split(idsStr, ",")
	numParts := strings.Split(numsStr, ",")
	if len(idParts) != len(numParts) {
		return nil, nil, errors.New("排序参数不匹配")
	}
	ids := make([]uint, 0, len(idParts))
	nums := make([]int, 0, len(numParts))
	for i := range idParts {
		idStr := strings.TrimSpace(idParts[i])
		numStr := strings.TrimSpace(numParts[i])
		if idStr == "" {
			continue
		}
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			return nil, nil, errors.New("无效的 ID")
		}
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return nil, nil, errors.New("无效的排序值")
		}
		ids = append(ids, uint(id))
		nums = append(nums, num)
	}
	return ids, nums, nil
}

// OperLogToMap 操作日志 VO。
func OperLogToMap(o *models.OperLog) map[string]any {
	if o == nil {
		return nil
	}
	return map[string]any{
		"operId":        o.ID,
		"title":         o.Title,
		"businessType":  o.BusinessType,
		"method":        o.Method,
		"requestMethod": o.RequestMethod,
		"operatorType":  o.OperatorType,
		"operName":      o.OperName,
		"deptName":      o.DeptName,
		"operUrl":       o.OperURL,
		"operIp":        o.OperIP,
		"operLocation":  o.OperLocation,
		"operParam":     o.OperParam,
		"jsonResult":    o.JSONResult,
		"status":        o.Status,
		"errorMsg":      o.ErrorMsg,
		"operTime":      FormatTime(o.OperTime),
		"costTime":      o.CostTime,
	}
}

// LogininforToMap 登录日志 VO。
func LogininforToMap(l *models.Logininfor) map[string]any {
	if l == nil {
		return nil
	}
	return map[string]any{
		"infoId":        l.ID,
		"userName":      l.UserName,
		"ipaddr":        l.IPAddr,
		"loginLocation": l.LoginLocation,
		"browser":       l.Browser,
		"os":            l.OS,
		"status":        l.Status,
		"msg":           l.Msg,
		"loginTime":     FormatTime(l.LoginTime),
	}
}

// JobToMap 定时任务 VO。
func JobToMap(j *models.Job) map[string]any {
	if j == nil {
		return nil
	}
	return map[string]any{
		"jobId":          j.ID,
		"jobName":        j.JobName,
		"jobGroup":       j.JobGroup,
		"invokeTarget":   j.InvokeTarget,
		"cronExpression": j.CronExpression,
		"misfirePolicy":  j.MisfirePolicy,
		"concurrent":     j.Concurrent,
		"status":         j.Status,
		"createdBy":       j.CreatedBy,
		"createdAt":     FormatTimeVal(j.CreatedAt),
		"updatedBy":       j.UpdatedBy,
		"updatedAt":     FormatTimeVal(j.UpdatedAt),
		"remark":         j.Remark,
	}
}

// JobLogToMap 任务日志 VO。
func JobLogToMap(j *models.JobLog) map[string]any {
	if j == nil {
		return nil
	}
	return map[string]any{
		"jobLogId":      j.ID,
		"jobName":       j.JobName,
		"jobGroup":      j.JobGroup,
		"invokeTarget":  j.InvokeTarget,
		"jobMessage":    j.JobMessage,
		"status":        j.Status,
		"exceptionInfo": j.ExceptionInfo,
		"startTime":     FormatTime(j.StartTime),
		"endTime":       FormatTime(j.EndTime),
	}
}

// GenTableToMap 代码生成表 VO。
func GenTableToMap(t *models.GenTable) map[string]any {
	if t == nil {
		return nil
	}
	return map[string]any{
		"tableId":        t.ID,
		"tableName":      t.Name,
		"tableComment":   t.TableComment,
		"subTableName":   t.SubTableName,
		"subTableFkName": t.SubTableFkName,
		"className":      t.ClassName,
		"tplCategory":    t.TplCategory,
		"tplWebType":     t.TplWebType,
		"packageName":    t.PackageName,
		"moduleName":     t.ModuleName,
		"businessName":   t.BusinessName,
		"functionName":   t.FunctionName,
		"functionAuthor": t.FunctionAuthor,
		"formColNum":     t.FormColNum,
		"genType":        t.GenType,
		"genPath":        t.GenPath,
		"options":        t.Options,
		"createdBy":       t.CreatedBy,
		"createdAt":     FormatTimeVal(t.CreatedAt),
		"updatedBy":       t.UpdatedBy,
		"updatedAt":     FormatTimeVal(t.UpdatedAt),
		"remark":         t.Remark,
	}
}

// GenTableColumnToMap 代码生成列 VO。
func GenTableColumnToMap(c *models.GenTableColumn) map[string]any {
	if c == nil {
		return nil
	}
	tableID := uint(0)
	if c.TableID != nil {
		tableID = *c.TableID
	}
	return map[string]any{
		"columnId":      c.ID,
		"tableId":       tableID,
		"columnName":    c.ColumnName,
		"columnComment": c.ColumnComment,
		"columnType":    c.ColumnType,
		"javaType":      c.JavaType,
		"javaField":     c.JavaField,
		"isPk":          c.IsPk,
		"isIncrement":   c.IsIncrement,
		"isRequired":    c.IsRequired,
		"isInsert":      c.IsInsert,
		"isEdit":        c.IsEdit,
		"isList":        c.IsList,
		"isQuery":       c.IsQuery,
		"queryType":     c.QueryType,
		"htmlType":      c.HTMLType,
		"dictType":      c.DictType,
		"sort":          c.Sort,
	}
}

// FlexInt 从 JSON 混合类型解析 int。
func FlexInt(v any, def int) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case string:
		if t == "" {
			return def
		}
		n, err := strconv.Atoi(t)
		if err != nil {
			return def
		}
		return n
	case int:
		return t
	case int64:
		return int(t)
	default:
		return def
	}
}
