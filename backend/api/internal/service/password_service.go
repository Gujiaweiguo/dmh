package service

import (
	"fmt"
	"regexp"
	"time"

	"dmh/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// PasswordService 密码服务
type PasswordService struct {
	db *gorm.DB
}

// NewPasswordService 创建密码服务
func NewPasswordService(db *gorm.DB) *PasswordService {
	return &PasswordService{
		db: db,
	}
}

// GetPasswordPolicy 获取密码策略
func (s *PasswordService) GetPasswordPolicy() (*model.PasswordPolicy, error) {
	var policies []model.PasswordPolicy

	// 获取所有策略，按ID倒序（最新的在前）
	err := s.db.Order("id DESC").Find(&policies).Error
	if err != nil {
		return nil, err
	}

	// 如果没有策略，返回默认策略
	if len(policies) == 0 {
		return s.getDefaultPasswordPolicy(), nil
	}

	// 返回最新的策略
	return &policies[0], nil
}

// getDefaultPasswordPolicy 获取默认密码策略
func (s *PasswordService) getDefaultPasswordPolicy() *model.PasswordPolicy {
	return &model.PasswordPolicy{
		MinLength:             8,
		RequireUppercase:      true,
		RequireLowercase:      true,
		RequireNumbers:        true,
		RequireSpecialChars:   true,
		MaxAge:                90,
		HistoryCount:          5,
		MaxLoginAttempts:      5,
		LockoutDuration:       30,
		SessionTimeout:        480,
		MaxConcurrentSessions: 3,
	}
}

// ValidatePassword 验证密码是否符合策略
func (s *PasswordService) ValidatePassword(password string, userID int64) error {
	policy, err := s.GetPasswordPolicy()
	if err != nil {
		return fmt.Errorf("获取密码策略失败: %v", err)
	}

	// 检查密码长度
	if len(password) < policy.MinLength {
		return fmt.Errorf("密码长度不能少于%d位", policy.MinLength)
	}

	// 检查大写字母
	if policy.RequireUppercase {
		if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
			return fmt.Errorf("密码必须包含至少一个大写字母")
		}
	}

	// 检查小写字母
	if policy.RequireLowercase {
		if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
			return fmt.Errorf("密码必须包含至少一个小写字母")
		}
	}

	// 检查数字
	if policy.RequireNumbers {
		if matched, _ := regexp.MatchString(`[0-9]`, password); !matched {
			return fmt.Errorf("密码必须包含至少一个数字")
		}
	}

	// 检查特殊字符
	if policy.RequireSpecialChars {
		if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\?]`, password); !matched {
			return fmt.Errorf("密码必须包含至少一个特殊字符")
		}
	}

	// 检查密码历史
	if userID > 0 {
		if err := s.checkPasswordHistory(password, userID, policy.HistoryCount); err != nil {
			return err
		}
	}

	return nil
}

// checkPasswordHistory 检查密码历史
func (s *PasswordService) checkPasswordHistory(password string, userID int64, historyCount int) error {
	if historyCount <= 0 {
		return nil
	}

	var histories []model.PasswordHistory
	err := s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(historyCount).
		Find(&histories).Error

	if err != nil {
		return fmt.Errorf("检查密码历史失败: %v", err)
	}

	for _, history := range histories {
		if err := bcrypt.CompareHashAndPassword([]byte(history.PasswordHash), []byte(password)); err == nil {
			return fmt.Errorf("不能使用最近%d次使用过的密码", historyCount)
		}
	}

	return nil
}

// SavePasswordHistory 保存密码历史
func (s *PasswordService) SavePasswordHistory(userID int64, passwordHash string) error {
	policy, err := s.GetPasswordPolicy()
	if err != nil {
		return err
	}

	// 保存新的密码历史
	history := &model.PasswordHistory{
		UserID:       userID,
		PasswordHash: passwordHash,
	}

	if err := s.db.Create(history).Error; err != nil {
		return fmt.Errorf("保存密码历史失败: %v", err)
	}

	// 清理超出限制的历史记录
	if policy.HistoryCount > 0 {
		var oldHistories []model.PasswordHistory
		err := s.db.Where("user_id = ?", userID).
			Order("created_at DESC").
			Offset(policy.HistoryCount).
			Find(&oldHistories).Error

		if err == nil && len(oldHistories) > 0 {
			var ids []int64
			for _, h := range oldHistories {
				ids = append(ids, h.ID)
			}
			s.db.Where("id IN ?", ids).Delete(&model.PasswordHistory{})
		}
	}

	return nil
}

// IsPasswordExpired 检查密码是否过期
func (s *PasswordService) IsPasswordExpired(userID int64) (bool, error) {
	policy, err := s.GetPasswordPolicy()
	if err != nil {
		return false, err
	}

	if policy.MaxAge <= 0 {
		return false, nil // 密码不过期
	}

	var user model.User
	err = s.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return false, err
	}

	// 检查最后一次密码更新时间
	var lastPasswordChange time.Time
	var history model.PasswordHistory
	err = s.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&history).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果没有密码历史，使用用户创建时间
			lastPasswordChange = user.CreatedAt
		} else {
			return false, err
		}
	} else {
		lastPasswordChange = history.CreatedAt
	}

	// 检查是否过期
	expireTime := lastPasswordChange.AddDate(0, 0, policy.MaxAge)
	return time.Now().After(expireTime), nil
}

// GeneratePasswordStrengthScore 生成密码强度评分
func (s *PasswordService) GeneratePasswordStrengthScore(password string) int {
	score := 0

	// 长度评分
	if len(password) >= 8 {
		score += 25
	}
	if len(password) >= 12 {
		score += 25
	}

	// 字符类型评分
	if matched, _ := regexp.MatchString(`[a-z]`, password); matched {
		score += 10
	}
	if matched, _ := regexp.MatchString(`[A-Z]`, password); matched {
		score += 10
	}
	if matched, _ := regexp.MatchString(`[0-9]`, password); matched {
		score += 10
	}
	if matched, _ := regexp.MatchString(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\?]`, password); matched {
		score += 20
	}

	// 复杂度评分
	if len(password) > 0 {
		uniqueChars := make(map[rune]bool)
		for _, char := range password {
			uniqueChars[char] = true
		}
		if len(uniqueChars) > len(password)/2 {
			score += 10
		}
	}

	// 确保评分在0-100之间
	if score > 100 {
		score = 100
	}

	return score
}

// GetPasswordStrengthLevel 获取密码强度等级
func (s *PasswordService) GetPasswordStrengthLevel(score int) string {
	if score >= 80 {
		return "强"
	} else if score >= 60 {
		return "中等"
	} else if score >= 40 {
		return "弱"
	} else {
		return "很弱"
	}
}

// HashPassword 加密密码
func (s *PasswordService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %v", err)
	}
	return string(hashedBytes), nil
}

// VerifyPassword 验证密码
func (s *PasswordService) VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// UpdatePasswordPolicy 更新密码策略
func (s *PasswordService) UpdatePasswordPolicy(policy *model.PasswordPolicy) error {
	var existingPolicy model.PasswordPolicy
	err := s.db.First(&existingPolicy).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新策略
			return s.db.Create(policy).Error
		}
		return err
	}

	// 更新现有策略
	policy.ID = existingPolicy.ID
	return s.db.Save(policy).Error
}
