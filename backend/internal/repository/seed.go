package repository

import (
	"log"

	"github.com/taskify/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// SeedData 填充种子数据（预定义角色和权限）
func SeedData() error {
	// 填充预定义权限
	for _, perm := range models.SystemPermissions {
		var existing models.Permission
		result := DB.Where("resource = ? AND action = ?", perm.Resource, perm.Action).First(&existing)
		if result.RowsAffected == 0 {
			if err := DB.Create(&perm).Error; err != nil {
				return err
			}
		}
	}
	log.Println("权限数据填充成功")

	// 填充预定义角色
	for _, role := range models.SystemRoles {
		var existing models.Role
		result := DB.Where("name = ?", role.Name).First(&existing)
		if result.RowsAffected == 0 {
			if err := DB.Create(&role).Error; err != nil {
				return err
			}
		}
	}
	log.Println("角色数据填充成功")

	// 为预定义角色分配权限
	if err := seedRolePermissions(); err != nil {
		return err
	}

	return nil
}

// seedRolePermissions 为预定义角色分配默认权限
func seedRolePermissions() error {
	// 管理员拥有所有权限
	adminRole := models.Role{}
	if err := DB.Where("name = ?", models.RoleAdmin).First(&adminRole).Error; err != nil {
		return err
	}

	// 获取所有权限
	var allPermissions []models.Permission
	DB.Find(&allPermissions)

	// 为管理员分配所有权限
	for _, perm := range allPermissions {
		var existing models.RolePermission
		result := DB.Where("role_id = ? AND permission_id = ?", adminRole.ID, perm.ID).First(&existing)
		if result.RowsAffected == 0 {
			rp := models.RolePermission{
				RoleID:       adminRole.ID,
				PermissionID: perm.ID,
			}
			if err := DB.Create(&rp).Error; err != nil {
				return err
			}
		}
	}

	// 项目所有者拥有项目相关所有权限
	ownerRole := models.Role{}
	if err := DB.Where("name = ?", models.RoleOwner).First(&ownerRole).Error; err != nil {
		return err
	}

	ownerPermissions := []struct {
		Resource string
		Action   string
	}{
		{models.ResourceProjects, models.ActionView},
		{models.ResourceProjects, models.ActionEdit},
		{models.ResourceProjects, models.ActionDelete},
		{models.ResourceProjects, models.ActionManage},
		{models.ResourceTasks, models.ActionView},
		{models.ResourceTasks, models.ActionCreate},
		{models.ResourceTasks, models.ActionEdit},
		{models.ResourceTasks, models.ActionMove},
		{models.ResourceTasks, models.ActionDelete},
		{models.ResourceComments, models.ActionView},
		{models.ResourceComments, models.ActionCreate},
		{models.ResourceComments, models.ActionEdit},
		{models.ResourceComments, models.ActionDelete},
	}

	for _, op := range ownerPermissions {
		var perm models.Permission
		if err := DB.Where("resource = ? AND action = ?", op.Resource, op.Action).First(&perm).Error; err != nil {
			continue
		}
		var existing models.RolePermission
		result := DB.Where("role_id = ? AND permission_id = ?", ownerRole.ID, perm.ID).First(&existing)
		if result.RowsAffected == 0 {
			rp := models.RolePermission{
				RoleID:       ownerRole.ID,
				PermissionID: perm.ID,
			}
			DB.Create(&rp)
		}
	}

	// 项目成员拥有创建和编辑权限
	memberRole := models.Role{}
	if err := DB.Where("name = ?", models.RoleMember).First(&memberRole).Error; err != nil {
		return err
	}

	memberPermissions := []struct {
		Resource string
		Action   string
	}{
		{models.ResourceProjects, models.ActionView},
		{models.ResourceTasks, models.ActionView},
		{models.ResourceTasks, models.ActionCreate},
		{models.ResourceTasks, models.ActionEdit},
		{models.ResourceTasks, models.ActionMove},
		{models.ResourceComments, models.ActionView},
		{models.ResourceComments, models.ActionCreate},
		{models.ResourceComments, models.ActionEdit},
	}

	for _, mp := range memberPermissions {
		var perm models.Permission
		if err := DB.Where("resource = ? AND action = ?", mp.Resource, mp.Action).First(&perm).Error; err != nil {
			continue
		}
		var existing models.RolePermission
		result := DB.Where("role_id = ? AND permission_id = ?", memberRole.ID, perm.ID).First(&existing)
		if result.RowsAffected == 0 {
			rp := models.RolePermission{
				RoleID:       memberRole.ID,
				PermissionID: perm.ID,
			}
			DB.Create(&rp)
		}
	}

	// 项目访客仅有查看权限
	guestRole := models.Role{}
	if err := DB.Where("name = ?", models.RoleGuest).First(&guestRole).Error; err != nil {
		return err
	}

	guestPermissions := []struct {
		Resource string
		Action   string
	}{
		{models.ResourceProjects, models.ActionView},
		{models.ResourceTasks, models.ActionView},
		{models.ResourceComments, models.ActionView},
	}

	for _, gp := range guestPermissions {
		var perm models.Permission
		if err := DB.Where("resource = ? AND action = ?", gp.Resource, gp.Action).First(&perm).Error; err != nil {
			continue
		}
		var existing models.RolePermission
		result := DB.Where("role_id = ? AND permission_id = ?", guestRole.ID, perm.ID).First(&existing)
		if result.RowsAffected == 0 {
			rp := models.RolePermission{
				RoleID:       guestRole.ID,
				PermissionID: perm.ID,
			}
			DB.Create(&rp)
		}
	}

	log.Println("角色权限分配成功")

	// 创建管理员账号（如不存在）
	if err := seedAdminUser(); err != nil {
		return err
	}

	return nil
}

// seedAdminUser 创建默认管理员账号
func seedAdminUser() error {
	adminEmail := "admin@taskify.local"

	var existingUser models.User
	result := DB.Where("email = ?", adminEmail).First(&existingUser)
	if result.RowsAffected > 0 {
		log.Println("管理员账号已存在，跳过创建")
		return nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{
		Email:         adminEmail,
		DisplayName:   "系统管理员",
		PasswordHash:  string(hashedPassword),
		EmailVerified: true,
		IsDisabled:    false,
	}

	if err := DB.Create(&user).Error; err != nil {
		return err
	}

	// 关联 admin 角色（system scope）
	var adminRole models.Role
	if err := DB.Where("name = ? AND scope = ?", models.RoleAdmin, "system").First(&adminRole).Error; err != nil {
		return err
	}

	membership := models.ProjectMembership{
		UserID:   user.ID,
		ProjectID: 0, // system-level role，不属于具体项目
		RoleID:   adminRole.ID,
	}

	if err := DB.Create(&membership).Error; err != nil {
		return err
	}

	log.Printf("管理员账号创建成功: %s / admin123\n", adminEmail)
	return nil
}