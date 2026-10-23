/*
Copyright 2017 The Kubernetes Authors.

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

package slice

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestSortInts64(t *testing.T) {
	src := []int64{10, 1, 2, 3, 4, 5, 6}
	expected := []int64{1, 2, 3, 4, 5, 6, 10}
	SortInts64(src)
	if !reflect.DeepEqual(src, expected) {
		t.Errorf("func Ints64 didnt sort correctly, %v !- %v", src, expected)
	}
}

func TestContainsString(t *testing.T) {
	testCases := []struct {
		name     string
		slice    []string
		s        string
		modifier func(string) string
		expected bool
	}{
		{
			name:     "nil modifier with match",
			slice:    []string{"foo", "bar", "baz"},
			s:        "bar",
			modifier: nil,
			expected: true,
		},
		{
			name:     "nil modifier with no match",
			slice:    []string{"foo", "bar", "baz"},
			s:        "qux",
			modifier: nil,
			expected: false,
		},
		{
			name:     "nil modifier with case mismatch",
			slice:    []string{"Foo", "bar"},
			s:        "foo",
			modifier: nil,
			expected: false,
		},
		{
			name:     "strings.ToLower modifier with case-insensitive match",
			slice:    []string{"Foo", "bar"},
			s:        "foo",
			modifier: strings.ToLower,
			expected: true,
		},
		{
			name:     "strings.ToLower modifier with no match",
			slice:    []string{"Foo", "bar"},
			s:        "baz",
			modifier: strings.ToLower,
			expected: false,
		},
		{
			name:     "empty slice",
			slice:    []string{},
			s:        "foo",
			modifier: nil,
			expected: false,
		},
		{
			name:     "nil slice",
			slice:    nil,
			s:        "foo",
			modifier: nil,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ContainsString(tc.slice, tc.s, tc.modifier)
			if actual != tc.expected {
				t.Errorf("ContainsString(%v, %q, modifier) = %v; want %v", tc.slice, tc.s, actual, tc.expected)
			}
		})
	}
}

func TestToSet(t *testing.T) {
	testCases := map[string]struct {
		input    [][]string
		expected []string
	}{
		"nil should be returned if no slices are passed to the function": {
			input:    [][]string{},
			expected: nil,
		},
		"empty slice should be returned if an empty slice is passed to the function": {
			input:    [][]string{{}},
			expected: []string{},
		},
		"a single slice with no duplicates should have the same values": {
			input:    [][]string{{"a", "b", "c"}},
			expected: []string{"a", "b", "c"},
		},
		"duplicates should be removed from a single slice": {
			input:    [][]string{{"a", "b", "a", "c", "b"}},
			expected: []string{"a", "b", "c"},
		},
		"multiple slices with no duplicates should be combined": {
			input:    [][]string{{"a", "b", "c"}, {"d", "e", "f"}},
			expected: []string{"a", "b", "c", "d", "e", "f"},
		},
		"duplicates should be removed from multiple slices": {
			input:    [][]string{{"a", "b", "c"}, {"d", "b", "e"}, {"e", "f", "a"}},
			expected: []string{"a", "b", "c", "d", "e", "f"},
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := ToSet(tc.input...)
			sort.Strings(actual) // Sort is needed to compare the output because ToSet is non-deterministic
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("wrong output. Actual=%v, Expected=%v", actual, tc.expected)
			}
		})
	}
}
