package entity

import "time"

// user.go (聚合根)
type User struct {
	UserID         string    `json:"user_id" orm:"id,primary" description:"用户ID"`
	Username       string    `json:"username" orm:"username" description:"用户名"`
	Email          string    `json:"email" orm:"email" description:"邮箱"`
	MobilePhone    string    `json:"mobile_phone" orm:"mobile_phone" description:"手机号"`
	PasswordHash   string    `json:"-" orm:"password_hash" description:"密码哈希"`
	IsLocked       bool      `json:"is_locked" orm:"is_locked"   description:"是否锁定"`
	LockedUntil    time.Time `json:"locked_until" orm:"locked_until" description:"锁定截止时间"`
	FailedAttempts int       `json:"-" orm:"failed_attempts" description:"失败尝试次数"`
}

func NewUser() *User {
	return &User{}
}
