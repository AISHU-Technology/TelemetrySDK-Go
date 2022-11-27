package eventsdk

import (
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

func contextWithDone() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func TestGetEventProvider(t *testing.T) {
	tests := []struct {
		name string
		want EventProvider
	}{
		{
			"",
			globalEventProvider,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetEventProvider(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetEventProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEventProvider(t *testing.T) {
	type args struct {
		opts []EventProviderOption
	}
	tests := []struct {
		name string
		args args
		want EventProvider
	}{
		{
			"",
			args{nil},
			&eventProvider{
				RWLock:    sync.RWMutex{},
				Exporters: make(map[string]EventExporter),
				Ticker:    time.NewTicker(5 * time.Second),
				Limit:     9,
				Events:    make([]Event, 0, 10),
				Sent:      make(chan bool, 3),
				stopOnce:  sync.Once{},
				stopCh:    make(chan struct{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEventProvider(tt.args.opts...); !reflect.DeepEqual(got.Shutdown(context.Background()), tt.want.Shutdown(context.Background())) {
				t.Errorf("NewEventProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetEventProvider(t *testing.T) {
	type args struct {
		ep EventProvider
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"",
			args{globalEventProvider},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetEventProvider(tt.args.ep)
		})
	}
}

func TestEventProviderForceFlash(t *testing.T) {
	exporterMap := make(map[string]EventExporter)
	exporterMap["DefaultExporter"] = GetDefaultExporter()
	type fields struct {
		RWLock    *sync.RWMutex
		Exporters map[string]EventExporter
		Ticker    *time.Ticker
		Limit     int
		Events    []Event
		Sent      chan bool
		stopOnce  *sync.Once
		stopCh    chan struct{}
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"",
			fields{
				RWLock:    &sync.RWMutex{},
				Exporters: nil,
				Ticker:    nil,
				Limit:     0,
				Events:    nil,
				Sent:      nil,
				stopOnce:  &sync.Once{},
				stopCh:    nil,
			},
			args{context.Background()},
			false,
		}, {
			"",
			fields{
				RWLock:    &sync.RWMutex{},
				Exporters: nil,
				Ticker:    nil,
				Limit:     0,
				Events:    make([]Event, 1, 2),
				Sent:      nil,
				stopOnce:  &sync.Once{},
				stopCh:    nil,
			},
			args{context.Background()},
			false,
		}, {
			"",
			fields{
				RWLock:    &sync.RWMutex{},
				Exporters: exporterMap,
				Ticker:    nil,
				Limit:     0,
				Events:    make([]Event, 1, 2),
				Sent:      nil,
				stopOnce:  &sync.Once{},
				stopCh:    nil,
			},
			args{contextWithDone()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &eventProvider{
				RWLock:    sync.RWMutex{},
				Exporters: tt.fields.Exporters,
				Ticker:    tt.fields.Ticker,
				Limit:     tt.fields.Limit,
				Events:    tt.fields.Events,
				Sent:      tt.fields.Sent,
				stopOnce:  sync.Once{},
				stopCh:    tt.fields.stopCh,
			}
			if err := ep.ForceFlash(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ForceFlash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventProviderShutdown(t *testing.T) {
	exporterMap := make(map[string]EventExporter)
	exporterMap["DefaultExporter"] = GetDefaultExporter()
	type fields struct {
		RWLock    *sync.RWMutex
		Exporters map[string]EventExporter
		Ticker    *time.Ticker
		Limit     int
		Events    []Event
		Sent      chan bool
		stopOnce  *sync.Once
		stopCh    chan struct{}
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"",
			fields{
				RWLock:    &sync.RWMutex{},
				Exporters: nil,
				Ticker:    nil,
				Limit:     0,
				Events:    nil,
				Sent:      nil,
				stopOnce:  &sync.Once{},
				stopCh:    make(chan struct{}),
			},
			args{context.Background()},
			false,
		}, {
			"",
			fields{
				RWLock:    &sync.RWMutex{},
				Exporters: exporterMap,
				Ticker:    nil,
				Limit:     0,
				Events:    nil,
				Sent:      nil,
				stopOnce:  &sync.Once{},
				stopCh:    make(chan struct{}),
			},
			args{contextWithDone()},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &eventProvider{
				RWLock:    sync.RWMutex{},
				Exporters: tt.fields.Exporters,
				Ticker:    tt.fields.Ticker,
				Limit:     tt.fields.Limit,
				Events:    tt.fields.Events,
				Sent:      tt.fields.Sent,
				stopOnce:  sync.Once{},
				stopCh:    tt.fields.stopCh,
			}
			if err := ep.Shutdown(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Shutdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEventProvideLoadEvent(t *testing.T) {
	type fields struct {
		RWLock    *sync.RWMutex
		Exporters map[string]EventExporter
		Ticker    *time.Ticker
		Limit     int
		Events    []Event
		Sent      chan bool
		stopOnce  *sync.Once
		stopCh    chan struct{}
	}
	type args struct {
		event Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{

		{
			"",
			fields{
				RWLock:    &sync.RWMutex{},
				Exporters: nil,
				Ticker:    nil,
				Limit:     0,
				Events:    []Event{},
				Sent:      make(chan bool, 1),
				stopOnce:  &sync.Once{},
				stopCh:    make(chan struct{}),
			},
			args{nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &eventProvider{
				RWLock:    sync.RWMutex{},
				Exporters: tt.fields.Exporters,
				Ticker:    tt.fields.Ticker,
				Limit:     tt.fields.Limit,
				Events:    tt.fields.Events,
				Sent:      tt.fields.Sent,
				stopOnce:  sync.Once{},
				stopCh:    tt.fields.stopCh,
			}
			ep.loadEvent(tt.args.event)
		})
	}
}

func TestEventProviderPrivate(t *testing.T) {
	type fields struct {
		RWLock    *sync.RWMutex
		Exporters map[string]EventExporter
		Ticker    *time.Ticker
		Limit     int
		Events    []Event
		Sent      chan bool
		stopOnce  *sync.Once
		stopCh    chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				RWLock:    &sync.RWMutex{},
				Exporters: nil,
				Ticker:    nil,
				Limit:     0,
				Events:    nil,
				Sent:      nil,
				stopOnce:  &sync.Once{},
				stopCh:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &eventProvider{
				RWLock:    sync.RWMutex{},
				Exporters: tt.fields.Exporters,
				Ticker:    tt.fields.Ticker,
				Limit:     tt.fields.Limit,
				Events:    tt.fields.Events,
				Sent:      tt.fields.Sent,
				stopOnce:  sync.Once{},
				stopCh:    tt.fields.stopCh,
			}
			ep.private()
		})
	}
}

func channelWithStop() chan struct{} {
	stopCh := make(chan struct{})
	close(stopCh)
	return stopCh
}

func TestEventProviderSendEvents(t *testing.T) {
	type fields struct {
		RWLock    *sync.RWMutex
		Exporters map[string]EventExporter
		Ticker    *time.Ticker
		Limit     int
		Events    []Event
		Sent      chan bool
		stopOnce  *sync.Once
		stopCh    chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				RWLock:    &sync.RWMutex{},
				Exporters: nil,
				Ticker:    time.NewTicker(time.Second),
				Limit:     0,
				Events:    nil,
				Sent:      nil,
				stopOnce:  &sync.Once{},
				stopCh:    channelWithStop(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &eventProvider{
				RWLock:    sync.RWMutex{},
				Exporters: tt.fields.Exporters,
				Ticker:    tt.fields.Ticker,
				Limit:     tt.fields.Limit,
				Events:    tt.fields.Events,
				Sent:      tt.fields.Sent,
				stopOnce:  sync.Once{},
				stopCh:    tt.fields.stopCh,
			}
			ep.sendEvents()
		})
	}
}

func TestEventProviderVerdictSent(t *testing.T) {
	type fields struct {
		RWLock    *sync.RWMutex
		Exporters map[string]EventExporter
		Ticker    *time.Ticker
		Limit     int
		Events    []Event
		Sent      chan bool
		stopOnce  *sync.Once
		stopCh    chan struct{}
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"",
			fields{
				RWLock:    &sync.RWMutex{},
				Exporters: nil,
				Ticker:    nil,
				Limit:     0,
				Events:    nil,
				Sent:      make(chan bool, 1),
				stopOnce:  &sync.Once{},
				stopCh:    nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ep := &eventProvider{
				RWLock:    sync.RWMutex{},
				Exporters: tt.fields.Exporters,
				Ticker:    tt.fields.Ticker,
				Limit:     tt.fields.Limit,
				Events:    tt.fields.Events,
				Sent:      tt.fields.Sent,
				stopOnce:  sync.Once{},
				stopCh:    tt.fields.stopCh,
			}
			ep.verdictSent()
		})
	}
}
