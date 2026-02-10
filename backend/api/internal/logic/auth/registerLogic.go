// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package auth

import (
	"context"
	"errors"

	"dmh/api/internal/middleware"
	"dmh/api/internal/svc"
	"dmh/api/internal/types"
	"dmh/model"
	"golang.org/x/crypto/bcrypt"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.LoginResp, err error) {
	if req.Username == "" {
		return nil, errors.New("用户名不能为空")
	}
	if req.Password == "" {
		return nil, errors.New("密码不能为空")
	}
	if req.Phone == "" {
		return nil, errors.New("手机号不能为空")
	}

	var existingUser model.User
	err = l.svcCtx.DB.Where("username = ?", req.Username).First(&existingUser).Error
	if err == nil {
		return nil, errors.New("用户名已存在")
	}

	err = l.svcCtx.DB.Where("phone = ?", req.Phone).First(&existingUser).Error
	if err == nil {
		return nil, errors.New("手机号已被注册")
	}

	if len(req.Password) < 6 {
		return nil, errors.New("密码长度不能少于6位")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		l.Errorf("密码加密失败: %v", err)
		return nil, errors.New("注册失败")
	}

	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Email:    req.Email,
		RealName: req.RealName,
		Status:   "active",
	}

	err = l.svcCtx.DB.Create(&user).Error
	if err != nil {
		l.Errorf("创建用户失败: %v", err)
		return nil, errors.New("注册失败")
	}

	var participantRole model.Role
	err = l.svcCtx.DB.Where("code = ?", "participant").First(&participantRole).Error
	if err != nil {
		l.Errorf("查询participant角色失败: %v", err)
		participantRole = model.Role{
			Code: "participant",
			Name: "参与者",
		}
		err = l.svcCtx.DB.Create(&participantRole).Error
		if err != nil {
			l.Errorf("创建participant角色失败: %v", err)
		}
	}

	userRole := model.UserRole{
		UserID: user.Id,
		RoleID: participantRole.ID,
	}
	err = l.svcCtx.DB.Create(&userRole).Error
	if err != nil {
		l.Errorf("分配用户角色失败: %v", err)
	}

	authMiddleware := middleware.NewAuthMiddleware(l.svcCtx.Config.Auth.AccessSecret)
	token, err := authMiddleware.GenerateToken(user.Id, user.Username, []string{"participant"}, []int64{})
	if err != nil {
		l.Errorf("生成token失败: %v", err)
		return nil, errors.New("注册失败")
	}

	resp = &types.LoginResp{
		Token:    token,
		UserId:   user.Id,
		Username: user.Username,
		Phone:    user.Phone,
		RealName: user.RealName,
		Roles:    []string{"participant"},
		BrandIds: []int64{},
	}

	return resp, nil
}
