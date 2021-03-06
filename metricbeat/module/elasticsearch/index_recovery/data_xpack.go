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

package index_recovery

import (
	"encoding/json"
	"time"

	"github.com/borischernov/beats/libbeat/common"
	"github.com/borischernov/beats/metricbeat/helper/elastic"
	"github.com/borischernov/beats/metricbeat/mb"
	"github.com/borischernov/beats/metricbeat/module/elasticsearch"
)

func eventsMappingXPack(r mb.ReporterV2, m *MetricSet, content []byte) error {
	var data map[string]interface{}
	err := json.Unmarshal(content, &data)
	if err != nil {
		return err
	}

	var results []map[string]interface{}
	for indexName, indexData := range data {
		indexData, ok := indexData.(map[string]interface{})
		if !ok {
			continue
		}

		shards, ok := indexData["shards"]
		if !ok {
			continue
		}

		shardsArr, ok := shards.([]interface{})
		if !ok {
			continue
		}

		for _, shard := range shardsArr {
			shard, ok := shard.(map[string]interface{})
			if !ok {
				continue
			}

			shard["index_name"] = indexName
			results = append(results, shard)
		}
	}

	indexRecovery := common.MapStr{}
	indexRecovery["shards"] = results

	info, err := elasticsearch.GetInfo(m.HTTP, m.HTTP.GetURI())
	if err != nil {
		return err
	}

	event := mb.Event{}
	event.RootFields = common.MapStr{
		"cluster_uuid":   info.ClusterID,
		"timestamp":      common.Time(time.Now()),
		"interval_ms":    m.Module().Config().Period / time.Millisecond,
		"type":           "index_recovery",
		"index_recovery": indexRecovery,
	}

	event.Index = elastic.MakeXPackMonitoringIndexName(elastic.Elasticsearch)
	r.Event(event)

	return nil
}
