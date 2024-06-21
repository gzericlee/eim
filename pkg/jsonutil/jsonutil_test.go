package jsonutil

import "testing"

func TestToStruct(t *testing.T) {
	t.Log(ToIndentAny[string]([]byte(`{"name":"lirui"}`)))
}

func TestDeepCopy(t *testing.T) {
	type demo struct {
		Name string
	}
	d1 := &demo{Name: "L"}
	d2, _ := DeepCopy[*demo](d1)
	d2.Name = "I"
	t.Log(&d1, &d2)
}
