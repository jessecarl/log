package log_test

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/PermissionData/log"
	"github.com/PermissionData/log/mock"
)

func TestWithLevels(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockLogger := mock_log.NewMockLogger(mockCtrl)

	lvlLogger := log.WithLevels(mockLogger)

	mockLogger.EXPECT().Log(gomock.Eq(log.FatalLevel), gomock.Eq(log.Data{})).Times(1)
	lvlLogger.Fatal(log.Data{})

	mockLogger.EXPECT().Log(gomock.Eq(log.ErrorLevel), gomock.Eq(log.Data{})).Times(1)
	lvlLogger.Error(log.Data{})

	mockLogger.EXPECT().Log(gomock.Eq(log.InfoLevel), gomock.Eq(log.Data{})).Times(1)
	lvlLogger.Info(log.Data{})

	mockLogger.EXPECT().Log(gomock.Eq(log.TraceLevel), gomock.Eq(log.Data{})).Times(1)
	lvlLogger.Trace(log.Data{})
}
