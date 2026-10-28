/*
Copyright 2016 The Kubernetes Authors.

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

package version

import (
	"fmt"
	"reflect"
	"testing"

	apimachineryversion "k8s.io/apimachinery/pkg/version"
)

type testItem struct {
	version    string
	unparsed   string
	equalsPrev bool
}

func testOne(v *Version, item, prev testItem) error {
	str := v.String()
	if item.unparsed == "" {
		if str != item.version {
			return fmt.Errorf("bad round-trip: %q -> %q", item.version, str)
		}
	} else {
		if str != item.unparsed {
			return fmt.Errorf("bad unparse: %q -> %q, expected %q", item.version, str, item.unparsed)
		}
	}

	if prev.version != "" {
		cmp, err := v.Compare(prev.version)
		if err != nil {
			return fmt.Errorf("unexpected parse error: %v", err)
		}
		rv, err := parse(prev.version, v.semver)
		if err != nil {
			return fmt.Errorf("unexpected parse error: %v", err)
		}
		rcmp, err := rv.Compare(item.version)
		if err != nil {
			return fmt.Errorf("unexpected parse error: %v", err)
		}

		switch {
		case cmp == -1:
			return fmt.Errorf("unexpected ordering %q < %q", item.version, prev.version)
		case cmp == 0 && !item.equalsPrev:
			return fmt.Errorf("unexpected comparison %q == %q", item.version, prev.version)
		case cmp == 1 && item.equalsPrev:
			return fmt.Errorf("unexpected comparison %q != %q", item.version, prev.version)
		case cmp != -rcmp:
			return fmt.Errorf("unexpected reverse comparison %q <=> %q %v %v %v %v", item.version, prev.version, cmp, rcmp, v.Components(), rv.Components())
		}
	}

	return nil
}

func TestSemanticVersions(t *testing.T) {
	tests := []testItem{
		// This is every version string that appears in the 2.0 semver spec,
		// sorted in strictly increasing order except as noted.
		{version: "0.1.0"},
		{version: "1.0.0-0.3.7"},
		{version: "1.0.0-alpha"},
		{version: "1.0.0-alpha+001", equalsPrev: true},
		{version: "1.0.0-alpha.1"},
		{version: "1.0.0-alpha.beta"},
		{version: "1.0.0-beta"},
		{version: "1.0.0-beta+exp.sha.5114f85", equalsPrev: true},
		{version: "1.0.0-beta.2"},
		{version: "1.0.0-beta.11"},
		{version: "1.0.0-rc.1"},
		{version: "1.0.0-x.7.z.92"},
		{version: "1.0.0"},
		{version: "1.0.0+20130313144700", equalsPrev: true},
		{version: "1.8.0-alpha.3"},
		{version: "1.8.0-alpha.3.673+73326ef01d2d7c"},
		{version: "1.9.0"},
		{version: "1.10.0"},
		{version: "1.11.0"},
		{version: "2.0.0"},
		{version: "2.1.0"},
		{version: "2.1.1"},
		{version: "42.0.0"},

		// We also allow whitespace and "v" prefix
		{version: "   42.0.0", unparsed: "42.0.0", equalsPrev: true},
		{version: "\t42.0.0  ", unparsed: "42.0.0", equalsPrev: true},
		{version: "43.0.0-1", unparsed: "43.0.0-1"},
		{version: "43.0.0-1  ", unparsed: "43.0.0-1", equalsPrev: true},
		{version: "v43.0.0-1", unparsed: "43.0.0-1", equalsPrev: true},
		{version: "  v43.0.0", unparsed: "43.0.0"},
		{version: " 43.0.0 ", unparsed: "43.0.0", equalsPrev: true},
	}

	var prev testItem
	for _, item := range tests {
		v, err := ParseSemantic(item.version)
		if err != nil {
			t.Errorf("unexpected parse error: %v", err)
			continue
		}
		err = testOne(v, item, prev)
		if err != nil {
			t.Errorf("%v", err)
		}
		prev = item
	}
}

func TestBadSemanticVersions(t *testing.T) {
	tests := []string{
		// "MUST take the form X.Y.Z"
		"1",
		"1.2",
		"1.2.3.4",
		".2.3",
		"1..3",
		"1.2.",
		"",
		"..",
		// "where X, Y, and Z are non-negative integers"
		"-1.2.3",
		"1.-2.3",
		"1.2.-3",
		"1a.2.3",
		"1.2a.3",
		"1.2.3a",
		"a1.2.3",
		"a.b.c",
		"1 .2.3",
		"1. 2.3",
		// "and MUST NOT contain leading zeroes."
		"01.2.3",
		"1.02.3",
		"1.2.03",
		// "[pre-release] identifiers MUST comprise only ASCII alphanumerics and hyphen"
		"1.2.3-/",
		// "[pre-release] identifiers MUST NOT be empty"
		"1.2.3-",
		"1.2.3-.",
		"1.2.3-foo.",
		"1.2.3-.foo",
		// "Numeric [pre-release] identifiers MUST NOT include leading zeroes"
		"1.2.3-01",
		// "[build metadata] identifiers MUST comprise only ASCII alphanumerics and hyphen"
		"1.2.3+/",
		// "[build metadata] identifiers MUST NOT be empty"
		"1.2.3+",
		"1.2.3+.",
		"1.2.3+foo.",
		"1.2.3+.foo",

		// whitespace/"v"-prefix checks
		"v 1.2.3",
		"vv1.2.3",
	}

	for i := range tests {
		_, err := ParseSemantic(tests[i])
		if err == nil {
			t.Errorf("unexpected success parsing invalid semver %q", tests[i])
		}
	}
}

func TestGenericVersions(t *testing.T) {
	tests := []testItem{
		// This is all of the strings from TestSemanticVersions, plus some strings
		// from TestBadSemanticVersions that should parse as generic versions,
		// plus some additional strings.
		{version: "0.1.0", unparsed: "0.1.0"},
		{version: "1.0.0-0.3.7", unparsed: "1.0.0"},
		{version: "1.0.0-alpha", unparsed: "1.0.0", equalsPrev: true},
		{version: "1.0.0-alpha+001", unparsed: "1.0.0", equalsPrev: true},
		{version: "1.0.0-alpha.1", unparsed: "1.0.0", equalsPrev: true},
		{version: "1.0.0-alpha.beta", unparsed: "1.0.0", equalsPrev: true},
		{version: "1.0.0.beta", unparsed: "1.0.0", equalsPrev: true},
		{version: "1.0.0-beta+exp.sha.5114f85", unparsed: "1.0.0", equalsPrev: true},
		{version: "1.0.0.beta.2", unparsed: "1.0.0", equalsPrev: true},
		{version: "1.0.0.beta.11", unparsed: "1.0.0", equalsPrev: true},
		{version: "1.0.0.rc.1", unparsed: "1.0.0", equalsPrev: true},
		{version: "1.0.0-x.7.z.92", unparsed: "1.0.0", equalsPrev: true},
		{version: "1.0.0", unparsed: "1.0.0", equalsPrev: true},
		{version: "1.0.0+20130313144700", unparsed: "1.0.0", equalsPrev: true},
		{version: "1.2", unparsed: "1.2"},
		{version: "1.2a.3", unparsed: "1.2", equalsPrev: true},
		{version: "1.2.3", unparsed: "1.2.3"},
		{version: "1.2.3.0", unparsed: "1.2.3.0", equalsPrev: true},
		{version: "1.2.3a", unparsed: "1.2.3", equalsPrev: true},
		{version: "1.2.3-foo.", unparsed: "1.2.3", equalsPrev: true},
		{version: "1.2.3-.foo", unparsed: "1.2.3", equalsPrev: true},
		{version: "1.2.3-01", unparsed: "1.2.3", equalsPrev: true},
		{version: "1.2.3+", unparsed: "1.2.3", equalsPrev: true},
		{version: "1.2.3+foo.", unparsed: "1.2.3", equalsPrev: true},
		{version: "1.2.3+.foo", unparsed: "1.2.3", equalsPrev: true},
		{version: "1.02.3", unparsed: "1.2.3", equalsPrev: true},
		{version: "1.2.03", unparsed: "1.2.3", equalsPrev: true},
		{version: "1.2.003", unparsed: "1.2.3", equalsPrev: true},
		{version: "1.2.3.4", unparsed: "1.2.3.4"},
		{version: "1.2.3.4b3", unparsed: "1.2.3.4", equalsPrev: true},
		{version: "1.2.3.4.5", unparsed: "1.2.3.4.5"},
		{version: "1.9.0", unparsed: "1.9.0"},
		{version: "1.9.0.0.0.0.0.0", unparsed: "1.9.0.0.0.0.0.0", equalsPrev: true},
		{version: "1.10.0", unparsed: "1.10.0"},
		{version: "1.11.0", unparsed: "1.11.0"},
		{version: "1.11.0.0.5", unparsed: "1.11.0.0.5"},
		{version: "2.0.0", unparsed: "2.0.0"},
		{version: "2.1.0", unparsed: "2.1.0"},
		{version: "2.1.1", unparsed: "2.1.1"},
		{version: "42.0.0", unparsed: "42.0.0"},
		{version: "   42.0.0", unparsed: "42.0.0", equalsPrev: true},
		{version: "\t42.0.0  ", unparsed: "42.0.0", equalsPrev: true},
		{version: "42.0.0-1", unparsed: "42.0.0", equalsPrev: true},
		{version: "42.0.0-1  ", unparsed: "42.0.0", equalsPrev: true},
		{version: "v42.0.0-1", unparsed: "42.0.0", equalsPrev: true},
		{version: "  v43.0.0", unparsed: "43.0.0"},
		{version: " 43.0.0 ", unparsed: "43.0.0", equalsPrev: true},
	}

	var prev testItem
	for _, item := range tests {
		v, err := ParseGeneric(item.version)
		if err != nil {
			t.Errorf("unexpected parse error: %v", err)
			continue
		}
		err = testOne(v, item, prev)
		if err != nil {
			t.Errorf("%v", err)
		}
		prev = item
	}
}

func TestBadGenericVersions(t *testing.T) {
	tests := []string{
		"1",
		"01.2.3",
		"-1.2.3",
		"1.-2.3",
		".2.3",
		"1..3",
		"1a.2.3",
		"a1.2.3",
		"1 .2.3",
		"1. 2.3",
		"1.bob",
		"bob",
		"v 1.2.3",
		"vv1.2.3",
		"",
		".",
	}

	for i := range tests {
		_, err := ParseGeneric(tests[i])
		if err == nil {
			t.Errorf("unexpected success parsing invalid version %q", tests[i])
		}
	}
}

func TestComponents(t *testing.T) {

	var tests = []struct {
		version               string
		semver                bool
		expectedComponents    []uint
		expectedMajor         uint
		expectedMinor         uint
		expectedPatch         uint
		expectedPreRelease    string
		expectedBuildMetadata string
	}{
		{
			version:            "1.0.2",
			semver:             true,
			expectedComponents: []uint{1, 0, 2},
			expectedMajor:      1,
			expectedMinor:      0,
			expectedPatch:      2,
		},
		{
			version:               "1.0.2-alpha+001",
			semver:                true,
			expectedComponents:    []uint{1, 0, 2},
			expectedMajor:         1,
			expectedMinor:         0,
			expectedPatch:         2,
			expectedPreRelease:    "alpha",
			expectedBuildMetadata: "001",
		},
		{
			version:            "1.2",
			semver:             false,
			expectedComponents: []uint{1, 2},
			expectedMajor:      1,
			expectedMinor:      2,
		},
		{
			version:               "1.0.2-beta+exp.sha.5114f85",
			semver:                true,
			expectedComponents:    []uint{1, 0, 2},
			expectedMajor:         1,
			expectedMinor:         0,
			expectedPatch:         2,
			expectedPreRelease:    "beta",
			expectedBuildMetadata: "exp.sha.5114f85",
		},
	}

	for _, test := range tests {
		version, _ := parse(test.version, test.semver)
		if !reflect.DeepEqual(test.expectedComponents, version.Components()) {
			t.Error("parse returned un'expected components")
		}
		if test.expectedMajor != version.Major() {
			t.Errorf("parse returned version.Major %d, expected %d", test.expectedMajor, version.Major())
		}
		if test.expectedMinor != version.Minor() {
			t.Errorf("parse returned version.Minor %d, expected %d", test.expectedMinor, version.Minor())
		}
		if test.expectedPatch != version.Patch() {
			t.Errorf("parse returned version.Patch %d, expected %d", test.expectedPatch, version.Patch())
		}
		if test.expectedPreRelease != version.PreRelease() {
			t.Errorf("parse returned version.PreRelease %s, expected %s", test.expectedPreRelease, version.PreRelease())
		}
		if test.expectedBuildMetadata != version.BuildMetadata() {
			t.Errorf("parse returned version.BuildMetadata %s, expected %s", test.expectedBuildMetadata, version.BuildMetadata())
		}
	}
}

func TestHighestSupportedVersion(t *testing.T) {
	testCases := []struct {
		versions                        []string
		expectedHighestSupportedVersion string
		shouldFail                      bool
	}{
		{
			versions:                        []string{"v1.0.0"},
			expectedHighestSupportedVersion: "1.0.0",
			shouldFail:                      false,
		},
		{
			versions:   []string{"0.3.0"},
			shouldFail: true,
		},
		{
			versions:   []string{"0.2.0"},
			shouldFail: true,
		},
		{
			versions:                        []string{"1.0.0"},
			expectedHighestSupportedVersion: "1.0.0",
			shouldFail:                      false,
		},
		{
			versions:   []string{"v0.3.0"},
			shouldFail: true,
		},
		{
			versions:   []string{"v0.2.0"},
			shouldFail: true,
		},
		{
			versions:   []string{"0.2.0", "v0.3.0"},
			shouldFail: true,
		},
		{
			versions:                        []string{"0.2.0", "v1.0.0"},
			expectedHighestSupportedVersion: "1.0.0",
			shouldFail:                      false,
		},
		{
			versions:                        []string{"0.2.0", "v1.2.3"},
			expectedHighestSupportedVersion: "1.2.3",
			shouldFail:                      false,
		},
		{
			versions:                        []string{"v1.2.3", "v0.3.0"},
			expectedHighestSupportedVersion: "1.2.3",
			shouldFail:                      false,
		},
		{
			versions:                        []string{"v1.2.3", "v0.3.0", "2.0.1"},
			expectedHighestSupportedVersion: "1.2.3",
			shouldFail:                      false,
		},
		{
			versions:                        []string{"v1.2.3", "4.9.12", "v0.3.0", "2.0.1"},
			expectedHighestSupportedVersion: "1.2.3",
			shouldFail:                      false,
		},
		{
			versions:                        []string{"4.9.12", "2.0.1"},
			expectedHighestSupportedVersion: "",
			shouldFail:                      true,
		},
		{
			versions:                        []string{"v1.2.3", "boo", "v0.3.0", "2.0.1"},
			expectedHighestSupportedVersion: "1.2.3",
			shouldFail:                      false,
		},
		{
			versions:                        []string{},
			expectedHighestSupportedVersion: "",
			shouldFail:                      true,
		},
		{
			versions:                        []string{"var", "boo", "foo"},
			expectedHighestSupportedVersion: "",
			shouldFail:                      true,
		},
	}

	for _, tc := range testCases {
		// Arrange & Act
		actual, err := HighestSupportedVersion(tc.versions)

		// Assert
		if tc.shouldFail && err == nil {
			t.Fatalf("expecting highestSupportedVersion to fail, but got nil error for testcase: %#v", tc)
		}
		if !tc.shouldFail && err != nil {
			t.Fatalf("unexpected error during ValidatePlugin for testcase: %#v\r\n err:%v", tc, err)
		}
		if tc.expectedHighestSupportedVersion != "" {
			result, err := actual.Compare(tc.expectedHighestSupportedVersion)
			if err != nil {
				t.Fatalf("comparison failed with %v for testcase %#v", err, tc)
			}
			if result != 0 {
				t.Fatalf("expectedHighestSupportedVersion %v, but got %v for tc: %#v", tc.expectedHighestSupportedVersion, actual, tc)
			}
		}
	}
}

func TestOffsetMinor(t *testing.T) {
	var tests = []struct {
		version            string
		diff               int
		expectedComponents []uint
	}{
		{
			version:            "1.0.2",
			diff:               -3,
			expectedComponents: []uint{1, 0},
		},
		{
			version:            "1.3.2-alpha+001",
			diff:               -2,
			expectedComponents: []uint{1, 1},
		},
		{
			version:            "1.3.2-alpha+001",
			diff:               -3,
			expectedComponents: []uint{1, 0},
		},
		{
			version:            "1.20",
			diff:               -5,
			expectedComponents: []uint{1, 15},
		},
		{
			version:            "1.20",
			diff:               5,
			expectedComponents: []uint{1, 25},
		},
	}

	for _, test := range tests {
		version, _ := ParseGeneric(test.version)
		if !reflect.DeepEqual(test.expectedComponents, version.OffsetMinor(test.diff).Components()) {
			t.Error("parse returned un'expected components")
		}
	}
}

func TestParse(t *testing.T) {

	var tests = []struct {
		version               string
		expectErr             bool
		expectedComponents    []uint
		expectedPreRelease    string
		expectedBuildMetadata string
	}{
		{
			version:            "1.0.2",
			expectedComponents: []uint{1, 0, 2},
		},
		{
			version:               "1.0.2-alpha+001",
			expectedComponents:    []uint{1, 0, 2},
			expectedPreRelease:    "alpha",
			expectedBuildMetadata: "001",
		},
		{
			version:            "1.2",
			expectedComponents: []uint{1, 2},
		},
		{
			version:               "1.0.2-beta+exp.sha.5114f85",
			expectedComponents:    []uint{1, 0, 2},
			expectedPreRelease:    "beta",
			expectedBuildMetadata: "exp.sha.5114f85",
		},
		{
			version:   "a.b.c",
			expectErr: true,
		},
	}

	for _, test := range tests {
		version, err := Parse(test.version)
		if test.expectErr {
			if err == nil {
				t.Fatalf("got no err, expected err")
			}
			continue
		}
		if !reflect.DeepEqual(test.expectedComponents, version.Components()) {
			t.Error("parse returned un'expected components")
		}
		if test.expectedPreRelease != version.PreRelease() {
			t.Errorf("parse returned version.PreRelease %s, expected %s", test.expectedPreRelease, version.PreRelease())
		}
		if test.expectedBuildMetadata != version.BuildMetadata() {
			t.Errorf("parse returned version.BuildMetadata %s, expected %s", test.expectedBuildMetadata, version.BuildMetadata())
		}
	}
}

func TestNilSafeOrdering(t *testing.T) {
	v120Semver := MustParseSemantic("1.2.0")
	v120Alpha := MustParseSemantic("1.2.0-alpha.1")
	v120Build := MustParseSemantic("1.2.0+build.1")
	v130Semver := MustParseSemantic("1.3.0")
	v12Generic := MustParseGeneric("1.2")
	v123Generic := MustParseGeneric("1.2.3")

	tests := []struct {
		name        string
		lhs         *Version
		rhs         *Version
		equalTo     bool
		lessThan    bool
		greaterThan bool
		atLeast     bool
	}{
		{
			name:        "nil and nil",
			lhs:         nil,
			rhs:         nil,
			equalTo:     true,
			lessThan:    false,
			greaterThan: false,
			atLeast:     true,
		},
		{
			name:        "nil and semantic version",
			lhs:         nil,
			rhs:         v120Semver,
			equalTo:     false,
			lessThan:    true,
			greaterThan: false,
			atLeast:     false,
		},
		{
			name:        "semantic version and nil",
			lhs:         v120Semver,
			rhs:         nil,
			equalTo:     false,
			lessThan:    false,
			greaterThan: true,
			atLeast:     true,
		},
		{
			name:        "nil and generic version",
			lhs:         nil,
			rhs:         v12Generic,
			equalTo:     false,
			lessThan:    true,
			greaterThan: false,
			atLeast:     false,
		},
		{
			name:        "generic version and nil",
			lhs:         v12Generic,
			rhs:         nil,
			equalTo:     false,
			lessThan:    false,
			greaterThan: true,
			atLeast:     true,
		},
		{
			name:        "pre-release and release semantic versions",
			lhs:         v120Alpha,
			rhs:         v120Semver,
			equalTo:     false,
			lessThan:    true,
			greaterThan: false,
			atLeast:     false,
		},
		{
			name:        "release and pre-release semantic versions",
			lhs:         v120Semver,
			rhs:         v120Alpha,
			equalTo:     false,
			lessThan:    false,
			greaterThan: true,
			atLeast:     true,
		},
		{
			name:        "semantic versions differing only by build metadata",
			lhs:         v120Semver,
			rhs:         v120Build,
			equalTo:     true,
			lessThan:    false,
			greaterThan: false,
			atLeast:     true,
		},
		{
			name:        "distinct semantic versions",
			lhs:         v120Semver,
			rhs:         v130Semver,
			equalTo:     false,
			lessThan:    true,
			greaterThan: false,
			atLeast:     false,
		},
		{
			name:        "generic 1.2 and semantic 1.2.0",
			lhs:         v12Generic,
			rhs:         v120Semver,
			equalTo:     true,
			lessThan:    false,
			greaterThan: false,
			atLeast:     true,
		},
		{
			name:        "generic 1.2 and generic 1.2.3",
			lhs:         v12Generic,
			rhs:         v123Generic,
			equalTo:     false,
			lessThan:    true,
			greaterThan: false,
			atLeast:     false,
		},
		{
			name:        "generic 1.2.3 and generic 1.2",
			lhs:         v123Generic,
			rhs:         v12Generic,
			equalTo:     false,
			lessThan:    false,
			greaterThan: true,
			atLeast:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.lhs.EqualTo(tc.rhs); got != tc.equalTo {
				t.Errorf("(%v).EqualTo(%v) = %v, want %v", tc.lhs, tc.rhs, got, tc.equalTo)
			}
			if got := tc.lhs.LessThan(tc.rhs); got != tc.lessThan {
				t.Errorf("(%v).LessThan(%v) = %v, want %v", tc.lhs, tc.rhs, got, tc.lessThan)
			}
			if got := tc.lhs.GreaterThan(tc.rhs); got != tc.greaterThan {
				t.Errorf("(%v).GreaterThan(%v) = %v, want %v", tc.lhs, tc.rhs, got, tc.greaterThan)
			}
			if got := tc.lhs.AtLeast(tc.rhs); got != tc.atLeast {
				t.Errorf("(%v).AtLeast(%v) = %v, want %v", tc.lhs, tc.rhs, got, tc.atLeast)
			}

			// Invariant: LessThan(other) == !AtLeast(other)
			if tc.lhs.LessThan(tc.rhs) != !tc.lhs.AtLeast(tc.rhs) {
				t.Errorf("invariant violated for (%v, %v): LessThan != !AtLeast", tc.lhs, tc.rhs)
			}
			// Antisymmetry: lhs.LessThan(rhs) == rhs.GreaterThan(lhs)
			if tc.lhs.LessThan(tc.rhs) != tc.rhs.GreaterThan(tc.lhs) {
				t.Errorf("antisymmetry violated for (%v, %v): lhs.LessThan(rhs) != rhs.GreaterThan(lhs)", tc.lhs, tc.rhs)
			}
		})
	}
}

func TestCore(t *testing.T) {
	if got := (*Version)(nil).Core(); got != nil {
		t.Fatalf("(*Version)(nil).Core() = %v, want nil", got)
	}

	tests := []struct {
		name               string
		input              *Version
		expectedString     string
		expectedComponents []uint
	}{
		{
			name:               "semantic version with pre-release and build metadata",
			input:              MustParseSemantic("1.2.3-alpha.1+build.456"),
			expectedString:     "1.2.3",
			expectedComponents: []uint{1, 2, 3},
		},
		{
			name:               "semantic version with pre-release only",
			input:              MustParseSemantic("1.2.3-beta.2"),
			expectedString:     "1.2.3",
			expectedComponents: []uint{1, 2, 3},
		},
		{
			name:               "semantic version with build metadata only",
			input:              MustParseSemantic("1.2.3+20130313"),
			expectedString:     "1.2.3",
			expectedComponents: []uint{1, 2, 3},
		},
		{
			name:               "plain semantic version",
			input:              MustParseSemantic("1.2.3"),
			expectedString:     "1.2.3",
			expectedComponents: []uint{1, 2, 3},
		},
		{
			name:               "generic 2-component version",
			input:              MustParseGeneric("1.2"),
			expectedString:     "1.2.0",
			expectedComponents: []uint{1, 2, 0},
		},
		{
			name:               "generic 4-component version",
			input:              MustParseGeneric("1.2.3.4"),
			expectedString:     "1.2.3",
			expectedComponents: []uint{1, 2, 3},
		},
		{
			name:               "single-component struct",
			input:              &Version{components: []uint{4}},
			expectedString:     "4.0.0",
			expectedComponents: []uint{4, 0, 0},
		},
		{
			name:               "empty-component struct",
			input:              &Version{},
			expectedString:     "0.0.0",
			expectedComponents: []uint{0, 0, 0},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			origComponents := append([]uint(nil), tc.input.Components()...)
			origPreRelease := tc.input.PreRelease()
			origBuildMetadata := tc.input.BuildMetadata()

			core := tc.input.Core()
			if core == nil {
				t.Fatalf("Core() returned nil for non-nil input %v", tc.input)
			}
			if !core.semver {
				t.Errorf("Core().semver = false, want true")
			}
			if got := core.String(); got != tc.expectedString {
				t.Errorf("Core().String() = %q, want %q", got, tc.expectedString)
			}
			if !reflect.DeepEqual(core.Components(), tc.expectedComponents) {
				t.Errorf("Core().Components() = %v, want %v", core.Components(), tc.expectedComponents)
			}
			if core.PreRelease() != "" {
				t.Errorf("Core().PreRelease() = %q, want empty", core.PreRelease())
			}
			if core.BuildMetadata() != "" {
				t.Errorf("Core().BuildMetadata() = %q, want empty", core.BuildMetadata())
			}

			// Verify slice isolation and non-mutation of original receiver.
			core.Components()[0] = 999
			if !reflect.DeepEqual(tc.input.Components(), origComponents) {
				t.Errorf("mutating Core().Components() mutated input components: got %v, want %v", tc.input.Components(), origComponents)
			}
			if tc.input.PreRelease() != origPreRelease || tc.input.BuildMetadata() != origBuildMetadata {
				t.Errorf("Core() mutated input metadata")
			}
		})
	}
}

func TestVersionHelpers(t *testing.T) {
	assertPanics := func(t *testing.T, name string, fn func()) {
		t.Helper()
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("%s did not panic on invalid input", name)
			}
		}()
		fn()
	}

	if got := MustParseGeneric("1.2.3.4").String(); got != "1.2.3.4" {
		t.Errorf("MustParseGeneric(\"1.2.3.4\").String() = %q, want \"1.2.3.4\"", got)
	}
	assertPanics(t, "MustParseGeneric", func() {
		MustParseGeneric("invalid")
	})

	if got := MustParse("1.2.3-alpha.1").String(); got != "1.2.3-alpha.1" {
		t.Errorf("MustParse(\"1.2.3-alpha.1\").String() = %q, want \"1.2.3-alpha.1\"", got)
	}
	assertPanics(t, "MustParse", func() {
		MustParse("invalid")
	})

	mm, err := ParseMajorMinor("1.28.4")
	if err != nil || mm.String() != "1.28" {
		t.Errorf("ParseMajorMinor(\"1.28.4\") = (%v, %v), want (\"1.28\", nil)", mm, err)
	}
	if _, err := ParseMajorMinor("invalid"); err == nil {
		t.Errorf("ParseMajorMinor(\"invalid\") expected error, got nil")
	}
	if got := MustParseMajorMinor("1.28.4").String(); got != "1.28" {
		t.Errorf("MustParseMajorMinor(\"1.28.4\").String() = %q, want \"1.28\"", got)
	}
	assertPanics(t, "MustParseMajorMinor", func() {
		MustParseMajorMinor("invalid")
	})

	assertPanics(t, "MustParseSemantic", func() {
		MustParseSemantic("1.2")
	})

	// Overflow component error branch in parse
	if _, err := ParseSemantic("999999999999999999999999999999.0.0"); err == nil {
		t.Errorf("expected overflow error from ParseSemantic, got nil")
	}

	base := MustParseSemantic("1.5.2-alpha.1+build.10")
	if got := base.WithMajor(2).String(); got != "2.5.2-alpha.1+build.10" {
		t.Errorf("WithMajor(2) = %q, want \"2.5.2-alpha.1+build.10\"", got)
	}
	if got := base.WithMinor(8).String(); got != "1.8.2-alpha.1+build.10" {
		t.Errorf("WithMinor(8) = %q, want \"1.8.2-alpha.1+build.10\"", got)
	}
	if got := base.WithPatch(9).String(); got != "1.5.9-alpha.1+build.10" {
		t.Errorf("WithPatch(9) = %q, want \"1.5.9-alpha.1+build.10\"", got)
	}
	if got := base.WithPreRelease(""); got != base {
		t.Errorf("WithPreRelease(\"\") should return original pointer %p, got %p", base, got)
	}
	if got := base.WithPreRelease("rc.1").String(); got != "1.5.2-rc.1+build.10" {
		t.Errorf("WithPreRelease(\"rc.1\") = %q, want \"1.5.2-rc.1+build.10\"", got)
	}
	if got := base.WithBuildMetadata("sha.999").String(); got != "1.5.2-alpha.1+sha.999" {
		t.Errorf("WithBuildMetadata(\"sha.999\") = %q, want \"1.5.2-alpha.1+sha.999\"", got)
	}
	if got := base.SubtractMinor(2).String(); got != "1.3" {
		t.Errorf("SubtractMinor(2) = %q, want \"1.3\"", got)
	}
	if got := base.AddMinor(3).String(); got != "1.8" {
		t.Errorf("AddMinor(3) = %q, want \"1.8\"", got)
	}

	if got := (*Version)(nil).String(); got != "<nil>" {
		t.Errorf("(*Version)(nil).String() = %q, want \"<nil>\"", got)
	}
	if _, err := base.Compare("not-a-version"); err == nil {
		t.Errorf("Compare(\"not-a-version\") expected error, got nil")
	}

	withInfo := base.WithInfo(apimachineryversion.Info{Major: "9", Minor: "9"})
	if !withInfo.EqualTo(base) {
		t.Errorf("WithInfo should preserve version value")
	}
	if got := (*Version)(nil).Info(); got != nil {
		t.Errorf("(*Version)(nil).Info() = %v, want nil", got)
	}
	info := MustParseSemantic("0.2.3").Info()
	if info == nil || info.Major != "" || info.Minor != "2" || info.GitVersion != "0.2.3" {
		t.Errorf("MustParseSemantic(\"0.2.3\").Info() = %+v, unexpected result", info)
	}
}

