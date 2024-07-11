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
	MongoDBDriver Driver = "mongodb"
)

var _ IDatabase = &mongodb.Repository{}

type IDatabase interface {
	IMessage
	IDevice
	IBiz
	IBizMember
	ISegment
	ITenant
}

type ISegment interface {
	GetSegment(bizId, tenantId string) (*model.Segment, error)
}

type ITenant interface {
	SaveTenant(tenant *model.Tenant) error
	GetTenant(tenantId string) (*model.Tenant, error)
	EnableTenant(tenantId string) error
	DisableTenant(tenantId string) error
}

type IBiz interface {
	SaveBiz(biz *model.Biz) error
	GetBiz(bizId, tenantId string) (*model.Biz, error)
	EnableBiz(bizId, tenantId string) error
	DisableBiz(bizId, tenantId string) error
	ListBizs(filter map[string]interface{}, limit, offset int64) ([]*model.Biz, int64, error)
}

type IBizMember interface {
	InsertBizMember(member *model.BizMember) error
	DeleteBizMember(bizId, tenantId, memberId string) error
	GetBizMembers(bizId, tenantId string) ([]*model.BizMember, error)
}

type IDevice interface {
	SaveDevice(device *model.Device) error
	GetDevicesByUser(userId, tenantId string) ([]*model.Device, error)
	GetDevice(userId, deviceId string) (*model.Device, error)
	DeleteDevice(userId, tenantId, deviceId string) error
	ListDevices(filter map[string]interface{}, limit, offset int64) ([]*model.Device, int64, error)
}

type IMessage interface {
	SaveMessage(message *model.Message) error
	SaveMessages(messages []*model.Message) error
	GetMessagesByIds(msgIds []int64) ([]*model.Message, error)
	ListHistoryMessages(filter map[string]interface{}, minSeq, maxSeq, limit, offset int64) ([]*model.Message, error)
}

func NewDatabase(driver Driver, connection, name string) (IDatabase, error) {
	switch driver {
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
