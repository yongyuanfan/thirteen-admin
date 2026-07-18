package pathx

import (
	"os"
	"path/filepath"
	"strings"
)

// ResolveProjectRoot 解析项目根目录，优先使用环境变量。
func ResolveProjectRoot() string {
	root := strings.TrimSpace(os.Getenv("PROJECT_ROOT"))
	if root == "" {
		cwd, err := os.Getwd()
		if err == nil {
			root = cwd
		}
	}
	if root == "" {
		root = "/"
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return filepath.Clean(root)
	}
	return filepath.Clean(absRoot)
}

// AbsoluteOrOriginal 返回绝对路径；若解析失败则回退原值。
func AbsoluteOrOriginal(path string) string {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return absPath
}

// IsUnderRoot 判断目标路径是否位于指定根目录内。
func IsUnderRoot(root, target string) bool {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	if rel == ".." {
		return false
	}
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// IsProtectedRepoPath 判断目标路径是否为受保护的仓库元数据目录。
func IsProtectedRepoPath(root, target string) bool {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return true
	}
	rel = filepath.ToSlash(rel)
	return rel == ".git" || strings.HasPrefix(rel, ".git/") || rel == ".gitmodules"
}
