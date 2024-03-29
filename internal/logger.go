package internal

// Logger is a type that is used to log errors and info to throughout Portal. Any logger implementation that
// implements Logger may be used by passing it to portal.New.
type Logger interface {
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})
}
