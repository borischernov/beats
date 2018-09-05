// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package mb

import (
	"time"

	"github.com/borischernov/beats/libbeat/beat"
	"github.com/borischernov/beats/libbeat/common"
)

// EventModifier is a function that can modifies an Event. This is typically
// used to apply transformations to an Event as it is converted to a
// beat.Event. An example is AddMetricSetInfo.
type EventModifier func(module, metricset string, event *Event)

// Event contains the data generated by a MetricSet.
type Event struct {
	RootFields      common.MapStr // Fields that will be added to the root of the event.
	ModuleFields    common.MapStr // Fields that will be namespaced under [module].
	MetricSetFields common.MapStr // Fields that will be namespaced under [module].[metricset].

	Index     string        // Index name prefix. If set overwrites the default prefix.
	Namespace string        // Fully qualified namespace to use for MetricSetFields.
	Timestamp time.Time     // Timestamp when the event data was collected.
	Error     error         // Error that occurred while collecting the event data.
	Host      string        // Host from which the data was collected.
	Took      time.Duration // Amount of time it took to collect the event data.
}

// BeatEvent returns a new beat.Event containing the data this Event. It does
// mutate the underlying data in the Event.
func (e *Event) BeatEvent(module, metricSet string, modifiers ...EventModifier) beat.Event {
	if e.RootFields == nil {
		e.RootFields = common.MapStr{}
	}

	for _, modify := range modifiers {
		modify(module, metricSet, e)
	}

	b := beat.Event{
		Timestamp: e.Timestamp,
		Fields:    e.RootFields,
	}

	if len(e.ModuleFields) > 0 {
		b.Fields.Put(module, e.ModuleFields)
		e.ModuleFields = nil
	}

	if len(e.MetricSetFields) > 0 {
		switch e.Namespace {
		case ".":
			// Add fields to root.
			b.Fields.DeepUpdate(e.MetricSetFields)
		case "":
			b.Fields.Put(module+"."+metricSet, e.MetricSetFields)
		default:
			b.Fields.Put(e.Namespace, e.MetricSetFields)
		}

		e.MetricSetFields = nil
	}

	// Set index prefix to overwrite default
	if e.Index != "" {
		b.Meta = common.MapStr{"index": e.Index}
	}

	if e.Error != nil {
		b.Fields["error"] = common.MapStr{
			"message": e.Error.Error(),
		}
	}

	return b
}

// AddMetricSetInfo is an EventModifier that adds information about the
// MetricSet that generated the event. It will always add the metricset and
// module names. And it will add the host, namespace, and rtt (round-trip time
// in microseconds) values if they are non-zero values.
//
//     "metricset": {
//       "host": "apache",
//       "module": "apache",
//       "name": "status",
//       "rtt": 115
//     }
func AddMetricSetInfo(module, metricset string, event *Event) {
	info := common.MapStr{
		"name":   metricset,
		"module": module,
	}
	if event.Host != "" {
		info["host"] = event.Host
	}
	if event.Took > 0 {
		info["rtt"] = event.Took / time.Microsecond
	}
	if event.Namespace != "" {
		info["namespace"] = event.Namespace
	}
	info = common.MapStr{
		"metricset": info,
	}

	if event.RootFields == nil {
		event.RootFields = info
	} else {
		event.RootFields.DeepUpdate(info)
	}
}

// TransformMapStrToEvent transforms a common.MapStr produced by MetricSet
// (like any MetricSet that does not natively produce a mb.Event). It accounts
// for the special key names and routes the data stored under those keys to the
// correct location in the event.
func TransformMapStrToEvent(module string, m common.MapStr, err error) Event {
	var (
		event = Event{RootFields: common.MapStr{}, Error: err}
	)

	for k, v := range m {
		switch k {
		case TimestampKey:
			switch ts := v.(type) {
			case time.Time:
				delete(m, TimestampKey)
				event.Timestamp = ts
			case common.Time:
				delete(m, TimestampKey)
				event.Timestamp = time.Time(ts)
			}
		case ModuleDataKey:
			delete(m, ModuleDataKey)
			event.ModuleFields, _ = tryToMapStr(v)
		case RTTKey:
			delete(m, RTTKey)
			if took, ok := v.(time.Duration); ok {
				event.Took = took
			}
		case NamespaceKey:
			delete(m, NamespaceKey)
			if ns, ok := v.(string); ok {
				// The _namespace value does not include the module name and
				// it is required in the mb.Event.Namespace value.
				event.Namespace = module + "." + ns
			}
		}
	}

	event.MetricSetFields = m
	return event
}

func tryToMapStr(v interface{}) (common.MapStr, bool) {
	switch m := v.(type) {
	case common.MapStr:
		return m, true
	case map[string]interface{}:
		return common.MapStr(m), true
	default:
		return nil, false
	}
}
