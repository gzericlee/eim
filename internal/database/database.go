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
	DriverMySQL   Driver = "mysql"
	DriverMongoDB Driver = "mongodb"
)

type IDatabase interface {
	SaveDevice(device *model.Device) error
	SaveMessage(msg *model.Message) error
}

func NewDatabase(driver Driver, connection, name string) (IDatabase, error) {
	switch driver {
	//case DriverMySQL:
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
	case DriverMongoDB:
		{
			client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connection))
			if err != nil {
				return nil, err
			}
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := client.Ping(ctx, readpref.Primary()); err != nil {
				return nil, err
			}
			return mongodb.NewRepository(client.Database(name)), err
		}
	}
	return nil, fmt.Errorf("unsupported driver: %s", driver)
}
