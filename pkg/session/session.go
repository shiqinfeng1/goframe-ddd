package session

import (
	"context"
	"errors"
	"io"
	"strings"
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
	cache *gcache.Cache
	mutex sync.Mutex
}

var (
	sessionMgr *SessionMgr
)

func init() {
	sessionMgr = &SessionMgr{
		cache: cache.KV(),
	}
	go sessionMgr.checkLiveness()
}

func (s *SessionMgr) checkLiveness() {
	ctx := gctx.New()
	// 定时3秒检查连通性
	ticker := time.NewTicker(time.Second * 3)
	for range ticker.C {
		s.mutex.Lock()
		keys, _ := s.cache.KeyStrings(ctx)
		for _, key := range keys {
			if strings.HasPrefix(key, "session") {
				val, _ := s.cache.Get(ctx, key)
				var sess *smux.Session
				if err := val.Scan(&sess); err != nil || sess == nil {
					g.Log().Warning(ctx, "nodeId=%v session is invalid:%v", key, err)
					continue
				}
				if sess.IsClosed() {
					_, err := s.cache.Remove(ctx, key)
					g.Log().Infof(ctx, "remove nodeId=%v session:%v", key, err)
				}
			}
		}
		s.mutex.Unlock()
	}
}

func GetSessionNodeList(ctx context.Context) ([]string, error) {
	sessionMgr.mutex.Lock()
	defer sessionMgr.mutex.Unlock()
	ids, err := sessionMgr.cache.Keys(ctx)
	if err != nil {
		return []string{}, gerror.Wrap(err, "get sesson nodeIds fail")
	}
	out := make([]string, 0)
	tmp := gconv.Strings(ids)
	for _, v := range tmp {
		if strings.HasPrefix(v, "session") {
			out = append(out, strings.TrimPrefix(v, "session"))
		}
	}
	return out, nil
}

func SaveSession(ctx context.Context, nodeId string, sess *smux.Session) error {
	key := "session" + nodeId
	// 检查之前是否有会话
	sessionMgr.mutex.Lock()
	defer sessionMgr.mutex.Unlock()
	if exist, _ := sessionMgr.cache.Contains(ctx, key); exist {
		oldVal, err := sessionMgr.cache.Remove(ctx, key)
		if err != nil {
			return gerror.Wrapf(err, "remove old session fail. key=%v", key)
		}
		var oldSess *smux.Session
		if err := oldVal.Scan(&oldSess); err != nil || oldSess == nil {
			return gerror.Wrapf(err, "scan old session fail. key=%v", key)
		}
		if err := oldSess.Close(); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				g.Log().Warningf(ctx, "server-side close old session of nodeid=%v:%v", key, err)
			} else {
				g.Log().Warningf(ctx, "server-side close old session of nodeid=%v fail:%v", key, err)
			}
		}
	}
	if err := sessionMgr.cache.Set(ctx, key, sess, 0); err != nil {
		return gerror.Wrapf(err, "set new session fail. nodeId=%v", key)
	}
	return nil
}

func GetSession(ctx context.Context, nodeId string) (*smux.Session, error) {
	key := "session" + nodeId
	sessionMgr.mutex.Lock()
	defer sessionMgr.mutex.Unlock()
	item, err := sessionMgr.cache.Get(ctx, key)
	if err != nil {
		return nil, gerror.Wrapf(err, "get session fail. key=%v", key)
	}
	var sess *smux.Session
	if err := item.Scan(&sess); err != nil || sess == nil {
		return nil, gerror.Wrapf(err, "scan session fail. key=%v", key)
	}
	return sess, nil
}
