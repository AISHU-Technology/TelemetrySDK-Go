package eventsdk

type EventProviderOption interface {
	apply(*eventProviderConfig) *eventProviderConfig
}
