package services

import (
	"encoding/csv"
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"thirteen-admin/app/models"
	"thirteen-admin/internal/base/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户业务。
type UserService struct {
	service.Service
	dept *DeptService
	post *PostService
}

// NewUserService 创建用户服务。
func NewUserService(db *gorm.DB, dept *DeptService, post *PostService) *UserService {
	return &UserService{
		Service: service.New(db),
		dept:    dept,
		post:    post,
	}
}

// List 分页查询用户。
func (s *UserService) List(ctx *gin.Context) (int64, []map[string]any, error) {
	page := ParsePage(ctx)
	q := s.DB.Model(&models.User{}).Where("del_flag = ?", "0")
	if name := ctx.Query("userName"); name != "" {
		q = q.Where("user_name LIKE ?", "%"+name+"%")
	}
	if phone := ctx.Query("phonenumber"); phone != "" {
		q = q.Where("phonenumber LIKE ?", "%"+phone+"%")
	}
	if status := ctx.Query("status"); status != "" {
		q = q.Where("status = ?", status)
	}
	if deptIDStr := ctx.Query("deptId"); deptIDStr != "" {
		deptID, _ := strconv.ParseUint(deptIDStr, 10, 64)
		if deptID > 0 {
			ids, err := s.dept.ChildDeptIDs(uint(deptID))
			if err != nil {
				return 0, nil, err
			}
			q = q.Where("dept_id IN ?", ids)
		}
	}
	q = ApplyTimeRange(q, ctx, "created_at")
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return 0, nil, err
	}
	q = ApplyOrder(q, ctx, "id asc")
	var users []models.User
	if err := ApplyPage(q, page).Find(&users).Error; err != nil {
		return 0, nil, err
	}
	rows := make([]map[string]any, 0, len(users))
	for i := range users {
		row, err := s.userDetailMap(&users[i], false)
		if err != nil {
			return 0, nil, err
		}
		rows = append(rows, row)
	}
	return total, rows, nil
}

// FormData 用户表单数据。
func (s *UserService) FormData(userID uint) (user map[string]any, posts, roles []map[string]any, postIDs, roleIDs []uint, err error) {
	posts, err = s.post.ListAll()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	roles, err = s.listRoleOptions()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	if userID == 0 {
		return nil, posts, roles, []uint{}, []uint{}, nil
	}
	var u models.User
	if err = s.DB.Where("id = ? AND del_flag = ?", userID, "0").First(&u).Error; err != nil {
		return nil, nil, nil, nil, nil, err
	}
	user, err = s.userDetailMap(&u, true)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	postIDs, _ = s.userPostIDs(userID)
	roleIDs, _ = s.userRoleIDs(userID)
	return user, posts, roles, postIDs, roleIDs, nil
}

// Get 获取用户详情。
func (s *UserService) Get(id uint) (map[string]any, error) {
	var u models.User
	if err := s.DB.Where("id = ? AND del_flag = ?", id, "0").First(&u).Error; err != nil {
		return nil, err
	}
	return s.userDetailMap(&u, true)
}

// Add 新增用户。
func (s *UserService) Add(u *models.User, roleIDs, postIDs []uint, operator string) error {
	var count int64
	if err := s.DB.Model(&models.User{}).Where("user_name = ? AND del_flag = ?", u.UserName, "0").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("新增用户'" + u.UserName + "'失败，登录账号已存在")
	}
	if u.Password == "" {
		pwd, _ := s.defaultPassword()
		u.Password = pwd
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	u.DelFlag = "0"
	u.CreatedBy = operator
	now := time.Now()
	u.PwdUpdateDate = &now
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(u).Error; err != nil {
			return err
		}
		return s.saveUserRelations(tx, u.ID, roleIDs, postIDs)
	})
}

// Edit 修改用户。
func (s *UserService) Edit(u *models.User, roleIDs, postIDs []uint, operator string) error {
	u.UpdatedBy = operator
	return s.DB.Transaction(func(tx *gorm.DB) error {
		updates := map[string]any{
			"dept_id":     u.DeptID,
			"nick_name":   u.NickName,
			"email":       u.Email,
			"phonenumber": u.Phonenumber,
			"sex":         u.Sex,
			"status":      u.Status,
			"remark":      u.Remark,
			"updated_by":  u.UpdatedBy,
		}
		if err := tx.Model(&models.User{}).Where("id = ?", u.ID).Updates(updates).Error; err != nil {
			return err
		}
		return s.saveUserRelations(tx, u.ID, roleIDs, postIDs)
	})
}

// Delete 删除用户。
func (s *UserService) Delete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("请选择要删除的数据")
	}
	for _, id := range ids {
		if id == 1 {
			return errors.New("不允许操作超级管理员用户")
		}
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("id IN ?", ids).Update("del_flag", "2").Error; err != nil {
			return err
		}
		if err := tx.Where("user_id IN ?", ids).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}
		return tx.Where("user_id IN ?", ids).Delete(&models.UserPost{}).Error
	})
}

// ResetPwd 重置密码。
func (s *UserService) ResetPwd(userID uint, password, operator string) error {
	if userID == 1 {
		return errors.New("不允许操作超级管理员用户")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	now := time.Now()
	return s.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]any{
		"password":        string(hash),
		"pwd_update_date": now,
		"updated_by":      operator,
	}).Error
}

// ChangeStatus 修改状态。
func (s *UserService) ChangeStatus(userID uint, status, operator string) error {
	if userID == 1 {
		return errors.New("不允许操作超级管理员用户")
	}
	return s.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]any{
		"status":     status,
		"updated_by": operator,
	}).Error
}

// Profile 个人中心资料。
func (s *UserService) Profile(userID uint) (map[string]any, string, string, error) {
	var u models.User
	if err := s.DB.First(&u, userID).Error; err != nil {
		return nil, "", "", err
	}
	roleGroup, err := s.roleGroupName(userID)
	if err != nil {
		return nil, "", "", err
	}
	postGroup, err := s.postGroupName(userID)
	if err != nil {
		return nil, "", "", err
	}
	return UserToMap(&u), roleGroup, postGroup, nil
}

// UpdateProfile 修改个人资料。
func (s *UserService) UpdateProfile(userID uint, nickName, email, phonenumber, sex string) error {
	return s.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]any{
		"nick_name":   nickName,
		"email":       email,
		"phonenumber": phonenumber,
		"sex":         sex,
	}).Error
}

// UpdatePwd 修改密码。
func (s *UserService) UpdatePwd(userID uint, oldPwd, newPwd string) error {
	var u models.User
	if err := s.DB.First(&u, userID).Error; err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(oldPwd)); err != nil {
		return errors.New("修改密码失败，旧密码错误")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	now := time.Now()
	return s.DB.Model(&u).Updates(map[string]any{
		"password":        string(hash),
		"pwd_update_date": now,
	}).Error
}

// UpdateAvatar 更新头像。
func (s *UserService) UpdateAvatar(userID uint, url string) error {
	return s.DB.Model(&models.User{}).Where("id = ?", userID).Update("avatar", url).Error
}

// AuthRole 查询用户可分配角色。
func (s *UserService) AuthRole(userID uint) (map[string]any, []map[string]any, error) {
	var u models.User
	if err := s.DB.Where("id = ? AND del_flag = ?", userID, "0").First(&u).Error; err != nil {
		return nil, nil, err
	}
	roleIDs, err := s.userRoleIDs(userID)
	if err != nil {
		return nil, nil, err
	}
	selected := map[uint]struct{}{}
	for _, id := range roleIDs {
		selected[id] = struct{}{}
	}
	var roles []models.Role
	if err := s.DB.Where("del_flag = ? AND status = ?", "0", "0").Order("role_sort asc").Find(&roles).Error; err != nil {
		return nil, nil, err
	}
	rows := make([]map[string]any, 0, len(roles))
	for i := range roles {
		_, ok := selected[roles[i].ID]
		rows = append(rows, map[string]any{
			"roleId":   roles[i].ID,
			"roleName": roles[i].RoleName,
			"flag":     ok,
		})
	}
	return map[string]any{
		"userId":   u.ID,
		"userName": u.UserName,
	}, rows, nil
}

// SaveAuthRole 保存用户角色。
func (s *UserService) SaveAuthRole(userID uint, roleIDs []uint) error {
	if userID == 1 {
		return errors.New("不允许操作超级管理员用户")
	}
	return s.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}
		for _, roleID := range roleIDs {
			if err := tx.Create(&models.UserRole{UserID: userID, RoleID: roleID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// DeptTree 部门树。
func (s *UserService) DeptTree() ([]map[string]any, error) {
	return s.dept.TreeSelect()
}

// Export 导出用户。
func (s *UserService) Export(ctx *gin.Context) ([][]string, error) {
	_, rows, err := s.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([][]string, 0, len(rows))
	for _, row := range rows {
		deptName := ""
		if dept, ok := row["dept"].(map[string]any); ok {
			deptName = toString(dept["deptName"])
		}
		out = append(out, []string{
			toString(row["userId"]),
			toString(row["userName"]),
			toString(row["nickName"]),
			deptName,
			toString(row["phonenumber"]),
			toString(row["status"]),
			toString(row["createdAt"]),
		})
	}
	return out, nil
}

// ImportTemplate 导入模板。
func (s *UserService) ImportTemplate() [][]string {
	return [][]string{
		{"登录名称", "用户名称", "部门编号", "手机号码", "邮箱", "用户性别", "账号状态", "备注"},
		{"ry", "若依", "103", "15888888888", "ry@qq.com", "0", "0", "测试用户"},
	}
}

// ImportData 导入用户。
func (s *UserService) ImportData(file multipart.File, updateSupport bool, operator string) (string, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return "", err
	}
	if len(records) <= 1 {
		return "导入数据为空", nil
	}
	success, failure := 0, 0
	var msgs []string
	for i, row := range records[1:] {
		if len(row) < 2 {
			failure++
			msgs = append(msgs, fmt.Sprintf("第%d行格式错误", i+2))
			continue
		}
		userName := strings.TrimSpace(row[0])
		nickName := strings.TrimSpace(row[1])
		if userName == "" {
			failure++
			continue
		}
		var exist models.User
		err := s.DB.Where("user_name = ? AND del_flag = ?", userName, "0").First(&exist).Error
		if err == nil {
			if !updateSupport {
				failure++
				msgs = append(msgs, fmt.Sprintf("账号 %s 已存在", userName))
				continue
			}
			exist.NickName = nickName
			if len(row) > 3 {
				exist.Phonenumber = strings.TrimSpace(row[3])
			}
			if len(row) > 4 {
				exist.Email = strings.TrimSpace(row[4])
			}
			if len(row) > 5 {
				exist.Sex = strings.TrimSpace(row[5])
			}
			if len(row) > 6 {
				exist.Status = strings.TrimSpace(row[6])
			}
			if len(row) > 7 {
				exist.Remark = strings.TrimSpace(row[7])
			}
			exist.UpdatedBy = operator
			if err := s.DB.Save(&exist).Error; err != nil {
				failure++
				continue
			}
			success++
			continue
		}
		u := &models.User{
			UserName: userName,
			NickName: nickName,
			Status:   "0",
			DelFlag:  "0",
		}
		if len(row) > 2 {
			if deptID, e := strconv.ParseUint(strings.TrimSpace(row[2]), 10, 64); e == nil && deptID > 0 {
				id := uint(deptID)
				u.DeptID = &id
			}
		}
		if len(row) > 3 {
			u.Phonenumber = strings.TrimSpace(row[3])
		}
		if len(row) > 4 {
			u.Email = strings.TrimSpace(row[4])
		}
		if len(row) > 5 {
			u.Sex = strings.TrimSpace(row[5])
		}
		if len(row) > 6 && strings.TrimSpace(row[6]) != "" {
			u.Status = strings.TrimSpace(row[6])
		}
		if len(row) > 7 {
			u.Remark = strings.TrimSpace(row[7])
		}
		pwd, _ := s.defaultPassword()
		u.Password = pwd
		if err := s.Add(u, nil, nil, operator); err != nil {
			failure++
			msgs = append(msgs, err.Error())
			continue
		}
		success++
	}
	msg := fmt.Sprintf("导入完成，成功 %d 条，失败 %d 条", success, failure)
	if len(msgs) > 0 {
		msg += "；" + strings.Join(msgs, "；")
	}
	return msg, nil
}

func (s *UserService) userDetailMap(u *models.User, withRelations bool) (map[string]any, error) {
	m := UserToMap(u)
	// 列表页需要展示部门名称，部门信息始终附带
	if u.DeptID != nil && *u.DeptID > 0 {
		var dept models.Dept
		if err := s.DB.Where("id = ? AND del_flag = ?", *u.DeptID, "0").First(&dept).Error; err == nil {
			m["dept"] = DeptToMap(&dept)
		}
	}
	if withRelations {
		roleIDs, err := s.userRoleIDs(u.ID)
		if err != nil {
			return nil, err
		}
		if len(roleIDs) > 0 {
			var roles []models.Role
			if err := s.DB.Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
				return nil, err
			}
			roleRows := make([]map[string]any, 0, len(roles))
			for i := range roles {
				roleRows = append(roleRows, RoleToMap(&roles[i]))
			}
			m["roles"] = roleRows
		}
	}
	return m, nil
}

func (s *UserService) saveUserRelations(tx *gorm.DB, userID uint, roleIDs, postIDs []uint) error {
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
		return err
	}
	for _, roleID := range roleIDs {
		if err := tx.Create(&models.UserRole{UserID: userID, RoleID: roleID}).Error; err != nil {
			return err
		}
	}
	if err := tx.Where("user_id = ?", userID).Delete(&models.UserPost{}).Error; err != nil {
		return err
	}
	for _, postID := range postIDs {
		if err := tx.Create(&models.UserPost{UserID: userID, PostID: postID}).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) userRoleIDs(userID uint) ([]uint, error) {
	var ids []uint
	err := s.DB.Model(&models.UserRole{}).Where("user_id = ?", userID).Pluck("role_id", &ids).Error
	if ids == nil {
		ids = []uint{}
	}
	return ids, err
}

func (s *UserService) userPostIDs(userID uint) ([]uint, error) {
	var ids []uint
	err := s.DB.Model(&models.UserPost{}).Where("user_id = ?", userID).Pluck("post_id", &ids).Error
	if ids == nil {
		ids = []uint{}
	}
	return ids, err
}

func (s *UserService) listRoleOptions() ([]map[string]any, error) {
	var roles []models.Role
	if err := s.DB.Where("del_flag = ? AND status = ?", "0", "0").Order("role_sort asc").Find(&roles).Error; err != nil {
		return nil, err
	}
	rows := make([]map[string]any, 0, len(roles))
	for i := range roles {
		rows = append(rows, RoleToMap(&roles[i]))
	}
	return rows, nil
}

func (s *UserService) roleGroupName(userID uint) (string, error) {
	var names []string
	err := s.DB.Table("sys_role r").
		Joins("JOIN sys_user_role ur ON ur.role_id = r.id").
		Where("ur.user_id = ? AND r.del_flag = ?", userID, "0").
		Pluck("r.role_name", &names).Error
	return strings.Join(names, ","), err
}

func (s *UserService) postGroupName(userID uint) (string, error) {
	var names []string
	err := s.DB.Table("sys_post p").
		Joins("JOIN sys_user_post up ON up.post_id = p.id").
		Where("up.user_id = ?", userID).
		Pluck("p.post_name", &names).Error
	return strings.Join(names, ","), err
}

func (s *UserService) defaultPassword() (string, error) {
	var cfg models.Config
	if err := s.DB.Where("config_key = ?", "sys.user.initPassword").First(&cfg).Error; err != nil {
		return "123456", nil
	}
	if cfg.ConfigValue == "" {
		return "123456", nil
	}
	return cfg.ConfigValue, nil
}
