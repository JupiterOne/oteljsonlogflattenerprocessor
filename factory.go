package oteljsonlogflattenerprocessor

import (
	"context"
	"errors"

	"github.com/jupiterone/oteljsonlogflattenerprocessor/internal/metadata"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/processor"
)

func NewFactory() processor.Factory {
	return processor.NewFactory(
		metadata.Type,
		createDefaultConfig,
		processor.WithLogs(createLogsProcessor, metadata.LogsStability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createLogsProcessor(_ context.Context, set processor.CreateSettings, cfg component.Config, nextConsumer consumer.Logs) (processor.Logs, error) {
	pCfg, ok := cfg.(*Config)
	if !ok {
		return nil, errors.New("could not initialize logs transform processor")
	}
	return newProcessor(pCfg, nextConsumer, set.Logger)
}
