package service

import (
	"context"
	"net/http"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/rs/xid"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/dto"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/entity"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/repository"
)

const (
	REFRESH_TOKEN = "refresh_token"
)

// auth_service.go
type authService struct {
	logger    application.Logger
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewAuthService(ctx context.Context, logger application.Logger, userRepo repository.UserRepository, tokenRepo repository.TokenRepository) application.AuthService {
	as := &authService{
		logger:    logger,
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}

	return as
}

// RequestSendVerifyCode 请求密码重置
func (s *authService) RequestSendVerifyCode(ctx context.Context, email, mobilePhone string) error {
	// 1. 验证用户存在
	user, err := s.userRepo.FindByEmailOrPhone(ctx, email, mobilePhone)
	if err != nil {
		return err
	}
	// 2. 生成验证码
	verifyCode := grand.Digits(6)
	// 3. 存储重置记录
	if err := s.tokenRepo.SaveVerifyCode(ctx, user.UserID, verifyCode, 5*time.Minute); err != nil {
		return err
	}
	// 4. 发送重置邮件或短信验证码
	// if user.Email != "" {
	// 	go s.sendResetEmail(user.Email, verifyCode)
	// }
	// if user.MobilePhone != "" {
	// 	go s.sendResetMobile(user.MobilePhone, verifyCode)
	// }

	return nil
}

// ResetPassword 重置密码
func (s *authService) ResetPassword(ctx context.Context, verifyCode, newPassword string) error {
	// 1. 验证重置令牌
	userId := s.tokenRepo.GetUserIdByVerifyCode(ctx, verifyCode)
	if userId == "" {
		return gerror.New("failed to save verify code")
	}

	// 2. 哈希新密码
	hashedPassword, err := entity.NewPassword().HashPassword(ctx, newPassword)
	if err != nil {
		return err
	}

	// 3. 更新用户密码
	if err := s.userRepo.UpdatePassword(ctx, userId, hashedPassword); err != nil {
		return err
	}

	// 4. 删除重置记录
	if err := s.tokenRepo.DeleteVerifyCode(ctx, verifyCode); err != nil {
		return err
	}
	return nil
}

// VerifyCredentials 验证用户凭证
func (s *authService) verifyCredentials(ctx context.Context, username, plainPassword string) (*entity.User, error) {
	// 查询用户
	user, err := s.userRepo.FindByName(ctx, username)
	if err != nil {
		return nil, err
	}
	// 2. 检查账户锁定状态
	if user.IsLocked {
		if user.LockedUntil.After(time.Now()) {
			return nil, gerror.Newf("account is locked until %s", user.LockedUntil.Format("2006-01-02 15:04"))
		}
		// 自动解锁过期锁定
		if err := s.userRepo.ResetFailedAttempts(ctx, user.UserID); err != nil {
			return nil, err
		}
	}
	// 验证密码
	match, err := entity.NewPassword().VerifyPassword(ctx, user.PasswordHash, plainPassword)
	if err != nil {
		return nil, err
	}
	if !match {
		// 记录失败尝试
		user2, err := s.userRepo.RecordFailedAttempt(ctx, username)
		if err != nil {
			return nil, err
		}
		if user2.IsLocked {
			return nil, gerror.Newf("account locked due to too many failed attempts")
		}
		maxAttempts := g.Cfg().MustGet(ctx, "password.maxAttempts").Int()
		return nil, gerror.Newf("invalid password (%d/%d attempts remaining)",
			maxAttempts-user2.FailedAttempts, maxAttempts)
	}
	// 4. 密码正确，重置失败计数
	if err := s.userRepo.ResetFailedAttempts(ctx, user.UserID); err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser 创建用户
func (s *authService) CreateUser(ctx context.Context, in *dto.CreateUserIn) error {
	// 密码哈希
	hashedPassword, err := entity.NewPassword().HashPassword(ctx, in.Password)
	if err != nil {
		return err
	}

	// 创建用户记录
	user := &entity.User{
		UserID:       xid.New().String(),
		Username:     in.Username,
		Email:        in.Email,
		MobilePhone:  in.MobilePhone,
		PasswordHash: hashedPassword,
	}

	if s.userRepo.SaveUser(ctx, user) != nil {
		return err
	}

	return nil
}

func (s *authService) Login(ctx context.Context, username, password string) (*dto.Token, error) {
	user, err := s.verifyCredentials(ctx, username, password)
	if err != nil {
		return nil, err
	}
	// 2. 生成Token对
	t := entity.NewToken()
	at, rt, err := t.GenerateTokenPair(ctx, user.UserID)
	if err != nil {
		return nil, err
	}
	// 存储Refresh Token记录 (实现轮换机制)
	if err := s.tokenRepo.SaveRefrsehToken(ctx, t); err != nil {
		return nil, err
	}
	// 3. 设置Refresh Token到Cookie
	ghttp.RequestFromCtx(ctx).Cookie.SetCookie(
		REFRESH_TOKEN, rt, "", "/",
		g.Cfg().MustGet(ctx, "jwt.refreshExpire").Duration(),
		ghttp.CookieOptions{
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	return &dto.Token{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

// 刷新Token
func (s *authService) RefreshToken(ctx context.Context, oldRefreshToken string) (*dto.Token, error) {

	// 验证旧Refresh Token
	token := entity.NewToken()
	err := token.ParseJWT(
		oldRefreshToken,
		g.Cfg().MustGet(ctx, "jwt.refreshSecret").String(),
	)
	if err != nil {
		return nil, err
	}

	// 检查Refresh Token是否有效
	tokenRecord, err := s.tokenRepo.FindByRefreshToken(ctx, token.UserID, token.RefreshID)
	if err != nil {
		return nil, err
	}
	if tokenRecord == nil || tokenRecord.RefreshID == "" {
		return nil, gerror.New("invalid refresh token")
	}

	// 使旧Token失效
	if err := s.tokenRepo.DeleteRefreshToken(ctx, token.RefreshID); err != nil {
		return nil, err
	}
	// 生成新Token对
	at, rt, err := token.GenerateTokenPair(ctx, token.UserID)
	if err != nil {
		return nil, err
	}
	// 存储Refresh Token记录 (实现轮换机制)
	if err := s.tokenRepo.SaveRefrsehToken(ctx, token); err != nil {
		return nil, err
	}
	// 3. 设置Refresh Token到Cookie
	ghttp.RequestFromCtx(ctx).Cookie.SetCookie(
		REFRESH_TOKEN, rt,
		"", "/",
		g.Cfg().MustGet(ctx, "jwt.refreshExpire").Duration(),
		ghttp.CookieOptions{
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
		})
	return &dto.Token{
		AccessToken:  at,
		RefreshToken: rt,
	}, nil
}

// Logout 登出接口
func (s *authService) Logout(ctx context.Context) (err error) {
	// 1. 获取当前用户的Refresh Token
	r := ghttp.RequestFromCtx(ctx)
	refreshToken := r.Cookie.Get(REFRESH_TOKEN).String()

	// 2. 使Token失效
	if refreshToken != "" {
		t := entity.NewToken()
		err := t.ParseJWT(refreshToken, g.Cfg().MustGet(ctx, "jwt.refreshSecret").String())
		if err == nil {
			// 使旧Token失效
			if err := s.tokenRepo.DeleteRefreshToken(ctx, t.RefreshID); err != nil {
				return err
			}
		}
	}
	// 3. 清除Cookie
	r.Cookie.Remove(REFRESH_TOKEN)
	return nil
}

func Auth(r *ghttp.Request) {
	// 1. 获取Access Token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		r.Response.WriteStatusExit(401, "authorization header missing")
		return
	}

	tokenString := authHeader[len("Bearer "):]
	if tokenString == "" {
		r.Response.WriteStatusExit(401, "invalid token format")
		return
	}

	// 2. 验证Token
	t := entity.NewToken()
	err := t.ParseJWT(
		tokenString,
		g.Cfg().MustGet(r.Context(), "jwt.accessSecret").String(),
	)
	if err != nil {
		r.Response.WriteStatusExit(401, "invalid or expired token")
		return
	}

	// 3. 将用户信息存入上下文
	r.SetCtxVar("userId", t.UserID)

	r.Middleware.Next()
}
