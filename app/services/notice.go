package services

import (
	"errors"
	"strconv"
	"time"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NoticeService 通知公告业务。
type NoticeService struct {
	service.Service
}

// NewNoticeService 创建公告服务。
func NewNoticeService(db *gorm.DB) *NoticeService {
	return &NoticeService{Service: service.New(db)}
}

// List 管理端分页列表。
func (s *NoticeService) List(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	q := s.DB.Model(&models.Notice{})
	if title := ctx.Query("noticeTitle"); title != "" {
		q = q.Where("notice_title LIKE ?", "%"+title+"%")
	}
	if createBy := ctx.Query("createBy"); createBy != "" {
		q = q.Where("created_by LIKE ?", "%"+createBy+"%")
	}
	if typ := ctx.Query("noticeType"); typ != "" {
		q = q.Where("notice_type = ?", typ)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	q = ApplyOrder(q, ctx, "id desc")
	var notices []models.Notice
	if err := ApplyPage(q, page).Find(&notices).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(notices))
	for i := range notices {
		rows = append(rows, NoticeToMap(&notices[i]))
	}
	return total, rows, nil
}

// Get 获取公告详情。
func (s *NoticeService) Get(id uint) (map[string]any, error) {
	var notice models.Notice
	if err := s.DB.First(&notice, id).Error; err != nil {
		return nil, err
	}
	return NoticeToMap(&notice), nil
}

// Add 新增公告。
func (s *NoticeService) Add(notice *models.Notice, operator string) error {
	notice.CreatedBy = operator
	return s.DB.Create(notice).Error
}

// Edit 修改公告。
func (s *NoticeService) Edit(notice *models.Notice, operator string) error {
	notice.UpdatedBy = operator
	return s.DB.Model(notice).Updates(map[string]any{
		"notice_title":   notice.NoticeTitle,
		"notice_type":    notice.NoticeType,
		"notice_content": notice.NoticeContent,
		"status":         notice.Status,
		"remark":         notice.Remark,
		"updated_by":     notice.UpdatedBy,
	}).Error
}

// Delete 删除公告。
func (s *NoticeService) Delete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("notice_id IN ?", ids).Delete(&models.NoticeRead{}).Error; err != nil {
			return err
		}
		return tx.Where("id IN ?", ids).Delete(&models.Notice{}).Error
	})
}

// ListTop 顶栏公告列表。
func (s *NoticeService) ListTop(userID uint) ([]map[string]any, int64, error) {
	var notices []models.Notice
	if err := s.DB.Where("status = ?", "0").Order("id desc").Limit(20).Find(&notices).Error; err != nil {
		return nil, 0, err
	}
	readIDs := map[uint]struct{}{}
	if userID > 0 {
		var ids []uint
		_ = s.DB.Model(&models.NoticeRead{}).Where("user_id = ?", userID).Pluck("notice_id", &ids)
		for _, id := range ids {
			readIDs[id] = struct{}{}
		}
	}
	rows := make([]map[string]any, 0, len(notices))
	var unread int64
	for i := range notices {
		_, read := readIDs[notices[i].ID]
		if !read {
			unread++
		}
		rows = append(rows, map[string]any{
			"noticeId":    notices[i].ID,
			"noticeTitle": notices[i].NoticeTitle,
			"noticeType":  notices[i].NoticeType,
			"isRead":      read,
		})
	}
	return rows, unread, nil
}

// MarkRead 标记单条已读。
func (s *NoticeService) MarkRead(userID, noticeID uint) error {
	var count int64
	if err := s.DB.Model(&models.NoticeRead{}).
		Where("user_id = ? AND notice_id = ?", userID, noticeID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return s.DB.Create(&models.NoticeRead{
		NoticeID: noticeID,
		UserID:   userID,
		ReadTime: time.Now(),
	}).Error
}

// MarkReadAll 批量标记已读。
func (s *NoticeService) MarkReadAll(userID uint, ids []uint) error {
	for _, id := range ids {
		if err := s.MarkRead(userID, id); err != nil {
			return err
		}
	}
	return nil
}

// ReadUsersList 已读用户分页。
func (s *NoticeService) ReadUsersList(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	noticeIDStr := ctx.Query("noticeId")
	noticeID, _ := strconv.ParseUint(noticeIDStr, 10, 64)
	q := s.DB.Table("sys_notice_read nr").
		Select("u.id as user_id, u.user_name, u.nick_name, u.phonenumber, d.dept_name, nr.read_time").
		Joins("JOIN sys_user u ON u.id = nr.user_id").
		Joins("LEFT JOIN sys_dept d ON d.id = u.dept_id").
		Where("nr.notice_id = ?", noticeID)
	if kw := ctx.Query("searchValue"); kw != "" {
		q = q.Where("u.user_name LIKE ? OR u.nick_name LIKE ?", "%"+kw+"%", "%"+kw+"%")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	type readRow struct {
		UserID      uint
		UserName    string
		NickName    string
		Phonenumber string
		DeptName    string
		ReadTime    time.Time
	}
	var rowsData []readRow
	if err := ApplyPage(q.Order("nr.read_time desc"), page).Scan(&rowsData).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(rowsData))
	for _, r := range rowsData {
		rows = append(rows, map[string]any{
			"userId":      r.UserID,
			"userName":    r.UserName,
			"nickName":    r.NickName,
			"deptName":    r.DeptName,
			"phonenumber": r.Phonenumber,
			"readTime":    FormatTimeVal(r.ReadTime),
		})
	}
	return total, rows, nil
}
