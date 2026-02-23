package model

import (
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupModelTestDB(t *testing.T) *gorm.DB {
	dsn := "root:Admin168@tcp(127.0.0.1:3306)/dmh_test?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	err = db.Exec("DROP TABLE IF EXISTS faq_items").Error
	if err != nil {
		t.Fatalf("Failed to drop faq_items table: %v", err)
	}

	err = db.AutoMigrate(
		&User{}, &Role{}, &UserRole{}, &Permission{}, &RolePermission{},
		&Brand{}, &UserBrand{}, &BrandAsset{},
		&Campaign{}, &Order{}, &Reward{}, &UserBalance{}, &SyncLog{},
		&Distributor{}, &DistributorApplication{}, &DistributorLevelReward{}, &DistributorReward{}, &DistributorLink{},
		&Menu{}, &RoleMenu{},
		&Withdrawal{},
		&Member{}, &MemberTag{}, &MemberProfile{}, &MemberTagLink{}, &MemberBrandLink{}, &MemberMergeRequest{}, &ExportRequest{},
		&PosterTemplate{}, &PosterRecord{}, &PosterTemplateConfig{},
		&PageConfig{}, &AuditLog{},
		&UserFeedback{}, &FAQItem{}, &FeatureSatisfactionSurvey{}, &FeatureUsageStat{}, &FeedbackTag{}, &FeedbackTagRelation{},
		&PasswordPolicy{}, &PasswordHistory{}, &LoginAttempt{}, &UserSession{}, &SecurityEvent{},
		&VerificationRecord{},
		&Promoter{}, &PromoterLink{}, &PromoterReward{},
		&Material{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate database: %v", err)
	}

	err = db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error
	if err != nil {
		t.Fatalf("Failed to disable foreign key checks: %v", err)
	}

	tables := []string{
		"promoter_rewards", "promoter_links", "promoters",
		"materials",
		"distributor_rewards", "distributor_links", "distributor_level_rewards", "distributor_applications", "distributors",
		"sync_logs", "rewards", "orders", "campaigns",
		"member_merge_requests", "member_brand_links", "member_tag_links", "member_profiles", "member_tags", "members",
		"export_requests", "withdrawals", "verification_records", "poster_template_configs", "poster_records", "poster_templates",
		"feedback_tag_relations", "feedback_tags", "feature_usage_stats", "feature_satisfaction_surveys", "faq_items", "user_feedback",
		"password_histories", "password_policies", "user_sessions", "login_attempts", "security_events", "audit_logs",
		"role_menus", "menus", "role_permissions", "permissions", "user_roles", "user_brands", "brand_assets", "brands", "roles", "users",
	}
	for _, table := range tables {
		if execErr := db.Exec("TRUNCATE TABLE " + table).Error; execErr != nil {
			t.Fatalf("Failed to truncate table %s: %v", table, execErr)
		}
	}

	err = db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error
	if err != nil {
		t.Fatalf("Failed to enable foreign key checks: %v", err)
	}

	return db
}

// ========== User Model Tests ==========

func TestUserCRUD(t *testing.T) {
	db := setupModelTestDB(t)

	user := &User{
		Username: "testuser",
		Password: "hashedpassword",
		Phone:    "13800138000",
		Email:    "test@example.com",
		RealName: "Test User",
		Role:     "participant",
		Status:   "active",
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	if user.Id == 0 {
		t.Fatal("User ID should not be 0")
	}

	var readUser User
	if err := db.First(&readUser, user.Id).Error; err != nil {
		t.Fatalf("Failed to read user: %v", err)
	}
	if readUser.Username != user.Username {
		t.Errorf("Username mismatch")
	}

	db.Model(user).Update("real_name", "Updated Name")
	var updatedUser User
	db.First(&updatedUser, user.Id)
	if updatedUser.RealName != "Updated Name" {
		t.Errorf("RealName not updated")
	}

	db.Delete(user)
	var deletedUser User
	result := db.First(&deletedUser, user.Id)
	if result.Error == nil {
		t.Error("User should be deleted")
	}
}

func TestUserUniqueConstraint(t *testing.T) {
	db := setupModelTestDB(t)

	user1 := &User{Username: "uniqueuser", Password: "pass", Phone: "13800138000", Status: "active"}
	if err := db.Create(user1).Error; err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}

	user2 := &User{Username: "uniqueuser", Password: "pass", Phone: "13800138001", Status: "active"}
	err := db.Create(user2).Error
	if err == nil {
		t.Error("Should fail due to unique constraint")
	}
}

// ========== Role & Permission Tests ==========

func TestRoleCRUD(t *testing.T) {
	db := setupModelTestDB(t)

	role := &Role{Name: "Test Role", Code: "test_role", Description: "Test"}
	if err := db.Create(role).Error; err != nil {
		t.Fatalf("Failed to create role: %v", err)
	}
	if role.ID == 0 {
		t.Fatal("Role ID should not be 0")
	}

	role2 := &Role{Name: "Another", Code: "test_role"}
	err := db.Create(role2).Error
	if err == nil {
		t.Error("Should fail due to unique code")
	}
}

// ========== Brand Model Tests ==========

func TestBrandCRUD(t *testing.T) {
	db := setupModelTestDB(t)

	brand := &Brand{Name: "Test Brand", Logo: "https://example.com/logo.png", Status: "active"}
	if err := db.Create(brand).Error; err != nil {
		t.Fatalf("Failed to create brand: %v", err)
	}
	if brand.Id == 0 {
		t.Fatal("Brand ID should not be 0")
	}

	var readBrand Brand
	db.First(&readBrand, brand.Id)
	if readBrand.Name != "Test Brand" {
		t.Errorf("Brand name mismatch")
	}
}

func TestBrandAssetCRUD(t *testing.T) {
	db := setupModelTestDB(t)

	brand := &Brand{Name: "Test", Status: "active"}
	db.Create(brand)

	asset := &BrandAsset{BrandID: brand.Id, Name: "Asset", Type: "image", FileUrl: "https://example.com/img.png", FileSize: 1024}
	if err := db.Create(asset).Error; err != nil {
		t.Fatalf("Failed to create asset: %v", err)
	}

	var assets []BrandAsset
	db.Where("brand_id = ?", brand.Id).Find(&assets)
	if len(assets) != 1 {
		t.Errorf("Expected 1 asset, got %d", len(assets))
	}
}

// ========== Campaign Model Tests ==========

func TestCampaignCRUD(t *testing.T) {
	db := setupModelTestDB(t)

	brand := &Brand{Name: "Test", Status: "active"}
	db.Create(brand)

	campaign := &Campaign{BrandId: brand.Id, Name: "Test Campaign", Status: "active", StartTime: time.Now(), EndTime: time.Now().Add(24 * time.Hour)}
	if err := db.Create(campaign).Error; err != nil {
		t.Fatalf("Failed to create campaign: %v", err)
	}
	if campaign.Id == 0 {
		t.Fatal("Campaign ID should not be 0")
	}
}

func TestCampaignSoftDelete(t *testing.T) {
	db := setupModelTestDB(t)

	brand := &Brand{Name: "Test", Status: "active"}
	db.Create(brand)

	campaign := &Campaign{BrandId: brand.Id, Name: "Test", Status: "active", StartTime: time.Now(), EndTime: time.Now().Add(24 * time.Hour)}
	db.Create(campaign)
	deletedAt := time.Now()
	if err := db.Model(&Campaign{}).Where("id = ?", campaign.Id).Update("deleted_at", deletedAt).Error; err != nil {
		t.Fatalf("Failed to soft delete campaign: %v", err)
	}

	var result Campaign
	if db.Where("id = ? AND deleted_at IS NULL", campaign.Id).First(&result).Error == nil {
		t.Error("Should be soft deleted")
	}
	if db.Unscoped().First(&result, campaign.Id).Error != nil {
		t.Error("Should find with Unscoped")
	}
	if result.DeletedAt == nil {
		t.Error("DeletedAt should be set")
	}
}

// ========== Order Model Tests ==========

func TestOrderCRUD(t *testing.T) {
	db := setupModelTestDB(t)

	brand := &Brand{Name: "Test", Status: "active"}
	db.Create(brand)

	campaign := &Campaign{BrandId: brand.Id, Name: "Test", Status: "active", StartTime: time.Now(), EndTime: time.Now().Add(24 * time.Hour)}
	db.Create(campaign)

	order := &Order{CampaignId: campaign.Id, Phone: "13800138000", FormData: "{}", Amount: 100.00, PayStatus: "paid", Status: "active"}
	if err := db.Create(order).Error; err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	if order.Id == 0 {
		t.Fatal("Order ID should not be 0")
	}
}

// ========== Distributor Model Tests ==========

func TestDistributorCRUD(t *testing.T) {
	db := setupModelTestDB(t)

	user := &User{Username: "test", Password: "pass", Phone: "13800138000", Status: "active"}
	brand := &Brand{Name: "Test", Status: "active"}
	db.Create(user)
	db.Create(brand)

	distributor := &Distributor{UserId: user.Id, BrandId: brand.Id, Level: 1, Status: "active"}
	if err := db.Create(distributor).Error; err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if distributor.Id == 0 {
		t.Fatal("ID should not be 0")
	}

	distributor2 := &Distributor{UserId: user.Id, BrandId: brand.Id, Level: 2}
	err := db.Create(distributor2).Error
	if err == nil {
		t.Error("Should fail unique constraint")
	}
}

// ========== Menu Model Tests ==========

func TestMenuCRUD(t *testing.T) {
	db := setupModelTestDB(t)

	menu := &Menu{Name: "Dashboard", Code: "dashboard", Path: "/dashboard", Sort: 1, Type: "menu", Platform: "admin", Status: "active"}
	if err := db.Create(menu).Error; err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if menu.ID == 0 {
		t.Fatal("Menu ID should not be 0")
	}
}

// ========== Member Model Tests ==========

func TestMemberCRUD(t *testing.T) {
	db := setupModelTestDB(t)

	member := &Member{UnionID: "union123", Nickname: "Test", Phone: "13800138000", Status: "active"}
	if err := db.Create(member).Error; err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if member.ID == 0 {
		t.Fatal("Member ID should not be 0")
	}

	var found Member
	if err := db.Where("unionid = ?", "union123").First(&found).Error; err != nil {
		t.Fatalf("Failed to find: %v", err)
	}
}

// ========== Feedback Model Tests ==========

func TestUserFeedbackCRUD(t *testing.T) {
	db := setupModelTestDB(t)

	user := &User{Username: "test", Password: "pass", Phone: "13800138000", Status: "active"}
	db.Create(user)

	feedback := &UserFeedback{UserID: user.Id, Category: "suggestion", Title: "Title", Content: "Content", Priority: "normal", Status: "pending"}
	if err := db.Create(feedback).Error; err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if feedback.ID == 0 {
		t.Fatal("ID should not be 0")
	}
}

// ========== Security Model Tests ==========

func TestPasswordPolicyCRUD(t *testing.T) {
	db := setupModelTestDB(t)

	policy := &PasswordPolicy{MinLength: 8, RequireUppercase: true, RequireLowercase: true, RequireNumbers: true, RequireSpecialChars: true}
	if err := db.Create(policy).Error; err != nil {
		t.Fatalf("Failed: %v", err)
	}
	if policy.ID == 0 {
		t.Fatal("ID should not be 0")
	}
}

// ========== Table Name Tests ==========

func TestAllTableNames(t *testing.T) {
	tests := []struct {
		name     string
		got      string
		expected string
	}{
		{"User", User{}.TableName(), "users"},
		{"Role", Role{}.TableName(), "roles"},
		{"UserRole", UserRole{}.TableName(), "user_roles"},
		{"Permission", Permission{}.TableName(), "permissions"},
		{"RolePermission", RolePermission{}.TableName(), "role_permissions"},
		{"Brand", Brand{}.TableName(), "brands"},
		{"BrandAsset", BrandAsset{}.TableName(), "brand_assets"},
		{"UserBrand", UserBrand{}.TableName(), "user_brands"},
		{"Campaign", (&Campaign{}).TableName(), "campaigns"},
		{"Order", (&Order{}).TableName(), "orders"},
		{"Reward", (&Reward{}).TableName(), "rewards"},
		{"UserBalance", (&UserBalance{}).TableName(), "user_balances"},
		{"SyncLog", (&SyncLog{}).TableName(), "sync_logs"},
		{"Distributor", Distributor{}.TableName(), "distributors"},
		{"DistributorApplication", DistributorApplication{}.TableName(), "distributor_applications"},
		{"DistributorLevelReward", DistributorLevelReward{}.TableName(), "distributor_level_rewards"},
		{"DistributorReward", DistributorReward{}.TableName(), "distributor_rewards"},
		{"DistributorLink", DistributorLink{}.TableName(), "distributor_links"},
		{"Menu", Menu{}.TableName(), "menus"},
		{"RoleMenu", RoleMenu{}.TableName(), "role_menus"},
		{"Withdrawal", Withdrawal{}.TableName(), "withdrawals"},
		{"Member", Member{}.TableName(), "members"},
		{"MemberTag", MemberTag{}.TableName(), "member_tags"},
		{"PosterTemplate", PosterTemplate{}.TableName(), "poster_templates"},
		{"PosterRecord", PosterRecord{}.TableName(), "poster_records"},
		{"PageConfig", (&PageConfig{}).TableName(), "page_configs"},
		{"UserFeedback", UserFeedback{}.TableName(), "user_feedback"},
		{"FAQItem", FAQItem{}.TableName(), "faq_items"},
		{"PasswordPolicy", PasswordPolicy{}.TableName(), "password_policies"},
		{"VerificationRecord", VerificationRecord{}.TableName(), "verification_records"},
		{"AuditLog", AuditLog{}.TableName(), "audit_logs"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("TableName() = %v, want %v", tt.got, tt.expected)
			}
		})
	}
}
