// Copyright (c) 2019-2026 Red Hat, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package workspace

import (
	"testing"
)

func TestMatchesToleratedKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		patterns []string
		expected bool
	}{
		{
			name:     "exact match",
			key:      "paas.redhat.com/appcode",
			patterns: []string{"paas.redhat.com/appcode"},
			expected: true,
		},
		{
			name:     "no match",
			key:      "paas.redhat.com/appcode",
			patterns: []string{"other.com/label"},
			expected: false,
		},
		{
			name:     "wildcard prefix match",
			key:      "paas.redhat.com/appcode",
			patterns: []string{"paas.redhat.com/*"},
			expected: true,
		},
		{
			name:     "wildcard prefix no match",
			key:      "other.com/label",
			patterns: []string{"paas.redhat.com/*"},
			expected: false,
		},
		{
			name:     "wildcard does not match prefix itself",
			key:      "paas.redhat.com",
			patterns: []string{"paas.redhat.com/*"},
			expected: false,
		},
		{
			name:     "multiple patterns with match",
			key:      "example.com/annotation",
			patterns: []string{"paas.redhat.com/*", "example.com/annotation"},
			expected: true,
		},
		{
			name:     "empty patterns",
			key:      "paas.redhat.com/appcode",
			patterns: []string{},
			expected: false,
		},
		{
			name:     "nil patterns",
			key:      "paas.redhat.com/appcode",
			patterns: nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesToleratedKey(tt.key, tt.patterns)
			if result != tt.expected {
				t.Errorf("matchesToleratedKey(%q, %v) = %v, want %v", tt.key, tt.patterns, result, tt.expected)
			}
		})
	}
}
