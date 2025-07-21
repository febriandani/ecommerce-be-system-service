package infra

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/febriandani/ecommerce-be-system-service/cmd/server/utils"

	"github.com/gofiber/fiber/v2/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBSystem struct {
	Username     string `json:",omitempty"`
	Password     string `json:",omitempty"`
	URL          string `json:",omitempty"`
	Port         string `json:",omitempty"`
	DBName       string `json:",omitempty"`
	MaxIdleConns int    `json:",omitempty"`
	MaxOpenConns int    `json:",omitempty"`
	MaxLifeTime  int    `json:",omitempty"`
	Timeout      string `json:",omitempty"`
	SSLMode      string `json:",omitempty"`
}

type DatabaseSystem struct {
	ReadUser  DBSystem `json:",omitempty"`
	WriteUser DBSystem `json:",omitempty"`
}

type DatabaseList struct {
	Backend DatabaseType
}

type DatabaseType struct {
	Read  Database
	Write Database
}

// IDatabase is interface for database
type Database interface {
	ConnectDB(db *DBSystem)
	Close()

	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Get(dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	// DriverName() string

	Begin() (*sql.Tx, error)
	In(query string, params ...interface{}) (string, []interface{}, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row
	// QueryRowSqlx(query string, args ...interface{}) *sqlx.Row
	// QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	// GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type DBHandler struct {
	DB  *sqlx.DB
	Err error
	log log.Logger
}

func NewDB(logger log.Logger) *DBHandler {
	return &DBHandler{
		log: logger,
	}
}

// ConnectDB - function for connect DB.
func (d *DBHandler) ConnectDB(dbBs *DBSystem) {
	dbs, err := sqlx.Open("postgres", "user="+dbBs.Username+" password="+dbBs.Password+" sslmode="+dbBs.SSLMode+" dbname="+dbBs.DBName+" host="+dbBs.URL+" port="+dbBs.Port+" connect_timeout="+dbBs.Timeout)
	if err != nil {
		log.Error(utils.ConnectDBFail + " | " + err.Error())
		d.Err = err
		return
	}

	d.DB = dbs

	err = d.DB.Ping()
	if err != nil {
		if strings.Contains(err.Error(), "Error: failed to connect to") {
			log.Error(utils.ConnectDBFail, err.Error())
			d.Err = err
		}
		return
	}

	d.log.Info(utils.ConnectDBSuccess)
	d.DB.SetConnMaxLifetime(time.Duration(dbBs.MaxLifeTime))
}

// Close - function for connection lost.
func (d *DBHandler) Close() {
	if err := d.DB.Close(); err != nil {
		d.log.Info(utils.ClosingDBFailed + " | " + err.Error())
	} else {
		d.log.Info(utils.ClosingDBSuccess)
	}
}

func (d *DBHandler) Exec(query string, args ...interface{}) (sql.Result, error) {
	result, err := d.DB.Exec(query, args...)
	return result, err
}

func (d *DBHandler) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := d.DB.ExecContext(ctx, query, args...)
	return result, err
}

func (d *DBHandler) Query(query string, args ...interface{}) (*sql.Rows, error) {
	result, err := d.DB.Query(query, args...)
	return result, err
}

func (d *DBHandler) Select(dest interface{}, query string, args ...interface{}) error {
	err := d.DB.Select(dest, query, args...)
	return err
}

func (d *DBHandler) Get(dest interface{}, query string, args ...interface{}) error {
	err := d.DB.Get(dest, query, args...)
	return err
}

func (d *DBHandler) Rebind(query string) string {
	return d.DB.Rebind(query)
}

func (d *DBHandler) In(query string, params ...interface{}) (string, []interface{}, error) {
	query, args, err := sqlx.In(query, params...)
	return query, args, err
}

func (d *DBHandler) Begin() (*sql.Tx, error) {
	return d.DB.Begin()
}

func (d *DBHandler) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return d.DB.QueryRowContext(ctx, query, args...)
}

func (d *DBHandler) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	err := d.DB.GetContext(ctx, dest, query, args...)
	return err
}
