package core

import (
	"encoding/json"
	"fmt"
	"log"
)

// MaxJSONLen is the max length of the JSON bytes of a message
const MaxJSONLen = 99999

// Message is
type Message struct {
	Name    string  `json:"msg_name"`
	Payload Payload `json:"msg_data"`
}

// Payload is
type Payload interface{}

// ParseMessage returns the parsed message object from bytes on the wire.
func ParseMessage(bytes []byte) (*Message, error) {
	msg := &Message{Payload: new(json.RawMessage)}
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		return nil, fmt.Errorf("message error - fail to unmarshal bytes to message, error: %v", err)
	}
	return msg, nil
}

// String returns the JSON string of the messsage.
func (m *Message) String() string {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("message error - fail to marshal message to bytes, error: %v", err)
	}
	return string(jsonBytes)
}

// Bytes returns the bytes of the messsage on the wire.
func (m *Message) Bytes() ([]byte, error) {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("message error - fail to marshal message to bytes, error: %v", err)
	}

	jsonLen := len(jsonBytes)
	if jsonLen > MaxJSONLen {
		return nil, fmt.Errorf("message error - message length %v is greather than %v", jsonLen, MaxJSONLen)
	}

	lenBytes := []byte(fmt.Sprintf("%05d", jsonLen))
	wireBytes := append(lenBytes, jsonBytes...)
	log.Printf("message sent ->\n%v", string(wireBytes))
	return wireBytes, nil
}

const (
	// RegistrationName is registration message's type name
	RegistrationName = "registration"

	// LegStartName is leg start message's type name
	LegStartName = "leg_start"

	// LegEndName is leg end message's type name
	LegEndName = "leg_end"

	// RoundName is round message's type name
	RoundName = "round"

	// ActionName is action message's type name
	ActionName = "action"

	// GameOverName is game over message's type name
	GameOverName = "game_over"
)

// Registration returns the pointer of the un-marshaled Registration message.
func (m *Message) Registration() (*Registration, error) {
	if m.Name != RegistrationName {
		return nil, fmt.Errorf("message error - unmatched type conversion, want %q, actual %q", m.Name, RegistrationName)
	}
	obj := new(Registration)
	err := json.Unmarshal(m.Payload.(json.RawMessage), obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// LegStart returns returns the pointer of the un-marshaled LegStart message.
func (m *Message) LegStart() (*LegStart, error) {
	if m.Name != LegStartName {
		return nil, fmt.Errorf("message error - unmatched type conversion, want %q, actual %q", m.Name, LegStartName)
	}

	obj := new(LegStart)
	err := json.Unmarshal(m.Payload.(json.RawMessage), obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// LegEnd returns returns the pointer of the un-marshaled LegEnd message.
func (m *Message) LegEnd() (*LegEnd, error) {
	if m.Name != LegEndName {
		return nil, fmt.Errorf("message error - unmatched type conversion, want %q, actual %q", m.Name, LegEndName)
	}

	obj := new(LegEnd)
	err := json.Unmarshal(m.Payload.(json.RawMessage), obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// Round returns returns the pointer of the un-marshaled Round message.
func (m *Message) Round() (*Round, error) {
	if m.Name != RoundName {
		return nil, fmt.Errorf("message error - unmatched type conversion, want %q, actual %q", m.Name, RoundName)
	}

	obj := new(Round)
	err := json.Unmarshal(m.Payload.(json.RawMessage), obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// Action returns returns the pointer of the un-marshaled Action message.
func (m *Message) Action() (*Action, error) {
	if m.Name != ActionName {
		return nil, fmt.Errorf("message error - unmatched type conversion, want %q, actual %q", m.Name, ActionName)
	}

	obj := new(Action)
	err := json.Unmarshal(m.Payload.(json.RawMessage), obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

// GameOver returns returns the pointer of the un-marshaled GameOver message.
func (m *Message) GameOver() (*GameOver, error) {
	if m.Name != GameOverName {
		return nil, fmt.Errorf("message error - unmatched type conversion, want %q, actual %q", m.Name, GameOverName)
	}
	return new(GameOver), nil
}
