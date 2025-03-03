package stream

import (
	"io"
	"net"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/xtaci/smux"
)

func newSessoinByServer(conn net.Conn) (*smux.Session, error) {
	return newSession(conn, smux.Server)
}

func newSessionByClient(conn net.Conn) (*smux.Session, error) {
	return newSession(conn, smux.Client)
}

func newSession(conn net.Conn, creator func(io.ReadWriteCloser, *smux.Config) (*smux.Session, error)) (*smux.Session, error) {
	smuxConfig := smux.DefaultConfig()
	smuxConfig.Version = 2
	smuxConfig.MaxReceiveBuffer = 16*1024*1024 + 1
	smuxConfig.MaxStreamBuffer = 1 * 1024 * 1024
	smuxConfig.MaxFrameSize = 65000
	smuxConfig.KeepAliveInterval = time.Duration(5) * time.Second
	smuxConfig.KeepAliveTimeout = time.Duration(15) * time.Second
	if err := smux.VerifyConfig(smuxConfig); err != nil {
		return nil, gerror.Wrapf(err, "set smux config fail")
	}
	client, err := creator(conn, smuxConfig)
	if err != nil {
		return nil, gerror.Wrapf(err, "create smux session fail")
	}
	return client, nil
}
