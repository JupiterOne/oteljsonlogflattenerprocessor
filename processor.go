package oteljsonlogflattenerprocessor

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/jupiterone/oteljsonlogflattenerprocessor/internal/errwrap"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor"
	"go.uber.org/zap"
)

var _ processor.Logs = (*jsonlogflattenerProcessor)(nil)

type jsonlogflattenerProcessor struct {
	config *Config
	next   consumer.Logs
	logger *zap.Logger
}

func newProcessor(config *Config, next consumer.Logs, logger *zap.Logger) (*jsonlogflattenerProcessor, error) {
	return &jsonlogflattenerProcessor{
		config: config,
		logger: logger,
		next:   next,
	}, nil
}

func (jlp *jsonlogflattenerProcessor) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func (jlp *jsonlogflattenerProcessor) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	if err := jlp.flatten(&ld); err != nil {
		jlp.logger.Error(err.Error())
	}
	return jlp.next.ConsumeLogs(ctx, ld)
}

func (jlp *jsonlogflattenerProcessor) Start(_ context.Context, host component.Host) error {
	return nil
}

func (jlp *jsonlogflattenerProcessor) Shutdown(context.Context) error {
	return nil
}

func (jlp *jsonlogflattenerProcessor) flatten(ld *plog.Logs) error {
	defer func() {
		// if we paniced for some reason, recover from the panic and log
		// the error message so we don't break the collector workflow
		if r := recover(); r != nil {
			jlp.logger.Error(r.(error).Error())
		}
	}()
	var rootErr error
	for i := 0; i < ld.ResourceLogs().Len(); i++ {
		for j := 0; j < ld.ResourceLogs().At(i).ScopeLogs().Len(); j++ {
			for k := 0; k < ld.ResourceLogs().At(i).ScopeLogs().At(j).LogRecords().Len(); k++ {
				for l, m := range ld.ResourceLogs().At(i).ScopeLogs().At(j).LogRecords().At(k).Attributes().AsRaw() {
					if m != nil && (reflect.TypeOf(m).Kind() == reflect.Map || reflect.TypeOf(m).Kind() == reflect.Slice) {
						payload, err := json.Marshal(m)
						if err != nil {
							jlp.logger.Error(err.Error())
							rootErr = errwrap.NewErrWrap(rootErr, err)
							continue
						}
						ld.ResourceLogs().At(i).ScopeLogs().At(j).LogRecords().At(k).Attributes().PutStr(l, string(payload))
					}
				}
			}
		}
	}
	return rootErr
}
