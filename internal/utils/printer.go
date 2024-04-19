package utils

import "github.com/sirupsen/logrus"

type PrinterLogger struct {
	log *logrus.Entry
}

func NewPrinterLogger(log *logrus.Entry) *PrinterLogger {
	return &PrinterLogger{
		log: log,
	}
}

func (p *PrinterLogger) Print(level, what, message string) {
	p.log.WithFields(logrus.Fields{
		"level": level,
		"what":  what,
	}).Debug(message)
}

type PrinterDevNull struct{}

func (p PrinterDevNull) Print(string, string, string) {
	// Do nothing
}
