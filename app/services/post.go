package services

import (
	"errors"
	"strconv"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PostService 岗位业务。
type PostService struct {
	service.Service
}

// NewPostService 创建岗位服务。
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{Service: service.New(db)}
}

// List 分页查询岗位。
func (s *PostService) List(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	q := s.DB.Model(&models.Post{})
	if code := ctx.Query("postCode"); code != "" {
		q = q.Where("post_code LIKE ?", "%"+code+"%")
	}
	if name := ctx.Query("postName"); name != "" {
		q = q.Where("post_name LIKE ?", "%"+name+"%")
	}
	if status := ctx.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	q = ApplyOrder(q, ctx, "post_sort asc")
	var posts []models.Post
	if err := ApplyPage(q, page).Find(&posts).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(posts))
	for i := range posts {
		rows = append(rows, PostToMap(&posts[i]))
	}
	return total, rows, nil
}

// Get 获取岗位详情。
func (s *PostService) Get(id uint) (map[string]any, error) {
	var post models.Post
	if err := s.DB.First(&post, id).Error; err != nil {
		return nil, err
	}
	return PostToMap(&post), nil
}

// Add 新增岗位。
func (s *PostService) Add(post *models.Post, operator string) error {
	post.CreatedBy = operator
	return s.DB.Create(post).Error
}

// Edit 修改岗位。
func (s *PostService) Edit(post *models.Post, operator string) error {
	post.UpdatedBy = operator
	return s.DB.Model(post).Updates(map[string]any{
		"post_code":  post.PostCode,
		"post_name":  post.PostName,
		"post_sort":  post.PostSort,
		"status":     post.Status,
		"remark":     post.Remark,
		"updated_by": post.UpdatedBy,
	}).Error
}

// Delete 删除岗位。
func (s *PostService) Delete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("post_id IN ?", ids).Delete(&models.UserPost{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", ids).Delete(&models.Post{}).Error
	})
}

// Export 导出岗位 CSV 行。
func (s *PostService) Export(ctx *gin.Context) ([][]string, error) {
	_, rows, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, []string{
			toString(row["postId"]),
			toString(row["postCode"]),
			toString(row["postName"]),
			toString(row["postSort"]),
			toString(row["status"]),
			toString(row["createdAt"]),
		})
	}
	return out, nil
}

// ListAll 全部岗位选项。
func (s *PostService) ListAll() ([]map[string]any, error) {
	var posts []models.Post
	if err := s.DB.Where("status = ?", "0").Order("post_sort asc").Find(&posts).Error; err != nil {
		return nil, err
	}
	rows := make([]map[string]any, 0, len(posts))
	for i := range posts {
		rows = append(rows, PostToMap(&posts[i]))
	}
	return rows, nil
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case uint:
		return strconv.FormatUint(uint64(t), 10)
	default:
		return ""
	}
}
