# 菜单管理

前端：`src/api/system/menu.ts`  
类型：`src/types/api/system/menu.ts`

动态路由见 [auth.md](../auth.md) 的 `GET /getRouters`。

---

## 数据模型 SysMenu

| 字段 | 类型 | 说明 |
|------|------|------|
| `menuId` | number | 菜单 ID |
| `parentId` | number | 父菜单 ID，根为 `0` |
| `menuName` | string | 名称 |
| `orderNum` | number | 显示顺序 |
| `path` | string | 路由地址 |
| `component` | string | 组件路径 |
| `query` | string | 路由参数 |
| `routeName` | string | 路由名称 |
| `perms` | string | 权限标识，如 `system:user:list` |
| `icon` | string | 图标 |
| `isFrame` | `"0"\|"1"` | 0外链 1否（注意：与库中 int 语义一致即可，JSON 用字符串） |
| `isCache` | `"0"\|"1"` | 0缓存 1不缓存 |
| `menuType` | `"M"\|"C"\|"F"` | 目录 / 菜单 / 按钮 |
| `visible` | `"0"\|"1"` | 0显示 1隐藏 |
| `status` | `"0"\|"1"` | 0正常 1停用 |
| `children` | SysMenu[] | 子节点（树形时） |
| + BaseEntity | | |

---

## GET `/system/menu/list`

查询菜单列表（通常返回**扁平列表**，前端 `handleTree` 建树）。

### Query

`menuName`、`status`

### 响应

`{ "code": 200, "data": SysMenu[] }`

---

## GET `/system/menu/{menuId}`

### 响应

`{ "code": 200, "data": SysMenu }`

---

## GET `/system/menu/treeselect`

菜单下拉树（TreeSelect）。

### 响应

`{ "code": 200, "data": TreeSelect[] }`

---

## GET `/system/menu/roleMenuTreeselect/{roleId}`

角色授权用：菜单树 + 已选中 keys。

### 响应

```json
{
  "code": 200,
  "menus": [{ "id": 1, "label": "系统管理", "disabled": false, "children": [] }],
  "checkedKeys": [100, 1000, 1001]
}
```

---

## POST `/system/menu` / PUT `/system/menu`

新增 / 修改。Body：`SysMenu`。

---

## PUT `/system/menu/updateSort`

批量更新排序。

### 请求体

```json
{
  "menuIds": "1,2,3",
  "orderNums": "1,2,3"
}
```

`menuIds` 与 `orderNums` 按位置一一对应（逗号分隔字符串）。

---

## DELETE `/system/menu/{menuId}`

删除菜单（通常需无子节点）。
