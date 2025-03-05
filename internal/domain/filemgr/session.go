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

func (s *SessionMgr) GetNodeList(ctx context.Context) ([]string, error) {
	ids, err := s.serverSess.Keys(ctx)
	if err != nil {
		return []string{}, gerror.Wrap(err, "get sesson nodeIds fail")
	}
	return gconv.Strings(ids), nil
}

func (s *SessionMgr) SaveSession(ctx context.Context, nodeId string, sess *smux.Session) error {
	// 检查之前是否有会话
	if exist, _ := s.serverSess.Contains(ctx, nodeId); exist {
		old, err := s.serverSess.Remove(ctx, nodeId)
		if err != nil {
			return gerror.Wrapf(err, "remove old session fail. nodeId=%v", nodeId)
		}
		var s *smux.Session
		if err := old.Scan(&s); err != nil {
			return gerror.Wrapf(err, "scan old session fail. nodeId=%v", nodeId)
		}
		if err := s.Close(); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				g.Log().Warningf(ctx, "server-side close old session of nodeid=%v:%v", nodeId, err)
			} else {
				g.Log().Warningf(ctx, "server-side close old session of nodeid=%v fail:%v", nodeId, err)
			}
		}
	}

	if err := s.serverSess.Set(ctx, nodeId, sess, 0); err != nil {
		return gerror.Wrapf(err, "set new session fail. nodeId=%v", nodeId)
	}
	return nil
}

func (s *SessionMgr) GetSession(ctx context.Context, nodeId string) (*smux.Session, error) {
	item, err := s.serverSess.Get(ctx, nodeId)
	if err != nil {
		return nil, gerror.Wrapf(err, "get session fail. nodeId=%v", nodeId)
	}
	var sess *smux.Session
	if err := item.Scan(&sess); err != nil {
		return nil, gerror.Wrapf(err, "scan session fail. nodeId=%v", nodeId)
	}
	return sess, nil
}
