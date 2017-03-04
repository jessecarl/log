package log

import (
	"encoding/json"
	"fmt"
	"os"
)

// Logger defines the bare minimum interface for logging structured data
// while specifying the level or priority.
type Logger interface {
	Log(Level, Data)
}

// Data provides an easily marshaled payload for structured logging. While an
// empty interface alone would satisfy the most basic requirements for
// structured logging, string keys on the first level allow better performance
// for basic filters without the need for reflection or type assertion.
type Data map[string]interface{}

// Encoder is used to safely prepare and send structured data for consumption.
// The standard package `json.Encoder` and `gob.Encoder` types are good
// implementations of this interface.
type Encoder interface {
	Encode(interface{}) error
}

// DefaultEncoder ensures that a New logger does not requre an explicit Encoder
var DefaultEncoder Encoder = json.NewEncoder(os.Stdout)

type logger struct {
	encoder   Encoder
	filters   []Filter
	threshold Level
}

// Config contains the values that will be used by a new Logger
type Config struct {
	Threshold Level
	Encoder   Encoder
	Filters   []Filter
}

// New provides a basic Logger using the provided configuration.
func New(config Config) Logger {
	lg := &logger{
		encoder:   config.Encoder,
		filters:   config.Filters,
		threshold: config.Threshold,
	}
	if len(lg.filters) == 0 {
		lg.filters = []Filter{DefaultFilter}
	}
	if config.Encoder == nil {
		lg.encoder = DefaultEncoder
	}
	return lg
}

func (lg *logger) Log(lvl Level, data Data) {
	for _, fn := range lg.filters {
		if data = fn(lvl, lg.threshold, data); data == nil {
			return
		}
	}

	if err := lg.encoder.Encode(data); err != nil {
		// I'm ambivalent on printing anything to stdout/stderr, but this should probably happen. (jallen)
		// I agree. (rrichardson)
		fmt.Fprintf(os.Stderr, "Error writing to log: %+v\n", err)
	}
}
