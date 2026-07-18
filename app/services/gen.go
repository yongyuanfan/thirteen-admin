package services

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type dbTableRow struct {
	TableName    string `gorm:"column:TABLE_NAME"`
	TableComment string `gorm:"column:TABLE_COMMENT"`
}

type dbColumnRow struct {
	ColumnName    string `gorm:"column:COLUMN_NAME"`
	ColumnComment string `gorm:"column:COLUMN_COMMENT"`
	ColumnType    string `gorm:"column:COLUMN_TYPE"`
	ColumnKey     string `gorm:"column:COLUMN_KEY"`
	Extra         string `gorm:"column:EXTRA"`
	IsNullable    string `gorm:"column:IS_NULLABLE"`
}

// GenService 代码生成业务。
type GenService struct {
	service.Service
}

// NewGenService 创建代码生成服务。
func NewGenService(db *gorm.DB) *GenService {
	return &GenService{Service: service.New(db)}
}

func (s *GenService) buildListQuery(ctx *gin.Context) *gorm.DB {
	q := s.DB.Model(&models.GenTable{})
	if name := ctx.Query("tableName"); name != "" {
		q = q.Where("table_name LIKE ?", "%"+name+"%")
	}
	if comment := ctx.Query("tableComment"); comment != "" {
		q = q.Where("table_comment LIKE ?", "%"+comment+"%")
	}
	return ApplyTimeRange(q, ctx, "created_at")
}

// List 已导入的生成表分页。
func (s *GenService) List(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	q := s.buildListQuery(ctx)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	q = ApplyOrder(q, ctx, "id desc")
	var tables []models.GenTable
	if err := ApplyPage(q, page).Find(&tables).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(tables))
	for i := range tables {
		rows = append(rows, GenTableToMap(&tables[i]))
	}
	return total, rows, nil
}

// DbList 数据库中尚未导入的表。
func (s *GenService) DbList(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	all, err := s.loadDbTables(ctx.Query("tableName"), ctx.Query("tableComment"))
	if err != nil {
		return 0, nil, err
	}
	imported := map[string]bool{}
	var genTables []models.GenTable
	if err := s.DB.Select("table_name").Find(&genTables).Error; err != nil {
		return 0, nil, err
	}
	for _, t := range genTables {
		imported[t.Name] = true
	}
	filtered := make([]dbTableRow, 0)
	for _, t := range all {
		if imported[t.TableName] {
			continue
		}
		filtered = append(filtered, t)
	}
	total := int64(len(filtered))
	start := (page.PageNum - 1) * page.PageSize
	if start > len(filtered) {
		return total, []map[string]any{}, nil
	}
	end := start + page.PageSize
	if end > len(filtered) {
		end = len(filtered)
	}
	rows := make([]map[string]any, 0, end-start)
	for _, t := range filtered[start:end] {
		rows = append(rows, map[string]any{
			"tableName":    t.TableName,
			"tableComment": t.TableComment,
		})
	}
	return total, rows, nil
}

func (s *GenService) loadDbTables(nameFilter, commentFilter string) ([]dbTableRow, error) {
	q := s.DB.Table("information_schema.TABLES").
		Select("TABLE_NAME, TABLE_COMMENT").
		Where("TABLE_SCHEMA = DATABASE() AND TABLE_TYPE = 'BASE TABLE'")
	if nameFilter != "" {
		q = q.Where("TABLE_NAME LIKE ?", "%"+nameFilter+"%")
	}
	if commentFilter != "" {
		q = q.Where("TABLE_COMMENT LIKE ?", "%"+commentFilter+"%")
	}
	var rows []dbTableRow
	if err := q.Order("TABLE_NAME asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// Get 表详情（含列与全部表选项）。
func (s *GenService) Get(tableID uint) (map[string]any, error) {
	var table models.GenTable
	if err := s.DB.First(&table, tableID).Error; err != nil {
		return nil, err
	}
	var columns []models.GenTableColumn
	if err := s.DB.Where("table_id = ?", tableID).Order("sort asc, id asc").Find(&columns).Error; err != nil {
		return nil, err
	}
	colRows := make([]map[string]any, 0, len(columns))
	for i := range columns {
		colRows = append(colRows, GenTableColumnToMap(&columns[i]))
	}
	var allTables []models.GenTable
	if err := s.DB.Select("id, table_name, table_comment").Order("table_name asc").Find(&allTables).Error; err != nil {
		return nil, err
	}
	tableOpts := make([]map[string]any, 0, len(allTables))
	for i := range allTables {
		tableOpts = append(tableOpts, GenTableToMap(&allTables[i]))
	}
	info := GenTableToMap(&table)
	info["columns"] = colRows
	return map[string]any{
		"info":   info,
		"rows":   colRows,
		"tables": tableOpts,
	}, nil
}

// Update 保存生成配置。
func (s *GenService) Update(table *models.GenTable, columns []models.GenTableColumn, operator string) error {
	table.UpdatedBy = operator
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(table).Updates(map[string]any{
			"table_comment":     table.TableComment,
			"sub_table_name":    table.SubTableName,
			"sub_table_fk_name": table.SubTableFkName,
			"class_name":        table.ClassName,
			"tpl_category":      table.TplCategory,
			"tpl_web_type":      table.TplWebType,
			"package_name":      table.PackageName,
			"module_name":       table.ModuleName,
			"business_name":     table.BusinessName,
			"function_name":     table.FunctionName,
			"function_author":   table.FunctionAuthor,
			"form_col_num":      table.FormColNum,
			"gen_type":          table.GenType,
			"gen_path":          table.GenPath,
			"options":           table.Options,
			"remark":            table.Remark,
			"updated_by":        table.UpdatedBy,
		}).Error; err != nil {
			return err
		}
		if err := tx.Where("table_id = ?", table.ID).Delete(&models.GenTableColumn{}).Error; err != nil {
			return err
		}
		for i := range columns {
			columns[i].TableID = &table.ID
			columns[i].UpdatedBy = operator
			if columns[i].ID == 0 {
				columns[i].CreatedBy = operator
				if err := tx.Create(&columns[i]).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Save(&columns[i]).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// ImportTable 从数据库导入表。
func (s *GenService) ImportTable(tableNames []string, tplWebType, operator string) error {
	if len(tableNames) == 0 {
		return errors.New("请选择要导入的表")
	}
	for _, name := range tableNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		var count int64
		if err := s.DB.Model(&models.GenTable{}).Where("table_name = ?", name).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		cols, err := s.loadDbColumns(name)
		if err != nil {
			return err
		}
		className := toClassName(name)
		table := &models.GenTable{
			Name:           name,
			TableComment:   "",
			ClassName:      className,
			TplCategory:    "crud",
			TplWebType:     tplWebType,
			PackageName:    "thirteen-admin/app",
			ModuleName:     strings.TrimPrefix(name, "sys_"),
			BusinessName:   strings.TrimPrefix(name, "sys_"),
			FunctionName:   className,
			FunctionAuthor: operator,
			GenType:        "0",
			GenPath:        "/",
			CreatedBy:      operator,
		}
		if err := s.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(table).Error; err != nil {
				return err
			}
			for i, c := range cols {
				col := buildGenColumn(table.ID, c, i+1, operator)
				if err := tx.Create(&col).Error; err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}

// CreateTable 根据 SQL 建表并导入。
func (s *GenService) CreateTable(sql, tplWebType, operator string) error {
	sql = strings.TrimSpace(sql)
	if sql == "" {
		return errors.New("SQL 不能为空")
	}
	if err := s.DB.Exec(sql).Error; err != nil {
		return err
	}
	tableName := parseCreateTableName(sql)
	if tableName == "" {
		return nil
	}
	return s.ImportTable([]string{tableName}, tplWebType, operator)
}

// Preview 预览生成代码。
func (s *GenService) Preview(tableID uint) (map[string]string, error) {
	detail, err := s.Get(tableID)
	if err != nil {
		return nil, err
	}
	info, _ := detail["info"].(map[string]any)
	rows, _ := detail["rows"].([]map[string]any)
	tableName := toString(info["tableName"])
	className := toString(info["className"])
	if className == "" {
		className = toClassName(tableName)
	}
	goCode, err := renderGoStruct(className, tableName, rows)
	if err != nil {
		return nil, err
	}
	vueCode, err := renderVueStub(className, tableName, rows)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"vm/go/domain.go.vm":  goCode,
		"vm/vue/index.vue.vm": vueCode,
	}, nil
}

// Delete 删除生成表配置。
func (s *GenService) Delete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("table_id IN ?", ids).Delete(&models.GenTableColumn{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", ids).Delete(&models.GenTable{}).Error
	})
}

// SynchDb 与数据库结构同步列信息。
func (s *GenService) SynchDb(tableName, operator string) error {
	var table models.GenTable
	if err := s.DB.Where("table_name = ?", tableName).First(&table).Error; err != nil {
		return err
	}
	cols, err := s.loadDbColumns(tableName)
	if err != nil {
		return err
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("table_id = ?", table.ID).Delete(&models.GenTableColumn{}).Error; err != nil {
			return err
		}
		for i, c := range cols {
			col := buildGenColumn(table.ID, c, i+1, operator)
			if err := tx.Create(&col).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GenCode 按自定义路径生成代码（占位实现）。
func (s *GenService) GenCode(tableName string) error {
	var count int64
	if err := s.DB.Model(&models.GenTable{}).Where("table_name = ?", tableName).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("生成表不存在")
	}
	return nil
}

// BatchGenCode 批量生成 zip。
func (s *GenService) BatchGenCode(tableNames []string) ([]byte, string, error) {
	if len(tableNames) == 0 {
		return nil, "", errors.New("请选择要生成的表")
	}
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	for _, name := range tableNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		var table models.GenTable
		if err := s.DB.Where("table_name = ?", name).First(&table).Error; err != nil {
			return nil, "", err
		}
		preview, err := s.Preview(table.ID)
		if err != nil {
			return nil, "", err
		}
		for path, code := range preview {
			f, err := zw.Create(fmt.Sprintf("%s/%s", name, path))
			if err != nil {
				return nil, "", err
			}
			if _, err := f.Write([]byte(code)); err != nil {
				return nil, "", err
			}
		}
	}
	if err := zw.Close(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), "thirteen-gen.zip", nil
}

func (s *GenService) loadDbColumns(tableName string) ([]dbColumnRow, error) {
	var cols []dbColumnRow
	err := s.DB.Table("information_schema.COLUMNS").
		Select("COLUMN_NAME, COLUMN_COMMENT, COLUMN_TYPE, COLUMN_KEY, EXTRA, IS_NULLABLE").
		Where("TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?", tableName).
		Order("ORDINAL_POSITION asc").
		Find(&cols).Error
	return cols, err
}

func buildGenColumn(tableID uint, c dbColumnRow, sort int, operator string) models.GenTableColumn {
	isPK := "0"
	isInc := "0"
	if c.ColumnKey == "PRI" {
		isPK = "1"
	}
	if strings.Contains(strings.ToLower(c.Extra), "auto_increment") {
		isInc = "1"
	}
	required := "0"
	if c.IsNullable == "NO" && isPK == "0" {
		required = "1"
	}
	field := toCamelLower(c.ColumnName)
	skip := isPK == "1" || field == "createdAt" || field == "updatedAt"
	return models.GenTableColumn{
		TableID:       &tableID,
		ColumnName:    c.ColumnName,
		ColumnComment: c.ColumnComment,
		ColumnType:    c.ColumnType,
		JavaType:      mapGoType(c.ColumnType),
		JavaField:     field,
		IsPk:          isPK,
		IsIncrement:   isInc,
		IsRequired:    required,
		IsInsert:      boolFlag(!skip),
		IsEdit:        boolFlag(!skip),
		IsList:        "1",
		IsQuery:       boolFlag(!skip && strings.Contains(c.ColumnName, "name")),
		QueryType:     "EQ",
		HTMLType:      "input",
		Sort:          sort,
		CreatedBy:     operator,
	}
}

func boolFlag(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

func toClassName(tableName string) string {
	parts := strings.Split(tableName, "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(p[1:])
		}
	}
	return b.String()
}

func toCamelLower(s string) string {
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}
	var b strings.Builder
	for i, p := range parts {
		if p == "" {
			continue
		}
		if i == 0 {
			b.WriteString(strings.ToLower(p))
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(strings.ToLower(p[1:]))
		}
	}
	return b.String()
}

func mapGoType(columnType string) string {
	t := strings.ToLower(columnType)
	switch {
	case strings.Contains(t, "bigint"):
		return "int64"
	case strings.Contains(t, "int"):
		return "int"
	case strings.Contains(t, "decimal"), strings.Contains(t, "numeric"), strings.Contains(t, "float"), strings.Contains(t, "double"):
		return "float64"
	case strings.Contains(t, "datetime"), strings.Contains(t, "timestamp"), strings.Contains(t, "date"):
		return "time.Time"
	default:
		return "string"
	}
}

func parseCreateTableName(sql string) string {
	upper := strings.ToUpper(sql)
	idx := strings.Index(upper, "CREATE TABLE")
	if idx < 0 {
		return ""
	}
	rest := strings.TrimSpace(sql[idx+len("CREATE TABLE"):])
	rest = strings.TrimLeft(rest, " IF NOT EXISTS")
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return ""
	}
	if rest[0] == '`' {
		if end := strings.Index(rest[1:], "`"); end >= 0 {
			return rest[1 : end+1]
		}
	}
	parts := strings.FieldsFunc(rest, func(r rune) bool {
		return r == ' ' || r == '(' || r == ';'
	})
	if len(parts) > 0 {
		return strings.Trim(parts[0], "`")
	}
	return ""
}

var goStructTmpl = template.Must(template.New("go").Parse(`package models

// {{.ClassName}} {{.TableComment}}。
type {{.ClassName}} struct {
{{range .Fields}}	{{.Name}} {{.Type}} ` + "`json:\"{{.JSON}}\"`" + `
{{end}}}
`))

var vueStubTmpl = template.Must(template.New("vue").Parse(`<template>
  <div class="app-container">
    <el-table :data="list">
{{range .Columns}}      <el-table-column prop="{{.Field}}" label="{{.Label}}" />
{{end}}    </el-table>
  </div>
</template>

<script setup lang="ts">
// {{.ClassName}} generated from {{.TableName}}
</script>
`))

func renderGoStruct(className, tableName string, rows []map[string]any) (string, error) {
	type field struct {
		Name, Type, JSON string
	}
	fields := []field{{Name: "ID", Type: "uint", JSON: "id"}}
	for _, row := range rows {
		if toString(row["isPk"]) == "1" {
			continue
		}
		fields = append(fields, field{
			Name: toExportName(toString(row["javaField"])),
			Type: toString(row["javaType"]),
			JSON: toString(row["javaField"]),
		})
	}
	data := map[string]any{
		"ClassName":    className,
		"TableComment": tableName,
		"Fields":       fields,
	}
	var buf bytes.Buffer
	if err := goStructTmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderVueStub(className, tableName string, rows []map[string]any) (string, error) {
	type col struct {
		Field, Label string
	}
	cols := make([]col, 0)
	for _, row := range rows {
		if toString(row["isList"]) != "1" {
			continue
		}
		cols = append(cols, col{
			Field: toString(row["javaField"]),
			Label: toString(row["columnComment"]),
		})
	}
	data := map[string]any{
		"ClassName": className,
		"TableName": tableName,
		"Columns":   cols,
	}
	var buf bytes.Buffer
	if err := vueStubTmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func toExportName(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}
