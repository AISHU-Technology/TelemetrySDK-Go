package common

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/instrumentation"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"time"
)

type ScopeMetric struct {
	Scope   *instrumentation.Scope `json:"Scope"`
	Metrics []*Metrics             `json:"Metrics"`
}

func AnyRobotScopeMetricFromScopeMetric(scopeMetric *metricdata.ScopeMetrics) *ScopeMetric {
	return &ScopeMetric{
		Scope:   &scopeMetric.Scope,
		Metrics: AnyRobotMetricsFromMetrics(scopeMetric.Metrics),
	}
}

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

type Metrics struct {
	Name        string          `json:"Name"`
	Description string          `json:"Description"`
	Unit        string          `json:"Unit"`
	GaugeInt    *Gauge[int64]   `json:"Gauge,omitempty"`
	GaugeFloat  *Gauge[float64] `json:"Gauge,omitempty"`
	SumInt      *Sum[int64]     `json:"Sum,omitempty"`
	SumFloat    *Sum[float64]   `json:"Sum,omitempty"`
	Histogram   *Histogram      `json:"Histogram,omitempty"`
}

func AnyRobotMetricFromMetric(metric *metricdata.Metrics) *Metrics {
	if gauge, ok := metric.Data.(metricdata.Gauge[int64]); ok {
		return &Metrics{
			Name:        metric.Name,
			Description: metric.Description,
			Unit:        string(metric.Unit),
			GaugeInt:    AnyRobotGaugeFromGaugeInt(gauge),
		}
	}
	if gauge, ok := metric.Data.(metricdata.Gauge[float64]); ok {
		return &Metrics{
			Name:        metric.Name,
			Description: metric.Description,
			Unit:        string(metric.Unit),
			GaugeFloat:  AnyRobotGaugeFromGaugeFloat(gauge),
		}
	}
	if sum, ok := metric.Data.(metricdata.Sum[int64]); ok {
		return &Metrics{
			Name:        metric.Name,
			Description: metric.Description,
			Unit:        string(metric.Unit),
			SumInt:      AnyRobotSumFromSumInt(sum),
		}
	}
	if sum, ok := metric.Data.(metricdata.Sum[float64]); ok {
		return &Metrics{
			Name:        metric.Name,
			Description: metric.Description,
			Unit:        string(metric.Unit),
			SumFloat:    AnyRobotSumFromSumFloat(sum),
		}
	}
	if histogram, ok := metric.Data.(metricdata.Histogram); ok {
		return &Metrics{
			Name:        metric.Name,
			Description: metric.Description,
			Unit:        string(metric.Unit),
			Histogram:   AnyRobotHistogramFromHistogram(histogram),
		}
	}
	return nil
}

func AnyRobotMetricsFromMetrics(metrics []metricdata.Metrics) []*Metrics {
	arMetrics := make([]*Metrics, 0, len(metrics))
	for _, value := range metrics {
		if arMetric := AnyRobotMetricFromMetric(&value); arMetric != nil {
			arMetrics = append(arMetrics, arMetric)
		}
	}
	return arMetrics
}

type Gauge[N int64 | float64] struct {
	DataPoints []*DataPoint[N] `json:"DataPoints"`
}
type Sum[N int64 | float64] struct {
	DataPoints  []*DataPoint[N]        `json:"DataPoints"`
	Temporality metricdata.Temporality `json:"Temporality"`
	IsMonotonic bool                   `json:"IsMonotonic"`
}

type DataPoint[N int64 | float64] struct {
	Attributes attribute.Set `json:"Attributes"`
	StartTime  time.Time     `json:"StartTime"`
	EndTime    time.Time     `json:"EndTime"`
	Int        int64         `json:"Int,omitempty"`
	Float      float64       `json:"Float,omitempty"`
}
type Histogram struct {
	DataPoints  []*HistogramDataPoint  `json:"DataPoints"`
	Temporality metricdata.Temporality `json:"Temporality"`
}
type HistogramDataPoint struct {
	Attributes   attribute.Set `json:"Attributes"`
	StartTime    time.Time     `json:"StartTime"`
	EndTime      time.Time     `json:"EndTime"`
	Count        uint64        `json:"Count"`
	Bounds       []float64     `json:"Bounds"`
	BucketCounts []uint64      `json:"BucketCounts"`
	Min          *float64      `json:"Min,omitempty"`
	Max          *float64      `json:"Max,omitempty"`
	Sum          float64       `json:"Sum"`
}

func IntDataPoint(dp metricdata.DataPoint[int64]) *DataPoint[int64] {
	return &DataPoint[int64]{
		Attributes: dp.Attributes,
		StartTime:  dp.StartTime,
		EndTime:    dp.Time,
		Int:        dp.Value,
	}
}

func FloatDataPoint(dp metricdata.DataPoint[float64]) *DataPoint[float64] {
	return &DataPoint[float64]{
		Attributes: dp.Attributes,
		StartTime:  dp.StartTime,
		EndTime:    dp.Time,
		Float:      dp.Value,
	}
}

func IntDataPoints(dps []metricdata.DataPoint[int64]) []*DataPoint[int64] {
	arDataPoint := make([]*DataPoint[int64], 0, len(dps))
	for _, value := range dps {
		arDataPoint = append(arDataPoint, IntDataPoint(value))
	}
	return arDataPoint
}

func FloatDataPoints(dps []metricdata.DataPoint[float64]) []*DataPoint[float64] {
	arDataPoint := make([]*DataPoint[float64], 0, len(dps))
	for _, value := range dps {
		arDataPoint = append(arDataPoint, FloatDataPoint(value))
	}
	return arDataPoint
}
func SingleHistogramDataPoint(dp metricdata.HistogramDataPoint) *HistogramDataPoint {
	return &HistogramDataPoint{
		Attributes:   dp.Attributes,
		StartTime:    dp.StartTime,
		EndTime:      dp.Time,
		Count:        dp.Count,
		Bounds:       dp.Bounds,
		BucketCounts: dp.BucketCounts,
		Min:          dp.Min,
		Max:          dp.Max,
		Sum:          dp.Sum,
	}
}
func HistogramDataPoints(dps []metricdata.HistogramDataPoint) []*HistogramDataPoint {
	arDataPoint := make([]*HistogramDataPoint, 0, len(dps))
	for _, value := range dps {
		arDataPoint = append(arDataPoint, SingleHistogramDataPoint(value))
	}
	return arDataPoint
}
func AnyRobotGaugeFromGaugeInt(gauge metricdata.Gauge[int64]) *Gauge[int64] {
	return &Gauge[int64]{
		DataPoints: IntDataPoints(gauge.DataPoints),
	}
}

func AnyRobotGaugeFromGaugeFloat(gauge metricdata.Gauge[float64]) *Gauge[float64] {
	return &Gauge[float64]{
		DataPoints: FloatDataPoints(gauge.DataPoints),
	}
}

func AnyRobotSumFromSumInt(sum metricdata.Sum[int64]) *Sum[int64] {
	return &Sum[int64]{
		DataPoints:  IntDataPoints(sum.DataPoints),
		Temporality: sum.Temporality,
		IsMonotonic: sum.IsMonotonic,
	}
}
func AnyRobotSumFromSumFloat(sum metricdata.Sum[float64]) *Sum[float64] {
	return &Sum[float64]{
		DataPoints:  FloatDataPoints(sum.DataPoints),
		Temporality: sum.Temporality,
		IsMonotonic: sum.IsMonotonic,
	}
}
func AnyRobotHistogramFromHistogram(histogram metricdata.Histogram) *Histogram {
	return &Histogram{
		DataPoints:  HistogramDataPoints(histogram.DataPoints),
		Temporality: histogram.Temporality,
	}
}
