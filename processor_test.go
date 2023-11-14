package oteljsonlogflattenerprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor/processortest"
)

func TestFlatten(t *testing.T) {
	factory := NewFactory()
	tln := new(consumertest.LogsSink)
	jlfp, err := factory.CreateLogsProcessor(context.Background(), processortest.NewNopCreateSettings(), &Config{}, tln)
	require.NoError(t, err)
	assert.True(t, jlfp.Capabilities().MutatesData)

	err = jlfp.Start(context.Background(), nil)
	require.NoError(t, err)

	ld := buildLogs(1)

	err = jlfp.ConsumeLogs(context.Background(), ld)
	require.NoError(t, err)

	logs := tln.AllLogs()
	require.Len(t, logs, 1)
	attrs := logs[0].ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(0).Attributes()
	require.Equal(t, attrs.Len(), 6)
	fooString, exists := attrs.Get("foo-string")
	require.True(t, exists)
	require.Equal(t, fooString.Str(), "bar")
	fooDouble, exists := attrs.Get("foo-double")
	require.True(t, exists)
	require.Equal(t, fooDouble.Double(), 1.0)
	fooInt, exists := attrs.Get("foo-int")
	require.True(t, exists)
	require.Equal(t, fooInt.Int(), int64(1))
	fooBool, exists := attrs.Get("foo-bool")
	require.True(t, exists)
	require.True(t, fooBool.Bool())
	fooMap, exists := attrs.Get("foo-map")
	require.True(t, exists)
	require.Equal(t, fooMap.Type(), pcommon.ValueTypeStr)
	require.Equal(t, fooMap.Str(), `{"foo-map-bool":true,"foo-map-double":1,"foo-map-int":1,"foo-map-string":"bar"}`)
	fooSlice, exists := attrs.Get("foo-slice")
	require.True(t, exists)
	require.Equal(t, fooSlice.Type(), pcommon.ValueTypeStr)
	require.Equal(t, fooSlice.Str(), `[1,2]`)
}

func BenchmarkFlatten10(b *testing.B) {
	factory := NewFactory()
	tln := new(consumertest.LogsSink)
	jlfp, _ := factory.CreateLogsProcessor(context.Background(), processortest.NewNopCreateSettings(), &Config{}, tln)
	ld := buildLogs(10)
	for n := 0; n < b.N; n++ {
		jlfp.ConsumeLogs(context.Background(), ld)
	}
}

func BenchmarkFlatten500(b *testing.B) {
	factory := NewFactory()
	tln := new(consumertest.LogsSink)
	jlfp, _ := factory.CreateLogsProcessor(context.Background(), processortest.NewNopCreateSettings(), &Config{}, tln)
	ld := buildLogs(500)
	for n := 0; n < b.N; n++ {
		jlfp.ConsumeLogs(context.Background(), ld)
	}
}

func BenchmarkFlatten1000(b *testing.B) {
	factory := NewFactory()
	tln := new(consumertest.LogsSink)
	jlfp, _ := factory.CreateLogsProcessor(context.Background(), processortest.NewNopCreateSettings(), &Config{}, tln)
	ld := buildLogs(1000)
	for n := 0; n < b.N; n++ {
		jlfp.ConsumeLogs(context.Background(), ld)
	}
}

func buildLogs(count int) plog.Logs {
	ld := plog.NewLogs()
	for i := 0; i < count; i++ {
		ld.ResourceLogs().AppendEmpty()
		scope := ld.ResourceLogs().At(i).ScopeLogs().AppendEmpty()
		log := scope.LogRecords().AppendEmpty()

		log.Attributes().PutStr("foo-string", "bar")
		log.Attributes().PutDouble("foo-double", 1.0)
		log.Attributes().PutInt("foo-int", 1)
		log.Attributes().PutBool("foo-bool", true)
		inputFooMap := log.Attributes().PutEmptyMap("foo-map")
		inputFooMap.PutStr("foo-map-string", "bar")
		inputFooMap.PutDouble("foo-map-double", 1.0)
		inputFooMap.PutInt("foo-map-int", 1)
		inputFooMap.PutBool("foo-map-bool", true)
		inputFooSlice := log.Attributes().PutEmptySlice("foo-slice")
		inputFooSlice.AppendEmpty().SetInt(1)
		inputFooSlice.AppendEmpty().SetInt(2)
	}
	return ld
}
