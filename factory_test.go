package oteljsonlogflattenerprocessor

import (
	"testing"

	"github.com/jupiterone/oteljsonlogflattenerprocessor/internal/metadata"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
)

func TestFactoryType(t *testing.T) {
	factory := NewFactory()
	assert.Equal(t, factory.Type(), component.Type(metadata.Type))
}

func TestFactoryCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.Equal(t, cfg, &Config{})
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestFactoryCreateProcessorEmpty(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	err := component.ValidateConfig(cfg)
	assert.NoError(t, err)
}
