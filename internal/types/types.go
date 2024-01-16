package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type Direction string

const (
	DirectionZero  Direction = "zero"
	DirectionNorth Direction = "north"
	DirectionWest  Direction = "west"
	DirectionSouth Direction = "south"
	DirectionEast  Direction = "east"
)

type MessageType string

const (
	MessageTypeGameEvent MessageType = "game"
	MessageTypePlayer    MessageType = "player"
	MessageTypeBroadcast MessageType = "broadcast"
)

type Message struct {
	Type    MessageType     `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type GameEventType string

const (
	GameEventTypeError   GameEventType = "error"
	GameEventTypeCreate  GameEventType = "create"
	GameEventTypeDelete  GameEventType = "delete"
	GameEventTypeUpdate  GameEventType = "update"
	GameEventTypeChecked GameEventType = "checked"
)

type GameEvent struct {
	Type    GameEventType   `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type PlayerEventType string

const (
	PlayerEventTypeSize      PlayerEventType = "size"
	PlayerEventTypeSnake     PlayerEventType = "snake"
	PlayerEventTypeNotice    PlayerEventType = "notice"
	PlayerEventTypeError     PlayerEventType = "error"
	PlayerEventTypeCountdown PlayerEventType = "countdown"
	PlayerEventTypeObjects   PlayerEventType = "objects"
)

type PlayerEvent struct {
	Type    PlayerEventType `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type Dot struct {
	X uint8
	Y uint8
}

func (d Dot) String() string {
	return fmt.Sprintf("[%d, %d]", d.X, d.Y)
}

const (
	dotFieldX = iota
	dotFieldY
)

const errDotUnmarshalJSONAnnotation = "cannot decode the dot"

var (
	ErrDotUnmarshalJSONInsufficientInput = errors.Errorf("%s: %s",
		errDotUnmarshalJSONAnnotation, "insufficient input")
	ErrDotUnmarshalJSONInvalidStructure = errors.Errorf("%s: %s",
		errDotUnmarshalJSONAnnotation, "invalid structure")
)

func (d *Dot) UnmarshalJSON(b []byte) error {
	if len(b) < 2 {
		return ErrDotUnmarshalJSONInsufficientInput
	}
	var op, cl byte
	op, b, cl = b[0], b[1:len(b)-1], b[len(b)-1]
	if op != '[' || cl != ']' {
		return ErrDotUnmarshalJSONInvalidStructure
	}
	axes := bytes.SplitN(b, []byte{','}, 2)
	if len(axes) != 2 {
		return ErrDotUnmarshalJSONInvalidStructure
	}
	x, err := strconv.Atoi(string(axes[dotFieldX]))
	if err != nil {
		return errors.Wrap(err, errDotUnmarshalJSONAnnotation)
	}
	y, err := strconv.Atoi(string(axes[dotFieldY]))
	if err != nil {
		return errors.Wrap(err, errDotUnmarshalJSONAnnotation)
	}
	*d = Dot{
		X: uint8(x),
		Y: uint8(y),
	}
	return nil
}

type ObjectType uint8

const (
	ObjectTypeUnknown ObjectType = iota
	ObjectTypeSnake
	ObjectTypeApple
	ObjectTypeCorpse
	ObjectTypeMouse
	ObjectTypeWatermelon
	ObjectTypeWall
)

var (
	ObjectTypeTokenSnake      = []byte(`"snake"`)
	ObjectTypeTokenApple      = []byte(`"apple"`)
	ObjectTypeTokenCorpse     = []byte(`"corpse"`)
	ObjectTypeTokenMouse      = []byte(`"mouse"`)
	ObjectTypeTokenWatermelon = []byte(`"watermelon"`)
	ObjectTypeTokenWall       = []byte(`"wall"`)
)

func (t *ObjectType) UnmarshalJSON(b []byte) error {
	switch {
	case bytes.Equal(b, ObjectTypeTokenSnake):
		*t = ObjectTypeSnake
	case bytes.Equal(b, ObjectTypeTokenWatermelon):
		*t = ObjectTypeWatermelon
	case bytes.Equal(b, ObjectTypeTokenCorpse):
		*t = ObjectTypeCorpse
	case bytes.Equal(b, ObjectTypeTokenWall):
		*t = ObjectTypeWall
	case bytes.Equal(b, ObjectTypeTokenMouse):
		*t = ObjectTypeMouse
	case bytes.Equal(b, ObjectTypeTokenApple):
		*t = ObjectTypeApple
	default:
		*t = ObjectTypeUnknown
	}
	return nil
}

type Object struct {
	Type      ObjectType `json:"type"`
	Id        uint32     `json:"id"`
	Dot       Dot        `json:"dot"`
	Dots      []Dot      `json:"dots"`
	Direction Direction  `json:"direction"`
}

func (o *Object) GetType() ObjectType {
	return o.Type
}

func (o *Object) GetDots() []Dot {
	switch o.Type {
	case ObjectTypeMouse, ObjectTypeApple:
		return []Dot{o.Dot}
	case ObjectTypeSnake, ObjectTypeCorpse, ObjectTypeWatermelon,
		ObjectTypeWall:
		return o.Dots
	}
	return []Dot{}
}

type Size struct {
	Width  uint8 `json:"width"`
	Height uint8 `json:"height"`
}

type MessageSnakeCommand int8

const (
	MessageSnakeCommandNorth = iota
	MessageSnakeCommandSouth
	MessageSnakeCommandWest
	MessageSnakeCommandEast
)

const (
	messageSnakeCommandLabelNorth = "north"
	messageSnakeCommandLabelSouth = "south"
	messageSnakeCommandLabelWest  = "west"
	messageSnakeCommandLabelEast  = "east"
)

var messageSnakeCommandLabelMapping = map[MessageSnakeCommand]string{
	MessageSnakeCommandNorth: messageSnakeCommandLabelNorth,
	MessageSnakeCommandSouth: messageSnakeCommandLabelSouth,
	MessageSnakeCommandWest:  messageSnakeCommandLabelWest,
	MessageSnakeCommandEast:  messageSnakeCommandLabelEast,
}

func (m MessageSnakeCommand) String() string {
	if label, ok := messageSnakeCommandLabelMapping[m]; ok {
		return label
	}
	return "unknown"
}

const messageSnakeCommandBuffer = 64

func (m MessageSnakeCommand) MarshalJSON() ([]byte, error) {
	b := bytes.NewBuffer(make([]byte, 0, messageSnakeCommandBuffer))
	b.Write([]byte(`{"type":"snake","payload":"`))
	b.WriteString(m.String())
	b.Write([]byte(`"}`))
	return b.Bytes(), nil
}

var directionMessageSnakeCommandMapping = map[Direction]MessageSnakeCommand{
	DirectionNorth: MessageSnakeCommandNorth,
	DirectionSouth: MessageSnakeCommandSouth,
	DirectionEast:  MessageSnakeCommandEast,
	DirectionWest:  MessageSnakeCommandWest,
}

func (d Direction) ToMessageSnakeCommand() MessageSnakeCommand {
	if command, ok := directionMessageSnakeCommandMapping[d]; ok {
		return command
	}
	return -1
}
