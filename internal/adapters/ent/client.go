// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/recvchunk"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/recvfile"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/sendchunk"
	"github.com/shiqinfeng1/goframe-ddd/internal/adapters/ent/sendfile"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// RecvChunk is the client for interacting with the RecvChunk builders.
	RecvChunk *RecvChunkClient
	// RecvFile is the client for interacting with the RecvFile builders.
	RecvFile *RecvFileClient
	// SendChunk is the client for interacting with the SendChunk builders.
	SendChunk *SendChunkClient
	// SendFile is the client for interacting with the SendFile builders.
	SendFile *SendFileClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.RecvChunk = NewRecvChunkClient(c.config)
	c.RecvFile = NewRecvFileClient(c.config)
	c.SendChunk = NewSendChunkClient(c.config)
	c.SendFile = NewSendFileClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:       ctx,
		config:    cfg,
		RecvChunk: NewRecvChunkClient(cfg),
		RecvFile:  NewRecvFileClient(cfg),
		SendChunk: NewSendChunkClient(cfg),
		SendFile:  NewSendFileClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:       ctx,
		config:    cfg,
		RecvChunk: NewRecvChunkClient(cfg),
		RecvFile:  NewRecvFileClient(cfg),
		SendChunk: NewSendChunkClient(cfg),
		SendFile:  NewSendFileClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		RecvChunk.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.RecvChunk.Use(hooks...)
	c.RecvFile.Use(hooks...)
	c.SendChunk.Use(hooks...)
	c.SendFile.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.RecvChunk.Intercept(interceptors...)
	c.RecvFile.Intercept(interceptors...)
	c.SendChunk.Intercept(interceptors...)
	c.SendFile.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *RecvChunkMutation:
		return c.RecvChunk.mutate(ctx, m)
	case *RecvFileMutation:
		return c.RecvFile.mutate(ctx, m)
	case *SendChunkMutation:
		return c.SendChunk.mutate(ctx, m)
	case *SendFileMutation:
		return c.SendFile.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// RecvChunkClient is a client for the RecvChunk schema.
type RecvChunkClient struct {
	config
}

// NewRecvChunkClient returns a client for the RecvChunk from the given config.
func NewRecvChunkClient(c config) *RecvChunkClient {
	return &RecvChunkClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `recvchunk.Hooks(f(g(h())))`.
func (c *RecvChunkClient) Use(hooks ...Hook) {
	c.hooks.RecvChunk = append(c.hooks.RecvChunk, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `recvchunk.Intercept(f(g(h())))`.
func (c *RecvChunkClient) Intercept(interceptors ...Interceptor) {
	c.inters.RecvChunk = append(c.inters.RecvChunk, interceptors...)
}

// Create returns a builder for creating a RecvChunk entity.
func (c *RecvChunkClient) Create() *RecvChunkCreate {
	mutation := newRecvChunkMutation(c.config, OpCreate)
	return &RecvChunkCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of RecvChunk entities.
func (c *RecvChunkClient) CreateBulk(builders ...*RecvChunkCreate) *RecvChunkCreateBulk {
	return &RecvChunkCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *RecvChunkClient) MapCreateBulk(slice any, setFunc func(*RecvChunkCreate, int)) *RecvChunkCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &RecvChunkCreateBulk{err: fmt.Errorf("calling to RecvChunkClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*RecvChunkCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &RecvChunkCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for RecvChunk.
func (c *RecvChunkClient) Update() *RecvChunkUpdate {
	mutation := newRecvChunkMutation(c.config, OpUpdate)
	return &RecvChunkUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *RecvChunkClient) UpdateOne(rc *RecvChunk) *RecvChunkUpdateOne {
	mutation := newRecvChunkMutation(c.config, OpUpdateOne, withRecvChunk(rc))
	return &RecvChunkUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *RecvChunkClient) UpdateOneID(id int) *RecvChunkUpdateOne {
	mutation := newRecvChunkMutation(c.config, OpUpdateOne, withRecvChunkID(id))
	return &RecvChunkUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for RecvChunk.
func (c *RecvChunkClient) Delete() *RecvChunkDelete {
	mutation := newRecvChunkMutation(c.config, OpDelete)
	return &RecvChunkDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *RecvChunkClient) DeleteOne(rc *RecvChunk) *RecvChunkDeleteOne {
	return c.DeleteOneID(rc.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *RecvChunkClient) DeleteOneID(id int) *RecvChunkDeleteOne {
	builder := c.Delete().Where(recvchunk.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &RecvChunkDeleteOne{builder}
}

// Query returns a query builder for RecvChunk.
func (c *RecvChunkClient) Query() *RecvChunkQuery {
	return &RecvChunkQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeRecvChunk},
		inters: c.Interceptors(),
	}
}

// Get returns a RecvChunk entity by its id.
func (c *RecvChunkClient) Get(ctx context.Context, id int) (*RecvChunk, error) {
	return c.Query().Where(recvchunk.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *RecvChunkClient) GetX(ctx context.Context, id int) *RecvChunk {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryRecvFile queries the recv_file edge of a RecvChunk.
func (c *RecvChunkClient) QueryRecvFile(rc *RecvChunk) *RecvFileQuery {
	query := (&RecvFileClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := rc.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(recvchunk.Table, recvchunk.FieldID, id),
			sqlgraph.To(recvfile.Table, recvfile.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, recvchunk.RecvFileTable, recvchunk.RecvFileColumn),
		)
		fromV = sqlgraph.Neighbors(rc.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *RecvChunkClient) Hooks() []Hook {
	return c.hooks.RecvChunk
}

// Interceptors returns the client interceptors.
func (c *RecvChunkClient) Interceptors() []Interceptor {
	return c.inters.RecvChunk
}

func (c *RecvChunkClient) mutate(ctx context.Context, m *RecvChunkMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&RecvChunkCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&RecvChunkUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&RecvChunkUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&RecvChunkDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown RecvChunk mutation op: %q", m.Op())
	}
}

// RecvFileClient is a client for the RecvFile schema.
type RecvFileClient struct {
	config
}

// NewRecvFileClient returns a client for the RecvFile from the given config.
func NewRecvFileClient(c config) *RecvFileClient {
	return &RecvFileClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `recvfile.Hooks(f(g(h())))`.
func (c *RecvFileClient) Use(hooks ...Hook) {
	c.hooks.RecvFile = append(c.hooks.RecvFile, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `recvfile.Intercept(f(g(h())))`.
func (c *RecvFileClient) Intercept(interceptors ...Interceptor) {
	c.inters.RecvFile = append(c.inters.RecvFile, interceptors...)
}

// Create returns a builder for creating a RecvFile entity.
func (c *RecvFileClient) Create() *RecvFileCreate {
	mutation := newRecvFileMutation(c.config, OpCreate)
	return &RecvFileCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of RecvFile entities.
func (c *RecvFileClient) CreateBulk(builders ...*RecvFileCreate) *RecvFileCreateBulk {
	return &RecvFileCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *RecvFileClient) MapCreateBulk(slice any, setFunc func(*RecvFileCreate, int)) *RecvFileCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &RecvFileCreateBulk{err: fmt.Errorf("calling to RecvFileClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*RecvFileCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &RecvFileCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for RecvFile.
func (c *RecvFileClient) Update() *RecvFileUpdate {
	mutation := newRecvFileMutation(c.config, OpUpdate)
	return &RecvFileUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *RecvFileClient) UpdateOne(rf *RecvFile) *RecvFileUpdateOne {
	mutation := newRecvFileMutation(c.config, OpUpdateOne, withRecvFile(rf))
	return &RecvFileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *RecvFileClient) UpdateOneID(id int) *RecvFileUpdateOne {
	mutation := newRecvFileMutation(c.config, OpUpdateOne, withRecvFileID(id))
	return &RecvFileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for RecvFile.
func (c *RecvFileClient) Delete() *RecvFileDelete {
	mutation := newRecvFileMutation(c.config, OpDelete)
	return &RecvFileDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *RecvFileClient) DeleteOne(rf *RecvFile) *RecvFileDeleteOne {
	return c.DeleteOneID(rf.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *RecvFileClient) DeleteOneID(id int) *RecvFileDeleteOne {
	builder := c.Delete().Where(recvfile.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &RecvFileDeleteOne{builder}
}

// Query returns a query builder for RecvFile.
func (c *RecvFileClient) Query() *RecvFileQuery {
	return &RecvFileQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeRecvFile},
		inters: c.Interceptors(),
	}
}

// Get returns a RecvFile entity by its id.
func (c *RecvFileClient) Get(ctx context.Context, id int) (*RecvFile, error) {
	return c.Query().Where(recvfile.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *RecvFileClient) GetX(ctx context.Context, id int) *RecvFile {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryRecvChunks queries the recv_chunks edge of a RecvFile.
func (c *RecvFileClient) QueryRecvChunks(rf *RecvFile) *RecvChunkQuery {
	query := (&RecvChunkClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := rf.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(recvfile.Table, recvfile.FieldID, id),
			sqlgraph.To(recvchunk.Table, recvchunk.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, recvfile.RecvChunksTable, recvfile.RecvChunksColumn),
		)
		fromV = sqlgraph.Neighbors(rf.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *RecvFileClient) Hooks() []Hook {
	return c.hooks.RecvFile
}

// Interceptors returns the client interceptors.
func (c *RecvFileClient) Interceptors() []Interceptor {
	return c.inters.RecvFile
}

func (c *RecvFileClient) mutate(ctx context.Context, m *RecvFileMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&RecvFileCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&RecvFileUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&RecvFileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&RecvFileDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown RecvFile mutation op: %q", m.Op())
	}
}

// SendChunkClient is a client for the SendChunk schema.
type SendChunkClient struct {
	config
}

// NewSendChunkClient returns a client for the SendChunk from the given config.
func NewSendChunkClient(c config) *SendChunkClient {
	return &SendChunkClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `sendchunk.Hooks(f(g(h())))`.
func (c *SendChunkClient) Use(hooks ...Hook) {
	c.hooks.SendChunk = append(c.hooks.SendChunk, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `sendchunk.Intercept(f(g(h())))`.
func (c *SendChunkClient) Intercept(interceptors ...Interceptor) {
	c.inters.SendChunk = append(c.inters.SendChunk, interceptors...)
}

// Create returns a builder for creating a SendChunk entity.
func (c *SendChunkClient) Create() *SendChunkCreate {
	mutation := newSendChunkMutation(c.config, OpCreate)
	return &SendChunkCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of SendChunk entities.
func (c *SendChunkClient) CreateBulk(builders ...*SendChunkCreate) *SendChunkCreateBulk {
	return &SendChunkCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *SendChunkClient) MapCreateBulk(slice any, setFunc func(*SendChunkCreate, int)) *SendChunkCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &SendChunkCreateBulk{err: fmt.Errorf("calling to SendChunkClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*SendChunkCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &SendChunkCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for SendChunk.
func (c *SendChunkClient) Update() *SendChunkUpdate {
	mutation := newSendChunkMutation(c.config, OpUpdate)
	return &SendChunkUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SendChunkClient) UpdateOne(sc *SendChunk) *SendChunkUpdateOne {
	mutation := newSendChunkMutation(c.config, OpUpdateOne, withSendChunk(sc))
	return &SendChunkUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SendChunkClient) UpdateOneID(id int) *SendChunkUpdateOne {
	mutation := newSendChunkMutation(c.config, OpUpdateOne, withSendChunkID(id))
	return &SendChunkUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for SendChunk.
func (c *SendChunkClient) Delete() *SendChunkDelete {
	mutation := newSendChunkMutation(c.config, OpDelete)
	return &SendChunkDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *SendChunkClient) DeleteOne(sc *SendChunk) *SendChunkDeleteOne {
	return c.DeleteOneID(sc.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *SendChunkClient) DeleteOneID(id int) *SendChunkDeleteOne {
	builder := c.Delete().Where(sendchunk.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SendChunkDeleteOne{builder}
}

// Query returns a query builder for SendChunk.
func (c *SendChunkClient) Query() *SendChunkQuery {
	return &SendChunkQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeSendChunk},
		inters: c.Interceptors(),
	}
}

// Get returns a SendChunk entity by its id.
func (c *SendChunkClient) Get(ctx context.Context, id int) (*SendChunk, error) {
	return c.Query().Where(sendchunk.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SendChunkClient) GetX(ctx context.Context, id int) *SendChunk {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QuerySendFile queries the send_file edge of a SendChunk.
func (c *SendChunkClient) QuerySendFile(sc *SendChunk) *SendFileQuery {
	query := (&SendFileClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := sc.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(sendchunk.Table, sendchunk.FieldID, id),
			sqlgraph.To(sendfile.Table, sendfile.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, sendchunk.SendFileTable, sendchunk.SendFileColumn),
		)
		fromV = sqlgraph.Neighbors(sc.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SendChunkClient) Hooks() []Hook {
	return c.hooks.SendChunk
}

// Interceptors returns the client interceptors.
func (c *SendChunkClient) Interceptors() []Interceptor {
	return c.inters.SendChunk
}

func (c *SendChunkClient) mutate(ctx context.Context, m *SendChunkMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&SendChunkCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&SendChunkUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&SendChunkUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&SendChunkDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown SendChunk mutation op: %q", m.Op())
	}
}

// SendFileClient is a client for the SendFile schema.
type SendFileClient struct {
	config
}

// NewSendFileClient returns a client for the SendFile from the given config.
func NewSendFileClient(c config) *SendFileClient {
	return &SendFileClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `sendfile.Hooks(f(g(h())))`.
func (c *SendFileClient) Use(hooks ...Hook) {
	c.hooks.SendFile = append(c.hooks.SendFile, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `sendfile.Intercept(f(g(h())))`.
func (c *SendFileClient) Intercept(interceptors ...Interceptor) {
	c.inters.SendFile = append(c.inters.SendFile, interceptors...)
}

// Create returns a builder for creating a SendFile entity.
func (c *SendFileClient) Create() *SendFileCreate {
	mutation := newSendFileMutation(c.config, OpCreate)
	return &SendFileCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of SendFile entities.
func (c *SendFileClient) CreateBulk(builders ...*SendFileCreate) *SendFileCreateBulk {
	return &SendFileCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *SendFileClient) MapCreateBulk(slice any, setFunc func(*SendFileCreate, int)) *SendFileCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &SendFileCreateBulk{err: fmt.Errorf("calling to SendFileClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*SendFileCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &SendFileCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for SendFile.
func (c *SendFileClient) Update() *SendFileUpdate {
	mutation := newSendFileMutation(c.config, OpUpdate)
	return &SendFileUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *SendFileClient) UpdateOne(sf *SendFile) *SendFileUpdateOne {
	mutation := newSendFileMutation(c.config, OpUpdateOne, withSendFile(sf))
	return &SendFileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *SendFileClient) UpdateOneID(id int) *SendFileUpdateOne {
	mutation := newSendFileMutation(c.config, OpUpdateOne, withSendFileID(id))
	return &SendFileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for SendFile.
func (c *SendFileClient) Delete() *SendFileDelete {
	mutation := newSendFileMutation(c.config, OpDelete)
	return &SendFileDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *SendFileClient) DeleteOne(sf *SendFile) *SendFileDeleteOne {
	return c.DeleteOneID(sf.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *SendFileClient) DeleteOneID(id int) *SendFileDeleteOne {
	builder := c.Delete().Where(sendfile.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &SendFileDeleteOne{builder}
}

// Query returns a query builder for SendFile.
func (c *SendFileClient) Query() *SendFileQuery {
	return &SendFileQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeSendFile},
		inters: c.Interceptors(),
	}
}

// Get returns a SendFile entity by its id.
func (c *SendFileClient) Get(ctx context.Context, id int) (*SendFile, error) {
	return c.Query().Where(sendfile.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *SendFileClient) GetX(ctx context.Context, id int) *SendFile {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QuerySendChunks queries the send_chunks edge of a SendFile.
func (c *SendFileClient) QuerySendChunks(sf *SendFile) *SendChunkQuery {
	query := (&SendChunkClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := sf.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(sendfile.Table, sendfile.FieldID, id),
			sqlgraph.To(sendchunk.Table, sendchunk.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, sendfile.SendChunksTable, sendfile.SendChunksColumn),
		)
		fromV = sqlgraph.Neighbors(sf.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *SendFileClient) Hooks() []Hook {
	return c.hooks.SendFile
}

// Interceptors returns the client interceptors.
func (c *SendFileClient) Interceptors() []Interceptor {
	return c.inters.SendFile
}

func (c *SendFileClient) mutate(ctx context.Context, m *SendFileMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&SendFileCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&SendFileUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&SendFileUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&SendFileDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown SendFile mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		RecvChunk, RecvFile, SendChunk, SendFile []ent.Hook
	}
	inters struct {
		RecvChunk, RecvFile, SendChunk, SendFile []ent.Interceptor
	}
)
