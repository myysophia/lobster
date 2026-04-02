package tui

type productItem struct {
	Key         string
	DisplayName string
	Summary     string
	Available   bool
	Badge       string
}

func defaultProducts() []productItem {
	return []productItem{
		{
			Key:         "workbuddy",
			DisplayName: "WorkBuddy",
			Summary:     "腾讯 WorkBuddy 安装与首次启动向导，包含校验与下一步建议。",
			Available:   true,
			Badge:       "可用",
		},
		{
			Key:         "arkclaw",
			DisplayName: "ArkClaw",
			Summary:     "统一入口已预留，安装能力正在打磨中，稍后上线。",
			Available:   false,
			Badge:       "On The Way",
		},
		{
			Key:         "kimi-claw",
			DisplayName: "Kimi Claw",
			Summary:     "统一入口已预留，安装能力正在路上，等待工程验证。",
			Available:   false,
			Badge:       "On The Way",
		},
		{
			Key:         "autoclaw",
			DisplayName: "AutoClaw",
			Summary:     "统一入口已预留，安装能力正在路上，后续版本会接入。",
			Available:   false,
			Badge:       "On The Way",
		},
		{
			Key:         "qoderwork",
			DisplayName: "QoderWork",
			Summary:     "统一入口已预留，当前版本先接入产品骨架，安装能力后续补齐。",
			Available:   false,
			Badge:       "On The Way",
		},
	}
}

func findProduct(items []productItem, key string) productItem {
	for _, item := range items {
		if item.Key == key {
			return item
		}
	}
	return items[0]
}
