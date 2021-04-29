package field

import "time"

type ExternalSpanField struct {
	// Method                            string
	// Url                               string
	// Target                            string
	// Host                              string
	// Status                            int
	// Scheme                            string
	// StatusCode                        int
	// Flavor                            int
	// UserAgent                         string
	// RequestContentLength              int
	// RequestContentLengthUncompressed  int
	// ResponeContentLength              int
	// ResponseContentLengthUncompressed int
	Attributes StructField
	StartTime  time.Time
	EndTime    time.Time
	// Attrs      [][2]string
	traceID          string
	id               string
	parentID         string
	internalParentID string
}

func (f *ExternalSpanField) TraceID() string {
	return f.traceID
}

func (f *ExternalSpanField) ID() string {
	return f.id
}

func (f *ExternalSpanField) ParentID() string {
	return f.parentID
}

func (f *ExternalSpanField) InternalParentID() string {
	return f.internalParentID
}
