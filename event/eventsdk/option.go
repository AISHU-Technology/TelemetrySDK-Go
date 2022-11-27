package eventsdk

// EventProviderOption EventProvider初始化选项。
type EventProviderOption interface {
	// apply 更改EventProvider默认配置。
	apply(eventProviderConfig) eventProviderConfig
}
