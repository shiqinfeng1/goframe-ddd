package filemgr

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shiqinfeng1/goframe-ddd/pkg/cache"
	"github.com/xtaci/smux"
)

type SessionMgr struct {
	serverSess *gcache.Cache
	mutex      sync.Mutex
}

var (
	sessionMgr *SessionMgr
	once       sync.Once
)

func Session() *SessionMgr {
	once.Do(func() {
		sessionMgr = &SessionMgr{
			serverSess: cache.Memory(),
		}
		go sessionMgr.checkLiveness()
	})
	return sessionMgr
}

func (s *SessionMgr) checkLiveness() {
	ctx := gctx.New()
	// 定时3秒检查连通性
	ticker := time.NewTicker(time.Second * 3)
	for range ticker.C {
		s.mutex.Lock()
		keys, _ := s.serverSess.KeyStrings(ctx)
		for _, key := range keys {
			val, _ := s.serverSess.Get(ctx, key)
			var sess *smux.Session
			if err := val.Scan(&sess); err != nil || sess == nil {
				g.Log().Warning(ctx, "nodeId=%v session is invalid:%v", key, err)
				continue
			}
			if sess.IsClosed() {
				_, err := s.serverSess.Remove(ctx, key)
				g.Log().Infof(ctx, "remove nodeId=%v session:%v", key, err)
			}
		}
		s.mutex.Unlock()
	}
}

func (s *SessionMgr) GetNodeList(ctx context.Context) ([]string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	ids, err := s.serverSess.Keys(ctx)
	if err != nil {
		return []string{}, gerror.Wrap(err, "get sesson nodeIds fail")
	}
	return gconv.Strings(ids), nil
}

func (s *SessionMgr) SaveSession(ctx context.Context, nodeId string, sess *smux.Session) error {
	// 检查之前是否有会话
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if exist, _ := s.serverSess.Contains(ctx, nodeId); exist {
		oldVal, err := s.serverSess.Remove(ctx, nodeId)
		if err != nil {
			return gerror.Wrapf(err, "remove old session fail. nodeId=%v", nodeId)
		}
		var oldSess *smux.Session
		if err := oldVal.Scan(&oldSess); err != nil || oldSess == nil {
			return gerror.Wrapf(err, "scan old session fail. nodeId=%v", nodeId)
		}
		if err := oldSess.Close(); err != nil {
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
	s.mutex.Lock()
	defer s.mutex.Unlock()
	item, err := s.serverSess.Get(ctx, nodeId)
	if err != nil {
		return nil, gerror.Wrapf(err, "get session fail. nodeId=%v", nodeId)
	}
	var sess *smux.Session
	if err := item.Scan(&sess); err != nil || sess == nil {
		return nil, gerror.Wrapf(err, "scan session fail. nodeId=%v", nodeId)
	}
	return sess, nil
}
