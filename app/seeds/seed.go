package seeds

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/utils/pathx"

	"gorm.io/gorm"
)

// SeedIfEmpty 当用户表为空时执行 ry.sql 初始化。
func SeedIfEmpty(db *gorm.DB) error {
	var count int64
	if err := db.Model(&models.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return nil
	}
	sqlPath, err := locateRySQL()
	if err != nil {
		return err
	}
	raw, err := os.ReadFile(sqlPath)
	if err != nil {
		return fmt.Errorf("read ry.sql: %w", err)
	}
	stmts := splitSQL(string(raw))
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if err := db.Exec(stmt).Error; err != nil {
			return fmt.Errorf("exec sql: %w\nstmt: %s", err, truncate(stmt, 120))
		}
	}
	return nil
}

func locateRySQL() (string, error) {
	candidates := []string{
		"docs/ry.sql",
		filepath.Join(pathx.ResolveProjectRoot(), "docs", "ry.sql"),
	}
	for _, c := range candidates {
		if st, err := os.Stat(c); err == nil && !st.IsDir() {
			return c, nil
		}
	}
	return "", fmt.Errorf("docs/ry.sql not found")
}

func splitSQL(s string) []string {
	var out []string
	var b strings.Builder
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "--") {
			continue
		}
		b.WriteString(line)
		b.WriteByte('\n')
		if strings.HasSuffix(trim, ";") {
			out = append(out, b.String())
			b.Reset()
		}
	}
	if rem := strings.TrimSpace(b.String()); rem != "" {
		out = append(out, rem)
	}
	return out
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
