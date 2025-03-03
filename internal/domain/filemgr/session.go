package filemgr

import (
	"context"
	"errors"
	"io"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/xtaci/smux"
)

var serverSess = gcache.NewWithAdapter(gcache.NewAdapterMemory())

func saveSession(ctx context.Context, clientId string, sess *smux.Session) error {
	if exist, _ := serverSess.Contains(ctx, clientId); exist {
		old, err := serverSess.Remove(ctx, clientId)
		if err != nil {
			return gerror.Wrapf(err, "remove old session fail. clientId=%v", clientId)
		}
		var s *smux.Session
		if err := old.Scan(&s); err != nil {
			return gerror.Wrapf(err, "scan old session fail. clientId=%v", clientId)
		}
		if err := s.Close(); err != nil {
			if errors.Is(err, io.ErrClosedPipe) {
				return nil
			}
		}
	}
	if err := serverSess.Set(ctx, clientId, sess, 0); err != nil {
		return gerror.Wrapf(err, "set new session fail. clientId=%v", clientId)
	}
	return nil
}

func GetSession(ctx context.Context, clientId string) (*smux.Session, error) {
	sess, err := serverSess.Get(ctx, clientId)
	if err != nil {
		return nil, gerror.Wrapf(err, "get session fail. clientId=%v", clientId)
	}
	var s *smux.Session
	if err := sess.Scan(&s); err != nil {
		return nil, gerror.Wrapf(err, "scan session fail. clientId=%v", clientId)
	}
	return s, nil
}
