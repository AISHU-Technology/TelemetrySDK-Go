package common

import (
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

type AnyRobotMetric struct {
	Resource     *Resource      `json:"Resource"`
	ScopeMetrics []*ScopeMetric `json:"ScopeMetrics"`
}

func AnyRobotMetricFromResourceMetric(resourceMetric *metricdata.ResourceMetrics) *AnyRobotMetric {
	return &AnyRobotMetric{
		Resource:     AnyRobotResourceFromResource(resourceMetric.Resource),
		ScopeMetrics: AnyRobotScopeMetricsFromScopeMetrics(resourceMetric.ScopeMetrics),
	}
}

func AnyRobotMetricsFromResourceMetrics(resourceMetrics []*metricdata.ResourceMetrics) []*AnyRobotMetric {
	arMetric := make([]*AnyRobotMetric, 0, len(resourceMetrics))
	for _, value := range resourceMetrics {
		arMetric = append(arMetric, AnyRobotMetricFromResourceMetric(value))
	}
	return arMetric
}
