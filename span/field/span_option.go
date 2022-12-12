package field

import "context"

type LogOptionFunc func(*logSpanV1)

func WithAttribute(attr *attribute) LogOptionFunc {
	return func(l *logSpanV1) {
		if attr == nil || attr.Message == nil {
			return
		}
		// record := MallocStructField(1)
		// record.Set(attr.Type, attr.Message)
		//record.Set("Type", StringField(attr.Type))
		mapAttributes := MapField(map[string]interface{}{attr.Type: attr.Message})
		l.attributes = mapAttributes
	}
}

func WithContext(ctx context.Context) LogOptionFunc {
	return func(l *logSpanV1) {
		if ctx != nil {
			l.ctx = ctx
		}
	}
}
