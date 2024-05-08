package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"eim/internal/database/mongodb"
	"eim/internal/model"
)

type Driver string

const (
	MySQLDriver   Driver = "mysql"
	MongoDBDriver Driver = "mongodb"
)

type IDatabase interface {
	SaveDevice(device *model.Device) error
	SaveMessage(msg *model.Message) error
	GetSegment(id string) (*model.Segment, error)
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
	//		_ = orm.AutoMigrate(&model.Device{}, &model.Message{})
	//
	//		return mysqldb.NewRepository(orm), nil
	//	}
	case MongoDBDriver:
		{
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			client, err := mongo.Connect(ctx, options.Client().ApplyURI(connection))
			if err != nil {
				return nil, err
			}

			if err := client.Ping(ctx, readpref.Primary()); err != nil {
				return nil, err
			}

			db := client.Database(name)

			return mongodb.NewRepository(db)
		}
	}
	return nil, fmt.Errorf("unsupported driver: %s", driver)
}
