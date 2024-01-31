package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/ryanadiputraa/unclatter/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	maxOpenConns    = 60
	connMaxLifeTime = 120
	maxIdleConn     = 30
	connMaxIdleTime = 20
)

func NewDB(c *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s TimeZone=UTC",
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.User,
		c.Postgres.DBName,
		c.Postgres.SSLMode,
		c.Postgres.Password,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxLifetime(connMaxLifeTime * time.Second)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)

	err = pingWithRetry(db, 3*time.Second, 5)
	if err != nil {
		return nil, err
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB, err
}

func pingWithRetry(db *sql.DB, interval time.Duration, maxRetries int) (err error) {
	for i := 0; i < maxRetries; i++ {
		err = db.Ping()
		if err == nil {
			break
		}
		time.Sleep(interval)
	}
	return
}
