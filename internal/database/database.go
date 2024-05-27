package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"eim/internal/database/mongodb"
	"eim/internal/database/mysql"
	"eim/internal/model"
)

type Driver string

const (
	MySQLDriver   Driver = "mysql"
	MongoDBDriver Driver = "mongodb"
)

var _ IDatabase = &mongodb.Repository{}
var _ IDatabase = &mysql.Repository{}

type IDatabase interface {
	SaveDevice(device *model.Device) error
	GetDevices(userId string) ([]*model.Device, error)
	GetDevice(userId, deviceId string) (*model.Device, error)

	SaveUser(user *model.User) error
	GetUser(userId, tenantId string) (*model.User, error)

	SaveMessage(msg *model.Message) error
	GetMessagesByIds(msgIds []int64) ([]*model.Message, error)

	GetSegment(bizId string) (*model.Segment, error)
}

func NewDatabase(driver Driver, connection, name string) (IDatabase, error) {
	switch driver {
	//case MySQLDriver:
	//	{
	//		orm, err := gorm.Open(mysql.Open(connection), &gorm.Config{SkipDefaultTransaction: true})
	//		if err != nil {
	//			return nil, err
	//		}
	//		db, err := orm.DB()
	//		if err != nil {
	//			return nil, err
	//		}
	//		db.SetConnMaxLifetime(time.Hour)
	//		db.SetMaxIdleConns(100)
	//		db.SetMaxOpenConns(200)
	//
	//		_ = orm.AutoMigrate(&model.device{}, &model.message{})
	//
	//		return mysqldb.NewRepository(orm), nil
	//	}
	case MongoDBDriver:
		{
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
			if err != nil {
				return nil, fmt.Errorf("connect mongodb -> %w", err)
			}

			if err := client.Ping(ctx, readpref.Primary()); err != nil {
				return nil, fmt.Errorf("ping mongodb -> %w", err)
			}

			db := client.Database(name)

			repo, err := mongodb.NewRepository(db)
			if err != nil {
				return nil, fmt.Errorf("new mongodb repository -> %w", err)
			}

			return repo, nil
		}
	}
	return nil, fmt.Errorf("unsupported driver: %s", driver)
}
