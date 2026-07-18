package app

import (
	// 触发 router 子包 init，与 handlers 目录结构对齐。
	_ "thirteen-admin/app/router/monitor"
	_ "thirteen-admin/app/router/system"
	_ "thirteen-admin/app/router/tool"
)
