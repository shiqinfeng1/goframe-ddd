package v1

import "github.com/gogf/gf/v2/frame/g"

type SendVerifyCodeReq struct {
	g.Meta      `path:"/auth/password/sendVerifyCode" tags:"认证管理" method:"post" summary:"发送验证码"`
	Email       string `p:"email"  dc:"邮箱"`
	MobilePhone string `p:"mobile_phone"  dc:"手机号"`
}
type SendVerifyCodeRes struct {
	g.Meta `status:"200"`
}

type ResetPasswordReq struct {
	g.Meta      `path:"/auth/password/reset" tags:"认证管理" method:"post" summary:"重置密码请求"`
	VerifyCode  string `p:"verify_code" v:"required" dc:"验证码"`
	NewPassword string `p:"new_password" v:"required" dc:"新密码"`
}
type ResetPasswordRes struct {
	g.Meta `status:"200"`
}
type RegisterUserReq struct {
	g.Meta      `path:"/auth/user/register" tags:"认证管理" method:"post" summary:"注册用户"`
	Username    string `p:"username" v:"required" dc:"用户名"`
	Email       string `p:"email" v:"required" dc:"邮箱"`
	MobilePhone string `p:"mobile_phone" v:"required" dc:"手机号"`
	Password    string `p:"password" v:"required" dc:"密码"`
}
type RegisterUserRes struct {
	g.Meta `status:"200"`
}

type LoginReq struct {
	g.Meta   `path:"/auth/user/login" tags:"认证管理" method:"post" summary:"登录"`
	Username string `p:"username" v:"required" dc:"用户名"`
	Password string `p:"password" v:"required" dc:"密码"`
}
type LoginRes struct {
	g.Meta       `status:"200"`
	AccessToken  string `json:"access_token" dc:"访问token"`
	RefreshToken string `json:"refresh_token" dc:"刷新token"`
}

type LogoutReq struct {
	g.Meta        `path:"/auth/user/logout" tags:"认证管理" method:"post" summary:"登出"`
	Authorization string `p:"Authorization" in:"header" v:"required" dc:"Bearer Token"` // 从Header获取
}
type LogoutRes struct {
	g.Meta `status:"200"`
}
