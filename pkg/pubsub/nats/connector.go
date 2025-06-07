// Package nats connector.go
package natsclient

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Connector represents the main Client connection.
type Connector interface {
	Connect(string, ...nats.Option) (natsConn, error)
}

// JetStreamCreator represents the main Client jStream Client.
type JetStreamCreator interface {
	New(conn natsConn) (jetstream.JetStream, error)
}

// ConnIntf represents the main Client connection.
type natsConn interface {
	Status() nats.Status
	Close()
	IsConnected() bool
	NatsConn() *nats.Conn
	NewJetStream() (jetstream.JetStream, error)
}

// natsConnWrapper wraps a nats.Conn to implement the ConnIntf.
type natsConnWrapper struct {
	conn *nats.Conn
}

func (w *natsConnWrapper) IsConnected() bool {
	return w.conn.IsConnected()
}
func (w *natsConnWrapper) Status() nats.Status {
	w.conn.IsConnected()
	return w.conn.Status()
}

func (w *natsConnWrapper) Close() {
	w.conn.Close()
}

func (w *natsConnWrapper) NatsConn() *nats.Conn {
	return w.conn
}

func (w *natsConnWrapper) NewJetStream() (jetstream.JetStream, error) {
	if w.conn == nil {
		return nil, gerror.New("invalid nats conn")
	}
	js, err := jetstream.New(w.conn)
	if err != nil {
		return nil, gerror.Wrap(err, "new jetstream fail")
	}
	return js, nil
}

type defaultConnector struct{}

func (*defaultConnector) Connect(serverURL string, opts ...nats.Option) (natsConn, error) {
	nc, err := nats.Connect(serverURL, opts...)
	if err != nil {
		return nil, gerror.Wrap(err, "nats connect fail")
	}
	return &natsConnWrapper{conn: nc}, nil
}
