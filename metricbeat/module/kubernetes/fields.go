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

package kubernetes

import (
	"github.com/borischernov/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("metricbeat", "kubernetes", Asset); err != nil {
		panic(err)
	}
}

// Asset returns asset data
func Asset() string {
	return "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"
}
