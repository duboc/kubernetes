/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fieldpath

import "testing"

func TestSplitMaybeSubscriptedPath(t *testing.T) {
	testCases := []struct {
		name          string
		fieldPath     string
		wantPath      string
		wantSubscript string
		wantOK        bool
	}{
		{
			name:          "subscripted annotation path",
			fieldPath:     "metadata.annotations['myKey']",
			wantPath:      "metadata.annotations",
			wantSubscript: "myKey",
			wantOK:        true,
		},
		{
			name:          "subscript containing brackets",
			fieldPath:     "metadata.annotations['a[b]c']",
			wantPath:      "metadata.annotations",
			wantSubscript: "a[b]c",
			wantOK:        true,
		},
		{
			name:          "empty subscript",
			fieldPath:     "metadata.labels['']",
			wantPath:      "metadata.labels",
			wantSubscript: "",
			wantOK:        true,
		},
		{
			name:          "unsubscripted path",
			fieldPath:     "metadata.labels",
			wantPath:      "metadata.labels",
			wantSubscript: "",
			wantOK:        false,
		},
		{
			name:          "no field before subscript",
			fieldPath:     "['myKey']",
			wantPath:      "['myKey']",
			wantSubscript: "",
			wantOK:        false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotPath, gotSubscript, gotOK := SplitMaybeSubscriptedPath(tc.fieldPath)
			if gotPath != tc.wantPath || gotSubscript != tc.wantSubscript || gotOK != tc.wantOK {
				t.Errorf("SplitMaybeSubscriptedPath(%q) = (%q, %q, %v); want (%q, %q, %v)",
					tc.fieldPath, gotPath, gotSubscript, gotOK, tc.wantPath, tc.wantSubscript, tc.wantOK)
			}
		})
	}
}

func TestFormatMap(t *testing.T) {
	testCases := []struct {
		name     string
		input    map[string]string
		expected string
	}{
		{
			name:     "empty map",
			input:    map[string]string{},
			expected: "",
		},
		{
			name:     "nil map",
			input:    nil,
			expected: "",
		},
		{
			name: "unsorted keys with quoted and escaped values",
			input: map[string]string{
				"z-key": "plain",
				"a-key": `quoted "value"`,
				"m-key": "line1\nline2",
			},
			expected: "a-key=\"quoted \\\"value\\\"\"\nm-key=\"line1\\nline2\"\nz-key=\"plain\"",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := FormatMap(tc.input)
			if got != tc.expected {
				t.Errorf("FormatMap(%v) = %q; want %q", tc.input, got, tc.expected)
			}
		})
	}
}
