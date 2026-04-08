package handler

import (
	"crypto/md5"
	"fmt"
	"go-source/internal/database"
	"go-source/internal/model"
	"go-source/internal/response"
	jwtpkg "go-source/pkg/jwt"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginReq struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	var admin model.Admin
	if err := database.DB.Where("username = ? AND status = 'normal'", req.UserName).First(&admin).Error; err != nil {
		response.Fail(c, "账号不存在或已禁用")
		return
	}

	// md5(md5(password) + salt)  与PHP一致
	inner := md5.Sum([]byte(req.Password))
	innerStr := fmt.Sprintf("%x", inner)
	hash := md5.Sum([]byte(innerStr + admin.Salt))
	pwd := fmt.Sprintf("%x", hash)
	if pwd != admin.Password {
		response.Fail(c, "密码错误")
		return
	}

	token, err := jwtpkg.GenerateToken(admin.ID, admin.Username)
	if err != nil {
		response.Fail(c, "生成token失败")
		return
	}

	// 更新登录时间
	database.DB.Model(&admin).Updates(map[string]interface{}{
		"logintime": time.Now().Unix(),
		"loginip":   c.ClientIP(),
		"token":     token,
	})

	response.OK(c, gin.H{
		"token":        token,
		"refreshToken": token,
	})
}

func GetUserInfo(c *gin.Context) {
	adminId, _ := c.Get("adminId")
	username, _ := c.Get("username")

	var admin model.Admin
	if err := database.DB.First(&admin, adminId).Error; err != nil {
		response.Unauthorized(c, "用户不存在")
		return
	}

	response.OK(c, gin.H{
		"userId":   admin.ID,
		"userName": username,
		"email":    admin.Email,
		"avatar":   admin.Avatar,
		"roles":    []string{"R_SUPER"},
		"buttons":  []string{"add", "edit", "delete"},
		"nickName": admin.Nickname,
	})
}

type ChangePasswordReq struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

func ChangePassword(c *gin.Context) {
	adminId, _ := c.Get("adminId")
	var req ChangePasswordReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	var admin model.Admin
	if err := database.DB.First(&admin, adminId).Error; err != nil {
		response.Fail(c, "用户不存在")
		return
	}

	// 验证旧密码
	inner := md5.Sum([]byte(req.OldPassword))
	innerStr := fmt.Sprintf("%x", inner)
	hash := md5.Sum([]byte(innerStr + admin.Salt))
	if fmt.Sprintf("%x", hash) != admin.Password {
		response.Fail(c, "当前密码错误")
		return
	}

	// 生成新密码
	newInner := md5.Sum([]byte(req.NewPassword))
	newInnerStr := fmt.Sprintf("%x", newInner)
	newHash := md5.Sum([]byte(newInnerStr + admin.Salt))
	newPwd := fmt.Sprintf("%x", newHash)

	database.DB.Model(&admin).Update("password", newPwd)
	response.OKMsg(c, "密码修改成功", nil)
}

type UpdateProfileReq struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

func UpdateProfile(c *gin.Context) {
	adminId, _ := c.Get("adminId")
	var req UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	database.DB.Model(&model.Admin{}).Where("id = ?", adminId).Updates(updates)
	response.OKMsg(c, "保存成功", nil)
}
