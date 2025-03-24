// Package nats connector.go
package nats

import (
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// natsConnWrapper wraps a nats.Conn to implement the ConnIntf.
type natsConnWrapper struct {
	conn *nats.Conn
}

func (w *natsConnWrapper) Status() nats.Status {
	return w.conn.Status()
}

func (w *natsConnWrapper) Close() {
	w.conn.Close()
}

func (w *natsConnWrapper) NATSConn() *nats.Conn {
	return w.conn
}

func (w *natsConnWrapper) JetStream() (jetstream.JetStream, error) {
	return jetstream.New(w.conn)
}

type defaultConnector struct{}

func (*defaultConnector) Connect(serverURL string, opts ...nats.Option) (ConnIntf, error) {
	nc, err := nats.Connect(serverURL, opts...)
	if err != nil {
		return nil, err
	}

	return &natsConnWrapper{nc}, nil
}

type defaultJetStreamCreator struct{}

func (*defaultJetStreamCreator) New(conn ConnIntf) (jetstream.JetStream, error) {
	return conn.JetStream()
}
