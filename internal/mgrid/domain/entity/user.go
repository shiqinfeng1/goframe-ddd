package entity

// user.go (聚合根)
type User struct {
	UserID         string `json:"user_id" orm:"user_id,primary" dc:"用户ID"`
	Username       string `json:"username" orm:"username,notnull,unique" dc:"用户名"`
	Email          string `json:"email" orm:"email,notnull,unique" dc:"邮箱"`
	MobilePhone    string `json:"mobile_phone" orm:"mobile_phone,notnull,unique" dc:"手机号"`
	PasswordHash   string `json:"-" orm:"password_hash" dc:"密码哈希"`
	IsLocked       bool   `json:"is_locked" orm:"is_locked"   dc:"是否锁定"`
	LockedUntil    string `json:"locked_until" orm:"locked_until,datetime" dc:"锁定截止时间"`
	FailedAttempts int    `json:"-" orm:"failed_attempts" dc:"失败尝试次数"`
}

func NewUser() *User {
	return &User{}
}
