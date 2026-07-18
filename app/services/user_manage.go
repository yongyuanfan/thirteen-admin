package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"clash-admin/app/models"
	"clash-admin/internal/base/model"
	"clash-admin/internal/base/service"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrManageNotFound   = errors.New("记录不存在")
	ErrManageConflict   = errors.New("记录已存在")
	ErrManageInvalid    = errors.New("参数无效")
	ErrManageForbidden  = errors.New("操作不被允许")
	ErrManageRoleInUse  = errors.New("角色仍被用户引用，无法删除")
	defaultUserPassword = "123456"
)

// UserVO 用户列表/详情响应。
type UserVO struct {
	ID         uint     `json:"id"`
	CreateBy   string   `json:"createBy"`
	CreateTime string   `json:"createTime"`
	UpdateBy   string   `json:"updateBy"`
	UpdateTime string   `json:"updateTime"`
	Status     string   `json:"status"`
	UserName   string   `json:"userName"`
	UserGender *string  `json:"userGender"`
	NickName   string   `json:"nickName"`
	UserPhone  string   `json:"userPhone"`
	UserEmail  string   `json:"userEmail"`
	UserRoles  []string `json:"userRoles"`
}

// UserListQuery 用户列表查询。
type UserListQuery struct {
	Current    int
	Size       int
	UserName   string
	UserGender string
	NickName   string
	UserPhone  string
	UserEmail  string
	Status     string
}

// UserUpsertInput 新增/更新用户入参。
type UserUpsertInput struct {
	ID         uint
	UserName   string
	UserGender *string
	NickName   string
	UserPhone  string
	UserEmail  string
	UserRoles  []string
	Status     string
	Operator   string
}

// UserManageService 用户管理业务。
type UserManageService struct {
	service.Service
}

// NewUserManageService 创建用户管理服务。
func NewUserManageService(db *gorm.DB) *UserManageService {
	return &UserManageService{Service: service.New(db)}
}

func (s *UserManageService) toVO(u *models.User) UserVO {
	phone := ""
	if u.Phone != nil {
		phone = *u.Phone
	}
	return UserVO{
		ID:         u.ID,
		CreateBy:   u.CreateBy,
		CreateTime: formatTime(u.CreatedAt),
		UpdateBy:   u.UpdateBy,
		UpdateTime: formatTime(u.UpdatedAt),
		Status:     statusString(u.Status),
		UserName:   u.UserName,
		UserGender: u.UserGender,
		NickName:   u.NickName,
		UserPhone:  phone,
		UserEmail:  u.UserEmail,
		UserRoles:  stringSliceOrEmpty(u.Roles),
	}
}

// GetUserList 分页查询用户。
func (s *UserManageService) GetUserList(ctx context.Context, q UserListQuery) (*PageResult[UserVO], error) {
	page, size, offset := parsePage(q.Current, q.Size)
	db := s.DB.WithContext(ctx).Model(&models.User{})
	if q.UserName != "" {
		db = db.Where("user_name LIKE ?", "%"+q.UserName+"%")
	}
	if q.NickName != "" {
		db = db.Where("nick_name LIKE ?", "%"+q.NickName+"%")
	}
	if q.UserPhone != "" {
		db = db.Where("phone LIKE ?", "%"+q.UserPhone+"%")
	}
	if q.UserEmail != "" {
		db = db.Where("user_email LIKE ?", "%"+q.UserEmail+"%")
	}
	if q.UserGender != "" {
		db = db.Where("user_gender = ?", q.UserGender)
	}
	if q.Status != "" {
		st, err := parseStatus(q.Status)
		if err != nil {
			return nil, err
		}
		db = db.Where("status = ?", st)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	var rows []models.User
	if err := db.Order("id DESC").Offset(offset).Limit(size).Find(&rows).Error; err != nil {
		return nil, err
	}
	records := make([]UserVO, 0, len(rows))
	for i := range rows {
		records = append(records, s.toVO(&rows[i]))
	}
	return &PageResult[UserVO]{Current: page, Size: size, Total: total, Records: records}, nil
}

// AddUser 新增用户。
func (s *UserManageService) AddUser(ctx context.Context, in UserUpsertInput) error {
	userName := strings.TrimSpace(in.UserName)
	if userName == "" {
		return fmt.Errorf("%w: userName 不能为空", ErrManageInvalid)
	}
	st, err := parseStatus(in.Status)
	if err != nil {
		return err
	}
	var count int64
	if err := s.DB.WithContext(ctx).Model(&models.User{}).Where("user_name = ?", userName).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrManageConflict
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultUserPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	var phone *string
	if p := strings.TrimSpace(in.UserPhone); p != "" {
		phone = &p
		var phoneCount int64
		if err := s.DB.WithContext(ctx).Model(&models.User{}).Where("phone = ?", p).Count(&phoneCount).Error; err != nil {
			return err
		}
		if phoneCount > 0 {
			return fmt.Errorf("%w: 手机号已存在", ErrManageConflict)
		}
	}
	roles := in.UserRoles
	if roles == nil {
		roles = []string{}
	}
	user := models.User{
		UserName:     userName,
		PasswordHash: string(hash),
		NickName:     strings.TrimSpace(in.NickName),
		UserGender:   in.UserGender,
		Phone:        phone,
		UserEmail:    strings.TrimSpace(in.UserEmail),
		Roles:        model.StringSlice(roles),
		Buttons:      model.StringSlice{},
		Status:       st,
		CreateBy:     in.Operator,
		UpdateBy:     in.Operator,
	}
	return s.DB.WithContext(ctx).Create(&user).Error
}

// UpdateUser 更新用户。
func (s *UserManageService) UpdateUser(ctx context.Context, in UserUpsertInput) error {
	if in.ID == 0 {
		return fmt.Errorf("%w: id 不能为空", ErrManageInvalid)
	}
	st, err := parseStatus(in.Status)
	if err != nil {
		return err
	}
	var user models.User
	if err := s.DB.WithContext(ctx).First(&user, in.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrManageNotFound
		}
		return err
	}
	userName := strings.TrimSpace(in.UserName)
	if userName == "" {
		return fmt.Errorf("%w: userName 不能为空", ErrManageInvalid)
	}
	var count int64
	if err := s.DB.WithContext(ctx).Model(&models.User{}).Where("user_name = ? AND id <> ?", userName, in.ID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrManageConflict
	}
	var phone *string
	if p := strings.TrimSpace(in.UserPhone); p != "" {
		phone = &p
		var phoneCount int64
		if err := s.DB.WithContext(ctx).Model(&models.User{}).Where("phone = ? AND id <> ?", p, in.ID).Count(&phoneCount).Error; err != nil {
			return err
		}
		if phoneCount > 0 {
			return fmt.Errorf("%w: 手机号已存在", ErrManageConflict)
		}
	}
	roles := in.UserRoles
	if roles == nil {
		roles = []string{}
	}
	updates := map[string]any{
		"user_name":   userName,
		"nick_name":   strings.TrimSpace(in.NickName),
		"user_gender": in.UserGender,
		"phone":       phone,
		"user_email":  strings.TrimSpace(in.UserEmail),
		"roles":       model.StringSlice(roles),
		"status":      st,
		"update_by":   in.Operator,
	}
	return s.DB.WithContext(ctx).Model(&user).Updates(updates).Error
}

// DeleteUsers 删除用户（禁止删除当前登录用户）。
func (s *UserManageService) DeleteUsers(ctx context.Context, ids []uint, currentUserID uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("%w: ids 不能为空", ErrManageInvalid)
	}
	for _, id := range ids {
		if id == currentUserID {
			return fmt.Errorf("%w: 不能删除当前登录用户", ErrManageForbidden)
		}
	}
	return s.DB.WithContext(ctx).Where("id IN ?", ids).Delete(&models.User{}).Error
}
