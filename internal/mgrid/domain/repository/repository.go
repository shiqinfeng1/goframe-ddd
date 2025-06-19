package repository

import (
	"context"
	"time"

	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/entity"
)

type PointdataRepository interface {
}
type UserRepository interface {
	SaveUser(ctx context.Context, user *entity.User) error
	FindByName(ctx context.Context, username string) (*entity.User, error)
	FindByEmailOrPhone(ctx context.Context, email, phone string) (*entity.User, error)
	UpdatePassword(ctx context.Context, userId, pwd string) error
	RecordFailedAttempt(ctx context.Context, username string, maxAttempts int, lockDuration time.Duration) (*entity.User, error)
	ResetFailedAttempts(ctx context.Context, userId string) error
	GetFailedAttempts(ctx context.Context, username string) (int, error)
	UserIsExisted(ctx context.Context, name, mobile, email string) (bool, error)
}

type TokenRepository interface {
	SaveRefrsehToken(ctx context.Context, token *entity.Token) error
	DeleteRefreshToken(ctx context.Context, refreshID string) error
	FindByRefreshToken(ctx context.Context, userID string, refreshID string) (*entity.Token, error)
	SaveVerifyCode(ctx context.Context, userId, verifyCode string, expired time.Duration) error
	GetUserIdByVerifyCode(ctx context.Context, verifyCode string) string
	DeleteVerifyCode(ctx context.Context, verifyCode string) error
}
