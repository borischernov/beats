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

// +build integration

package bgwriter

import (
	"testing"

	"github.com/borischernov/beats/libbeat/common"
	"github.com/borischernov/beats/libbeat/tests/compose"
	mbtest "github.com/borischernov/beats/metricbeat/mb/testing"
	"github.com/borischernov/beats/metricbeat/module/postgresql"

	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	compose.EnsureUp(t, "postgresql")

	f := mbtest.NewEventFetcher(t, getConfig())
	event, err := f.Fetch()
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Logf("%s/%s event: %+v", f.Module().Name(), f.Name(), event)

	assert.Contains(t, event, "checkpoints")
	assert.Contains(t, event, "buffers")
	assert.Contains(t, event, "stats_reset")

	checkpoints := event["checkpoints"].(common.MapStr)
	assert.Contains(t, checkpoints, "scheduled")
	assert.Contains(t, checkpoints, "requested")
	assert.Contains(t, checkpoints, "times")

	buffers := event["buffers"].(common.MapStr)
	assert.Contains(t, buffers, "checkpoints")
	assert.Contains(t, buffers, "clean")
	assert.Contains(t, buffers, "clean_full")
	assert.Contains(t, buffers, "backend")
	assert.Contains(t, buffers, "backend_fsync")
	assert.Contains(t, buffers, "allocated")
}

func TestData(t *testing.T) {
	compose.EnsureUp(t, "postgresql")

	f := mbtest.NewEventFetcher(t, getConfig())

	err := mbtest.WriteEvent(f, t)
	if err != nil {
		t.Fatal("write", err)
	}
}

func getConfig() map[string]interface{} {
	return map[string]interface{}{
		"module":     "postgresql",
		"metricsets": []string{"bgwriter"},
		"hosts":      []string{postgresql.GetEnvDSN()},
		"username":   postgresql.GetEnvUsername(),
		"password":   postgresql.GetEnvPassword(),
	}
}
