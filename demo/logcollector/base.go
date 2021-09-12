package logcollector

import "context"

type LogLevel int

const (
	LogLevelDebug LogLevel = 0
	LogLevelInfo  LogLevel = 1
	LogLevelError LogLevel = 2
)

type LogCollector interface {
	CtxLog(ctx context.Context, logLevel LogLevel, fmt string, args ...interface{})
}