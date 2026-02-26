package factory

import (
	"dmh/model"
	"time"

	"gorm.io/gorm"
)

// UserFactory 用户测试数据工厂
type UserFactory struct {
	BaseFactory
}

// NewUserFactory 创建用户工厂实例
func NewUserFactory() *UserFactory {
	return &UserFactory{}
}

// UserOverrides 定义用户字段覆盖选项
type UserOverrides struct {
	ID            *int64
	Username      *string
	Password      *string
	Phone         *string
	Email         *string
	Avatar        *string
	RealName      *string
	Role          *string
	Status        *string
	LoginAttempts *int
}

// parseTimePointer 解析时间值，支持 time.Time 和 *time.Time 两种类型
func parseTimePtr(v any) *time.Time {
	if pt, ok := v.(*time.Time); ok {
		return pt
	}
	if t, ok := v.(time.Time); ok {
		return &t
	}
	return nil
}

// Build 创建一个用户实例，使用合理的默认值
func (f *UserFactory) Build() *model.User {
	return f.BuildWith(nil)
}

// BuildWith 创建一个用户实例，允许覆盖指定字段
// 使用 map[string]any 类型提供灵活的字段覆盖
func (f *UserFactory) BuildWith(overrides map[string]any) *model.User {
	user := &model.User{
		Username:      f.RandomUsername("testuser"),
		Password:      "$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrD1pFpL/qVzLqBxJMdOqYQe1V5q6S", // "123456" 的 bcrypt hash
		Phone:         f.RandomPhone(),
		Email:         f.RandomEmail(),
		Avatar:        "",
		RealName:      "测试用户",
		Role:          RoleParticipant,
		Status:        StatusActive,
		LoginAttempts: 0,
	}

	// 应用覆盖值
	if overrides != nil {
		if v, ok := overrides["id"]; ok {
			user.Id = v.(int64)
		}
		if v, ok := overrides["username"]; ok {
			user.Username = v.(string)
		}
		if v, ok := overrides["password"]; ok {
			user.Password = v.(string)
		}
		if v, ok := overrides["phone"]; ok {
			user.Phone = v.(string)
		}
		if v, ok := overrides["email"]; ok {
			user.Email = v.(string)
		}
		if v, ok := overrides["avatar"]; ok {
			user.Avatar = v.(string)
		}
		if v, ok := overrides["real_name"]; ok {
			user.RealName = v.(string)
		}
		if v, ok := overrides["realName"]; ok {
			user.RealName = v.(string)
		}
		if v, ok := overrides["role"]; ok {
			user.Role = v.(string)
		}
		if v, ok := overrides["status"]; ok {
			user.Status = v.(string)
		}
		if v, ok := overrides["login_attempts"]; ok {
			user.LoginAttempts = v.(int)
		}
		if v, ok := overrides["loginAttempts"]; ok {
			user.LoginAttempts = v.(int)
		}
		if v, ok := overrides["locked_until"]; ok {
			user.LockedUntil = parseTimePtr(v)
		}
		if v, ok := overrides["lockedUntil"]; ok {
			user.LockedUntil = parseTimePtr(v)
		}
	}

	return user
}

// BuildList 创建多个用户实例
func (f *UserFactory) BuildList(count int) []*model.User {
	users := make([]*model.User, count)
	for i := 0; i < count; i++ {
		users[i] = f.Build()
	}
	return users
}

// Create 创建并持久化用户到数据库
func (f *UserFactory) Create(db *gorm.DB) (*model.User, error) {
	return f.CreateWith(db, nil)
}

// CreateWith 创建并持久化用户到数据库，允许覆盖指定字段
func (f *UserFactory) CreateWith(db *gorm.DB, overrides map[string]any) (*model.User, error) {
	user := f.BuildWith(overrides)
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// CreateList 创建并持久化多个用户到数据库
func (f *UserFactory) CreateList(db *gorm.DB, count int) ([]*model.User, error) {
	users := f.BuildList(count)
	for _, user := range users {
		if err := db.Create(user).Error; err != nil {
			return nil, err
		}
	}
	return users, nil
}

// ========== 便捷方法：创建特定类型用户 ==========

// BuildPlatformAdmin 创建平台管理员用户（不持久化）
func (f *UserFactory) BuildPlatformAdmin() *model.User {
	return f.BuildWith(map[string]any{
		"role":      RolePlatformAdmin,
		"username":  f.RandomUsername("admin"),
		"real_name": "平台管理员",
	})
}

// CreatePlatformAdmin 创建并持久化平台管理员用户
func (f *UserFactory) CreatePlatformAdmin(db *gorm.DB) (*model.User, error) {
	return f.CreateWith(db, map[string]any{
		"role":      RolePlatformAdmin,
		"username":  f.RandomUsername("admin"),
		"real_name": "平台管理员",
	})
}

// BuildBrandAdmin 创建品牌管理员用户（不持久化）
func (f *UserFactory) BuildBrandAdmin() *model.User {
	return f.BuildWith(map[string]any{
		"role":      RoleBrandAdmin,
		"username":  f.RandomUsername("brand_admin"),
		"real_name": "品牌管理员",
	})
}

// CreateBrandAdmin 创建并持久化品牌管理员用户
func (f *UserFactory) CreateBrandAdmin(db *gorm.DB) (*model.User, error) {
	return f.CreateWith(db, map[string]any{
		"role":      RoleBrandAdmin,
		"username":  f.RandomUsername("brand_admin"),
		"real_name": "品牌管理员",
	})
}

// BuildParticipant 创建普通参与者用户（不持久化）
func (f *UserFactory) BuildParticipant() *model.User {
	return f.BuildWith(map[string]any{
		"role": RoleParticipant,
	})
}

// CreateParticipant 创建并持久化普通参与者用户
func (f *UserFactory) CreateParticipant(db *gorm.DB) (*model.User, error) {
	return f.CreateWith(db, map[string]any{
		"role": RoleParticipant,
	})
}

// BuildLocked 创建被锁定的用户（不持久化）
func (f *UserFactory) BuildLocked() *model.User {
	lockedUntil := f.FutureTime(30 * 24 * time.Hour) // 30 天后解锁
	return f.BuildWith(map[string]any{
		"status":        StatusLocked,
		"loginAttempts": 5,
		"locked_until":  &lockedUntil,
	})
}
