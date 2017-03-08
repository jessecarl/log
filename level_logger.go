package log

// LevelLogger extends the Logger with convenience methods for common Levels
type LevelLogger interface {
	Logger
	Fatal(Data)
	Error(Data)
	Info(Data)
	Trace(Data)
}

var _ LevelLogger = &logWithLevels{}

type logWithLevels struct {
	Logger
}

// WithLevels wraps a Logger to make it into a LevelLogger
func WithLevels(log Logger) LevelLogger {
	return &logWithLevels{log}
}

func (wl *logWithLevels) Fatal(data Data) { wl.Log(FatalLevel, data) }
func (wl *logWithLevels) Error(data Data) { wl.Log(ErrorLevel, data) }
func (wl *logWithLevels) Info(data Data)  { wl.Log(InfoLevel, data) }
func (wl *logWithLevels) Trace(data Data) { wl.Log(TraceLevel, data) }
