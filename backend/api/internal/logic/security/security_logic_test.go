package security

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

func setupSecurityLogicTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	err = db.AutoMigrate(&model.PasswordPolicy{}, &model.UserSession{}, &model.User{}, &model.SecurityEvent{}, &model.AuditLog{}, &model.LoginAttempt{})
	if err != nil {
		t.Fatalf("failed to migrate sqlite db: %v", err)
	}

	return db
}

func TestSecurityLogicConstructors(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}

	assert.NotNil(t, NewCheckPasswordStrengthLogic(ctx, svcCtx))
	assert.NotNil(t, NewForceLogoutUserLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetAuditLogsLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetLoginAttemptsLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetPasswordPolicyLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetSecurityEventsLogic(ctx, svcCtx))
	assert.NotNil(t, NewGetUserSessionsLogic(ctx, svcCtx))
	assert.NotNil(t, NewHandleSecurityEventLogic(ctx, svcCtx))
	assert.NotNil(t, NewRevokeSessionLogic(ctx, svcCtx))
	assert.NotNil(t, NewUpdatePasswordPolicyLogic(ctx, svcCtx))
}

func TestSecurityLogicCurrentMethodBehavior(t *testing.T) {
	ctx := context.Background()
	svcCtx := &svc.ServiceContext{}

	check := NewCheckPasswordStrengthLogic(ctx, svcCtx)
	checkResp, checkErr := check.CheckPasswordStrength(&types.ChangePasswordReq{NewPassword: "Password123!"})
	assert.NoError(t, checkErr)
	assert.NotNil(t, checkResp)
	assert.NotZero(t, checkResp.Score)

	force := NewForceLogoutUserLogic(ctx, svcCtx)
	forceResp, forceErr := force.ForceLogoutUser(1, &types.ForceLogoutReq{})
	assert.Error(t, forceErr)
	assert.Nil(t, forceResp)

	audit := NewGetAuditLogsLogic(ctx, svcCtx)
	auditResp, auditErr := audit.GetAuditLogs()
	assert.NoError(t, auditErr)
	assert.NotNil(t, auditResp)
	assert.Equal(t, int64(0), auditResp.Total)
	assert.Len(t, auditResp.Logs, 0)

	attempts := NewGetLoginAttemptsLogic(ctx, svcCtx)
	attemptsResp, attemptsErr := attempts.GetLoginAttempts()
	assert.NoError(t, attemptsErr)
	assert.NotNil(t, attemptsResp)
	assert.Equal(t, int64(0), attemptsResp.Total)
	assert.Len(t, attemptsResp.Attempts, 0)

	policy := NewGetPasswordPolicyLogic(ctx, svcCtx)
	policyResp, policyErr := policy.GetPasswordPolicy()
	assert.NoError(t, policyErr)
	assert.NotNil(t, policyResp)
	assert.Equal(t, 8, policyResp.MinLength)
	assert.Equal(t, 5, policyResp.MaxLoginAttempts)

	events := NewGetSecurityEventsLogic(ctx, svcCtx)
	eventsResp, eventsErr := events.GetSecurityEvents()
	assert.NoError(t, eventsErr)
	assert.NotNil(t, eventsResp)
	assert.Equal(t, int64(0), eventsResp.Total)
	assert.Len(t, eventsResp.Events, 0)

	sessions := NewGetUserSessionsLogic(ctx, svcCtx)
	sessionsResp, sessionsErr := sessions.GetUserSessions()
	assert.NoError(t, sessionsErr)
	assert.NotNil(t, sessionsResp)
	assert.Equal(t, int64(0), sessionsResp.Total)
	assert.Len(t, sessionsResp.Sessions, 0)

	handle := NewHandleSecurityEventLogic(ctx, svcCtx)
	handleResp, handleErr := handle.HandleSecurityEvent(1, &types.HandleSecurityEventReq{})
	assert.Error(t, handleErr)
	assert.Nil(t, handleResp)

	revoke := NewRevokeSessionLogic(ctx, svcCtx)
	revokeResp, revokeErr := revoke.RevokeSession("session-id")
	assert.Error(t, revokeErr)
	assert.Nil(t, revokeResp)

	update := NewUpdatePasswordPolicyLogic(ctx, svcCtx)
	updateResp, updateErr := update.UpdatePasswordPolicy(&types.UpdatePasswordPolicyReq{MinLength: 8})
	assert.Error(t, updateErr)
	assert.Nil(t, updateResp)
}

func TestUpdatePasswordPolicyLogic_Success(t *testing.T) {
	db := setupSecurityLogicTestDB(t)

	logic := NewUpdatePasswordPolicyLogic(context.WithValue(context.Background(), "userId", int64(1)), &svc.ServiceContext{DB: db})
	resp, err := logic.UpdatePasswordPolicy(&types.UpdatePasswordPolicyReq{
		MinLength:             10,
		RequireUppercase:      true,
		RequireLowercase:      true,
		RequireNumbers:        true,
		RequireSpecialChars:   false,
		MaxAge:                120,
		HistoryCount:          6,
		MaxLoginAttempts:      6,
		LockoutDuration:       20,
		SessionTimeout:        180,
		MaxConcurrentSessions: 4,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 10, resp.MinLength)

	var saved model.PasswordPolicy
	err = db.First(&saved).Error
	assert.NoError(t, err)
	assert.Equal(t, 10, saved.MinLength)
	assert.Equal(t, 180, saved.SessionTimeout)
}

func TestForceLogoutUserLogic_Success(t *testing.T) {
	db := setupSecurityLogicTestDB(t)

	user := &model.User{Username: "target", Password: "hashed", Phone: "13800138000", Status: "active"}
	assert.NoError(t, db.Create(user).Error)
	assert.NoError(t, db.Create(&model.UserSession{ID: "s1", UserID: user.Id, ClientIP: "127.0.0.1", UserAgent: "ua", Status: "active"}).Error)
	assert.NoError(t, db.Create(&model.UserSession{ID: "s2", UserID: user.Id, ClientIP: "127.0.0.1", UserAgent: "ua", Status: "active"}).Error)

	logic := NewForceLogoutUserLogic(context.WithValue(context.Background(), "userId", int64(99)), &svc.ServiceContext{DB: db})
	resp, err := logic.ForceLogoutUser(user.Id, &types.ForceLogoutReq{Reason: "security check"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "用户已强制下线", resp.Message)

	var activeCount int64
	err = db.Model(&model.UserSession{}).Where("user_id = ? AND status = ?", user.Id, "active").Count(&activeCount).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(0), activeCount)
}

func validPasswordPolicyReq() *types.UpdatePasswordPolicyReq {
	return &types.UpdatePasswordPolicyReq{
		MinLength:             8,
		RequireUppercase:      true,
		RequireLowercase:      true,
		RequireNumbers:        true,
		RequireSpecialChars:   true,
		MaxAge:                90,
		HistoryCount:          5,
		MaxLoginAttempts:      5,
		LockoutDuration:       30,
		SessionTimeout:        120,
		MaxConcurrentSessions: 3,
	}
}

func TestValidatePasswordPolicy_Errors(t *testing.T) {
	cases := []struct {
		name    string
		mutate  func(*types.UpdatePasswordPolicyReq)
		errText string
	}{
		{
			name: "min length too small",
			mutate: func(req *types.UpdatePasswordPolicyReq) {
				req.MinLength = 5
			},
			errText: "最小密码长度必须在6到50之间",
		},
		{
			name: "max age too large",
			mutate: func(req *types.UpdatePasswordPolicyReq) {
				req.MaxAge = 366
			},
			errText: "密码有效期必须在0到365天之间",
		},
		{
			name: "history count too large",
			mutate: func(req *types.UpdatePasswordPolicyReq) {
				req.HistoryCount = 21
			},
			errText: "历史密码记录数量必须在0到20之间",
		},
		{
			name: "max login attempts too small",
			mutate: func(req *types.UpdatePasswordPolicyReq) {
				req.MaxLoginAttempts = 0
			},
			errText: "最大登录尝试次数必须在1到20之间",
		},
		{
			name: "lockout duration too large",
			mutate: func(req *types.UpdatePasswordPolicyReq) {
				req.LockoutDuration = 1441
			},
			errText: "锁定时长必须在1到1440分钟之间",
		},
		{
			name: "session timeout too small",
			mutate: func(req *types.UpdatePasswordPolicyReq) {
				req.SessionTimeout = 4
			},
			errText: "会话超时时间必须在5到1440分钟之间",
		},
		{
			name: "max concurrent sessions too large",
			mutate: func(req *types.UpdatePasswordPolicyReq) {
				req.MaxConcurrentSessions = 11
			},
			errText: "最大并发会话数必须在1到10之间",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := validPasswordPolicyReq()
			tc.mutate(req)

			err := validatePasswordPolicy(req)
			assert.EqualError(t, err, tc.errText)
		})
	}
}

func TestUserIDFromContext_Conversions(t *testing.T) {
	t.Run("int64", func(t *testing.T) {
		id, ok := userIDFromContext(context.WithValue(context.Background(), "userId", int64(7)))
		assert.True(t, ok)
		assert.Equal(t, int64(7), id)
	})

	t.Run("int", func(t *testing.T) {
		id, ok := userIDFromContext(context.WithValue(context.Background(), "userId", 8))
		assert.True(t, ok)
		assert.Equal(t, int64(8), id)
	})

	t.Run("float64", func(t *testing.T) {
		id, ok := userIDFromContext(context.WithValue(context.Background(), "userId", float64(9)))
		assert.True(t, ok)
		assert.Equal(t, int64(9), id)
	})

	t.Run("string", func(t *testing.T) {
		id, ok := userIDFromContext(context.WithValue(context.Background(), "userId", "10"))
		assert.True(t, ok)
		assert.Equal(t, int64(10), id)
	})

	t.Run("invalid string", func(t *testing.T) {
		id, ok := userIDFromContext(context.WithValue(context.Background(), "userId", "abc"))
		assert.False(t, ok)
		assert.Equal(t, int64(0), id)
	})
}

func TestForceLogoutUserLogic_NoActiveSessions(t *testing.T) {
	db := setupSecurityLogicTestDB(t)
	user := &model.User{Username: "no-session", Password: "hashed", Phone: "13800138000", Status: "active"}
	assert.NoError(t, db.Create(user).Error)

	logic := NewForceLogoutUserLogic(context.Background(), &svc.ServiceContext{DB: db})
	resp, err := logic.ForceLogoutUser(user.Id, &types.ForceLogoutReq{Reason: "manual"})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "用户当前无活跃会话", resp.Message)
}

func TestForceLogoutUserLogic_ErrorBranches(t *testing.T) {
	t.Run("invalid user id", func(t *testing.T) {
		logic := NewForceLogoutUserLogic(context.Background(), &svc.ServiceContext{})
		resp, err := logic.ForceLogoutUser(0, nil)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "用户ID无效")
	})

	t.Run("db not initialized", func(t *testing.T) {
		logic := NewForceLogoutUserLogic(context.Background(), &svc.ServiceContext{})
		resp, err := logic.ForceLogoutUser(1, nil)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "数据库未初始化")
	})

	t.Run("user not found", func(t *testing.T) {
		db := setupSecurityLogicTestDB(t)
		logic := NewForceLogoutUserLogic(context.Background(), &svc.ServiceContext{DB: db})
		resp, err := logic.ForceLogoutUser(999, nil)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "用户不存在")
	})
}

func TestRevokeSessionLogic_Branches(t *testing.T) {
	t.Run("invalid session id", func(t *testing.T) {
		logic := NewRevokeSessionLogic(context.Background(), &svc.ServiceContext{})
		resp, err := logic.RevokeSession("")
		assert.Nil(t, resp)
		assert.EqualError(t, err, "会话ID无效")
	})

	t.Run("db not initialized", func(t *testing.T) {
		logic := NewRevokeSessionLogic(context.Background(), &svc.ServiceContext{})
		resp, err := logic.RevokeSession("sid")
		assert.Nil(t, resp)
		assert.EqualError(t, err, "数据库未初始化")
	})

	t.Run("session not found", func(t *testing.T) {
		db := setupSecurityLogicTestDB(t)
		logic := NewRevokeSessionLogic(context.Background(), &svc.ServiceContext{DB: db})
		resp, err := logic.RevokeSession("missing")
		assert.Nil(t, resp)
		assert.EqualError(t, err, "会话不存在")
	})

	t.Run("already revoked", func(t *testing.T) {
		db := setupSecurityLogicTestDB(t)
		session := &model.UserSession{ID: "revoked-1", UserID: 1, ClientIP: "127.0.0.1", UserAgent: "ua", Status: "revoked"}
		assert.NoError(t, db.Create(session).Error)

		logic := NewRevokeSessionLogic(context.Background(), &svc.ServiceContext{DB: db})
		resp, err := logic.RevokeSession("revoked-1")
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "会话已吊销", resp.Message)
	})
}

func TestHandleSecurityEventLogic_Branches(t *testing.T) {
	t.Run("invalid event id", func(t *testing.T) {
		logic := NewHandleSecurityEventLogic(context.Background(), &svc.ServiceContext{})
		resp, err := logic.HandleSecurityEvent(0, nil)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "事件ID无效")
	})

	t.Run("db not initialized", func(t *testing.T) {
		logic := NewHandleSecurityEventLogic(context.Background(), &svc.ServiceContext{})
		resp, err := logic.HandleSecurityEvent(1, nil)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "数据库未初始化")
	})

	t.Run("event not found", func(t *testing.T) {
		db := setupSecurityLogicTestDB(t)
		logic := NewHandleSecurityEventLogic(context.Background(), &svc.ServiceContext{DB: db})
		resp, err := logic.HandleSecurityEvent(999, nil)
		assert.Nil(t, resp)
		assert.EqualError(t, err, "安全事件不存在")
	})

	t.Run("already handled", func(t *testing.T) {
		db := setupSecurityLogicTestDB(t)
		event := &model.SecurityEvent{EventType: "login_failed", Severity: "medium", ClientIP: "127.0.0.1", Description: "desc", Handled: true}
		assert.NoError(t, db.Create(event).Error)

		logic := NewHandleSecurityEventLogic(context.Background(), &svc.ServiceContext{DB: db})
		resp, err := logic.HandleSecurityEvent(event.ID, nil)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "事件已处理", resp.Message)
	})

	t.Run("append note and set handler", func(t *testing.T) {
		db := setupSecurityLogicTestDB(t)
		event := &model.SecurityEvent{EventType: "login_failed", Severity: "high", ClientIP: "127.0.0.1", Description: "desc", Details: "initial"}
		assert.NoError(t, db.Create(event).Error)

		ctx := context.WithValue(context.Background(), "userId", int64(42))
		logic := NewHandleSecurityEventLogic(ctx, &svc.ServiceContext{DB: db})
		resp, err := logic.HandleSecurityEvent(event.ID, &types.HandleSecurityEventReq{Note: "processed"})
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "处理成功", resp.Message)

		var saved model.SecurityEvent
		assert.NoError(t, db.First(&saved, event.ID).Error)
		assert.True(t, saved.Handled)
		assert.NotNil(t, saved.HandledBy)
		assert.Equal(t, int64(42), *saved.HandledBy)
		assert.Equal(t, "initial\nprocessed", saved.Details)
	})
}

func TestUpdatePasswordPolicyLogic_ErrorBranches(t *testing.T) {
	logic := NewUpdatePasswordPolicyLogic(context.Background(), &svc.ServiceContext{})

	resp, err := logic.UpdatePasswordPolicy(nil)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "请求参数不能为空")

	resp, err = logic.UpdatePasswordPolicy(&types.UpdatePasswordPolicyReq{MinLength: 5})
	assert.Nil(t, resp)
	assert.EqualError(t, err, "最小密码长度必须在6到50之间")

	resp, err = logic.UpdatePasswordPolicy(validPasswordPolicyReq())
	assert.Nil(t, resp)
	assert.EqualError(t, err, "数据库未初始化")
}

func TestCheckPasswordStrengthLogic_ErrorBranches(t *testing.T) {
	logic := NewCheckPasswordStrengthLogic(context.Background(), &svc.ServiceContext{})

	resp, err := logic.CheckPasswordStrength(nil)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "请求参数不能为空")

	resp, err = logic.CheckPasswordStrength(&types.ChangePasswordReq{NewPassword: "   "})
	assert.Nil(t, resp)
	assert.EqualError(t, err, "新密码不能为空")
}

func TestGetPasswordPolicyLogic_FromDB(t *testing.T) {
	db := setupSecurityLogicTestDB(t)
	policy := &model.PasswordPolicy{
		MinLength:             12,
		RequireUppercase:      true,
		RequireLowercase:      true,
		RequireNumbers:        true,
		RequireSpecialChars:   false,
		MaxAge:                180,
		HistoryCount:          8,
		MaxLoginAttempts:      6,
		LockoutDuration:       60,
		SessionTimeout:        240,
		MaxConcurrentSessions: 2,
	}
	assert.NoError(t, db.Create(policy).Error)

	logic := NewGetPasswordPolicyLogic(context.Background(), &svc.ServiceContext{DB: db})
	resp, err := logic.GetPasswordPolicy()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, policy.ID, resp.Id)
	assert.Equal(t, 12, resp.MinLength)
	assert.Equal(t, 240, resp.SessionTimeout)
	assert.Equal(t, 2, resp.MaxConcurrentSessions)
}

func TestGetAuditLogsLogic_WithData(t *testing.T) {
	db := setupSecurityLogicTestDB(t)

	userID := int64(101)
	log1 := &model.AuditLog{UserID: &userID, Username: "u1", Action: "a1", Resource: "r1", ResourceID: "1", Status: "success"}
	log2 := &model.AuditLog{UserID: nil, Username: "u2", Action: "a2", Resource: "r2", ResourceID: "2", Status: "failed"}
	assert.NoError(t, db.Create(log1).Error)
	assert.NoError(t, db.Create(log2).Error)

	logic := NewGetAuditLogsLogic(context.Background(), &svc.ServiceContext{DB: db})
	resp, err := logic.GetAuditLogs()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(2), resp.Total)
	assert.Len(t, resp.Logs, 2)
	assert.Equal(t, log2.ID, resp.Logs[0].Id)
	assert.Equal(t, int64(0), resp.Logs[0].UserId)
	assert.Equal(t, log1.ID, resp.Logs[1].Id)
	assert.Equal(t, userID, resp.Logs[1].UserId)
}

func TestGetLoginAttemptsLogic_WithData(t *testing.T) {
	db := setupSecurityLogicTestDB(t)

	userID := int64(202)
	attempt1 := &model.LoginAttempt{UserID: &userID, Username: "admin", ClientIP: "127.0.0.1", Success: true}
	attempt2 := &model.LoginAttempt{UserID: nil, Username: "ghost", ClientIP: "127.0.0.2", Success: false, FailReason: "bad password"}
	assert.NoError(t, db.Create(attempt1).Error)
	assert.NoError(t, db.Create(attempt2).Error)

	logic := NewGetLoginAttemptsLogic(context.Background(), &svc.ServiceContext{DB: db})
	resp, err := logic.GetLoginAttempts()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(2), resp.Total)
	assert.Len(t, resp.Attempts, 2)
	assert.Equal(t, attempt2.ID, resp.Attempts[0].Id)
	assert.Equal(t, int64(0), resp.Attempts[0].UserId)
	assert.Equal(t, "bad password", resp.Attempts[0].FailReason)
	assert.Equal(t, attempt1.ID, resp.Attempts[1].Id)
	assert.Equal(t, userID, resp.Attempts[1].UserId)
}

func TestGetSecurityEventsLogic_WithData(t *testing.T) {
	db := setupSecurityLogicTestDB(t)

	userID := int64(303)
	handledBy := int64(404)
	event1 := &model.SecurityEvent{EventType: "risk", Severity: "high", UserID: &userID, ClientIP: "127.0.0.1", Description: "event-1"}
	event2 := &model.SecurityEvent{EventType: "login_failed", Severity: "medium", ClientIP: "127.0.0.2", Description: "event-2", Handled: true, HandledBy: &handledBy}
	assert.NoError(t, db.Create(event1).Error)
	assert.NoError(t, db.Create(event2).Error)

	logic := NewGetSecurityEventsLogic(context.Background(), &svc.ServiceContext{DB: db})
	resp, err := logic.GetSecurityEvents()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(2), resp.Total)
	assert.Len(t, resp.Events, 2)
	assert.Equal(t, event2.ID, resp.Events[0].Id)
	assert.True(t, resp.Events[0].Handled)
	assert.Equal(t, handledBy, resp.Events[0].HandledBy)
	assert.Equal(t, event1.ID, resp.Events[1].Id)
	assert.Equal(t, userID, resp.Events[1].UserId)
}

func TestGetUserSessionsLogic_WithData(t *testing.T) {
	db := setupSecurityLogicTestDB(t)

	s1 := &model.UserSession{ID: "s1", UserID: 1, ClientIP: "127.0.0.1", UserAgent: "ua1", Status: "active"}
	s2 := &model.UserSession{ID: "s2", UserID: 2, ClientIP: "127.0.0.2", UserAgent: "ua2", Status: "revoked"}
	assert.NoError(t, db.Create(s1).Error)
	assert.NoError(t, db.Create(s2).Error)

	logic := NewGetUserSessionsLogic(context.Background(), &svc.ServiceContext{DB: db})
	resp, err := logic.GetUserSessions()
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(2), resp.Total)
	assert.Len(t, resp.Sessions, 2)
	assert.Equal(t, "s2", resp.Sessions[0].Id)
	assert.Equal(t, "s1", resp.Sessions[1].Id)
}
