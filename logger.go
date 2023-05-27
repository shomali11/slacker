package slacker

import "log"

var logDebugMode = false

func infof(format string, v ...any) {
	log.Printf(format, v...)
}

func debugf(format string, v ...any) {
	if logDebugMode {
		log.Printf(format, v...)
	}
}

func setLogDebugMode(debug bool) {
	logDebugMode = debug
}
