package auth

import (
	"context"
	"testing"

	"dmh/api/internal/config"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestRegisterLogic_Register_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	participantRole := &model.Role{Code: "participant", Name: "参与者"}
	db.Create(participantRole)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{
		DB: db,
		Config: config.Config{
			Auth: struct {
				AccessSecret string
				AccessExpire int64
			}{
				AccessSecret: "test-secret-key-for-testing",
				AccessExpire: 86400,
			},
		},
	}
	logic := NewRegisterLogic(ctx, svcCtx)

	req := &types.RegisterReq{
		Username: "newuser",
		Password: "password123",
		Phone:    "13900139000",
		Email:    "test@example.com",
		RealName: "测试用户",
	}

	resp, err := logic.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "newuser", resp.Username)
	assert.Equal(t, "participant", resp.Roles[0])
}

func TestRegisterLogic_UsernameAlreadyExists(t *testing.T) {
	db := setupAuthTestDB(t)
	createTestUserWithPassword(t, db, "existinguser", "password123", "active")
	participantRole := &model.Role{Code: "participant", Name: "参与者"}
	db.Create(participantRole)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{
		DB: db,
		Config: config.Config{
			Auth: struct {
				AccessSecret string
				AccessExpire int64
			}{
				AccessSecret: "test-secret-key",
				AccessExpire: 86400,
			},
		},
	}
	logic := NewRegisterLogic(ctx, svcCtx)

	req := &types.RegisterReq{
		Username: "existinguser",
		Password: "password123",
		Phone:    "13900139001",
	}

	resp, err := logic.Register(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户名已存在")
}

func TestRegisterLogic_PhoneAlreadyExists(t *testing.T) {
	db := setupAuthTestDB(t)
	createTestUserWithPassword(t, db, "newuser1", "password123", "active")
	participantRole := &model.Role{Code: "participant", Name: "参与者"}
	db.Create(participantRole)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{
		DB: db,
		Config: config.Config{
			Auth: struct {
				AccessSecret string
				AccessExpire int64
			}{
				AccessSecret: "test-secret-key",
				AccessExpire: 86400,
			},
		},
	}
	logic := NewRegisterLogic(ctx, svcCtx)

	req := &types.RegisterReq{
		Username: "newuser2",
		Password: "password123",
		Phone:    "13800138000",
	}

	resp, err := logic.Register(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "手机号已被注册")
}

func TestRegisterLogic_EmptyUsername(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewRegisterLogic(ctx, svcCtx)

	req := &types.RegisterReq{
		Username: "",
		Password: "password123",
		Phone:    "13900139002",
	}

	resp, err := logic.Register(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户名不能为空")
}

func TestRegisterLogic_EmptyPassword(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewRegisterLogic(ctx, svcCtx)

	req := &types.RegisterReq{
		Username: "testuser",
		Password: "",
		Phone:    "13900139003",
	}

	resp, err := logic.Register(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "密码不能为空")
}

func TestRegisterLogic_EmptyPhone(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewRegisterLogic(ctx, svcCtx)

	req := &types.RegisterReq{
		Username: "testuser",
		Password: "password123",
		Phone:    "",
	}

	resp, err := logic.Register(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "手机号不能为空")
}

func TestRegisterLogic_PasswordTooShort(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewRegisterLogic(ctx, svcCtx)

	req := &types.RegisterReq{
		Username: "testuser",
		Password: "12345",
		Phone:    "13900139004",
	}

	_, err := logic.Register(req)

	assert.Contains(t, err.Error(), "密码长度不能少于6位")
}

func TestRegisterLogic_CreateRoleIfNotExist(t *testing.T) {
	db := setupAuthTestDB(t)
	// 注意：不预先创建 participant 角色

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{
		DB: db,
		Config: config.Config{
			Auth: struct {
				AccessSecret string
				AccessExpire int64
			}{
				AccessSecret: "test-secret-key-for-testing",
				AccessExpire: 86400,
			},
		},
	}
	logic := NewRegisterLogic(ctx, svcCtx)

	req := &types.RegisterReq{
		Username: "newuser_norole",
		Password: "password123",
		Phone:    "13900139999",
		Email:    "test_norole@example.com",
		RealName: "测试用户无角色",
	}

	resp, err := logic.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)

	// 验证角色被自动创建
	var role model.Role
	err = db.Where("code = ?", "participant").First(&role).Error
	assert.NoError(t, err)
	assert.Equal(t, "participant", role.Code)
}

func TestChangePasswordLogic_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	user := createTestUserWithPassword(t, db, "testuser", "oldpassword123", "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewChangePasswordLogic(ctx, svcCtx)

	req := &types.ChangePasswordReq{
		OldPassword: "oldpassword123",
		NewPassword: "newpassword456",
	}

	resp, err := logic.ChangePassword(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, resp.Message, "成功")

	var updatedUser model.User
	db.First(&updatedUser, user.Id)
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("newpassword456"))
	assert.NoError(t, err)
}

func TestChangePasswordLogic_WrongOldPassword(t *testing.T) {
	db := setupAuthTestDB(t)
	user := createTestUserWithPassword(t, db, "testuser", "oldpassword123", "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewChangePasswordLogic(ctx, svcCtx)

	req := &types.ChangePasswordReq{
		OldPassword: "wrongpassword",
		NewPassword: "newpassword456",
	}

	resp, err := logic.ChangePassword(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "旧密码错误")
}

func TestChangePasswordLogic_NewPasswordTooShort(t *testing.T) {
	db := setupAuthTestDB(t)
	user := createTestUserWithPassword(t, db, "testuser", "oldpassword123", "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewChangePasswordLogic(ctx, svcCtx)

	req := &types.ChangePasswordReq{
		OldPassword: "oldpassword123",
		NewPassword: "12345",
	}

	resp, err := logic.ChangePassword(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "新密码长度不能少于6位")
}

func TestChangePasswordLogic_NotLoggedIn(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewChangePasswordLogic(ctx, svcCtx)

	req := &types.ChangePasswordReq{
		OldPassword: "oldpassword123",
		NewPassword: "newpassword456",
	}

	resp, err := logic.ChangePassword(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "无法从context中获取用户ID")
}

func TestLogoutLogic_Success(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}
	logic := NewLogoutLogic(ctx, svcCtx)

	resp, err := logic.Logout(&types.CommonResp{})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "登出成功", resp.Message)
}

func TestRefreshTokenLogic_EmptyToken(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{
		Config: config.Config{
			Auth: struct {
				AccessSecret string
				AccessExpire int64
			}{
				AccessSecret: "test-secret-key-for-testing",
				AccessExpire: 86400,
			},
		},
	}
	logic := NewRefreshTokenLogic(ctx, svcCtx)

	req := &types.RefreshTokenReq{
		Token: "",
	}

	resp, err := logic.RefreshToken(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "token不能为空")
}

func TestBindPhoneLogic_EmptyPhone(t *testing.T) {
	db := setupAuthTestDB(t)
	user := createTestUserWithPassword(t, db, "testuser", "password123", "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewBindPhoneLogic(ctx, svcCtx)

	req := &types.BindPhoneReq{
		Phone: "",
		Code:  "123456",
	}

	resp, err := logic.BindPhone(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "手机号不能为空")
}

func TestBindPhoneLogic_UserNotExists(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.WithValue(context.Background(), "userId", int64(99999))
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewBindPhoneLogic(ctx, svcCtx)

	req := &types.BindPhoneReq{
		Phone: "13912345678",
		Code:  "123456",
	}

	resp, err := logic.BindPhone(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")
}

func TestBindEmailLogic_EmptyEmail(t *testing.T) {
	db := setupAuthTestDB(t)
	user := createTestUserWithPassword(t, db, "testuser", "password123", "active")

	ctx := context.WithValue(context.Background(), "userId", user.Id)
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewBindEmailLogic(ctx, svcCtx)

	req := &types.BindEmailReq{
		Email: "",
		Code:  "123456",
	}

	resp, err := logic.BindEmail(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "邮箱不能为空")
}

func TestBindEmailLogic_UserNotExists(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.WithValue(context.Background(), "userId", int64(99999))
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewBindEmailLogic(ctx, svcCtx)

	req := &types.BindEmailReq{
		Email: "test@example.com",
		Code:  "123456",
	}

	resp, err := logic.BindEmail(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")
}

func TestSendPhoneCodeLogic_EmptyTarget(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}
	logic := NewSendPhoneCodeLogic(ctx, svcCtx)

	req := &types.SendCodeReq{
		Target: "",
		Type:   "phone",
	}

	resp, err := logic.SendPhoneCode(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "手机号不能为空")
}

func TestSendEmailCodeLogic_EmptyTarget(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}
	logic := NewSendEmailCodeLogic(ctx, svcCtx)

	req := &types.SendCodeReq{
		Target: "",
		Type:   "email",
	}

	resp, err := logic.SendEmailCode(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "邮箱不能为空")
}

func TestUpdateProfileLogic_UserNotExists(t *testing.T) {
	db := setupAuthTestDB(t)
	ctx := context.WithValue(context.Background(), "userId", int64(99999))
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateProfileLogic(ctx, svcCtx)

	req := &types.UpdateProfileReq{
		RealName: "测试姓名",
	}

	resp, err := logic.UpdateProfile(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")
}

func TestChangePasswordLogic_UserNotFound(t *testing.T) {
	db := setupAuthTestDB(t)

	// 使用一个不存在于数据库中的 userId
	ctx := context.WithValue(context.Background(), "userId", int64(99999))
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewChangePasswordLogic(ctx, svcCtx)

	req := &types.ChangePasswordReq{
		OldPassword: "oldpassword123",
		NewPassword: "newpassword456",
	}

	resp, err := logic.ChangePassword(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户不存在")
}
