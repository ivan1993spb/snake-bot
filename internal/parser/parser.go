package parser

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ivan1993spb/snake-bot/internal/types"
)

type Countdown interface {
	Countdown(sec int)
}

type Me interface {
	Me(id uint32)
}

type Size interface {
	Size(width, height uint8)
}

type Game interface {
	Create(object *types.Object)
	Update(object *types.Object)
	Delete(object *types.Object)
}

type Printer interface {
	Print(level, what, message string)
}

type Parser struct {
	Countdown
	Me
	Size
	Game
	Printer
}

const errParseAnnotation = "parsing error"

var ErrParseUnknownMessage = errors.Errorf("%s: %s", errParseAnnotation,
	"unknown message")

func (p *Parser) Parse(data []byte) error {
	var message *types.Message
	if err := json.Unmarshal(data, &message); err != nil {
		return errors.Wrap(err, errParseAnnotation)
	}
	switch message.Type {
	case types.MessageTypeGameEvent:
		if err := p.ParseGameEvent(message.Payload); err != nil {
			return errors.Wrap(err, errParseAnnotation)
		}
	case types.MessageTypePlayer:
		if err := p.ParsePlayerEvent(message.Payload); err != nil {
			return errors.Wrap(err, errParseAnnotation)
		}
	case types.MessageTypeBroadcast:
		if err := p.ParseBroadcast(message.Payload); err != nil {
			return errors.Wrap(err, errParseAnnotation)
		}
	default:
		return errors.WithMessagef(ErrParseUnknownMessage, "unkown type %q",
			message.Type)
	}
	return nil
}

const errParseGameEventAnnotation = "cannot parse game event"

func (p *Parser) ParseGameEvent(data []byte) error {
	var event *types.GameEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return errors.Wrap(err, errParseGameEventAnnotation)
	}
	if event.Type == types.GameEventTypeError {
		var errorMessage string
		if err := json.Unmarshal(event.Payload, &errorMessage); err != nil {
			return errors.Wrap(err, errParseGameEventAnnotation)
		}
		p.Printer.Print("game", "error", errorMessage)
		return nil
	}
	var object *types.Object
	if err := json.Unmarshal(event.Payload, &object); err != nil {
		return errors.Wrap(err, errParseGameEventAnnotation)
	}
	switch event.Type {
	case types.GameEventTypeCreate:
		p.Game.Create(object)
	case types.GameEventTypeDelete:
		p.Game.Delete(object)
	case types.GameEventTypeUpdate:
		p.Game.Update(object)
	case types.GameEventTypeChecked:
		// ignore deprecated feature.
	}
	return nil
}

const errParsePlayerEventAnnotation = "cannot parse player event"

var ErrParseUnknownPlayerEvent = errors.Errorf("%s: %s", errParsePlayerEventAnnotation,
	"unknown player event")

func (p *Parser) ParsePlayerEvent(data []byte) error {
	var event *types.PlayerEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return errors.Wrap(err, errParsePlayerEventAnnotation)
	}
	switch event.Type {
	case types.PlayerEventTypeSize:
		var size *types.Size
		if err := json.Unmarshal(event.Payload, &size); err != nil {
			return errors.Wrap(err, errParsePlayerEventAnnotation)
		}
		p.Size.Size(size.Width, size.Height)
	case types.PlayerEventTypeSnake:
		var me uint32
		if err := json.Unmarshal(event.Payload, &me); err != nil {
			return errors.Wrap(err, errParsePlayerEventAnnotation)
		}
		p.Me.Me(me)
	case types.PlayerEventTypeNotice:
		var notice string
		if err := json.Unmarshal(event.Payload, &notice); err != nil {
			return errors.Wrap(err, errParsePlayerEventAnnotation)
		}
		p.Printer.Print("player", "notice", notice)
	case types.PlayerEventTypeError:
		var errorMessage string
		if err := json.Unmarshal(event.Payload, &errorMessage); err != nil {
			return errors.Wrap(err, errParsePlayerEventAnnotation)
		}
		p.Printer.Print("player", "error", errorMessage)
	case types.PlayerEventTypeCountdown:
		var countdown int
		if err := json.Unmarshal(event.Payload, &countdown); err != nil {
			return errors.Wrap(err, errParsePlayerEventAnnotation)
		}
		p.Countdown.Countdown(countdown)
	case types.PlayerEventTypeObjects:
		var objects []*types.Object
		if err := json.Unmarshal(event.Payload, &objects); err != nil {
			return errors.Wrap(err, errParsePlayerEventAnnotation)
		}
		for _, object := range objects {
			p.Game.Create(object)
		}
	default:
		return errors.WithMessagef(ErrParseUnknownPlayerEvent, "unkown type %q",
			event.Type)
	}
	return nil
}

const errParseBroadcastAnnotation = "cannot parse broadcast message"

func (p *Parser) ParseBroadcast(data []byte) error {
	var message string
	if err := json.Unmarshal(data, &message); err != nil {
		return errors.Wrap(err, errParseBroadcastAnnotation)
	}
	p.Printer.Print("broadcast", "message", message)
	return nil
}
