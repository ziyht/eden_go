package etimer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetRunPend(t *testing.T) {

	var run, pend int32

	js := &JobState{}

	pend, run = js.addPendRun(-1, -1)
	assert.Equal(t, int32(-1), pend)
	assert.Equal(t, int32(-1), run)
	
	pend, run = js.addPendRun(1, 1)
	assert.Equal(t, int32(0), pend)
	assert.Equal(t, int32(0), run)
	
	pend, run = js.addPendRun(1, 0)
	assert.Equal(t, int32(1), pend)
	assert.Equal(t, int32(0), run)
	
	pend, run = js.addPendRun(-1, 1)
	assert.Equal(t, int32(0), pend)
	assert.Equal(t, int32(1), run)
}