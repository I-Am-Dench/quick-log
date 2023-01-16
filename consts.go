package quicklog

type LogLevel int
type LogType string

const (
	LEVEL_DEBUG LogLevel = iota
	LEVEL_TRACE
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
	NUM_LOG_LEVELS
)

var LOG_PREFIX = [NUM_LOG_LEVELS]string{"D", "T", "I", "W", "E", "F"}
