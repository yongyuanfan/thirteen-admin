# 接口清单（速查）

路径均为去掉 `/dev-api` 等前缀后的后端路径。Auth：`✓` 需 Token，`-` 不需。

## 认证与路由

| Method | Path | Auth |
|--------|------|------|
| POST | `/login` | - |
| POST | `/register` | - |
| GET | `/captchaImage` | - |
| GET | `/getInfo` | ✓ |
| POST | `/logout` | ✓ |
| POST | `/unlockscreen` | ✓ |
| GET | `/getRouters` | ✓ |

## 公共

| Method | Path | Auth |
|--------|------|------|
| POST | `/common/upload` | ✓ |
| GET | `/common/download` | ✓ |
| GET | `/common/download/resource` | ✓ |

## 系统 · 用户

| Method | Path | Auth |
|--------|------|------|
| GET | `/system/user/list` | ✓ |
| GET | `/system/user/` | ✓ |
| GET | `/system/user/{userId}` | ✓ |
| POST | `/system/user` | ✓ |
| PUT | `/system/user` | ✓ |
| DELETE | `/system/user/{userId}` | ✓ |
| PUT | `/system/user/resetPwd` | ✓ |
| PUT | `/system/user/changeStatus` | ✓ |
| GET | `/system/user/profile` | ✓ |
| PUT | `/system/user/profile` | ✓ |
| PUT | `/system/user/profile/updatePwd` | ✓ |
| POST | `/system/user/profile/avatar` | ✓ |
| GET | `/system/user/authRole/{userId}` | ✓ |
| PUT | `/system/user/authRole` | ✓ |
| GET | `/system/user/deptTree` | ✓ |
| POST | `/system/user/export` | ✓ |
| POST | `/system/user/importData` | ✓ |
| POST | `/system/user/importTemplate` | ✓ |

## 系统 · 角色

| Method | Path | Auth |
|--------|------|------|
| GET | `/system/role/list` | ✓ |
| GET | `/system/role/{roleId}` | ✓ |
| POST | `/system/role` | ✓ |
| PUT | `/system/role` | ✓ |
| PUT | `/system/role/dataScope` | ✓ |
| PUT | `/system/role/changeStatus` | ✓ |
| DELETE | `/system/role/{roleId}` | ✓ |
| GET | `/system/role/authUser/allocatedList` | ✓ |
| GET | `/system/role/authUser/unallocatedList` | ✓ |
| PUT | `/system/role/authUser/cancel` | ✓ |
| PUT | `/system/role/authUser/cancelAll` | ✓ |
| PUT | `/system/role/authUser/selectAll` | ✓ |
| GET | `/system/role/deptTree/{roleId}` | ✓ |
| POST | `/system/role/export` | ✓ |

## 系统 · 菜单

| Method | Path | Auth |
|--------|------|------|
| GET | `/system/menu/list` | ✓ |
| GET | `/system/menu/{menuId}` | ✓ |
| GET | `/system/menu/treeselect` | ✓ |
| GET | `/system/menu/roleMenuTreeselect/{roleId}` | ✓ |
| POST | `/system/menu` | ✓ |
| PUT | `/system/menu` | ✓ |
| PUT | `/system/menu/updateSort` | ✓ |
| DELETE | `/system/menu/{menuId}` | ✓ |

## 系统 · 部门

| Method | Path | Auth |
|--------|------|------|
| GET | `/system/dept/list` | ✓ |
| GET | `/system/dept/list/exclude/{deptId}` | ✓ |
| GET | `/system/dept/{deptId}` | ✓ |
| POST | `/system/dept` | ✓ |
| PUT | `/system/dept` | ✓ |
| PUT | `/system/dept/updateSort` | ✓ |
| DELETE | `/system/dept/{deptId}` | ✓ |

## 系统 · 岗位

| Method | Path | Auth |
|--------|------|------|
| GET | `/system/post/list` | ✓ |
| GET | `/system/post/{postId}` | ✓ |
| POST | `/system/post` | ✓ |
| PUT | `/system/post` | ✓ |
| DELETE | `/system/post/{postId}` | ✓ |
| POST | `/system/post/export` | ✓ |

## 系统 · 字典

| Method | Path | Auth |
|--------|------|------|
| GET | `/system/dict/type/list` | ✓ |
| GET | `/system/dict/type/{dictId}` | ✓ |
| POST | `/system/dict/type` | ✓ |
| PUT | `/system/dict/type` | ✓ |
| DELETE | `/system/dict/type/{dictId}` | ✓ |
| DELETE | `/system/dict/type/refreshCache` | ✓ |
| GET | `/system/dict/type/optionselect` | ✓ |
| POST | `/system/dict/type/export` | ✓ |
| GET | `/system/dict/data/list` | ✓ |
| GET | `/system/dict/data/{dictCode}` | ✓ |
| GET | `/system/dict/data/type/{dictType}` | ✓ |
| POST | `/system/dict/data` | ✓ |
| PUT | `/system/dict/data` | ✓ |
| DELETE | `/system/dict/data/{dictCode}` | ✓ |
| POST | `/system/dict/data/export` | ✓ |

## 系统 · 参数

| Method | Path | Auth |
|--------|------|------|
| GET | `/system/config/list` | ✓ |
| GET | `/system/config/{configId}` | ✓ |
| GET | `/system/config/configKey/{configKey}` | ✓ |
| POST | `/system/config` | ✓ |
| PUT | `/system/config` | ✓ |
| DELETE | `/system/config/{configId}` | ✓ |
| DELETE | `/system/config/refreshCache` | ✓ |
| POST | `/system/config/export` | ✓ |

## 系统 · 公告

| Method | Path | Auth |
|--------|------|------|
| GET | `/system/notice/list` | ✓ |
| GET | `/system/notice/{noticeId}` | ✓ |
| POST | `/system/notice` | ✓ |
| PUT | `/system/notice` | ✓ |
| DELETE | `/system/notice/{noticeId}` | ✓ |
| GET | `/system/notice/listTop` | ✓ |
| POST | `/system/notice/markRead` | ✓ |
| POST | `/system/notice/markReadAll` | ✓ |
| GET | `/system/notice/readUsers/list` | ✓ |

## 监控

| Method | Path | Auth |
|--------|------|------|
| GET | `/monitor/job/list` | ✓ |
| GET | `/monitor/job/{jobId}` | ✓ |
| POST | `/monitor/job` | ✓ |
| PUT | `/monitor/job` | ✓ |
| DELETE | `/monitor/job/{jobId}` | ✓ |
| PUT | `/monitor/job/changeStatus` | ✓ |
| PUT | `/monitor/job/run` | ✓ |
| POST | `/monitor/job/export` | ✓ |
| GET | `/monitor/jobLog/list` | ✓ |
| DELETE | `/monitor/jobLog/{jobLogId}` | ✓ |
| DELETE | `/monitor/jobLog/clean` | ✓ |
| POST | `/monitor/jobLog/export` | ✓ |
| GET | `/monitor/operlog/list` | ✓ |
| DELETE | `/monitor/operlog/{operId}` | ✓ |
| DELETE | `/monitor/operlog/clean` | ✓ |
| POST | `/monitor/operlog/export` | ✓ |
| GET | `/monitor/logininfor/list` | ✓ |
| DELETE | `/monitor/logininfor/{infoId}` | ✓ |
| GET | `/monitor/logininfor/unlock/{userName}` | ✓ |
| DELETE | `/monitor/logininfor/clean` | ✓ |
| POST | `/monitor/logininfor/export` | ✓ |
| GET | `/monitor/online/list` | ✓ |
| DELETE | `/monitor/online/{tokenId}` | ✓ |
| GET | `/monitor/cache` | ✓ |
| GET | `/monitor/cache/getNames` | ✓ |
| GET | `/monitor/cache/getKeys/{cacheName}` | ✓ |
| GET | `/monitor/cache/getValue/{cacheName}/{cacheKey}` | ✓ |
| DELETE | `/monitor/cache/clearCacheName/{cacheName}` | ✓ |
| DELETE | `/monitor/cache/clearCacheKey/{cacheKey}` | ✓ |
| DELETE | `/monitor/cache/clearCacheAll` | ✓ |
| GET | `/monitor/server` | ✓ |

## 代码生成

| Method | Path | Auth |
|--------|------|------|
| GET | `/tool/gen/list` | ✓ |
| GET | `/tool/gen/db/list` | ✓ |
| GET | `/tool/gen/{tableId}` | ✓ |
| PUT | `/tool/gen` | ✓ |
| POST | `/tool/gen/importTable` | ✓ |
| POST | `/tool/gen/createTable` | ✓ |
| GET | `/tool/gen/preview/{tableId}` | ✓ |
| DELETE | `/tool/gen/{tableId}` | ✓ |
| GET | `/tool/gen/genCode/{tableName}` | ✓ |
| GET | `/tool/gen/synchDb/{tableName}` | ✓ |
| GET | `/tool/gen/batchGenCode` | ✓ |
