package auth

import (
	"context"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/testutil"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupAuthTestDB(t *testing.T) *gorm.DB {
	db, _ := testutil.SetupMySQLTestDB(t)
	return db
}

func createTestUserWithPassword(t *testing.T, db *gorm.DB, username, password, status string) *model.User {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
		Phone:    "13800138000",
		Status:   status,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

func TestLoginLogic_Login_Success(t *testing.T) {
	db := setupAuthTestDB(t)
	createTestUserWithPassword(t, db, "testuser", "password123", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewLoginLogic(ctx, svcCtx)

	req := &types.LoginReq{
		Username: "testuser",
		Password: "password123",
	}

	resp, err := logic.Login(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, "testuser", resp.Username)
}

func TestLoginLogic_InvalidUsername(t *testing.T) {
	db := setupAuthTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewLoginLogic(ctx, svcCtx)

	req := &types.LoginReq{
		Username: "nonexistent",
		Password: "password123",
	}

	resp, err := logic.Login(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户名或密码错误")
}

func TestLoginLogic_InvalidPassword(t *testing.T) {
	db := setupAuthTestDB(t)
	createTestUserWithPassword(t, db, "testuser2", "password123", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewLoginLogic(ctx, svcCtx)

	req := &types.LoginReq{
		Username: "testuser2",
		Password: "wrongpassword",
	}

	resp, err := logic.Login(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "用户名或密码错误")
}

func TestLoginLogic_DisabledUser(t *testing.T) {
	db := setupAuthTestDB(t)
	createTestUserWithPassword(t, db, "disableduser", "password123", "disabled")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewLoginLogic(ctx, svcCtx)

	req := &types.LoginReq{
		Username: "disableduser",
		Password: "password123",
	}

	resp, err := logic.Login(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "账号已被禁用")
}
