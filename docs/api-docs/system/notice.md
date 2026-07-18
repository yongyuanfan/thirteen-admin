# 通知公告

前端：`src/api/system/notice.ts`  
类型：`src/types/api/system/notice.ts`  
页面：公告管理、顶栏铃铛 `HeaderNotice`

---

## 数据模型

### SysNotice

| 字段 | 类型 | 说明 |
|------|------|------|
| `noticeId` | number | 公告 ID |
| `noticeTitle` | string | 标题 |
| `noticeType` | `"1"\|"2"` | 1通知 2公告 |
| `noticeContent` | string | 内容（HTML） |
| `status` | `"0"\|"1"` | 0正常 1关闭 |
| `isRead` | boolean | 当前用户是否已读（顶栏列表） |
| + BaseEntity | | |

### NoticeReadUser

| 字段 | 类型 | 说明 |
|------|------|------|
| `userId` | number | 用户 ID |
| `userName` | string | 账号 |
| `nickName` | string | 昵称 |
| `deptName` | string | 部门名 |
| `phonenumber` | string | 手机 |
| `readTime` | string | 阅读时间 |

---

## GET `/system/notice/list`

管理端分页列表。

### Query

`noticeTitle`、`createdBy`、`noticeType`、分页。

### 响应

`TableDataInfo<SysNotice>`

---

## GET `/system/notice/{noticeId}`

`{ "code": 200, "data": SysNotice }`

---

## POST `/system/notice` / PUT `/system/notice`

新增 / 修改。Body：`SysNotice`。

---

## DELETE `/system/notice/{noticeId}`

支持多 ID。

---

## GET `/system/notice/listTop`

顶栏公告（含已读状态与未读数）。**需要 Token。**

### 响应

```json
{
  "code": 200,
  "msg": "操作成功",
  "data": [
    {
      "noticeId": 1,
      "noticeTitle": "温馨提醒",
      "noticeType": "2",
      "isRead": false
    }
  ],
  "unreadCount": 2
}
```

---

## POST `/system/notice/markRead`

标记单条已读。**参数在 Query。**

### Query

`noticeId`（number）

---

## POST `/system/notice/markReadAll`

批量标记已读。**参数在 Query。**

### Query

`ids`：公告 ID 逗号分隔，如 `1,2,3`

---

## GET `/system/notice/readUsers/list`

某公告的已读用户分页。

### Query

`noticeId`、`searchValue`（账号/昵称关键字）、分页。

### 响应

`TableDataInfo<NoticeReadUser>`
