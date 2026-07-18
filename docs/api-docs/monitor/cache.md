# 缓存监控

前端：`src/api/monitor/cache.ts`  
类型：`src/types/api/monitor/cache.ts`  
页面：`views/monitor/cache/index.vue`、`list.vue`

---

## 数据模型 SysCache

| 字段 | 类型 | 说明 |
|------|------|------|
| `cacheName` | string | 缓存名称 |
| `cacheKey` | string | 缓存键 |
| `cacheValue` | string | 缓存内容 |
| `remark` | string | 备注 |

---

## GET `/monitor/cache`

Redis 监控概览（监控首页）。

### 响应 `data` 形状（前端使用）

```json
{
  "code": 200,
  "data": {
    "info": {
      "redis_version": "7.0.0",
      "redis_mode": "standalone",
      "tcp_port": 6379,
      "connected_clients": 1,
      "uptime_in_days": 1,
      "used_memory_human": "1M",
      "used_cpu_user_children": "0.01",
      "maxmemory_human": "0B",
      "aof_enabled": "0",
      "rdb_last_bgsave_status": "ok",
      "instantaneous_input_kbps": "0.00",
      "instantaneous_output_kbps": "0.00"
    },
    "dbSize": 10,
    "commandStats": [
      { "name": "get", "value": "100" },
      { "name": "set", "value": "50" }
    ]
  }
}
```

字段名需与 Redis INFO 及前端展示兼容；`commandStats` 为命令统计列表。

---

## GET `/monitor/cache/getNames`

缓存名称列表。

### 响应

`{ "code": 200, "data": SysCache[] }`（常用 `cacheName`、`remark`）

---

## GET `/monitor/cache/getKeys/{cacheName}`

某缓存名下的全部 key。

### 响应

`{ "code": 200, "data": ["key1", "key2"] }`

---

## GET `/monitor/cache/getValue/{cacheName}/{cacheKey}`

查询缓存值。

### 响应

`{ "code": 200, "data": SysCache }`

---

## DELETE `/monitor/cache/clearCacheName/{cacheName}`

按缓存名清空。

---

## DELETE `/monitor/cache/clearCacheKey/{cacheKey}`

按 key 删除。

---

## DELETE `/monitor/cache/clearCacheAll`

清空全部缓存。
