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

	if _, err := registry.Get("unknown"); err == nil {
		t.Fatalf("Get(unknown) 应返回错误")
	}
}

func TestRegistryKeys(t *testing.T) {
	registry := NewRegistry()
	keys := registry.Keys()

	if len(keys) != 1 {
		t.Fatalf("当前 keys 数量不符合预期，实际：%v", keys)
	}
	if keys[0] != "workbuddy" {
		t.Fatalf("keys[0] 不符合预期，实际：%s", keys[0])
	}
}
