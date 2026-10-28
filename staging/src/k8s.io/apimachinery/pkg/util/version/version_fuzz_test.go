/*
Copyright 2024 The Kubernetes Authors.

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
	"strings"
	"testing"
)

func verifyCompareLaws(t *testing.T, v1, v2 *Version) {
	t.Helper()

	c12, err := v1.Compare(v2.String())
	if err != nil {
		t.Fatalf("(%v).Compare(%q) unexpected error: %v", v1, v2.String(), err)
	}
	c21, err := v2.Compare(v1.String())
	if err != nil {
		t.Fatalf("(%v).Compare(%q) unexpected error: %v", v2, v1.String(), err)
	}
	if c12 != -c21 {
		t.Fatalf("Compare antisymmetry violated for %q and %q: c12=%d, c21=%d", v1.String(), v2.String(), c12, c21)
	}

	if v1.LessThan(v2) != !v1.AtLeast(v2) {
		t.Fatalf("LessThan/AtLeast complement violated for (%v, %v)", v1, v2)
	}
	if v1.LessThan(v2) != v2.GreaterThan(v1) {
		t.Fatalf("LessThan/GreaterThan antisymmetry violated for (%v, %v)", v1, v2)
	}
	if v1.EqualTo(v2) != v2.EqualTo(v1) {
		t.Fatalf("EqualTo symmetry violated for (%v, %v)", v1, v2)
	}
	if v1.EqualTo(v2) != (c12 == 0) {
		t.Fatalf("EqualTo inconsistent with Compare for (%v, %v): EqualTo=%v, c12=%d", v1, v2, v1.EqualTo(v2), c12)
	}
	if v1.LessThan(v2) != (c12 == -1) {
		t.Fatalf("LessThan inconsistent with Compare for (%v, %v): LessThan=%v, c12=%d", v1, v2, v1.LessThan(v2), c12)
	}
	if v1.GreaterThan(v2) != (c12 == 1) {
		t.Fatalf("GreaterThan inconsistent with Compare for (%v, %v): GreaterThan=%v, c12=%d", v1, v2, v1.GreaterThan(v2), c12)
	}
}

func FuzzParseAndCompare(f *testing.F) {
	seeds := []string{
		// Corpus from TestSemanticVersions
		"0.1.0",
		"1.0.0-0.3.7",
		"1.0.0-alpha",
		"1.0.0-alpha+001",
		"1.0.0-alpha.1",
		"1.0.0-alpha.beta",
		"1.0.0-beta",
		"1.0.0-beta+exp.sha.5114f85",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0-x.7.z.92",
		"1.0.0",
		"1.0.0+20130313144700",
		"1.8.0-alpha.3",
		"1.8.0-alpha.3.673+73326ef01d2d7c",
		"1.9.0",
		"1.10.0",
		"1.11.0",
		"2.0.0",
		"2.1.0",
		"2.1.1",
		"42.0.0",
		"   42.0.0",
		"\t42.0.0  ",
		"43.0.0-1",
		"v43.0.0-1",
		"  v43.0.0",
		// Corpus from TestGenericVersions
		"1.0.0.beta",
		"1.0.0.beta.2",
		"1.2",
		"1.2a.3",
		"1.2.3",
		"1.2.3.0",
		"1.2.3a",
		"1.02.3",
		"1.2.03",
		"1.2.3.4",
		"1.2.3.4b3",
		"1.2.3.4.5",
		"1.9.0.0.0.0.0.0",
		"1.11.0.0.5",
		// Representative invalid inputs from TestBadSemanticVersions / TestBadGenericVersions
		"",
		"1",
		".2.3",
		"1..3",
		"01.2.3",
		"1.2.3-01",
		"1.2.3-",
		"1.2.3+",
		"a.b.c",
		"v 1.2.3",
		// Pair seeds
		"1.0.0-alpha.1|1.0.0-alpha.2",
		"1.2.3-alpha+001|1.2.3+002",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	semanticAnchors := []*Version{
		MustParseSemantic("0.1.0"),
		MustParseSemantic("1.0.0-0.3.7"),
		MustParseSemantic("1.0.0-alpha"),
		MustParseSemantic("1.0.0-alpha.1"),
		MustParseSemantic("1.0.0-beta.2"),
		MustParseSemantic("1.0.0-rc.1"),
		MustParseSemantic("1.0.0"),
		MustParseSemantic("1.28.0"),
		MustParseSemantic("2.0.0"),
	}

	genericAnchors := []*Version{
		MustParseGeneric("0.1.0"),
		MustParseGeneric("1.0.0"),
		MustParseGeneric("1.2"),
		MustParseGeneric("1.2.0"),
		MustParseGeneric("1.2.3"),
		MustParseGeneric("1.2.3.4"),
		MustParseGeneric("2.0.0"),
	}

	f.Fuzz(func(t *testing.T, str string) {
		// Ensure top-level helpers never panic on arbitrary input.
		_, _ = Parse(str)
		_, _ = ParseMajorMinor(str)
		_, _ = HighestSupportedVersion([]string{str, "v1.0.0"})

		if v, err := ParseSemantic(str); err == nil {
			if v == nil {
				t.Fatalf("ParseSemantic(%q) returned nil Version without error", str)
			}

			// Semantic round-trip law: ParseSemantic(v.String()).EqualTo(v)
			rt, err := ParseSemantic(v.String())
			if err != nil {
				t.Fatalf("ParseSemantic(%q) round-trip ParseSemantic(%q) failed: %v", str, v.String(), err)
			}
			if !rt.EqualTo(v) || !v.EqualTo(rt) {
				t.Fatalf("ParseSemantic(%q) round-trip EqualTo failed: v=%v, rt=%v", str, v, rt)
			}
			if rt.String() != v.String() {
				t.Fatalf("ParseSemantic(%q) round-trip String mismatch: got %q, want %q", str, rt.String(), v.String())
			}

			// Core() helper law
			core := v.Core()
			if core == nil {
				t.Fatalf("(%v).Core() returned nil", v)
			}
			if core.PreRelease() != "" || core.BuildMetadata() != "" {
				t.Fatalf("(%v).Core() retained metadata: preRelease=%q, buildMetadata=%q", v, core.PreRelease(), core.BuildMetadata())
			}
			coreRT, err := ParseSemantic(core.String())
			if err != nil || !coreRT.EqualTo(core) {
				t.Fatalf("ParseSemantic(v.Core().String()) failed round-trip: core=%v, err=%v", core, err)
			}

			// Nil-ordering invariants
			if !v.GreaterThan(nil) || !v.AtLeast(nil) || v.LessThan(nil) || v.EqualTo(nil) {
				t.Fatalf("non-nil receiver %v violated nil ordering invariants", v)
			}
			if !(*Version)(nil).LessThan(v) || (*Version)(nil).AtLeast(v) || (*Version)(nil).GreaterThan(v) || (*Version)(nil).EqualTo(v) {
				t.Fatalf("nil receiver violated ordering invariants against %v", v)
			}

			// Compare antisymmetry & total-order predicate consistency
			verifyCompareLaws(t, v, v)
			verifyCompareLaws(t, v, rt)
			verifyCompareLaws(t, v, core)
			verifyCompareLaws(t, v, v.WithMajor(v.Major()+1))
			verifyCompareLaws(t, v, v.WithMinor(v.Minor()+1))
			verifyCompareLaws(t, v, v.WithPatch(v.Patch()+1))
			verifyCompareLaws(t, v, v.WithPreRelease("alpha.1"))
			for _, anchor := range semanticAnchors {
				verifyCompareLaws(t, v, anchor)
			}
		}

		if g, err := ParseGeneric(str); err == nil {
			if g == nil {
				t.Fatalf("ParseGeneric(%q) returned nil Version without error", str)
			}

			grt, err := ParseGeneric(g.String())
			if err != nil {
				t.Fatalf("ParseGeneric(%q) round-trip ParseGeneric(%q) failed: %v", str, g.String(), err)
			}
			if !grt.EqualTo(g) || !g.EqualTo(grt) {
				t.Fatalf("ParseGeneric(%q) round-trip EqualTo failed: g=%v, grt=%v", str, g, grt)
			}
			if grt.String() != g.String() {
				t.Fatalf("ParseGeneric(%q) round-trip String mismatch: got %q, want %q", str, grt.String(), g.String())
			}

			verifyCompareLaws(t, g, g)
			verifyCompareLaws(t, g, grt)
			verifyCompareLaws(t, g, g.OffsetMinor(1))
			verifyCompareLaws(t, g, g.OffsetMinor(-1))
			for _, anchor := range genericAnchors {
				verifyCompareLaws(t, g, anchor)
			}
		}

		if lhs, rhs, ok := strings.Cut(str, "|"); ok {
			if v1, err1 := ParseSemantic(lhs); err1 == nil {
				if v2, err2 := ParseSemantic(rhs); err2 == nil {
					verifyCompareLaws(t, v1, v2)
				}
			}
			if g1, err1 := ParseGeneric(lhs); err1 == nil {
				if g2, err2 := ParseGeneric(rhs); err2 == nil {
					verifyCompareLaws(t, g1, g2)
				}
			}
		}
	})
}
