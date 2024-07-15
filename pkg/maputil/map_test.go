package maputil

import (
	"testing"

	"github.com/gzericlee/eim/pkg/jsonutil"
)

func TestGetInterface(t *testing.T) {
	jsonStr := `{
    "code": 0,
    "data": {
        "currentPage": 1,
        "from": 1,
        "items": [
            {
                "addr": "10.200.21.118:3306",
                "backupInsAddr": 1
            }
        ],
        "pageSize": 10,
        "to": 1,
        "total": 1,
        "totalPage": 1
    },
    "message": "OK",
    "taskId": null
}`
	var tmp map[string]interface{}
	tmp, err := jsonutil.ToAny[map[string]interface{}](jsonStr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(GetAny[string](tmp, "data", "items"))
	t.Log(GetAny[int](tmp, "data", "pageSize"))
}

func TestMerge(t *testing.T) {
	map1 := map[string]string{"Li": "Rui"}
	map2 := map[string]string{"Cui": "YanZhu"}
	map3 := Merge[string, string](map1, map2)
	t.Log(map3)
}
