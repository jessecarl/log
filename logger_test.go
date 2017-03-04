package log_test

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/PermissionData/log"
	mock_log "github.com/PermissionData/log/mock"
)

func TestNewUsesAssignedFilters(t *testing.T) {
	testCases := []struct {
		name     string
		filters  []log.Filter
		wantData log.Data
	}{
		{"single filter",
			[]log.Filter{Pi},
			log.Data{
				"pi": 3.14,
			},
		},
		{"multiple filters",
			[]log.Filter{Pi, Phi},
			log.Data{
				"pi":  3.14,
				"phi": 1.618,
			},
		},
		{"multiple filters with conflict",
			[]log.Filter{
				Pi,
				Phi,
				func(lvl, t log.Level, d log.Data) log.Data {
					d["pi"] = "yum"
					return d
				},
			},
			log.Data{
				"pi":  "yum",
				"phi": 1.618,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockEncoder := mock_log.NewMockEncoder(mockCtrl)
			mockEncoder.EXPECT().Encode(gomock.Eq(tc.wantData)).Times(1)

			lg := log.New(log.Config{
				Encoder: mockEncoder,
				Filters: tc.filters,
			})
			lg.Log(log.InfoLevel, log.Data{})
		})
	}
}

func TestNewUsesDefaultFilter(t *testing.T) {
	oldFilter := log.DefaultFilter
	defer func() { log.DefaultFilter = oldFilter }()
	log.DefaultFilter = Pi

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockEncoder := mock_log.NewMockEncoder(mockCtrl)
	mockEncoder.EXPECT().Encode(gomock.Eq(log.Data{"pi": 3.14})).Times(1)

	lg := log.New(log.Config{Encoder: mockEncoder})
	lg.Log(log.FatalLevel, log.Data{})
}

func TestNewUsesAssignedEncoder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockEncoder := mock_log.NewMockEncoder(mockCtrl)

	lg := log.New(log.Config{
		Filters: []log.Filter{Pi},
		Encoder: mockEncoder,
	})
	mockEncoder.EXPECT().Encode(gomock.Eq(log.Data{"pi": 3.14})).Times(1)
	lg.Log(log.FatalLevel, log.Data{})
	lg.Log(log.FatalLevel, nil) // should not encode
}

func TestNewUsesDefaultEncoder(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	oldEncoder := log.DefaultEncoder
	defer func() { log.DefaultEncoder = oldEncoder }()
	log.DefaultEncoder = mock_log.NewMockEncoder(mockCtrl)
	log.DefaultEncoder.(*mock_log.MockEncoder).EXPECT().Encode(gomock.Any()).Times(1)

	lg := log.New(log.Config{Filters: []log.Filter{nopFilter}})
	lg.Log(log.InfoLevel, log.Data{})
}

func TestNewUsesThreshold(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockEncoder := mock_log.NewMockEncoder(mockCtrl)

	for i := log.FatalLevel; i <= log.TraceLevel+1; i++ {
		mockEncoder.EXPECT().Encode(gomock.Eq(log.Data{"pi": 3.14})).Times(1)

		lg := log.New(log.Config{
			Encoder: mockEncoder,
			Filters: []log.Filter{
				Pi,
				func(lvl, threshold log.Level, data log.Data) log.Data {
					if threshold != i {
						t.Errorf("filter called with threshold of %v, expected %v", threshold, i)
					}
					return data
				},
			},
			Threshold: i,
		})
		lg.Log(log.FatalLevel, log.Data{})
	}
}

func Pi(lvl, threshold log.Level, data log.Data) log.Data {
	if data == nil {
		return nil
	}
	data["pi"] = 3.14
	return data
}

func Phi(lvl, threshold log.Level, data log.Data) log.Data {
	if data == nil {
		return nil
	}
	data["phi"] = 1.618
	return data
}

func nopFilter(lvl, threshold log.Level, data log.Data) log.Data {
	return data
}
