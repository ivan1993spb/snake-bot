package utils

import "github.com/sirupsen/logrus"

type PrinterLog struct {
	log logrus.FieldLogger
}

const printerGameKey = "print_game"

func NewPrinterLog(log logrus.FieldLogger, gameId int) *PrinterLog {
	return &PrinterLog{
		log: log.WithField(printerGameKey, gameId),
	}
}

func (p *PrinterLog) Print(level, what, message string) {
	p.log.WithFields(logrus.Fields{
		"level": level,
		"what":  what,
	}).Debug(message)
}

type PrinterDevNull struct{}

func (p PrinterDevNull) Print(string, string, string) {
	// Do nothing
}
