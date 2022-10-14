package maindb

import (
	"errors"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"eim/internal/types"
)

const (
	driverTidb = "mysql"
)

var (
	tidb *gorm.DB
)

type tidbRepository struct {
}

type Db interface {
	SaveDevice(device *types.Device) error
	SaveMessage(msg *types.Message) error
}

func NewMainDB() Db {
	switch driverTidb {
	case driverTidb:
		return &tidbRepository{}
	}
	return nil
}

func InitDBEngine(driver, dns string) error {
	var err error
	switch driver {
	case driverTidb:
		{
			tidb, err = gorm.Open(mysql.Open(dns), &gorm.Config{SkipDefaultTransaction: true})
			if err != nil {
				return err
			}
			db, err := tidb.DB()
			if err != nil {
				return err
			}
			db.SetConnMaxLifetime(time.Hour)
			db.SetMaxIdleConns(100)
			db.SetMaxOpenConns(200)

			_ = tidb.AutoMigrate(&types.Device{}, &types.Message{})

			return nil
		}
	}
	return errors.New("unsupported driver")
}
