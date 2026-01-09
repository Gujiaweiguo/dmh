// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"dmh/api/internal/config"
	"dmh/api/internal/service"
	
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config          config.Config
	DB              *gorm.DB
	PasswordService *service.PasswordService
	AuditService    *service.AuditService
	SessionService  *service.SessionService
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化GORM数据库连接
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		logx.Errorf("连接数据库失败: %v", err)
		// 不中断服务，继续运行（使用mock数据）
	} else {
		// 测试连接
		sqlDB, err := db.DB()
		if err != nil {
			logx.Errorf("获取数据库实例失败: %v", err)
		} else if err := sqlDB.Ping(); err != nil {
			logx.Errorf("数据库ping失败: %v", err)
		} else {
			logx.Info("数据库连接成功")
		}
	}
	
	// 初始化安全服务
	passwordService := service.NewPasswordService(db)
	auditService := service.NewAuditService(db)
	sessionService := service.NewSessionService(db, passwordService)
	
	return &ServiceContext{
		Config:          c,
		DB:              db,
		PasswordService: passwordService,
		AuditService:    auditService,
		SessionService:  sessionService,
	}
}