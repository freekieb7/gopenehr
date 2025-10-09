package v2

type EHR struct {
	ID        string
	SubjectID string
	Folders   []Folder
}

type Folder struct {
	Name    string
	Folders []Folder
}

// Meta data about the service emmitting the data (e.g., service name, environment, host).
type Resource struct {
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Context from the instrumentation scope (e.g., which logger or module produced it).
type Scope struct {
	Name       string            `json:"name"`
	Version    string            `json:"version,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type EventData struct {
	ResourceEvents []ResourceEvent `json:"resource_events"`
}

type ResourceEvent struct {
	Resource  Resource     `json:"resource"`
	Scope     []ScopeEvent `json:"scope_events"`
	SchemaURL string       `json:"schema_url,omitempty"` // The Schema URL, if known. This is the identifier of the Schema that the resource data is recorded in.
}

type ScopeEvent struct {
	Scope     Scope   `json:"scope"`
	Events    []Event `json:"events"`
	SchemaURL string  `json:"schema_url,omitempty"` // The Schema URL, if known. This is the identifier of the Schema that the event data is recorded in.
}

type Event struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Data        isEventData       `json:"data"`
	MetaData    map[string]string `json:"meta_data,omitempty"`
}

type isEventData interface {
	isEventData()
}

type Observation struct {
	TimeUnixNano int64             `json:"time_unix_nano"` // The time at which the event occurred. The value is in Unix Epoch time in nanoseconds since 1970-01-01 00:00:00 +0000 UTC.
	Attributes   map[string]string `json:"attributes,omitempty"`
}

func (Observation) isEventData() {}

type Evaluation struct {
	Attributes map[string]string `json:"attributes,omitempty"`
}

func (Evaluation) isEventData() {}

type Instruction struct {
	Activities []Activity        `json:"activities,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

func (Instruction) isEventData() {}

type Activity struct {
	Description string            `json:"description,omitempty"`
	Attributes  map[string]string `json:"attributes,omitempty"`
}

type Action struct {
	Attributes map[string]string `json:"attributes,omitempty"`
}

func (Action) isEventData() {}

// type MetricsData struct {
// 	ResourceMetrics []ResourceMetrics `json:"resource_metrics"`
// }

// type ResourceMetrics struct {
// 	Resource  Resource       `json:"resource"`
// 	Scope     []ScopeMetrics `json:"scope_metrics"`
// 	SchemaURL string         `json:"schema_url,omitempty"` // The Schema URL, if known. This is the identifier of the Schema that the resource data is recorded in.
// }

// type ScopeMetrics struct {
// 	Scope     Scope    `json:"scope"`
// 	Metrics   []Metric `json:"metrics"`
// 	SchemaURL string   `json:"schema_url,omitempty"` // The Schema URL, if known. This is the identifier of the Schema that the metric data is recorded in.
// }

// type Metric struct {
// 	Name        string            `json:"name"`
// 	Description string            `json:"description,omitempty"`
// 	Unit        string            `json:"unit,omitempty"`
// 	Data        isMetricData      `json:"data"`
// 	MetaData    map[string]string `json:"meta_data,omitempty"` // Additional metadata attributes that describe the metric. [Optional].
// }

// type isMetricData interface {
// 	isMetricData()
// }

// type Gauge struct {
// 	DataPoints []NumberDataPoint `json:"data_points"`
// }

// func (Gauge) isMetricData() {}

// type Sum struct {
// 	DataPoints             []NumberDataPoint `json:"data_points"`
// 	AggregationTemporality string            `json:"aggregation_temporality"` // The aggregation temporality of this sum. Valid values are: "AGGREGATION_TEMPORALITY_UNSPECIFIED", "AGGREGATION_TEMPORALITY_DELTA", "AGGREGATION_TEMPORALITY_CUMULATIVE"
// 	IsMonotonic            bool              `json:"is_monotonic"`            // True if the sum is monotonic.
// }

// func (Sum) isMetricData() {}

// type Histogram struct {
// 	DataPoints             []NumberDataPoint `json:"data_points"`
// 	AggregationTemporality string            `json:"aggregation_temporality"` // The aggregation temporality of this histogram. Valid values are: "AGGREGATION_TEMPORALITY_UNSPECIFIED", "AGGREGATION_TEMPORALITY_DELTA", "AGGREGATION_TEMPORALITY_CUMULATIVE"
// }

// func (Histogram) isMetricData() {}

// type NumberDataPoint struct {
// 	Attributes    map[string]string `json:"attributes,omitempty"`
// 	Value         float64           `json:"value"`
// 	TimeUnixNano  int64             `json:"time_unix_nano"` // The time at which the data point was recorded. The value is in Unix Epoch time in nanoseconds since 1970-01-01 00:00:00 +0000 UTC.
// 	StartUnixNano int64             `json:"start_unix_nano,omitempty"`
// }
