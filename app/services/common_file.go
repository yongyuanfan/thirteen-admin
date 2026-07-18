package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"thirteen-admin/internal/utils/pathx"
)

// FileService 本地文件上传下载。
type FileService struct {
	root string
}

// NewFileService 创建文件服务，根目录为项目下 profile。
func NewFileService() *FileService {
	root := filepath.Join(pathx.ResolveProjectRoot(), "profile")
	_ = os.MkdirAll(filepath.Join(root, "upload"), 0o755)
	return &FileService{root: root}
}

// SaveUpload 保存上传文件，返回相对路径与访问 URL 路径。
func (s *FileService) SaveUpload(filename string, r io.Reader) (fileName, newFileName, urlPath string, err error) {
	ext := filepath.Ext(filename)
	newFileName = fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	relDir := filepath.Join("upload", time.Now().Format("2006/01/02"))
	absDir := filepath.Join(s.root, relDir)
	if err = os.MkdirAll(absDir, 0o755); err != nil {
		return "", "", "", err
	}
	absPath := filepath.Join(absDir, newFileName)
	f, err := os.Create(absPath)
	if err != nil {
		return "", "", "", err
	}
	defer f.Close()
	if _, err = io.Copy(f, r); err != nil {
		return "", "", "", err
	}
	fileName = filepath.ToSlash(filepath.Join("/profile", relDir, newFileName))
	urlPath = fileName
	return fileName, newFileName, urlPath, nil
}

// Resolve 解析可下载路径（防目录穿越）。
func (s *FileService) Resolve(name string) (string, error) {
	name = strings.TrimPrefix(name, "/profile/")
	name = strings.TrimPrefix(name, "profile/")
	abs := filepath.Join(s.root, filepath.Clean(name))
	if !pathx.IsUnderRoot(s.root, abs) {
		return "", fmt.Errorf("非法路径")
	}
	return abs, nil
}

// Root 返回文件根目录。
func (s *FileService) Root() string { return s.root }
