package repositories

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/entity"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/repository"
)

type userRepo struct {
}

// NewTrainingRepo .
func NewUserRepo() repository.UserRepository {
	return &userRepo{}
}

func (s *userRepo) SaveUser(ctx context.Context, t *entity.User) error {
	_, err := g.DB().Model("users").Ctx(ctx).Insert(t)
	if err != nil {
		return gerror.Wrapf(err, "failed to save user %s", t.Username)
	}
	return nil
}
func (s *userRepo) DeleteUser(ctx context.Context, userId string) error {
	_, err := g.DB().Model("users").Ctx(ctx).Where("user_id", userId).Delete()
	if err != nil {
		return gerror.Wrapf(err, "failed to delete user %s", userId)
	}
	return nil
}

func (s *userRepo) FindByID(ctx context.Context, userId string) (*entity.User, error) {
	var user *entity.User
	err := g.DB().Model("users").Ctx(ctx).Where("user_id", userId).Scan(&user)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to find user %s", userId)
	}
	if user == nil {
		return nil, gerror.Newf("user not found by id %s", userId)
	}
	return user, nil
}

func (s *userRepo) FindByName(ctx context.Context, username string) (*entity.User, error) {
	var user *entity.User
	err := g.DB().Model("users").Ctx(ctx).Where("username", username).Scan(&user)
	if err != nil {
		return nil, gerror.Wrapf(err, "failed to find user %s", username)
	}
	if user == nil {
		return nil, gerror.Newf("user not found by username %s", username)
	}
	return user, nil
}
func (s *userRepo) FindByEmailOrPhone(ctx context.Context, email, phone string) (*entity.User, error) {
	var user *entity.User
	err := g.DB().Model("users").Ctx(ctx).Where("email", email).Where("mobile_phone", phone).Scan(&user)
	if err != nil {
		return nil, gerror.Wrapf(err, "user not found by email %s or phone %s", email, phone)
	}
	if user == nil {
		return nil, gerror.Newf("user not found by email %s or phone %s", email, phone)
	}
	return user, err
}

func (s *userRepo) UpdatePassword(ctx context.Context, userId, pwd string) error {
	_, err := g.DB().Model("users").Ctx(ctx).Where("user_id", userId).Update("password_hash", pwd)
	if err != nil {
		return gerror.Wrapf(err, "failed to update password for user %s", userId)
	}
	return nil
}

// RecordFailedAttempt 记录登录失败尝试
func (s *userRepo) RecordFailedAttempt(ctx context.Context, username string) (*entity.User, error) {
	maxAttempts := g.Cfg().MustGet(ctx, "password.maxAttempts").Int()
	lockDuration := g.Cfg().MustGet(ctx, "password.lockDuration").Duration()

	var user *entity.User

	// 使用事务确保原子操作
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {

		// 递增失败计数
		_, err := tx.Model("users").
			Where("username", username).
			Data(g.Map{
				"failed_attempts": gdb.Raw("failed_attempts + 1"),
			}).
			Update()
		if err != nil {
			return err
		}

		// 获取更新后的状态
		err = tx.Model("users").
			Where("username", username).
			Scan(&user)
		if err != nil {
			return err
		}

		// 检查是否需要锁定
		if user.FailedAttempts >= maxAttempts {
			_, err = tx.Model("users").
				Where("username", username).
				Data(g.Map{
					"is_locked":    true,
					"locked_until": time.Now().Add(lockDuration),
				}).
				Update()
			if err != nil {
				return err
			}
			user.IsLocked = true
			user.LockedUntil = time.Now().Add(lockDuration)
		}

		return nil
	})

	if err != nil {
		return nil, gerror.Wrapf(err, "failed to record attempt for user %s", username)
	}

	return user, nil
}

// ResetFailedAttempts 重置失败计数
func (s *userRepo) ResetFailedAttempts(ctx context.Context, userId string) error {
	_, err := g.DB().Model("users").Ctx(ctx).
		Where("user_id", userId).
		Data(g.Map{
			"failed_attempts": 0,
			"is_locked":       false,
			"locked_until":    nil,
		}).
		Update()

	if err != nil {
		return gerror.Wrapf(err, "failed to reset attempts for user %s", userId)
	}
	return nil
}

// GetFailedAttempts 获取当前失败尝试次数
func (s *userRepo) GetFailedAttempts(ctx context.Context, username string) (int, error) {
	var attempts int
	_, err := g.DB().Model("users").
		Ctx(ctx).
		Where("username", username).
		Value("failed_attempts", &attempts)

	if err != nil {
		return 0, gerror.Wrapf(err, "failed to get attempts for user %s", username)
	}

	return attempts, nil
}
