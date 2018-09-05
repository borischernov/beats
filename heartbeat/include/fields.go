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

package include

import (
	"github.com/borischernov/beats/libbeat/asset"
)

func init() {
	if err := asset.SetFields("heartbeat", "fields.yml", Asset); err != nil {
		panic(err)
	}
}

// Asset returns asset data
func Asset() string {
	return "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"
}
