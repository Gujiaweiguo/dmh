package member

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupMemberTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.AutoMigrate(&model.Member{}, &model.MemberProfile{})
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	return db
}

func createTestMember(t *testing.T, db *gorm.DB, unionID, nickname, phone, status string) *model.Member {
	member := &model.Member{
		UnionID:  unionID,
		Nickname: nickname,
		Phone:    phone,
		Status:   status,
		Gender:   1,
		Source:   "wechat",
	}
	if err := db.Create(member).Error; err != nil {
		t.Fatalf("Failed to create test member: %v", err)
	}
	return member
}

func TestGetMembersLogic_GetMembers_Success(t *testing.T) {
	db := setupMemberTestDB(t)

	createTestMember(t, db, "wx001", "张三", "13800138001", "active")
	createTestMember(t, db, "wx002", "李四", "13800138002", "active")
	createTestMember(t, db, "wx003", "王五", "13800138003", "disabled")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMembersLogic(ctx, svcCtx)

	req := &types.GetMembersReq{
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetMembers(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(3), resp.Total)
	assert.Len(t, resp.Members, 3)
}

func TestGetMembersLogic_WithStatusFilter(t *testing.T) {
	db := setupMemberTestDB(t)

	createTestMember(t, db, "wx001", "张三", "13800138001", "active")
	createTestMember(t, db, "wx002", "李四", "13800138002", "active")
	createTestMember(t, db, "wx003", "王五", "13800138003", "disabled")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMembersLogic(ctx, svcCtx)

	req := &types.GetMembersReq{
		Status:   "active",
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetMembers(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(2), resp.Total)
	assert.Len(t, resp.Members, 2)
}

func TestGetMembersLogic_WithKeywordFilter(t *testing.T) {
	db := setupMemberTestDB(t)

	createTestMember(t, db, "wx001", "张三", "13800138001", "active")
	createTestMember(t, db, "wx002", "李四", "13800138002", "active")
	createTestMember(t, db, "wx003", "王五", "13800138003", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMembersLogic(ctx, svcCtx)

	req := &types.GetMembersReq{
		Keyword:  "张三",
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetMembers(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Total)
	assert.Len(t, resp.Members, 1)
	assert.Equal(t, "张三", resp.Members[0].Nickname)
}

func TestGetMembersLogic_WithPhoneFilter(t *testing.T) {
	db := setupMemberTestDB(t)

	createTestMember(t, db, "wx001", "张三", "13800138001", "active")
	createTestMember(t, db, "wx002", "李四", "13800138002", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMembersLogic(ctx, svcCtx)

	req := &types.GetMembersReq{
		Keyword:  "13800138001",
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetMembers(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, "13800138001", resp.Members[0].Phone)
}

func TestGetMembersLogic_WithSourceFilter(t *testing.T) {
	db := setupMemberTestDB(t)

	member1 := createTestMember(t, db, "wx001", "张三", "13800138001", "active")
	member1.Source = "wechat"
	db.Save(member1)

	member2 := createTestMember(t, db, "wx002", "李四", "13800138002", "active")
	member2.Source = "alipay"
	db.Save(member2)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMembersLogic(ctx, svcCtx)

	req := &types.GetMembersReq{
		Source:   "wechat",
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetMembers(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Total >= 1)
}

func TestGetMembersLogic_WithGenderFilter(t *testing.T) {
	db := setupMemberTestDB(t)

	member1 := createTestMember(t, db, "wx001", "张三", "13800138001", "active")
	member1.Gender = 1
	db.Save(member1)

	member2 := createTestMember(t, db, "wx002", "李四", "13800138002", "active")
	member2.Gender = 2
	db.Save(member2)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMembersLogic(ctx, svcCtx)

	req := &types.GetMembersReq{
		Gender:   1,
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetMembers(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Total >= 1)
}

func TestGetMembersLogic_EmptyResult(t *testing.T) {
	db := setupMemberTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMembersLogic(ctx, svcCtx)

	req := &types.GetMembersReq{
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetMembers(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(0), resp.Total)
	assert.Len(t, resp.Members, 0)
}

func TestGetMembersLogic_Pagination(t *testing.T) {
	db := setupMemberTestDB(t)

	for i := 0; i < 25; i++ {
		createTestMember(t, db, fmt.Sprintf("wx%03d", i), fmt.Sprintf("用户%d", i), fmt.Sprintf("13800138%03d", i), "active")
	}

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMembersLogic(ctx, svcCtx)

	req := &types.GetMembersReq{
		Page:     1,
		PageSize: 10,
	}

	resp, err := logic.GetMembers(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(25), resp.Total)
	assert.Len(t, resp.Members, 10)
}

func TestGetMemberLogic_GetMember_Success(t *testing.T) {
	db := setupMemberTestDB(t)

	member := createTestMember(t, db, "wx001", "张三", "13800138001", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMemberLogic(ctx, svcCtx)

	req := &types.GetMemberReq{Id: member.ID}

	resp, err := logic.GetMember(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, member.ID, resp.Id)
	assert.Equal(t, member.UnionID, resp.UnionID)
	assert.Equal(t, member.Nickname, resp.Nickname)
	assert.Equal(t, member.Phone, resp.Phone)
	assert.Equal(t, member.Status, resp.Status)
}

func TestGetMemberLogic_MemberNotFound(t *testing.T) {
	db := setupMemberTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMemberLogic(ctx, svcCtx)

	req := &types.GetMemberReq{Id: 99999}

	resp, err := logic.GetMember(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, strings.ToLower(err.Error()), "not found")
}

func TestUpdateMemberLogic_UpdateMember_Success(t *testing.T) {
	db := setupMemberTestDB(t)

	member := createTestMember(t, db, "wx001", "张三", "13800138001", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateMemberLogic(ctx, svcCtx)

	req := &types.UpdateMemberReq{
		Nickname: "张三丰",
		Avatar:   "http://example.com/avatar.jpg",
		Gender:   2,
	}

	resp, err := logic.UpdateMember(member.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "张三丰", resp.Nickname)
	assert.Equal(t, "http://example.com/avatar.jpg", resp.Avatar)
	assert.Equal(t, 2, resp.Gender)

	var updatedMember model.Member
	db.First(&updatedMember, member.ID)
	assert.Equal(t, "张三丰", updatedMember.Nickname)
	assert.Equal(t, "http://example.com/avatar.jpg", updatedMember.Avatar)
	assert.Equal(t, 2, updatedMember.Gender)
}

func TestUpdateMemberLogic_MemberNotFound(t *testing.T) {
	db := setupMemberTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateMemberLogic(ctx, svcCtx)

	req := &types.UpdateMemberReq{
		Nickname: "新昵称",
	}

	resp, err := logic.UpdateMember(99999, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpdateMemberStatusLogic_UpdateStatus_Success(t *testing.T) {
	db := setupMemberTestDB(t)

	member := createTestMember(t, db, "wx001", "张三", "13800138001", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateMemberStatusLogic(ctx, svcCtx)

	req := &types.UpdateMemberStatusReq{
		Status: "disabled",
	}

	resp, err := logic.UpdateMemberStatus(member.ID, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)

	var updatedMember model.Member
	db.First(&updatedMember, member.ID)
	assert.Equal(t, "disabled", updatedMember.Status)
}

func TestUpdateMemberStatusLogic_MemberNotFound(t *testing.T) {
	db := setupMemberTestDB(t)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewUpdateMemberStatusLogic(ctx, svcCtx)

	req := &types.UpdateMemberStatusReq{
		Status: "disabled",
	}

	resp, err := logic.UpdateMemberStatus(99999, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetMemberProfileLogic_GetProfile_Success(t *testing.T) {
	db := setupMemberTestDB(t)

	member := createTestMember(t, db, "wx001", "张三", "13800138001", "active")

	profile := &model.MemberProfile{
		MemberID:              member.ID,
		TotalOrders:           10,
		TotalPayment:          5000.00,
		TotalReward:           200.00,
		ParticipatedCampaigns: 5,
	}
	db.Create(profile)

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMemberProfileLogic(ctx, svcCtx)

	req := &types.GetMemberReq{Id: member.ID}

	resp, err := logic.GetMemberProfile(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, member.ID, resp.MemberId)
	assert.Equal(t, 10, resp.TotalOrders)
	assert.Equal(t, 5000.00, resp.TotalPayment)
	assert.Equal(t, 200.00, resp.TotalReward)
	assert.Equal(t, 5, resp.ParticipatedCampaigns)
}

func TestGetMemberProfileLogic_ProfileNotFound(t *testing.T) {
	db := setupMemberTestDB(t)

	member := createTestMember(t, db, "wx001", "张三", "13800138001", "active")

	ctx := context.Background()
	svcCtx := &svc.ServiceContext{DB: db}
	logic := NewGetMemberProfileLogic(ctx, svcCtx)

	req := &types.GetMemberReq{Id: member.ID}

	resp, err := logic.GetMemberProfile(req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}
