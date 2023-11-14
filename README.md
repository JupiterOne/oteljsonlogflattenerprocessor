# OpenTelemetry JSON Log Flattener 
NewRelic doesn't support nested JSON attributes in OpenTelemetry logs. To help alleviate these issues, this processor will flatten JSON at the top level. It turns messages like this:
```json
{
  "name": "Test User",
  "address": {
    "street": "First Ave",
    "house": 1234
  },
  "occupants": [
    "Test User",
    "Test User 2",
    "Test User 3"
  ]
}
```
into
```json
{
  "name": "Test User",
  "address.street": "First Ave",
  "address.house": 1234,
  "occupants": "[\"Test User\", \"Test User 2\", \"Test User 3\"]"
}
```


## Using this processor
To add this processor to your OpenTelemetry Collector, follow the instructions for building a collector here: https://opentelemetry.io/docs/collector/custom-collector/

In your config, add the following:
```yaml
processors:
  - github.com/JupiterOne/otel-jsonlogflattenerprocessor
```

In your collector config, add the following processor:
```yaml
receivers:
  ...
exporters:
  ...
processors:
  - oteljsonlogflattenerprocessor
pipelines:
  logs:
    receivers: [...]
    processors: [oteljsonlogflattenerprocessor]
    exporters: [...]
```