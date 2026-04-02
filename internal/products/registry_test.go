package products

import "testing"

func TestRegistryGet(t *testing.T) {
	registry := NewRegistry()

	product, err := registry.Get("workbuddy")
	if err != nil {
		t.Fatalf("Get(workbuddy) 不应报错，实际：%v", err)
	}
	if product.Key() != "workbuddy" {
		t.Fatalf("产品 key 不符合预期，实际：%s", product.Key())
	}

	autoClaw, err := registry.Get("autoclaw")
	if err != nil {
		t.Fatalf("Get(autoclaw) 不应报错，实际：%v", err)
	}
	if autoClaw.Key() != "autoclaw" {
		t.Fatalf("AutoClaw 产品 key 不符合预期，实际：%s", autoClaw.Key())
	}

	qoderWork, err := registry.Get("qoderwork")
	if err != nil {
		t.Fatalf("Get(qoderwork) 不应报错，实际：%v", err)
	}
	if qoderWork.Key() != "qoderwork" {
		t.Fatalf("QoderWork 产品 key 不符合预期，实际：%s", qoderWork.Key())
	}

	if _, err := registry.Get("unknown"); err == nil {
		t.Fatalf("Get(unknown) 应返回错误")
	}
}

func TestRegistryKeys(t *testing.T) {
	registry := NewRegistry()
	keys := registry.Keys()

	if len(keys) != 3 {
		t.Fatalf("当前 keys 数量不符合预期，实际：%v", keys)
	}
	want := []string{"autoclaw", "qoderwork", "workbuddy"}
	for index, key := range want {
		if keys[index] != key {
			t.Fatalf("keys[%d] 不符合预期，实际：%v", index, keys)
		}
	}
}
