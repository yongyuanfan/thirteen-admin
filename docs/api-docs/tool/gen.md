# 代码生成

前端：`src/api/tool/gen.ts`  
类型：`src/types/api/tool/gen.ts`  
页面：`views/tool/gen/**`  
说明：`tool/build` 为纯前端表单设计器，**无后端接口**。

---

## 数据模型

### GenTable

| 字段 | 类型 | 说明 |
|------|------|------|
| `tableId` | number | 主键 |
| `tableName` | string | 表名称 |
| `tableComment` | string | 表描述 |
| `subTableName` | string | 子表名 |
| `subTableFkName` | string | 子表外键名 |
| `className` | string | 实体类名 |
| `tplCategory` | `"crud"\|"tree"\|"sub"` | 模板类型 |
| `tplWebType` | `"element-ui"\|"element-plus"` | 前端模板（导入时前端还会传 `element-plus-typescript`） |
| `packageName` | string | 生成包路径 |
| `moduleName` | string | 模块名 |
| `businessName` | string | 业务名 |
| `functionName` | string | 功能名 |
| `functionAuthor` | string | 作者 |
| `genType` | `"0"\|"1"` | 0zip 1自定义路径 |
| `genPath` | string | 生成路径 |
| `options` | string | 其它选项 JSON |
| `treeCode` / `treeParentCode` / `treeName` | string | 树表字段 |
| `parentMenuId` / `parentMenuName` | string | 上级菜单 |
| `columns` | GenTableColumn[] | 列配置 |
| + BaseEntity | | |

### GenTableColumn

| 字段 | 类型 | 说明 |
|------|------|------|
| `columnId` | number | 主键 |
| `tableId` | number | 归属表 |
| `columnName` | string | 列名 |
| `columnComment` | string | 列描述 |
| `columnType` | string | 列类型 |
| `javaType` | string | Java/语言类型名（前端字段名历史遗留） |
| `javaField` | string | 字段名 |
| `isPk` / `isIncrement` / `isRequired` / `isInsert` / `isEdit` / `isList` / `isQuery` | `"0"\|"1"` | 标志位 |
| `queryType` | string | EQ/NE/GT/LT/LIKE/BETWEEN |
| `htmlType` | string | input/textarea/select/... |
| `dictType` | string | 字典类型 |
| `sort` | number | 排序 |

---

## GET `/tool/gen/list`

已导入的生成表分页。

### Query

`tableName`、`tableComment`、`params[beginTime/endTime]`、分页。

### 响应

`TableDataInfo<GenTable>`

---

## GET `/tool/gen/db/list`

数据库中尚未导入的表分页。Query 同 list。

---

## GET `/tool/gen/{tableId}`

表详情（含列与全部表选项）。

### 响应

```json
{
  "code": 200,
  "data": {
    "info": { "tableId": 1, "tableName": "sys_notice", "...": "..." },
    "rows": [{ "columnId": 1, "columnName": "id", "...": "..." }],
    "tables": [{ "tableId": 1, "tableName": "sys_notice" }]
  }
}
```

类型：`GenTableInfoResult` → `{ info, rows, tables }`

---

## PUT `/tool/gen`

保存生成配置。Body：`GenTable`（含 `columns`）。

---

## POST `/tool/gen/importTable`

从数据库导入表。**参数在 Query。**

| 参数 | 说明 |
|------|------|
| `tables` | 表名，逗号分隔 |
| `tplWebType` | 前端模板，如 `element-plus-typescript` |

---

## POST `/tool/gen/createTable`

根据 SQL 创建表并导入。**参数在 Query。**

| 参数 | 说明 |
|------|------|
| `sql` |建表 SQL |
| `tplWebType` | 前端模板类型 |

---

## GET `/tool/gen/preview/{tableId}`

预览生成代码。

### 响应

```json
{
  "code": 200,
  "data": {
    "vm/java/domain.java.vm": "package ...",
    "vm/vue/index.vue.vm": "<template>...</template>"
  }
}
```

`data` 为 `Record<path, codeString>`。

---

## DELETE `/tool/gen/{tableId}`

删除生成表配置。支持多 ID。

---

## GET `/tool/gen/genCode/{tableName}`

按自定义路径生成代码到服务器。

---

## GET `/tool/gen/synchDb/{tableName}`

与数据库结构同步列信息。

---

## GET `/tool/gen/batchGenCode`

批量下载 zip。

### Query

`tables`：表名逗号分隔。

### 响应

`application/octet-stream`（zip blob）。前端：`$download.zip`。
