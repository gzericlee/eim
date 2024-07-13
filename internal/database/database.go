package database

import (
	"context"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"xorm.io/xorm"

	"eim/internal/database/mongodb"
	"eim/internal/database/pgsql"
	"eim/internal/model"
)

type Driver string

const (
	MongoDBDriver  Driver = "mongodb"
	PostgresDriver Driver = "postgres"
)

var _ IDatabase = &mongodb.Repository{}
var _ IDatabase = &pgsql.Repository{}

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
	InsertTenant(tenant *model.Tenant) error
	UpdateTenant(tenant *model.Tenant) error
	GetTenant(tenantId string) (*model.Tenant, error)
	EnableTenant(tenantId string) error
	DisableTenant(tenantId string) error
}

type IBiz interface {
	InsertBiz(biz *model.Biz) error
	UpdateBiz(biz *model.Biz) error
	GetBiz(bizId, tenantId string) (*model.Biz, error)
	EnableBiz(bizId, tenantId string) error
	DisableBiz(bizId, tenantId string) error
	ListBizs(filter map[string]interface{}, order []string, limit, offset int64) ([]*model.Biz, int64, error)
}

type IBizMember interface {
	InsertBizMember(member *model.BizMember) error
	DeleteBizMember(bizId, tenantId, memberId string) error
	GetBizMembers(bizId, tenantId string) ([]*model.BizMember, error)
}

type IDevice interface {
	InsertDevice(device *model.Device) error
	UpdateDevice(device *model.Device) error
	GetDevicesByUser(userId, tenantId string) ([]*model.Device, error)
	GetDevice(userId, tenantId, deviceId string) (*model.Device, error)
	DeleteDevice(userId, tenantId, deviceId string) error
	ListDevices(filter map[string]interface{}, order []string, limit, offset int64) ([]*model.Device, int64, error)
}

type IMessage interface {
	InsertMessage(message *model.Message) error
	InsertMessages(messages []*model.Message) error
	GetMessagesByIds(msgIds []int64) ([]*model.Message, error)
	ListHistoryMessages(filter map[string]interface{}, order []string, minSeq, maxSeq, limit, offset int64) ([]*model.Message, error)
}

func NewDatabase(driver Driver, connections []string, name string) (IDatabase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	switch driver {
	case PostgresDriver:
		{
			if len(connections) == 0 {
				return nil, fmt.Errorf("empty connections")
			}

			eg, err := xorm.NewEngineGroup(string(PostgresDriver), connections, xorm.RoundRobinPolicy())
			if err != nil {
				return nil, fmt.Errorf("new xorm engine group -> %w", err)
			}

			eg.SetMaxIdleConns(10)
			eg.SetMaxOpenConns(50)
			eg.SetConnMaxLifetime(1 * time.Hour)

			repo, err := pgsql.NewRepository(eg)
			if err != nil {
				return nil, fmt.Errorf("new pgsql repository -> %w", err)
			}

			return repo, nil
		}
	case MongoDBDriver:
		{
			client, err := mongo.Connect(ctx, options.Client().ApplyURI(connections[0]))
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
