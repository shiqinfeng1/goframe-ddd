package v1

import "github.com/gogf/gf/v2/frame/g"

type RefreshTokenReq struct {
	g.Meta       `path:"/auth/token/refresh" tags:"认证管理" method:"post" summary:"刷新token"`
	RefreshToken string `p:"refresh_token" v:"required" in:"cookie" dc:"刷新令牌"` // 从Cookie获取
}

type RefreshTokenRes struct {
	g.Meta       `status:"200"`
	AccessToken  string `json:"access_token" dc:"访问token"`
	RefreshToken string `json:"refresh_token" dc:"刷新token"`
}
