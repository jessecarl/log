package log

// Logger defines the bare minimum interface for logging structured data
// while specifying the level or priority.
type Logger interface {
	Log(Level, Data)
}

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

// Data provides an easily marshaled payload for structured logging. While an
// empty interface alone would satisfy the most basic requirements for
// structured logging, string keys on the first level allow better performance
// for basic filters without the need for reflection or type assertion.
type Data map[string]interface{}
