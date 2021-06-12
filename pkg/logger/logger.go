package logger

import "github.com/ipfs/go-log/v2"

func New(system string) *log.ZapEventLogger {
	logger := log.Logger(system)
	if err := log.SetLogLevel(system, "debug"); err != nil {
		// this should never happen
		panic(err)
	}
	return logger
}
