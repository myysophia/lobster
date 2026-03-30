package tui

import "testing"

func TestNewModelForLobsterEntersProductSelect(t *testing.T) {
	model := newModel("")
	if model.screen != screenProductSelect {
		t.Fatalf("lobster 默认应进入产品选择页，实际：%s", model.screen)
	}
	if model.selectedProduct.Key != "workbuddy" {
		t.Fatalf("默认选中产品应为 workbuddy，实际：%s", model.selectedProduct.Key)
	}
}

func TestNewModelForDirectProductEntersWorkBuddyWelcome(t *testing.T) {
	model := newModel("workbuddy")
	if model.screen != screenWorkBuddyWelcome {
		t.Fatalf("指定 workbuddy 时应直接进入 WorkBuddy 欢迎页，实际：%s", model.screen)
	}
	if model.selectedProduct.Key != "workbuddy" {
		t.Fatalf("指定 workbuddy 时应绑定 workbuddy，实际：%s", model.selectedProduct.Key)
	}
}
