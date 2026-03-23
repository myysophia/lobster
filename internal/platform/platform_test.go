package platform

import "testing"

func TestNormalizeOS(t *testing.T) {
	cases := map[string]OS{
		"windows": Windows,
		"darwin":  Darwin,
		"linux":   Linux,
		"plan9":   Unknown,
	}

	for input, want := range cases {
		if got := normalizeOS(input); got != want {
			t.Fatalf("normalizeOS(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestInfoString(t *testing.T) {
	info := Info{OS: Darwin, Arch: "arm64"}
	if got := info.String(); got != "darwin/arm64" {
		t.Fatalf("Info.String() = %q, want %q", got, "darwin/arm64")
	}
}
