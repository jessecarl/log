package log_test

import (
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
