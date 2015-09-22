package main

import (
	log "github.com/Sirupsen/logrus"
)

func MemoryAlarm(actual int, spec []int) (level string) {
	crit := spec[0]
	warn := spec[1]
	ok := spec[2]
	switch {
	case actual <= ok:
		log.Debug("ALARM: OK")
		return "OK"
	case actual <= warn:
		log.Debug("ALARM: WARN")
		return "WARN"
	case actual <= crit:
		log.Debug("ALARM: CRIT")
		return "CRITICAL"
	}
	return "NO ALARM MATCH FOUND"
}
