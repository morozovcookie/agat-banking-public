package percona

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	// DriverName is the name of driver which will be used for work with sql data storage.
	DriverName = "mysql"

	// ConnectTimeout is the maximum time for waiting until connect operation will be finished.
	ConnectTimeout = time.Second

	// PingTimeout is the maximum time for waiting until ping operation will be finished.
	PingTimeout = time.Millisecond * 100
)

var (
	_ TxBeginner = (*Client)(nil)
	_ Preparer   = (*Client)(nil)
)

// Client represents an object for basic manipulation with Percona MySQL Database System.
type Client struct {
	db  *sql.DB
	dsn string

	connMaxLifetime time.Duration
	maxIdleConns    int
	connMaxIdleTime time.Duration
	maxOpenConns    int

	dbName string
	dbUser string
}

// NewClient returns a new Client instance.
func NewClient(dsn string) *Client {
	return &Client{
		db:  nil,
		dsn: dsn,

		connMaxLifetime: DefaultConnMaxLifetime,
		maxIdleConns:    DefaultMaxIdleConns,
		connMaxIdleTime: DefaultConnMaxIdleTime,
		maxOpenConns:    DefaultMaxOpenConns,
	}
}

// DBName returns name of database which client are connected.
func (c *Client) DBName() string {
	return c.dbName
}

// DBUser returns name of user which connected to the database.
func (c *Client) DBUser() string {
	return c.dbUser
}

// Connect opens and set up a database.
func (c *Client) Connect(ctx context.Context) (err error) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(ConnectTimeout))
	defer cancel()

	cfg, err := mysql.ParseDSN(c.dsn)
	if err != nil {
		return errors.Wrap(err, "percona connect")
	}

	c.dbName, c.dbUser = cfg.DBName, cfg.User

	if c.db, err = sql.Open(DriverName, c.dsn); err != nil {
		return errors.Wrap(err, "percona connect")
	}

	c.db.SetConnMaxLifetime(c.connMaxLifetime)
	c.db.SetMaxIdleConns(c.maxIdleConns)
	c.db.SetConnMaxIdleTime(c.connMaxIdleTime)
	c.db.SetMaxOpenConns(c.maxOpenConns)

	if err = c.PingContext(ctx); err != nil {
		return errors.Wrap(err, "percona connect")
	}

	return nil
}

// PingContext verifies a connection to the database is still alive.
func (c *Client) PingContext(ctx context.Context) error {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(PingTimeout))
	defer cancel()

	if err := c.db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "percona ping")
	}

	return nil
}

// BeginTx starts a transaction.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	res, err := c.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, errors.Wrap(err, "percona begin")
	}

	return &tx{
		Tx: res,
	}, nil
}

// PrepareContext returns prepared statement.
func (c *Client) PrepareContext(ctx context.Context, query string) (Stmt, error) {
	res, err := c.db.PrepareContext(ctx, query) // nolint:sqlclosecheck
	if err != nil {
		return nil, errors.Wrap(err, "percona prepare")
	}

	return &stmt{
		Stmt: res,
	}, nil
}

// Close closes database connection.
func (c *Client) Close(_ context.Context) error {
	if err := c.db.Close(); err != nil {
		return errors.Wrap(err, "percona close")
	}

	return nil
}
