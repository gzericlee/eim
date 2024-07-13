package database

import (
	"fmt"
	"sync"
	"testing"

	"github.com/panjf2000/ants"

	"eim/internal/model"
	"eim/internal/model/consts"
)

var db IDatabase

func init() {
	var err error
	//db, err = NewDatabase(MongoDBDriver, []string{"mongodb://admin:pass%40word1@127.0.0.1:27017/?authSource=admin&connect=direct"}, "eim")
	db, err = NewDatabase(PostgresDriver, []string{
		"postgres://eim:pass@word1@127.0.0.1:15430/eim?sslmode=disable",
		"postgres://eim:pass@word1@127.0.0.1:15431/eim?sslmode=disable",
		"postgres://eim:pass@word1@127.0.0.1:15432/eim?sslmode=disable",
	}, "eim")
	if err != nil {
		panic(err)
	}
}

func TestInsertTenant(t *testing.T) {
	err := db.InsertTenant(&model.Tenant{
		TenantId:   "bingo",
		TenantName: "广州市品高软件股份有限公司",
		State:      consts.StatusEnabled,
		Attributes: map[string]string{
			consts.FileflexEnabled: fmt.Sprintf("%v", true),
			consts.FileflexBucket:  "bingo",
			consts.FileflexUser:    "bingo",
			consts.FileflexPasswd:  "pass@word1",
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestInsertBiz(t *testing.T) {
	wg := sync.WaitGroup{}
	pool, _ := ants.NewPool(100)
	for i := 1; i <= 1000000; i++ {
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()
			err := db.InsertBiz(&model.Biz{
				BizId:      fmt.Sprintf("user-%d", i),
				BizType:    consts.BizUser,
				BizName:    fmt.Sprintf("用户-%d", i),
				TenantId:   "bingo",
				TenantName: "品高软件",
				State:      consts.StatusEnabled,
				Attributes: map[string]string{
					"password": "pass@word1",
				},
			})
			if err != nil {
				t.Fatal(err)
			}
		})
	}
	wg.Wait()
}

func TestInsertDevice(t *testing.T) {
	wg := sync.WaitGroup{}
	pool, _ := ants.NewPool(1000)
	for i := 1; i <= 1000000; i++ {
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()
			err := db.InsertDevice(&model.Device{
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
		})
	}
	wg.Wait()
}

func TestInsertGroup(t *testing.T) {
	wg := sync.WaitGroup{}
	pool, _ := ants.NewPool(1000)
	for i := 1; i <= 10000; i++ {
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()
			err := db.InsertBiz(&model.Biz{
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
		})
	}
	wg.Wait()
}

func TestAppendBizMember(t *testing.T) {
	wg := sync.WaitGroup{}
	pool, _ := ants.NewPool(1000)
	for i := 1; i <= 1000; i++ {
		for j := 1; j <= 100; j++ {
			wg.Add(1)
			pool.Submit(func() {
				defer wg.Done()
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
			})
		}
		for j := 1; j <= 100; j++ {
			wg.Add(1)
			pool.Submit(func() {
				defer wg.Done()
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
			})
		}
		for j := 1; j <= 1000; j++ {
			wg.Add(1)
			pool.Submit(func() {
				defer wg.Done()
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
			})
		}
	}
	wg.Wait()
}

func TestGetTenant(t *testing.T) {
	tenant, err := db.GetTenant("bingo")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", tenant)
}

func TestGetBiz(t *testing.T) {
	biz, err := db.GetBiz("user-1", "bingo")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", biz)
}

func TestGetDevice(t *testing.T) {
	device, err := db.GetDevice("user-1", "bingo", "device-1")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", device)
}

func TestUpdateDevice(t *testing.T) {
	device, err := db.GetDevice("user-1", "bingo", "device-1")
	if err != nil {
		t.Fatal(err)
	}
	device.State = consts.StatusOffline
	t.Logf("%+v", device)
}

func TestGetBizMembers(t *testing.T) {
	members, err := db.GetBizMembers("group-1", "bingo")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", members)
}

func TestListBizs(t *testing.T) {
	filter := map[string]interface{}{"tenant_id": "bingo", "biz_type": consts.BizGroup}
	bizs, total, err := db.ListBizs(filter, []string{"created_at DESC"}, 10, 0)
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
	bizs, total, err := db.ListBizs(filter, []string{}, 10, 0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(total)
	for _, biz := range bizs {
		t.Log(biz.BizId)
	}
}

func TestListDevices(t *testing.T) {
	filter := map[string]interface{}{"tenant_id": "bingo", "state": consts.StatusOffline}
	devices, total, err := db.ListDevices(filter, nil, 10, 0)
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
