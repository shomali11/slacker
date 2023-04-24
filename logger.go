package slacker

import "log"

var logDebugMode = false

func infof(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func debugf(format string, v ...interface{}) {
	if logDebugMode {
		log.Printf(format, v...)
	}
}

func setLogDebugMode(debug bool) {
	logDebugMode = debug
}
