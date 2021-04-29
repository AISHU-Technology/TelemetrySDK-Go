package field

import (
	"testing"

	"gotest.tools/assert"
)

func TestExternalsSpanGetID(t *testing.T) {
	ID := GenTraceID()
	es := &ExternalSpanField{
		traceID:          ID,
		id:               ID,
		parentID:         ID,
		internalParentID: ID,
	}

	assert.Equal(t, ID, es.TraceID())
	assert.Equal(t, ID, es.ID())
	assert.Equal(t, ID, es.ParentID())
	assert.Equal(t, ID, es.InternalParentID())

}
