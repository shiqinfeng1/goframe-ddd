package auth

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	v1 "github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/auth/v1"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/application/dto"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
)

func (c *ControllerV1) SendVerifyCode(ctx context.Context, req *v1.SendVerifyCodeReq) (res *v1.SendVerifyCodeRes, err error) {
	lang := ghttp.RequestFromCtx(ctx).GetCtxVar("lang").String()
	err = c.app.Auth().RequestSendVerifyCode(ctx, req.Email, req.MobilePhone)
	if err != nil {
		c.logger.Error(ctx, err)
		return nil, errors.ErrSendVCodeFail(lang)
	}
	return &v1.SendVerifyCodeRes{}, nil
}
func (c *ControllerV1) ResetPassword(ctx context.Context, req *v1.ResetPasswordReq) (res *v1.ResetPasswordRes, err error) {
	lang := ghttp.RequestFromCtx(ctx).GetCtxVar("lang").String()
	err = c.app.Auth().ResetPassword(ctx, req.VerifyCode, req.NewPassword)
	if err != nil {
		c.logger.Error(ctx, err)
		return nil, errors.ErrRestPwdFail(lang)
	}
	return &v1.ResetPasswordRes{}, nil
}
func (c *ControllerV1) RegisterUser(ctx context.Context, req *v1.RegisterUserReq) (res *v1.RegisterUserRes, err error) {
	lang := ghttp.RequestFromCtx(ctx).GetCtxVar("lang").String()
	var exist bool
	exist, err = c.app.Auth().UserIsExisted(ctx, req.Username, req.MobilePhone, req.Email)
	if err != nil {
		c.logger.Error(ctx, err)
		return nil, errors.ErrRegisterUserFail(lang)
	}
	if exist {
		return nil, errors.ErrUserExisted(lang)
	}

	in := &dto.CreateUserIn{
		Username:    req.Username,
		Email:       req.Email,
		MobilePhone: req.MobilePhone,
		Password:    req.Password,
	}
	err = c.app.Auth().CreateUser(ctx, in)
	if err != nil {
		c.logger.Error(ctx, err)
		return nil, errors.ErrRegisterUserFail(lang)
	}
	return &v1.RegisterUserRes{}, nil
}
func (c *ControllerV1) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	lang := ghttp.RequestFromCtx(ctx).GetCtxVar("lang").String()
	user, err := c.app.Auth().VerifyCredentials(ctx, lang, req.Username, req.Password)
	if err != nil {
		code := gerror.Code(err).Code()
		if code == errors.CodeUserIsLockdBefore ||
			code == errors.CodeUserIsLockdTooManyAttempts ||
			code == errors.CodeUserVerifyAttemptsRemain {
			return nil, err
		}
		c.logger.Error(ctx, err)
		return nil, errors.ErrLoginFail(lang)
	}

	tokens, err := c.app.Auth().Login(ctx, user)
	if err != nil {
		c.logger.Error(ctx, err)
		return nil, errors.ErrLoginFail(lang)
	}
	return &v1.LoginRes{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}
func (c *ControllerV1) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	lang := ghttp.RequestFromCtx(ctx).GetCtxVar("lang").String()
	err = c.app.Auth().Logout(ctx)
	if err != nil {
		c.logger.Error(ctx, err)
		return nil, errors.ErrLogoutFail(lang)
	}
	return &v1.LogoutRes{}, nil
}
