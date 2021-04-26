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
	TraceID string
	Id      string
}
