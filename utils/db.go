package utils

import (
	"database/sql"
	"time"

	"github.com/canteen_management/logger"
	_ "github.com/go-sql-driver/mysql"
)

const (
	utilsLogTag = "Utils"
)

// Second is an integer type representing duration in second
type Second int

// Duration converts Second to time.Duration
func (i Second) Duration() time.Duration {
	return time.Second * time.Duration(i)
}

// Config defines configuration for mysql
type Config struct {
	Dsn             string `json:"dsn"`
	MaxIdle         int    `json:"maxIdle"`         // zero means to use default value; negative means 0
	MaxOpen         int    `json:"maxOpen"`         // <= 0 means unlimited
	ConnMaxLifetime Second `json:"connMaxLifetime"` // maximum amount of time a connection may be reused. if ConnMaxLifetime <= 0, no idle connections are retained.
}

func NewMysqlClient(config Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.Dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(config.ConnMaxLifetime.Duration())
	db.SetMaxIdleConns(config.MaxIdle)
	db.SetMaxOpenConns(config.MaxOpen)
	err = db.Ping()
	return db, err
}

func Begin(sqlCli *sql.DB) (*sql.Tx, error) {
	tx, err := sqlCli.Begin()
	if err != nil {
		logger.Warn(utilsLogTag, "Begin Failed|Err:%v", err)
		return nil, err
	}
	return tx, nil
}

func End(tx *sql.Tx, err error) error {
	if err != nil {
		rollErr := tx.Rollback()
		if rollErr != nil {
			logger.Warn(utilsLogTag, "Rollback Failed|Err:%v", err)
			return rollErr
		}
	} else {
		err = tx.Commit()
		if err != nil {
			logger.Warn(utilsLogTag, "Commit Failed|Err:%v", err)
			tx.Rollback()
			return err
		}
	}
	return nil
}
