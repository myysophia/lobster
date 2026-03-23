package products

import (
	"fmt"
	"sort"
)

type Registry struct {
	products map[string]Product
}

func NewRegistry() *Registry {
	items := []Product{
		NewWorkBuddy(),
	}

	m := make(map[string]Product, len(items))
	for _, item := range items {
		m[item.Key()] = item
	}

	return &Registry{products: m}
}

func (r *Registry) Get(key string) (Product, error) {
	product, ok := r.products[key]
	if !ok {
		return nil, fmt.Errorf("不支持的产品: %s", key)
	}
	return product, nil
}

func (r *Registry) Keys() []string {
	keys := make([]string, 0, len(r.products))
	for key := range r.products {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
