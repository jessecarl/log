package log_test

import (
	"reflect"
	"testing"

	"github.com/PermissionData/log"
)

func TestLevelMarshalText(t *testing.T) {
	testCases := []struct {
		lvl       log.Level
		wantBytes []byte
	}{
		{
			log.FatalLevel,
			[]byte("Fatal"),
		},
		{
			log.ErrorLevel,
			[]byte("Error"),
		},
		{
			log.InfoLevel,
			[]byte("Info"),
		},
		{
			log.TraceLevel,
			[]byte("Trace"),
		},
		{
			log.TraceLevel + 1,
			[]byte("4"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(string(tc.wantBytes), func(t *testing.T) {
			b, err := tc.lvl.MarshalText()
			if err != nil {
				t.Fatalf("Unexpected error, there should not be an error: %+v", err)
			}
			if !reflect.DeepEqual(b, tc.wantBytes) {
				t.Errorf("%v.MarshalText() = %q, expected %q", tc.lvl, b, tc.wantBytes)
			}
		})
	}
}

func TestLevelUnmarshalText(t *testing.T) {
	testCases := []struct {
		raw       []byte
		wantLvl   log.Level
		wantError bool
	}{
		{
			[]byte("Fatal"),
			log.FatalLevel,
			false,
		},
		{
			[]byte("fatal"),
			log.FatalLevel,
			false,
		},
		{
			[]byte("FATAL"),
			log.FatalLevel,
			false,
		},
		{
			[]byte("0"),
			log.FatalLevel,
			false,
		},
		{
			[]byte("Error"),
			log.ErrorLevel,
			false,
		},
		{
			[]byte("Info"),
			log.InfoLevel,
			false,
		},
		{
			[]byte("Trace"),
			log.TraceLevel,
			false,
		},
		{
			[]byte("3"),
			log.TraceLevel,
			false,
		},
		{
			[]byte("4"),
			log.TraceLevel + 1,
			false,
		},
		{
			nil,
			log.FatalLevel,
			true,
		},
		{
			[]byte("unknown"),
			log.FatalLevel,
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(string(tc.raw), func(t *testing.T) {
			var lvl = new(log.Level)
			err := lvl.UnmarshalText(tc.raw)
			if !tc.wantError && err != nil {
				t.Fatalf("Unexpected error, got: %v", err)
			}
			if tc.wantError && err == nil {
				t.Fatalf("lvl.UnmarshalText(%q), did not get expected error", tc.raw)
			}
			if *lvl != tc.wantLvl {
				t.Fatalf("lvl.UnmarshalText(%q), got %v, expected %v", tc.raw, lvl, tc.wantLvl)
			}
		})
	}
}
