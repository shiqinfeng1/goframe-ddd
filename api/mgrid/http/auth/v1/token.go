package v1

import "github.com/gogf/gf/v2/frame/g"

type RefreshTokenReq struct {
	g.Meta `path:"/token/refresh" tags:"认证管理" method:"post" summary:"刷新token"`
}

type RefreshTokenRes struct {
	g.Meta       `status:"200"`
	AccessToken  string `json:"access_token" dc:"访问token"`
	RefreshToken string `json:"refresh_token" dc:"刷新token"`
}
