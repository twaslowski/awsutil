package pkg

import (
	"errors"
	"reflect"
	"testing"
)

type TestCase struct {
	name string
	err  error
	val  any
	want any
}

func TestRequire(t *testing.T) {
	tests := []TestCase{
		{
			name: "no error returns value",
			val:  42,
			err:  nil,
			want: 42,
		},
		{
			name: "error panics",
			val:  "ignored",
			err:  errors.New("something went wrong"),
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if r != nil && tt.err == nil {
					t.Errorf("Require() panicked unexpectedly: %v", r)
				}
				if r == nil && tt.err != nil {
					t.Errorf("Require() did not panic with error: %v", tt.err)
				}
			}()

			got := Require(tt.val, tt.err)
			if tt.err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Require() = %v, want %v", got, tt.want)
			}
		})
	}
}

type CheckErrTestCase struct {
	name string
	err  error
	want any
}

func TestCheckErr(t *testing.T) {
	tests := []CheckErrTestCase{
		{
			name: "Should panic if err != nil",
			err:  nil,
			want: nil,
		},
	}

	for _, test := range tests {
		CheckErr(test.err)

		r := recover()
		if test.err != nil && r == nil {
			t.Errorf("CheckErr() did not panic with error=%v", test.err)
		}
		if test.err == nil && r != nil {
			t.Errorf("CheckErr() panicked unexpectedly: %v", r)
		}
	}
}
