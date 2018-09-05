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

package mysql

import (
	"github.com/borischernov/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "mysql", Asset); err != nil {
		panic(err)
	}
}

// Asset returns asset data
func Asset() string {
	return "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"
}
