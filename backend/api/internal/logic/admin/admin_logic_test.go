package admin

import (
	"context"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AdminLogicTestSuite struct {
	suite.Suite
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (suite *AdminLogicTestSuite) SetupTest() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.User{}, &model.Role{}, &model.UserRole{}, &model.UserBrand{})

	suite.ctx = context.Background()
	suite.svcCtx = &svc.ServiceContext{
		DB: db,
	}
}

// ==================== GetUsers 测试 ====================

func (suite *AdminLogicTestSuite) TestGetUsers() {
	users := []model.User{
		{Id: 1, Username: "admin", Phone: "13800138000", Status: "active"},
		{Id: 2, Username: "user1", Phone: "13800138001", Status: "active"},
		{Id: 3, Username: "user2", Phone: "13800138002", Status: "disabled"},
	}
	for _, user := range users {
		suite.svcCtx.DB.Create(&user)
	}

	logic := NewGetUsersLogic(suite.ctx, suite.svcCtx)

	tests := []struct {
		name     string
		req      *types.AdminGetUsersReq
		validate func(resp *types.AdminUserListResp)
	}{
		{
			name: "查询所有用户",
			req: &types.AdminGetUsersReq{
				Page:     1,
				PageSize: 10,
			},
			validate: func(resp *types.AdminUserListResp) {
				assert.Equal(suite.T(), int64(3), resp.Total)
				assert.Len(suite.T(), resp.Users, 3)
			},
		},
		{
			name: "按状态筛选",
			req: &types.AdminGetUsersReq{
				Page:     1,
				PageSize: 10,
				Status:   "active",
			},
			validate: func(resp *types.AdminUserListResp) {
				assert.Equal(suite.T(), int64(2), resp.Total)
			},
		},
		{
			name: "关键词搜索",
			req: &types.AdminGetUsersReq{
				Page:     1,
				PageSize: 10,
				Keyword:  "admin",
			},
			validate: func(resp *types.AdminUserListResp) {
				assert.Equal(suite.T(), int64(1), resp.Total)
			},
		},
		{
			name: "分页查询",
			req: &types.AdminGetUsersReq{
				Page:     1,
				PageSize: 2,
			},
			validate: func(resp *types.AdminUserListResp) {
				assert.Equal(suite.T(), int64(3), resp.Total)
				assert.Len(suite.T(), resp.Users, 2)
			},
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			resp, err := logic.GetUsers(tt.req)
			assert.NoError(t, err)
			tt.validate(resp)
		})
	}
}

// ==================== CreateUser 测试 ====================

func (suite *AdminLogicTestSuite) TestCreateUser_Success() {
	logic := NewCreateUserLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminCreateUserReq{
		Username: "newuser",
		Password: "password123",
		Phone:    "13900139000",
		Email:    "newuser@test.com",
		RealName: "测试用户",
	}

	resp, err := logic.CreateUser(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "newuser", resp.Username)
	assert.Equal(suite.T(), "13900139000", resp.Phone)
	assert.Equal(suite.T(), "newuser@test.com", resp.Email)
	assert.Equal(suite.T(), "测试用户", resp.RealName)
	assert.Equal(suite.T(), "active", resp.Status)

	// 验证密码已加密
	var user model.User
	suite.svcCtx.DB.Where("username = ?", "newuser").First(&user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	assert.NoError(suite.T(), err)
}

func (suite *AdminLogicTestSuite) TestCreateUser_WithBrandIds() {
	brand1 := model.Role{ID: 100, Name: "品牌1", Code: "brand_1"}
	brand2 := model.Role{ID: 101, Name: "品牌2", Code: "brand_2"}
	suite.svcCtx.DB.Create(&brand1)
	suite.svcCtx.DB.Create(&brand2)

	logic := NewCreateUserLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminCreateUserReq{
		Username: "branduser",
		Password: "password123",
		Phone:    "13900139001",
		BrandIds: []int64{100, 101},
	}

	resp, err := logic.CreateUser(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)

	// 验证品牌关系已创建
	var userRoles []model.UserRole
	suite.svcCtx.DB.Where("user_id = ?", resp.Id).Find(&userRoles)
	assert.Len(suite.T(), userRoles, 2)
}

func (suite *AdminLogicTestSuite) TestCreateUser_EmptyUsername() {
	logic := NewCreateUserLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminCreateUserReq{
		Username: "",
		Password: "password123",
		Phone:    "13900139002",
	}

	resp, err := logic.CreateUser(req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "用户名不能为空", err.Error())
}

func (suite *AdminLogicTestSuite) TestCreateUser_EmptyPassword() {
	logic := NewCreateUserLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminCreateUserReq{
		Username: "testuser",
		Password: "",
		Phone:    "13900139003",
	}

	resp, err := logic.CreateUser(req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "密码不能为空", err.Error())
}

func (suite *AdminLogicTestSuite) TestCreateUser_EmptyPhone() {
	logic := NewCreateUserLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminCreateUserReq{
		Username: "testuser",
		Password: "password123",
		Phone:    "",
	}

	resp, err := logic.CreateUser(req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "手机号不能为空", err.Error())
}

func (suite *AdminLogicTestSuite) TestCreateUser_DuplicateUsername() {
	// 先创建一个用户
	existingUser := model.User{
		Username: "existinguser",
		Password: "hashedpassword",
		Phone:    "13900139004",
		Status:   "active",
	}
	suite.svcCtx.DB.Create(&existingUser)

	logic := NewCreateUserLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminCreateUserReq{
		Username: "existinguser",
		Password: "password123",
		Phone:    "13900139005",
	}

	resp, err := logic.CreateUser(req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "用户名或手机号已存在", err.Error())
}

func (suite *AdminLogicTestSuite) TestCreateUser_DuplicatePhone() {
	// 先创建一个用户
	existingUser := model.User{
		Username: "user1",
		Password: "hashedpassword",
		Phone:    "13900139006",
		Status:   "active",
	}
	suite.svcCtx.DB.Create(&existingUser)

	logic := NewCreateUserLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminCreateUserReq{
		Username: "newuser",
		Password: "password123",
		Phone:    "13900139006",
	}

	resp, err := logic.CreateUser(req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "用户名或手机号已存在", err.Error())
}

// ==================== UpdateUser 测试 ====================

func (suite *AdminLogicTestSuite) TestUpdateUser_Success() {
	// 创建测试用户
	user := model.User{
		Username: "updateuser",
		Password: "hashedpassword",
		Phone:    "13900139010",
		Email:    "old@test.com",
		RealName: "旧名字",
		Status:   "active",
		Role:     "participant",
	}
	suite.svcCtx.DB.Create(&user)

	logic := NewUpdateUserLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminUpdateUserReq{
		RealName: "新名字",
		Email:    "new@test.com",
		Status:   "disabled",
		Role:     "brand_admin",
	}

	resp, err := logic.UpdateUser(user.Id, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "新名字", resp.RealName)
	assert.Equal(suite.T(), "new@test.com", resp.Email)
	assert.Equal(suite.T(), "disabled", resp.Status)
}

func (suite *AdminLogicTestSuite) TestUpdateUser_NoUserIdInContext() {
	logic := NewUpdateUserLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminUpdateUserReq{
		RealName: "新名字",
	}

	resp, err := logic.UpdateUser(0, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "用户ID无效", err.Error())
}

func (suite *AdminLogicTestSuite) TestUpdateUser_UserNotFound() {
	logic := NewUpdateUserLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminUpdateUserReq{
		RealName: "新名字",
	}

	resp, err := logic.UpdateUser(99999, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "用户不存在", err.Error())
}

func (suite *AdminLogicTestSuite) TestUpdateUser_PartialUpdate() {
	// 创建测试用户
	user := model.User{
		Username: "partialuser",
		Password: "hashedpassword",
		Phone:    "13900139011",
		Email:    "partial@test.com",
		RealName: "原始名字",
		Status:   "active",
	}
	suite.svcCtx.DB.Create(&user)

	logic := NewUpdateUserLogic(suite.ctx, suite.svcCtx)

	// 只更新邮箱
	req := &types.AdminUpdateUserReq{
		Email: "updated@test.com",
	}

	resp, err := logic.UpdateUser(user.Id, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "updated@test.com", resp.Email)
	assert.Equal(suite.T(), "原始名字", resp.RealName) // 名字应该保持不变
}

// ==================== DeleteUser 测试 ====================

func (suite *AdminLogicTestSuite) TestDeleteUser_Success() {
	// 创建测试用户
	user := model.User{
		Username: "deleteuser",
		Password: "hashedpassword",
		Phone:    "13900139020",
		Status:   "active",
	}
	suite.svcCtx.DB.Create(&user)

	logic := NewDeleteUserLogic(suite.ctx, suite.svcCtx)

	resp, err := logic.DeleteUser(user.Id)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "删除成功", resp.Message)

	// 验证用户已被删除
	var deletedUser model.User
	result := suite.svcCtx.DB.Where("id = ?", user.Id).First(&deletedUser)
	assert.Error(suite.T(), result.Error)
}

func (suite *AdminLogicTestSuite) TestDeleteUser_NoUserIdInContext() {
	logic := NewDeleteUserLogic(suite.ctx, suite.svcCtx)

	resp, err := logic.DeleteUser(0)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "用户ID无效", err.Error())
}

func (suite *AdminLogicTestSuite) TestDeleteUser_UserNotFound() {
	logic := NewDeleteUserLogic(suite.ctx, suite.svcCtx)

	resp, err := logic.DeleteUser(99999)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "用户不存在", err.Error())
}

// ==================== GetUser 测试 ====================

func (suite *AdminLogicTestSuite) TestGetUser_Success() {
	// 创建测试用户
	user := model.User{
		Username: "getuser",
		Password: "hashedpassword",
		Phone:    "13900139030",
		Email:    "getuser@test.com",
		RealName: "查询用户",
		Status:   "active",
	}
	suite.svcCtx.DB.Create(&user)

	role1 := model.Role{ID: 1, Name: "管理员", Code: "admin"}
	role2 := model.Role{ID: 2, Name: "品牌管理员", Code: "brand_admin"}
	suite.svcCtx.DB.Create(&role1)
	suite.svcCtx.DB.Create(&role2)

	// 创建用户角色关系
	userRole1 := model.UserRole{UserID: user.Id, RoleID: 1}
	userRole2 := model.UserRole{UserID: user.Id, RoleID: 2}
	suite.svcCtx.DB.Create(&userRole1)
	suite.svcCtx.DB.Create(&userRole2)

	ctx := context.WithValue(suite.ctx, "userId", user.Id)
	logic := NewGetUserLogic(ctx, suite.svcCtx)

	resp, err := logic.GetUser()

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), user.Id, resp.Id)
	assert.Equal(suite.T(), "getuser", resp.Username)
	assert.Equal(suite.T(), "查询用户", resp.RealName)
	assert.Contains(suite.T(), resp.Roles, "admin")
	assert.Contains(suite.T(), resp.Roles, "brand_admin")
}

func (suite *AdminLogicTestSuite) TestGetUser_NoUserIdInContext() {
	logic := NewGetUserLogic(suite.ctx, suite.svcCtx)

	resp, err := logic.GetUser()

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "无法从context中获取用户ID", err.Error())
}

func (suite *AdminLogicTestSuite) TestGetUser_UserNotFound() {
	ctx := context.WithValue(suite.ctx, "userId", int64(99999))
	logic := NewGetUserLogic(ctx, suite.svcCtx)

	resp, err := logic.GetUser()

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "用户不存在", err.Error())
}

func (suite *AdminLogicTestSuite) TestGetUser_NoRoles() {
	// 创建测试用户
	user := model.User{
		Username: "noroleuser",
		Password: "hashedpassword",
		Phone:    "13900139031",
		Status:   "active",
	}
	suite.svcCtx.DB.Create(&user)

	ctx := context.WithValue(suite.ctx, "userId", user.Id)
	logic := NewGetUserLogic(ctx, suite.svcCtx)

	resp, err := logic.GetUser()

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Empty(suite.T(), resp.Roles)
}

// ==================== ResetUserPassword 测试 ====================

func (suite *AdminLogicTestSuite) TestResetUserPassword_Success() {
	// 创建测试用户
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("oldpassword"), bcrypt.DefaultCost)
	user := model.User{
		Username: "resetuser",
		Password: string(hashedPassword),
		Phone:    "13900139040",
		Status:   "active",
	}
	suite.svcCtx.DB.Create(&user)

	logic := NewResetUserPasswordLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminResetPasswordReq{
		NewPassword: "newpassword123",
	}

	resp, err := logic.ResetUserPassword(user.Id, req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "重置成功", resp.Message)

	// 验证新密码生效
	var updatedUser model.User
	suite.svcCtx.DB.Where("id = ?", user.Id).First(&updatedUser)
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("newpassword123"))
	assert.NoError(suite.T(), err)

	// 验证旧密码失效
	err = bcrypt.CompareHashAndPassword([]byte(updatedUser.Password), []byte("oldpassword"))
	assert.Error(suite.T(), err)
}

func (suite *AdminLogicTestSuite) TestResetUserPassword_NoUserIdInContext() {
	logic := NewResetUserPasswordLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminResetPasswordReq{
		NewPassword: "newpassword123",
	}

	resp, err := logic.ResetUserPassword(0, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "用户ID无效", err.Error())
}

func (suite *AdminLogicTestSuite) TestResetUserPassword_UserNotFound() {
	logic := NewResetUserPasswordLogic(suite.ctx, suite.svcCtx)

	req := &types.AdminResetPasswordReq{
		NewPassword: "newpassword123",
	}

	resp, err := logic.ResetUserPassword(99999, req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "用户不存在", err.Error())
}

// ==================== ManageBrandAdminRelation 测试 ====================

func (suite *AdminLogicTestSuite) TestManageBrandAdminRelation_Success() {
	// 创建测试用户
	user := model.User{
		Username: "brandadmin",
		Password: "hashedpassword",
		Phone:    "13900139050",
		Status:   "active",
	}
	suite.svcCtx.DB.Create(&user)

	brand1 := model.Role{ID: 100, Name: "品牌1", Code: "brand_1"}
	brand2 := model.Role{ID: 101, Name: "品牌2", Code: "brand_2"}
	suite.svcCtx.DB.Create(&brand1)
	suite.svcCtx.DB.Create(&brand2)

	logic := NewManageBrandAdminRelationLogic(suite.ctx, suite.svcCtx)

	req := &types.BrandAdminRelationReq{
		UserId:   user.Id,
		BrandIds: []int64{100, 101},
	}

	resp, err := logic.ManageBrandAdminRelation(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), "操作成功", resp.Message)

	// 验证关系已创建
	var userRoles []model.UserRole
	suite.svcCtx.DB.Where("user_id = ?", user.Id).Find(&userRoles)
	assert.Len(suite.T(), userRoles, 2)
}

func (suite *AdminLogicTestSuite) TestManageBrandAdminRelation_ReplaceExisting() {
	// 创建测试用户
	user := model.User{
		Username: "replaceadmin",
		Password: "hashedpassword",
		Phone:    "13900139051",
		Status:   "active",
	}
	suite.svcCtx.DB.Create(&user)

	oldRole1 := model.Role{ID: 200, Name: "旧品牌1", Code: "old_brand_1"}
	oldRole2 := model.Role{ID: 201, Name: "旧品牌2", Code: "old_brand_2"}
	suite.svcCtx.DB.Create(&oldRole1)
	suite.svcCtx.DB.Create(&oldRole2)

	oldUserRole1 := model.UserRole{UserID: user.Id, RoleID: 200}
	oldUserRole2 := model.UserRole{UserID: user.Id, RoleID: 201}
	suite.svcCtx.DB.Create(&oldUserRole1)
	suite.svcCtx.DB.Create(&oldUserRole2)

	newBrand := model.Role{ID: 300, Name: "新品牌", Code: "new_brand"}
	suite.svcCtx.DB.Create(&newBrand)

	logic := NewManageBrandAdminRelationLogic(suite.ctx, suite.svcCtx)

	req := &types.BrandAdminRelationReq{
		UserId:   user.Id,
		BrandIds: []int64{300},
	}

	resp, err := logic.ManageBrandAdminRelation(req)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), resp)

	// 验证旧关系已删除，新关系已创建
	var userRoles []model.UserRole
	suite.svcCtx.DB.Where("user_id = ?", user.Id).Find(&userRoles)
	assert.Len(suite.T(), userRoles, 1)
	assert.Equal(suite.T(), int64(300), userRoles[0].RoleID)
}

func (suite *AdminLogicTestSuite) TestManageBrandAdminRelation_InvalidUserId() {
	logic := NewManageBrandAdminRelationLogic(suite.ctx, suite.svcCtx)

	req := &types.BrandAdminRelationReq{
		UserId:   0,
		BrandIds: []int64{100},
	}

	resp, err := logic.ManageBrandAdminRelation(req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "用户ID无效", err.Error())
}

func (suite *AdminLogicTestSuite) TestManageBrandAdminRelation_EmptyBrandIds() {
	logic := NewManageBrandAdminRelationLogic(suite.ctx, suite.svcCtx)

	req := &types.BrandAdminRelationReq{
		UserId:   1,
		BrandIds: []int64{},
	}

	resp, err := logic.ManageBrandAdminRelation(req)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), resp)
	assert.Equal(suite.T(), "品牌ID列表不能为空", err.Error())
}

func TestAdminLogicTestSuite(t *testing.T) {
	suite.Run(t, new(AdminLogicTestSuite))
}
