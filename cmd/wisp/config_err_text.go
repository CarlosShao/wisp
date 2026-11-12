package main

import "github.com/CarlosShao/wisp/internal/config"

// configErrText keeps the warn record honest about a read that produced no
// value at all (err non-nil) versus a read that returned nothing usable.
//
// CROSS-PLATFORM ON PURPOSE: it is a pure function with no Windows dependency,
// but non-windows files (logsink.go, memory_retention_180k.go) call it, so it
// moved here out of resident_audio_windows.go (424-fix-b) to keep the Linux
// type-check whole. The body is unchanged.
func configErrText(err error, c *config.Config) string {
	if err != nil {
		return err.Error()
	}
	if c == nil {
		return "config.LoadFile returned a nil config with a nil error"
	}
	return ""
}
