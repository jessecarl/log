package log

import (
	"runtime/debug"
	"time"
)

// Filter is a function that is used to manipulate the Data passed to a Logger.
// This could be adding fields, removing fields, or even setting the Data to
// nil. Filters need to be able to accept nil Data, as that is the most
// effective way of stopping Data from being logged.
type Filter func(lvl, threshold Level, data Data) Data

var (
	// DefaultTimestampFormat is used by the BaseFilter to add timestamps to logs.
	DefaultTimestampFormat = "2006-01-02T15:04:05.000Z"
	// DefaultFilter sets the BaseFilter to be enabled by default.
	DefaultFilter = BaseFilter()
)

// BaseFilter provides a Filter that prevents logs above the set threshold from
// being logged, and adds a timestamp, version, and log level to the log Data.
func BaseFilter() Filter {
	return func(lvl, threshold Level, data Data) Data {
		if data == nil {
			return nil
		}

		if lvl <= TraceLevel && lvl > threshold {
			return nil
		}

		data["@timestamp"] = time.Now().UTC().Format(DefaultTimestampFormat)
		data["@version"] = "1"
		data["log_level"] = lvl
		return data
	}
}

// StackFilter adds a stack trace to log data when the log level exceeds the threshold.
func StackFilter(stackLevel Level) Filter {
	return func(lvl, threshold Level, data Data) Data {
		if data == nil {
			return nil
		}
		if lvl <= stackLevel {
			data["_stack"] = string(debug.Stack())
		}
		return data
	}
}
