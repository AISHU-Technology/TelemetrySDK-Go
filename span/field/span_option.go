package field

import "context"

type LogOptionFunc func(*logSpanV1)

func WithAttribute(attr *attribute) LogOptionFunc {
	return func(l *logSpanV1) {
		if attr != nil {
			if attr == nil || attr.Message == nil {
				return
			}
			record := MallocStructField(2)
			record.Set(attr.Type, attr.Message)
			record.Set("Type", StringField(attr.Type))
			l.attributes = record
		}
	}
}

func WithContext(ctx context.Context) LogOptionFunc {
	return func(l *logSpanV1) {
		if ctx != nil {
			l.ctx = ctx
		}
	}
}
