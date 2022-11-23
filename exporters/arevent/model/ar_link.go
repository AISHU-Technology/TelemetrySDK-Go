package model

type ARLink interface {
	GetTraceID() string
	GetSpanID() string

	//private()
}
