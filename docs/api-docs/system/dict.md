# 字典管理

前端：`src/api/system/dict/type.ts`、`src/api/system/dict/data.ts`  
类型：`src/types/api/system/dict.ts`  
全局：`useDict` → `GET /system/dict/data/type/{dictType}`

---

## 数据模型

### SysDictType

| 字段 | 类型 | 说明 |
|------|------|------|
| `dictId` | number | 字典主键 |
| `dictName` | string | 字典名称 |
| `dictType` | string | 字典类型（唯一） |
| `status` | `"0"\|"1"` | 0正常 1停用 |
| + BaseEntity | | |

### SysDictData

| 字段 | 类型 | 说明 |
|------|------|------|
| `dictCode` | number | 字典编码（主键） |
| `dictSort` | number | 排序 |
| `dictLabel` | string | 标签 |
| `dictValue` | string | 键值 |
| `dictType` | string | 所属类型 |
| `cssClass` | string | 样式属性 |
| `listClass` | string | 表格回显样式（如 `primary`） |
| `isDefault` | `"Y"\|"N"` | 是否默认 |
| `status` | `"0"\|"1"` | 状态 |
| + BaseEntity | | |

---

## 字典类型

### GET `/system/dict/type/list`

Query：`dictName`、`dictType`、`status`、`params[beginTime/endTime]`、分页。  
响应：`TableDataInfo<SysDictType>`

### GET `/system/dict/type/{dictId}`

`{ "code": 200, "data": SysDictType }`

### POST `/system/dict/type` / PUT `/system/dict/type`

新增 / 修改。Body：`SysDictType`。

### DELETE `/system/dict/type/{dictId}`

支持多 ID。

### DELETE `/system/dict/type/refreshCache`

刷新字典缓存。

### GET `/system/dict/type/optionselect`

下拉选项：`{ "code": 200, "data": SysDictType[] }`

### POST `/system/dict/type/export`

导出 Excel。

---

## 字典数据

### GET `/system/dict/data/list`

Query：`dictType`、`dictLabel`、`status`、分页等。  
响应：`TableDataInfo<SysDictData>`

### GET `/system/dict/data/{dictCode}`

`{ "code": 200, "data": SysDictData }`

### GET `/system/dict/data/type/{dictType}`

按类型获取全部字典项（前端缓存用）。

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": [
    { "dictCode": 1, "dictLabel": "男", "dictValue": "0", "dictType": "sys_user_sex", "listClass": "", "isDefault": "Y" }
  ]
}
```

### POST `/system/dict/data` / PUT `/system/dict/data`

新增 / 修改。Body：`SysDictData`。

### DELETE `/system/dict/data/{dictCode}`

支持多 ID。

### POST `/system/dict/data/export`

导出 Excel。
