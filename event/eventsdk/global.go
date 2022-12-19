package eventsdk

import (
	environment "go.opentelemetry.io/otel/sdk/resource"
	"strings"
)

var (
	globalEventProvider = NewEventProvider()

	globalServiceName     = defaultServiceName()
	globalServiceVersion  = "UnknownServiceVersion"
	globalServiceInstance = "UnknownServiceInstance"
)

func defaultServiceName() string {
	attributes := environment.Default().Attributes()
	if len(attributes) > 0 {
		if attributes[0].Key == "service.name" {
			if v := strings.Split(attributes[0].Value.AsString(), "___"); len(v) >= 2 {
				return strings.Split(attributes[0].Value.AsString(), "___")[1]
			}
		}
	}
	return "UnknownServiceName"
}
