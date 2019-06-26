package tests

import (
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

// CompareResults compare tests  results
func CompareResults(t *testing.T, expected interface{}, got interface{}) {
	if diff := pretty.Compare(expected, got); diff != "" {
		t.Fatalf("diff: (-got +want)\n%s", diff)
	}
}
