package log_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/PermissionData/log"
)

func TestBaseFilter(t *testing.T) {
	testCases := []struct {
		name      string
		lvl       log.Level
		threshold log.Level
		inData    log.Data
		wantData  log.Data
	}{
		{"nil data",
			log.FatalLevel,
			log.FatalLevel,
			nil,
			nil,
		},
		{"no data, ErrorLevel, ErrorLevel",
			log.ErrorLevel,
			log.ErrorLevel,
			log.Data{},
			log.Data{
				"@timestamp": log.DefaultTimestampFormat, // assertion will work with the format
				"@version":   "1",
				"log_level":  log.ErrorLevel,
			},
		},
		{"no data, InfoLevel, ErrorLevel",
			log.InfoLevel,
			log.ErrorLevel,
			log.Data{},
			nil,
		},
		{"no data, FatalLevel, ErrorLevel",
			log.FatalLevel,
			log.ErrorLevel,
			log.Data{},
			log.Data{
				"@timestamp": log.DefaultTimestampFormat,
				"@version":   "1",
				"log_level":  log.FatalLevel,
			},
		},
		{"no data, Custom Level, ErrorLevel",
			log.TraceLevel + 1,
			log.ErrorLevel,
			log.Data{},
			log.Data{
				"@timestamp": log.DefaultTimestampFormat,
				"@version":   "1",
				"log_level":  log.TraceLevel + 1,
			},
		},
		{"no data, ErrorLevel, Custom Level",
			log.ErrorLevel,
			log.TraceLevel + 1,
			log.Data{},
			log.Data{
				"@timestamp": log.DefaultTimestampFormat,
				"@version":   "1",
				"log_level":  log.ErrorLevel,
			},
		},
		{"misc data, InfoLevel, InfoLevel",
			log.InfoLevel,
			log.InfoLevel,
			log.Data{
				"pi": 3.14,
			},
			log.Data{
				"@timestamp": log.DefaultTimestampFormat, // assertion will work with the format
				"@version":   "1",
				"log_level":  log.InfoLevel,
				"pi":         3.14,
			},
		},
		{"conflicting data",
			log.InfoLevel,
			log.InfoLevel,
			log.Data{
				"@version":   3.14,
				"log_level":  "foo",
				"@timestamp": "4:20pm",
			},
			log.Data{
				"@timestamp": log.DefaultTimestampFormat,
				"@version":   "1",
				"log_level":  log.InfoLevel,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			inData := log.Data{}
			if tc.inData == nil {
				inData = nil
			}
			for k, v := range tc.inData {
				inData[k] = v
			}
			gotData := log.BaseFilter()(tc.lvl, tc.threshold, inData)

			if gotData != nil && tc.wantData == nil {
				t.Fatalf("log.BaseFilter()(%v, %v, %+v) = %+v, expected %+v", tc.lvl, tc.threshold, tc.inData, gotData, tc.wantData)
			}
			if gotData == nil && tc.wantData != nil {
				t.Fatalf("log.BaseFilter()(%v, %v, %+v) = %+v, expected %+v", tc.lvl, tc.threshold, tc.inData, gotData, tc.wantData)
			}
			if len(gotData) != len(tc.wantData) {
				t.Fatalf("log.BaseFilter()(%v, %v, %+v) = %+v, expected %+v", tc.lvl, tc.threshold, tc.inData, gotData, tc.wantData)
			}

			for k, v := range tc.wantData {
				got, ok := gotData[k]
				if !ok {
					t.Fatalf("log.BaseFilter()(%v, %v, %+v)[%q] == <nil>, expected, %+v", tc.lvl, tc.threshold, tc.inData, k, v)
				}
				if k == "@timestamp" {
					if _, err := time.Parse(v.(string), got.(string)); err != nil {
						t.Fatalf("log.BaseFilter()(%v, %v, %+v)[\"@timestamp\"] = %+v, expected like %q without parsing error %+v", tc.lvl, tc.threshold, tc.inData, got, v, err)
					}
					continue
				}
				if !reflect.DeepEqual(got, v) {
					t.Fatalf("log.BaseFilter()(%v, %v, %+v)[%q] == %+v, expected, %+v", tc.lvl, tc.threshold, tc.inData, k, got, v)
				}
			}
		})
	}
}

func TestStackFilter(t *testing.T) {
	testCases := []struct {
		name        string
		inData      log.Data
		inLevel     log.Level
		inThreshold log.Level
		wantData    log.Data
	}{
		{"nil data",
			nil,
			log.FatalLevel,
			log.FatalLevel,
			nil,
		},
		{"no data",
			log.Data{},
			log.FatalLevel,
			log.FatalLevel,
			log.Data{
				"_stack": "foo",
			},
		},
		{"misc data",
			log.Data{
				"pi": 3.14,
			},
			log.ErrorLevel,
			log.ErrorLevel,
			log.Data{
				"_stack": "foo",
				"pi":     3.14,
			},
		},
		{"conflicting data",
			log.Data{
				"_stack": 1.618,
			},
			log.InfoLevel,
			log.InfoLevel,
			log.Data{
				"_stack": "foo",
			},
		},
		{"below threshold",
			log.Data{
				"pi": 3.14,
			},
			log.InfoLevel,
			log.ErrorLevel,
			log.Data{
				"pi": 3.14,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotData := log.StackFilter(tc.inThreshold)(tc.inLevel, log.InfoLevel, tc.inData)

			if gotData != nil && tc.wantData == nil {
				t.Fatalf("StackFilter(%+v)(%+v, InfoLevel, %+v) = %+v, expected %+v", tc.inThreshold, tc.inLevel, tc.inData, gotData, tc.wantData)
			}
			if gotData == nil && tc.wantData != nil {
				t.Fatalf("StackFilter(%+v)(%+v, InfoLevel, %+v) = %+v, expected %+v", tc.inThreshold, tc.inLevel, tc.inData, gotData, tc.wantData)
			}

			if len(gotData) != len(tc.wantData) {
				t.Fatalf("StackFilter(%+v)(%+v, InfoLevel, %+v) = %+v, expected %+v", tc.inThreshold, tc.inLevel, tc.inData, gotData, tc.wantData)
			}

			for k, v := range tc.wantData {
				got, ok := gotData[k]
				if !ok {
					t.Fatalf("StackFilter(%+v)(%+v, InfoLevel, %+v)[%q] == <nil>, expected, %+v", tc.inThreshold, tc.inLevel, tc.inData, k, v)
				}
				if k == "_stack" {
					if _, ok := got.(string); !ok {
						t.Fatalf("StackFilter(%+v)(%+v, InfoLevel, %+v)[\"_stack\"] = %+v, expected a string", tc.inThreshold, tc.inLevel, tc.inData, got)
					}
					continue
				}
				if !reflect.DeepEqual(got, v) {
					t.Fatalf(
						"StackFilter(%+v)(%+v, InfoLevel, %+v)[%q] == %+v, expected, %+v",
						tc.inThreshold,
						tc.inLevel,
						tc.inData,
						k,
						got,
						v,
					)
				}
			}
		})
	}
}

func TestErrorFilter(t *testing.T) {
	testCases := []struct {
		name      string
		lvl       log.Level
		threshold log.Level
		errorKeys []string
		inData    log.Data
		wantData  log.Data
	}{
		{"nil data, no keys",
			log.FatalLevel,
			log.FatalLevel,
			nil,
			nil,
			nil,
		},
		{"nil data, one key",
			log.InfoLevel,
			log.InfoLevel,
			[]string{"Error"},
			nil,
			nil,
		},
		{"no data, ErrorLevel, ErrorLevel, no keys",
			log.ErrorLevel,
			log.ErrorLevel,
			[]string{},
			log.Data{},
			log.Data{},
		},
		{"misc data, no keys",
			log.TraceLevel,
			log.TraceLevel,
			[]string{},
			log.Data{
				"pi": 3.14,
			},
			log.Data{
				"pi": 3.14,
			},
		},
		{"misc data, with keys, no matches",
			log.TraceLevel,
			log.TraceLevel,
			[]string{"Error", "foo"},
			log.Data{
				"pi": 3.14,
			},
			log.Data{
				"pi": 3.14,
			},
		},
		{"misc data, with keys, with matches, no errors",
			log.TraceLevel,
			log.TraceLevel,
			[]string{"Error", "foo"},
			log.Data{
				"pi":  3.14,
				"foo": "bar",
			},
			log.Data{
				"pi":  3.14,
				"foo": "bar",
			},
		},
		{"misc data, with keys, with matches, including errors",
			log.TraceLevel,
			log.TraceLevel,
			[]string{"Error", "foo"},
			log.Data{
				"pi":    3.14,
				"foo":   "bar",
				"Error": fmt.Errorf("baz"),
			},
			log.Data{
				"pi":    3.14,
				"foo":   "bar",
				"Error": "baz",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			gotData := log.ErrorFilter(tc.errorKeys...)(tc.lvl, tc.threshold, tc.inData)

			if gotData != nil && tc.wantData == nil {
				t.Fatalf("ErrorFilter(%+v)(%v, %v, %+v) = %+v, expected %+v", tc.errorKeys, tc.lvl, tc.threshold, tc.inData, gotData, tc.wantData)
			}
			if gotData == nil && tc.wantData != nil {
				t.Fatalf("ErrorFilter(%+v)(%v, %v, %+v) = %+v, expected %+v", tc.errorKeys, tc.lvl, tc.threshold, tc.inData, gotData, tc.wantData)
			}

			if len(gotData) != len(tc.wantData) {
				t.Fatalf("ErrorFilter(%+v)(%v, %v, %+v) = %+v, expected %+v", tc.errorKeys, tc.lvl, tc.threshold, tc.inData, gotData, tc.wantData)
			}

			for k, v := range tc.wantData {
				got, ok := gotData[k]
				if !ok {
					t.Fatalf("ErrorFilter(%+v)(%v, %v, %+v)[%q] == <nil>, expected %+v", tc.errorKeys, tc.lvl, tc.threshold, tc.inData, k, v)
				}
				if k == "@timestamp" {
					if _, err := time.Parse(v.(string), got.(string)); err != nil {
						t.Fatalf("ErrorFilter(%+v)(%v, %v, %+v)[\"@timestamp\"] = %+v, expected like %q without parsing error %+v", tc.errorKeys, tc.lvl, tc.threshold, tc.inData, got, v, err)
					}
					continue
				}
				if !reflect.DeepEqual(got, v) {
					t.Fatalf("ErrorFilter(%+v)(%v, %v, %+v)[%q] == %+v, expected %+v", tc.errorKeys, tc.lvl, tc.threshold, tc.inData, k, got, v)
				}
			}
		})
	}
}
