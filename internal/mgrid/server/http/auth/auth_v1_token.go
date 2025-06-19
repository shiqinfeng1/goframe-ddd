package auth

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
	v1 "github.com/shiqinfeng1/goframe-ddd/api/mgrid/http/auth/v1"
	"github.com/shiqinfeng1/goframe-ddd/pkg/errors"
)

func (c *ControllerV1) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error) {
	lang := ghttp.RequestFromCtx(ctx).GetCtxVar("lang").String()

	newTokens, err := c.app.Auth().RefreshToken(ctx)
	if err != nil {
		c.logger.Error(ctx, err)
		return nil, errors.ErrAuthRefreshTokenFail(lang)
	}
	return &v1.RefreshTokenRes{
		AccessToken:  newTokens.AccessToken,
		RefreshToken: newTokens.RefreshToken,
	}, nil
}
