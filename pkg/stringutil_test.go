package pkg

import (
	"reflect"
	"testing"
)

type testcase struct {
	original   string
	truncated  string
	truncateAt int
}

func TestTruncateString(t *testing.T) {
	cases := []testcase{
		{original: "test", truncated: "tes", truncateAt: 3},
		{original: "test", truncated: "test", truncateAt: 5},
		{original: "", truncated: "", truncateAt: 3},
	}

	for _, c := range cases {
		result := TruncateString(c.original, c.truncateAt)
		if !reflect.DeepEqual(result, c.truncated) {
			t.Errorf("Expected: %s, got: %s", c.truncated, result)
		}
	}
}
