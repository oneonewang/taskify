package handlers

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/taskify/backend/internal/middleware"
	"github.com/taskify/backend/internal/models"
	"github.com/taskify/backend/pkg/response"
	"golang.org/x/crypto/bcrypt"
	"github.com/xuri/excelize/v2"
)

const DefaultImportPassword = "admin123"

// ImportResult 导入结果
type ImportResult struct {
	Total    int      `json:"total"`
	Success  int      `json:"success"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors,omitempty"`
}

// ImportUsers 批量导入用户
// POST /api/admin/users/import
func (h *AdminUserHandler) ImportUsers(c *gin.Context) {
	adminID, ok := middleware.GetUserID(c)
	if !ok {
		response.Unauthorized(c, "请先登录")
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请上传Excel文件")
		return
	}

	// 验证文件类型
	if !strings.HasSuffix(file.Filename, ".xlsx") {
		response.BadRequest(c, "只支持 .xlsx 格式的Excel文件")
		return
	}

	result, err := h.parseAndImportUsers(file)
	if err != nil {
		response.InternalError(c, "导入失败: "+err.Error())
		return
	}

	// 记录审计日志
	details := fmt.Sprintf("批量导入用户: 成功%d，失败%d", result.Success, result.Failed)
	logAdminAudit(adminID, models.EventUserImport, details, middleware.GetClientIP(c))

	response.Success(c, result)
}

// GetImportTemplate 下载导入模板
// GET /api/admin/users/import/template
func (h *AdminUserHandler) GetImportTemplate(c *gin.Context) {
	// 创建Excel文件
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "用户导入"
	sheetIndex, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(sheetIndex)

	// 设置表头
	headers := []string{"邮箱", "显示名称"}
	f.SetSheetRow(sheetName, "A1", headers)

	// 添加示例数据
	exampleData := [][]interface{}{
		{"user1@example.com", "张三"},
		{"user2@example.com", "李四"},
	}
	for i, row := range exampleData {
		rowNum := i + 2
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), row[0])
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), row[1])
	}

	// 设置列宽
	f.SetColWidth(sheetName, "A", "A", 30)
	f.SetColWidth(sheetName, "B", "B", 20)
	f.SetColWidth(sheetName, "C", "C", 15)

	// 设置HTTP头
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=用户导入模板.xlsx")

	if err := f.Write(c.Writer); err != nil {
		response.InternalError(c, "生成模板失败")
		return
	}
}

// parseAndImportUsers 解析并导入用户
func (h *AdminUserHandler) parseAndImportUsers(file *multipart.FileHeader) (*ImportResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	f, err := excelize.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows("用户导入")
	if err != nil {
		// 尝试默认sheet名
		rows, err = f.GetRows("Sheet1")
		if err != nil {
			return nil, err
		}
	}

	result := &ImportResult{
		Total:  len(rows) - 1, // 减去表头行
		Errors: []string{},
	}

	if len(rows) < 2 {
		result.Errors = append(result.Errors, "Excel文件为空或格式不正确")
		result.Failed = result.Total
		return result, nil
	}

	// 跳过表头，从第二行开始
	for i, row := range rows[1:] {
		if len(row) < 2 {
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 数据不完整", i+2))
			result.Failed++
			continue
		}

		email := strings.TrimSpace(row[0])
		displayName := strings.TrimSpace(row[1])

		// 验证邮箱格式
		if email == "" || !isValidImportEmail(email) {
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 无效的邮箱地址 '%s'", i+2, email))
			result.Failed++
			continue
		}

		// 检查邮箱是否已存在
		existingUser, _ := h.userRepo.FindByEmail(email)
		if existingUser != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 邮箱 '%s' 已存在，跳过", i+2, email))
			result.Failed++
			continue
		}

		// 哈希密码
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(DefaultImportPassword), bcrypt.DefaultCost)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 密码加密失败", i+2))
			result.Failed++
			continue
		}

		// 创建用户（不分配全局角色，角色通过项目管理页面分配）
		user := &models.User{
			Email:        email,
			DisplayName:  displayName,
			PasswordHash: string(hashedPassword),
			EmailVerified: false,
			IsDisabled:    false,
		}

		if err := h.userRepo.Create(user); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行: 创建用户失败 '%s'", i+2, email))
			result.Failed++
			continue
		}

		result.Success++
	}

	return result, nil
}

// isValidImportEmail 验证邮箱格式（用于导入）
func isValidImportEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	atIndex := strings.Index(email, "@")
	if atIndex < 1 || atIndex > len(email)-3 {
		return false
	}
	domain := email[atIndex+1:]
	if !strings.Contains(domain, ".") {
		return false
	}
	return true
}
