package log

import (
	"bytes"
	"fmt"
	"strconv"
)

// Level is used to indicate priority or threshold
type Level int

const (
	// FatalLevel should be used to communicate when the application has failed
	// and is left in an unpredictable state.
	FatalLevel Level = iota
	// ErrorLevel should be used to communicate when something went wrong in the
	// application, but the application can continue.
	ErrorLevel
	// InfoLevel should be used to communicate when something happened that is
	// worth noting.
	InfoLevel
	// TraceLevel should be used to communicate when something happened.
	TraceLevel
)

// LogLevelToString maps log levels to string representations
var LogLevelToString = map[Level]string{
	FatalLevel: "Fatal",
	ErrorLevel: "Error",
	InfoLevel:  "Info",
	TraceLevel: "Trace",
}

// StringToLogLevel maps string representations to log levels.
var StringToLogLevel = map[string]Level{
	"FATAL": FatalLevel,
	"ERROR": ErrorLevel,
	"INFO":  InfoLevel,
	"TRACE": TraceLevel,
}

// String represents a Level as a human-readable string instead of the integer value.
func (lvl Level) String() string {
	txt, ok := LogLevelToString[lvl]
	if !ok {
		txt = strconv.Itoa(int(lvl))
	}
	return txt
}

// MarshalText returns a human-readable text representation of a Level.
func (lvl Level) MarshalText() ([]byte, error) {
	return []byte(lvl.String()), nil
}

// UnmarshalText assigns a Level according to a text representation of either
// the human-readable string or integer value of a Level.
func (lvl *Level) UnmarshalText(raw []byte) error {
	level, ok := StringToLogLevel[string(bytes.ToUpper(raw))]
	if !ok {
		lvlInt, err := strconv.Atoi(string(raw))
		if err != nil {
			return fmt.Errorf("level not found: %+v", err)
		}
		level = Level(lvlInt)
	}
	*lvl = level
	return nil
}
