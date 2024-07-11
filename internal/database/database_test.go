package database

import (
	"fmt"
	"testing"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"eim/internal/model"
	"eim/internal/model/consts"
)

var db IDatabase

func init() {
	var err error
	db, err = NewDatabase(MongoDBDriver, "mongodb://admin:pass%40word1@127.0.0.1:27017/?authSource=admin&connect=direct", "eim")
	if err != nil {
		panic(err)
	}
}

func TestSaveTenant(t *testing.T) {
	boolValue, _ := anypb.New(&wrapperspb.BoolValue{Value: true})

	err := db.SaveTenant(&model.Tenant{
		TenantId:   "bingo",
		TenantName: "广州市品高软件股份有限公司",
		State:      consts.StatusEnabled,
		Attributes: map[string]*anypb.Any{
			consts.FileflexEnabled: boolValue,
			consts.FileflexBucket:  {Value: []byte("bingo")},
			consts.FileflexUser:    {Value: []byte("bingo")},
			consts.FileflexPasswd:  {Value: []byte("pass@word1")},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestSaveBiz(t *testing.T) {
	for i := 1; i <= 1000000; i++ {
		err := db.SaveBiz(&model.Biz{
			BizId:      fmt.Sprintf("user-%d", i),
			BizType:    consts.BizUser,
			BizName:    fmt.Sprintf("用户-%d", i),
			TenantId:   "bingo",
			TenantName: "品高软件",
			State:      consts.StatusEnabled,
			Attributes: map[string]*anypb.Any{"password": {Value: []byte("pass@word1")}},
		})
		if err != nil {
			t.Fatal(err)
			return
		}
	}
}

func TestSaveDevice(t *testing.T) {
	for i := 1; i <= 1000000; i++ {
		err := db.SaveDevice(&model.Device{
			DeviceId:      fmt.Sprintf("device-%d", i),
			UserId:        fmt.Sprintf("user-%d", i),
			TenantId:      "bingo",
			DeviceType:    consts.LinuxDevice,
			DeviceVersion: "1.0.0",
			State:         consts.StatusOffline,
		})
		if err != nil {
			t.Fatal(err)
			return
		}
	}
}

func TestSaveGroup(t *testing.T) {
	for i := 1; i <= 10000; i++ {
		err := db.SaveBiz(&model.Biz{
			BizId:      fmt.Sprintf("group-%d", i),
			BizType:    consts.BizGroup,
			State:      consts.StatusEnabled,
			BizName:    fmt.Sprintf("群组-%d", i),
			TenantId:   "bingo",
			TenantName: "品高软件",
		})
		if err != nil {
			t.Fatal(err)
			return
		}
	}
}

func TestAppendBizMember(t *testing.T) {
	for i := 1; i <= 1000; i++ {
		for j := 1; j <= 100; j++ {
			err := db.InsertBizMember(&model.BizMember{
				BizId:          fmt.Sprintf("group-%d", i),
				MemberId:       fmt.Sprintf("user-%d", j),
				MemberType:     consts.BizUser,
				BizTenantId:    "bingo",
				MemberTenantId: "bingo",
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		for j := 1; j <= 100; j++ {
			err := db.InsertBizMember(&model.BizMember{
				BizId:          fmt.Sprintf("group-%d", i+1000),
				MemberId:       fmt.Sprintf("user-%d", j),
				MemberType:     consts.BizUser,
				BizTenantId:    "bingo",
				MemberTenantId: "bingo",
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		for j := 1; j <= 1000; j++ {
			err := db.InsertBizMember(&model.BizMember{
				BizId:          fmt.Sprintf("group-%d", i+2000),
				MemberId:       fmt.Sprintf("user-%d", j),
				MemberType:     consts.BizUser,
				BizTenantId:    "bingo",
				MemberTenantId: "bingo",
			})
			if err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestListBizs(t *testing.T) {
	filter := map[string]interface{}{"tenant_id": "bingo", "biz_type": consts.BizUser}
	bizs, total, err := db.ListBizs(filter, 10, 1000)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(total)
	for _, biz := range bizs {
		t.Log(biz.BizId)
	}
}

func TestListGroups(t *testing.T) {
	filter := map[string]interface{}{"tenant_id": "bingo", "biz_type": consts.BizGroup}
	bizs, total, err := db.ListBizs(filter, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(total)
	for _, biz := range bizs {
		t.Log(biz.BizId)
	}
}

func TestListDevices(t *testing.T) {
	filter := map[string]interface{}{"tenant_id": "bingo", "state": consts.StatusOnline}
	devices, total, err := db.ListDevices(filter, 20, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(total)
	for _, device := range devices {
		t.Logf("%+v", device)
	}
}

func TestGetDevicesByUser(t *testing.T) {
	devices, err := db.GetDevicesByUser("user-1", "bingo")
	if err != nil {
		t.Fatal(err)
	}
	for _, device := range devices {
		t.Logf("%+v", device)
	}
}
