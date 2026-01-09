package brand

import (
	"context"
	"encoding/json"
	"fmt"
)

// UserInfo 用户信息结构
type UserInfo struct {
	UserID int64
	Role   string
}

// GetUserInfoFromContext 从go-zero JWT context中获取用户信息
func GetUserInfoFromContext(ctx context.Context) (*UserInfo, error) {
	// 获取用户ID
	userID, ok := ctx.Value("userId").(json.Number)
	if !ok {
		return nil, fmt.Errorf("无法从context获取用户ID")
	}
	
	userIDInt64, err := userID.Int64()
	if err != nil {
		return nil, fmt.Errorf("用户ID转换失败: %v", err)
	}

	// 获取用户角色
	roles, ok := ctx.Value("roles").([]interface{})
	if !ok || len(roles) == 0 {
		return nil, fmt.Errorf("无法从context获取用户角色")
	}
	
	// 转换角色为字符串
	var roleStrings []string
	for _, role := range roles {
		if roleStr, ok := role.(string); ok {
			roleStrings = append(roleStrings, roleStr)
		}
	}
	
	if len(roleStrings) == 0 {
		return nil, fmt.Errorf("用户角色信息无效")
	}
	
	return &UserInfo{
		UserID: userIDInt64,
		Role:   roleStrings[0], // 取第一个角色
	}, nil
}