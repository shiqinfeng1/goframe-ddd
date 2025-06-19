package entity

import (
	"context"
	"errors"
	"unicode"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/crypto/bcrypt"
)

type Password struct{}

func NewPassword() *Password {
	return &Password{}
}

// HashPassword 生成密码哈希
func (s *Password) HashPassword(ctx context.Context, plainPassword string) (string, error) {
	// 密码强度校验
	if err := s.validatePasswordStrength(ctx, plainPassword); err != nil {
		return "", err
	}

	// 生成哈希
	hashCost := g.Cfg().MustGet(ctx, "password.hashCost").Int()
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), hashCost)
	if err != nil {
		return "", gerror.Wrap(err, "failed to hash password")
	}

	return string(hashedBytes), nil
}

// VerifyPassword 验证密码
func (s *Password) VerifyPassword(ctx context.Context, hashedPassword, plainPassword string) (bool, error) {
	// 空密码检查
	if hashedPassword == "" || plainPassword == "" {
		return false, gerror.New("password cannot be empty")
	}

	// 使用 bcrypt 比较哈希
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, gerror.Wrap(err, "failed to compare passwords")
	}

	return true, nil
}

// validatePasswordStrength 密码强度验证
func (s *Password) validatePasswordStrength(ctx context.Context, password string) error {
	// 长度检查
	minLen := g.Cfg().MustGet(ctx, "password.minLength").Int()
	maxLen := g.Cfg().MustGet(ctx, "password.maxLength").Int()
	if len(password) < minLen || len(password) > maxLen {
		return gerror.Newf("password length must be between %d and %d characters", minLen, maxLen)
	}

	// 复杂度检查
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if g.Cfg().MustGet(ctx, "password.requireMixed").Bool() && !(hasUpper && hasLower) {
		return gerror.New("password must contain both uppercase and lowercase letters")
	}

	if g.Cfg().MustGet(ctx, "password.requireNumber").Bool() && !hasNumber {
		return gerror.New("password must contain at least one number")
	}

	if g.Cfg().MustGet(ctx, "password.requireSpecial").Bool() && !hasSpecial {
		return gerror.New("password must contain at least one special character")
	}

	// 常见弱密码检查
	if s.isCommonPassword(ctx, password) {
		return gerror.New("password is too common")
	}

	return nil
}

// isCommonPassword 检查是否为常见弱密码
func (s *Password) isCommonPassword(ctx context.Context, password string) bool {
	// commonPasswords := []string{
	// 	"password", "123456", "12345678", "11111111", "88888888", "qwerty", "admin", "admin123", "welcome",
	// 	// 可以扩展更多常见弱密码
	// }
	commonPasswords := g.Cfg().MustGet(ctx, "commonPassword").Strings()
	for _, p := range commonPasswords {
		if password == p {
			return true
		}
	}

	// 检查连续字符
	// if match, _ := regexp.MatchString(`(?:0123|1234|2345|3456|4567|5678|6789|7890)`, password); match {
	// 	return true
	// }

	return false
}
