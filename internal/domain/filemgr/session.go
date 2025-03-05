package filemgr

import (
	"context"
	"errors"
	"io"
	"sync"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/xtaci/smux"
)

type SessionMgr struct {
	serverSess *gcache.Cache
}

var (
	sessionMgr *SessionMgr
	once       sync.Once
)

func Session() *SessionMgr {
	once.Do(func() {
		sessionMgr = &SessionMgr{
			serverSess: gcache.NewWithAdapter(gcache.NewAdapterMemory()),
		}
	})
	return sessionMgr
}

func (s *SessionMgr) GetSessionList(ctx context.Context) ([]string, error) {
	ids, err := s.serverSess.Keys(ctx)
	if err != nil {
		return []string{}, gerror.Wrap(err, "get sesson clientIds fail")
	}
	return gconv.Strings(ids), nil
}

func (s *SessionMgr) SaveSession(ctx context.Context, clientId string, sess *smux.Session) error {
	// 检查之前是否有会话
	if exist, _ := s.serverSess.Contains(ctx, clientId); exist {
		old, err := s.serverSess.Remove(ctx, clientId)
		if err != nil {
			return gerror.Wrapf(err, "remove old session fail. clientId=%v", clientId)
		}
		var s *smux.Session
		if err := old.Scan(&s); err != nil {
			return gerror.Wrapf(err, "scan old session fail. clientId=%v", clientId)
		}
		if err := s.Close(); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				g.Log().Warningf(ctx, "delete old session of clientid=%v:%v", clientId, err)
			} else {
				g.Log().Warningf(ctx, "delete old session of clientid=%v fail:%v", clientId, err)
			}
		}
	}

	if err := s.serverSess.Set(ctx, clientId, sess, 0); err != nil {
		return gerror.Wrapf(err, "set new session fail. clientId=%v", clientId)
	}
	return nil
}

func (s *SessionMgr) GetSession(ctx context.Context, clientId string) (*smux.Session, error) {
	item, err := s.serverSess.Get(ctx, clientId)
	if err != nil {
		return nil, gerror.Wrapf(err, "get session fail. clientId=%v", clientId)
	}
	var sess *smux.Session
	if err := item.Scan(&sess); err != nil {
		return nil, gerror.Wrapf(err, "scan session fail. clientId=%v", clientId)
	}
	return sess, nil
}
