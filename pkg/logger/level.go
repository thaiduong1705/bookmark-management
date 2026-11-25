package logger

import "github.com/rs/zerolog"

// SetLogLevel sets the global zerolog level; invalid or empty values fall back to info.
func SetLogLevel(level string) {
	lv, err := zerolog.ParseLevel(level)
	if err != nil || lv == zerolog.NoLevel {
		lv = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lv)
}
