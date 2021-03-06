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

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package rabbitmq

import (
	"github.com/borischernov/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "rabbitmq", Asset); err != nil {
		panic(err)
	}
}

// Asset returns asset data
func Asset() string {
	return "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"
}
