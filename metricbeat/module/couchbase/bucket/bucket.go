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

package bucket

import (
	"github.com/borischernov/beats/libbeat/common"
	"github.com/borischernov/beats/libbeat/common/cfgwarn"
	"github.com/borischernov/beats/metricbeat/helper"
	"github.com/borischernov/beats/metricbeat/mb"
	"github.com/borischernov/beats/metricbeat/mb/parse"
)

const (
	defaultScheme = "http"
	defaultPath   = "/pools/default/buckets"
)

var (
	hostParser = parse.URLHostParserBuilder{
		DefaultScheme: defaultScheme,
		DefaultPath:   defaultPath,
	}.Build()
)

// init registers the MetricSet with the central registry.
// The New method will be called after the setup of the module and before starting to fetch data
func init() {
	mb.Registry.MustAddMetricSet("couchbase", "bucket", New,
		mb.WithHostParser(hostParser),
		mb.DefaultMetricSet(),
	)
}

// MetricSet type defines all fields of the MetricSet
type MetricSet struct {
	mb.BaseMetricSet
	http *helper.HTTP
}

// New create a new instance of the MetricSet
func New(base mb.BaseMetricSet) (mb.MetricSet, error) {
	cfgwarn.Beta("The couchbase bucket metricset is beta")

	http, err := helper.NewHTTP(base)
	if err != nil {
		return nil, err
	}
	return &MetricSet{
		BaseMetricSet: base,
		http:          http,
	}, nil
}

// Fetch methods implements the data gathering and data conversion to the right format
// It returns the event which is then forward to the output. In case of an error, a
// descriptive error must be returned.
func (m *MetricSet) Fetch() ([]common.MapStr, error) {
	content, err := m.http.FetchContent()
	if err != nil {
		return nil, err
	}
	return eventsMapping(content), nil
}
