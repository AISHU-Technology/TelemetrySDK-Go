package common

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"time"
)

// ScopeMetric 自定义 ScopeMetric，改造了Metrics。
type ScopeMetric struct {
	Scope   *instrumentation.Scope `json:"Scope"`
	Metrics []*Metrics             `json:"Metrics"`
}

// AnyRobotScopeMetricFromScopeMetric 单条 *metricdata.ScopeMetrics 转换为 *ScopeMetric 。
func AnyRobotScopeMetricFromScopeMetric(scopeMetric *metricdata.ScopeMetrics) *ScopeMetric {
	return &ScopeMetric{
		Scope:   &scopeMetric.Scope,
		Metrics: AnyRobotMetricsFromMetrics(scopeMetric.Metrics),
	}
}

// AnyRobotScopeMetricsFromScopeMetrics 批量 []metricdata.ScopeMetrics 转换为 []*ScopeMetric 。
func AnyRobotScopeMetricsFromScopeMetrics(scopeMetrics []metricdata.ScopeMetrics) []*ScopeMetric {
	if scopeMetrics == nil {
		return make([]*ScopeMetric, 0)
	}
	arScopeMetrics := make([]*ScopeMetric, 0, len(scopeMetrics))
	for _, value := range scopeMetrics {
		arScopeMetrics = append(arScopeMetrics, AnyRobotScopeMetricFromScopeMetric(&value))
	}
	return arScopeMetrics
}

// Metrics 自定义 Metrics，改造了Data->Gauge/Sum/Histogram。
type Metrics struct {
	Name        string                `json:"Name"`
	Description string                `json:"Description"`
	Unit        string                `json:"Unit"`
	Gauge       *Gauge                `json:"Gauge,omitempty"`
	Sum         *Sum                  `json:"Sum,omitempty"`
	Histogram   *metricdata.Histogram `json:"Histogram,omitempty"`
}

// AnyRobotMetricFromMetric 单条 *metricdata.Metrics 转换为 *Metrics 。
func AnyRobotMetricFromMetric(metric *metricdata.Metrics) *Metrics {
	if gauge, ok := metric.Data.(metricdata.Gauge[int64]); ok {
		return &Metrics{
			Name:        metric.Name,
			Description: metric.Description,
			Unit:        string(metric.Unit),
			Gauge:       AnyRobotGaugeFromGaugeInt(gauge),
		}
	}
	if gauge, ok := metric.Data.(metricdata.Gauge[float64]); ok {
		return &Metrics{
			Name:        metric.Name,
			Description: metric.Description,
			Unit:        string(metric.Unit),
			Gauge:       AnyRobotGaugeFromGaugeFloat(gauge),
		}
	}
	if sum, ok := metric.Data.(metricdata.Sum[int64]); ok {
		return &Metrics{
			Name:        metric.Name,
			Description: metric.Description,
			Unit:        string(metric.Unit),
			Sum:         AnyRobotSumFromSumInt(sum),
		}
	}
	if sum, ok := metric.Data.(metricdata.Sum[float64]); ok {
		return &Metrics{
			Name:        metric.Name,
			Description: metric.Description,
			Unit:        string(metric.Unit),
			Sum:         AnyRobotSumFromSumFloat(sum),
		}
	}
	if histogram, ok := metric.Data.(metricdata.Histogram); ok {
		return &Metrics{
			Name:        metric.Name,
			Description: metric.Description,
			Unit:        string(metric.Unit),
			Histogram:   &histogram,
		}
	}
	return nil
}

// AnyRobotMetricsFromMetrics 批量 []metricdata.Metrics 转换为 []*Metrics 。
func AnyRobotMetricsFromMetrics(metrics []metricdata.Metrics) []*Metrics {
	arMetrics := make([]*Metrics, 0, len(metrics))
	for _, value := range metrics {
		if arMetric := AnyRobotMetricFromMetric(&value); arMetric != nil {
			arMetrics = append(arMetrics, arMetric)
		}
	}
	return arMetrics
}

// Gauge 自定义 Gauge，改造了DataPoints。
type Gauge struct {
	DataPoints []*DataPoint `json:"DataPoints"`
}

// Sum 自定义 Sum，改造了DataPoints。
type Sum struct {
	DataPoints  []*DataPoint           `json:"DataPoints"`
	Temporality metricdata.Temporality `json:"Temporality"`
	IsMonotonic bool                   `json:"IsMonotonic"`
}

// DataPoint 自定义 DataPoint，改造了Value->Int/Float。
type DataPoint struct {
	Attributes attribute.Set `json:"Attributes"`
	StartTime  time.Time     `json:"StartTime"`
	Time       time.Time     `json:"Time"`
	Int        *int64        `json:"Int,omitempty"`
	Float      *float64      `json:"Float,omitempty"`
}

// IntDataPoint 单条 metricdata.DataPoint[int64] 转换为 *DataPoint 。
func IntDataPoint(dp metricdata.DataPoint[int64]) *DataPoint {
	return &DataPoint{
		Attributes: dp.Attributes,
		StartTime:  dp.StartTime,
		Time:       dp.Time,
		Int:        &dp.Value,
	}
}

// FloatDataPoint 单条 metricdata.DataPoint[float64] 转换为 *DataPoint 。
func FloatDataPoint(dp metricdata.DataPoint[float64]) *DataPoint {
	return &DataPoint{
		Attributes: dp.Attributes,
		StartTime:  dp.StartTime,
		Time:       dp.Time,
		Float:      &dp.Value,
	}
}

// IntDataPoints 批量 []metricdata.DataPoint[int64] 转换为 []*DataPoint 。
func IntDataPoints(dps []metricdata.DataPoint[int64]) []*DataPoint {
	arDataPoint := make([]*DataPoint, 0, len(dps))
	for _, value := range dps {
		arDataPoint = append(arDataPoint, IntDataPoint(value))
	}
	return arDataPoint
}

// FloatDataPoints 批量 []metricdata.DataPoint[float64] 转换为 []*DataPoint 。
func FloatDataPoints(dps []metricdata.DataPoint[float64]) []*DataPoint {
	arDataPoint := make([]*DataPoint, 0, len(dps))
	for _, value := range dps {
		arDataPoint = append(arDataPoint, FloatDataPoint(value))
	}
	return arDataPoint
}

// AnyRobotGaugeFromGaugeInt 单条 metricdata.Gauge[int64] 转换为 *Gauge 。
func AnyRobotGaugeFromGaugeInt(gauge metricdata.Gauge[int64]) *Gauge {
	return &Gauge{
		DataPoints: IntDataPoints(gauge.DataPoints),
	}
}

// AnyRobotGaugeFromGaugeFloat 单条 metricdata.Gauge[float64] 转换为 *Gauge 。
func AnyRobotGaugeFromGaugeFloat(gauge metricdata.Gauge[float64]) *Gauge {
	return &Gauge{
		DataPoints: FloatDataPoints(gauge.DataPoints),
	}
}

// AnyRobotSumFromSumInt 单条 metricdata.Sum[int64] 转换为 *Sum 。
func AnyRobotSumFromSumInt(sum metricdata.Sum[int64]) *Sum {
	return &Sum{
		DataPoints:  IntDataPoints(sum.DataPoints),
		Temporality: sum.Temporality,
		IsMonotonic: sum.IsMonotonic,
	}
}

// AnyRobotSumFromSumFloat 单条 metricdata.Sum[float64] 转换为 *Sum 。
func AnyRobotSumFromSumFloat(sum metricdata.Sum[float64]) *Sum {
	return &Sum{
		DataPoints:  FloatDataPoints(sum.DataPoints),
		Temporality: sum.Temporality,
		IsMonotonic: sum.IsMonotonic,
	}
}
