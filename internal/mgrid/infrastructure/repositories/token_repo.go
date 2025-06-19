package repositories

import (
	"context"
	"time"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/entity"
	"github.com/shiqinfeng1/goframe-ddd/internal/mgrid/domain/repository"
	"github.com/shiqinfeng1/goframe-ddd/pkg/cache"
)

type tokenRepo struct {
	labelCache *gcache.Cache
}

// NewTrainingRepo .
func NewTokenRepo() repository.TokenRepository {
	return &tokenRepo{
		labelCache: cache.KV(),
	}
}

func (s *tokenRepo) SaveRefrsehToken(ctx context.Context, t *entity.Token) error {
	_, err := g.DB().Model("user_tokens").Ctx(ctx).OnConflict("user_id").Save(t)
	if err != nil {
		return gerror.Wrapf(err, "insert user_token fail: userid=%v refreshid=%v", t.UserID, t.RefreshID)
	}
	return nil
}
func (s *tokenRepo) DeleteRefreshToken(ctx context.Context, refreshID string) error {
	_, err := g.DB().Model("user_tokens").Ctx(ctx).Where("refresh_id", refreshID).Delete()
	if err != nil {
		return gerror.Wrapf(err, "delete user_token fail: refreshid=%v", refreshID)
	}
	return nil
}

func (s *tokenRepo) FindByRefreshToken(ctx context.Context, userID, refreshID string) (*entity.Token, error) {
	var token *entity.Token
	err := g.DB().Model("user_tokens").Ctx(ctx).Where("refresh_id", refreshID).Scan(&token)
	if err != nil {
		return nil, gerror.Wrapf(err, "find user_token by userid and refreshid fail: userid=%v refreshid=%v", userID, refreshID)
	}
	return token, nil
}

func (s *tokenRepo) SaveVerifyCode(ctx context.Context, userId, verifyCode string, expired time.Duration) error {
	// 注意： code作为key， userid作为value
	err := s.labelCache.Set(ctx, verifyCode, userId, expired)
	if err != nil {
		return gerror.Wrapf(err, "save verify code fail: userid=%v verifycode=%v", userId, verifyCode)
	}
	return nil
}

func (s *tokenRepo) GetUserIdByVerifyCode(ctx context.Context, verifyCode string) string {
	v, _ := s.labelCache.Get(ctx, verifyCode)
	return v.String()
}
func (s *tokenRepo) DeleteVerifyCode(ctx context.Context, verifyCode string) error {
	_, err := s.labelCache.Remove(ctx, verifyCode)
	if err != nil {
		return gerror.Wrapf(err, "delete verify code fail: verifycode=%v", verifyCode)
	}
	return nil
}
