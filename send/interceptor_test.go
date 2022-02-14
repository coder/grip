package send

import (
	"testing"

	"cdr.dev/grip/level"
	"cdr.dev/grip/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterceptor(t *testing.T) {
	base, err := NewInternalLogger("test", LevelInfo{Threshold: level.Info, Default: level.Debug})
	require.NoError(t, err)

	var count int
	filter := func(m message.Composer) { count++ }

	icept := NewInterceptor(base, filter)

	assert.Equal(t, 0, base.Len())
	icept.Send(message.NewSimpleStringMessage(level.Info, "hello"))
	assert.Equal(t, 1, base.Len())
	assert.Equal(t, 1, count)

	icept.Send(message.NewSimpleStringMessage(level.Trace, "hello"))
	assert.Equal(t, 1, base.Len())
	assert.Equal(t, 2, count)
}
