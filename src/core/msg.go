package core

import (
	"encoding/json"
	"fmt"
)

// MaxJSONLen is the max length of the JSON bytes of a message
const MaxJSONLen = 100000 - 1

// MsgHeadLen is the length of the message header on wire
const MsgHeadLen = 5

const (

	// InvitationName is invitation message's type name
	InvitationName = "invitation"

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

// Message is
type Message struct {
	Name    string  `json:"msg_name"`
	Payload Payload `json:"msg_data"`
	Raw     bool    `json:"-"`
}

// Payload is
type Payload interface{}

// ParseMessageOnWire returns the parsed message object from bytes on the wire.
func ParseMessageOnWire(bytes []byte) (*Message, error) {
	msg := &Message{Payload: new(json.RawMessage), Raw: true}
	err := json.Unmarshal(bytes[MsgHeadLen:], &msg)
	if err != nil {
		return nil, fmt.Errorf("message error - fail to unmarshal bytes to message, error: %v\n%v", err, string(bytes))
	}
	return msg, msg.substantiatePayload()
}

func (m *Message) substantiatePayload() (err error) {
	switch m.Name {
	case InvitationName:
		_, err = m.Invitation()
	case RegistrationName:
		_, err = m.Registration()
	case LegStartName:
		_, err = m.LegStart()
	case LegEndName:
		_, err = m.LegEnd()
	case RoundName:
		_, err = m.Round()
	case ActionName:
		_, err = m.Action()
	case GameOverName:
		_, err = m.GameOver()
	default:
		err = fmt.Errorf("message error - unknown message type %v", m.Name)
	}
	return err
}

// BytesOnWire returns the bytes of the messsage on the wire.
func (m *Message) BytesOnWire() ([]byte, error) {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("message error - fail to marshal message to bytes, error: %v", err)
	}

	jsonLen := len(jsonBytes)
	if jsonLen > MaxJSONLen {
		return nil, fmt.Errorf("message error - message length %v is greather than %v", jsonLen, MaxJSONLen)
	}

	lenBytes := []byte(fmt.Sprintf("%05d", jsonLen))
	return append(lenBytes, jsonBytes...), nil
}

// String returns the JSON string of the messsage.
func (m *Message) String() string {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return fmt.Sprintf("message error - fail to marshal message to bytes, error: %v", err)
	}
	return string(jsonBytes)
}

// Equal equals
// func (m *Message) Equal(n *Message) bool {
// 	if m == nil || n == nil {
// 		return false
// 	}
// 	if m.Name != n.Name {
// 		return false
// 	}
// 	v1 := reflect.ValueOf(m.Payload)
// 	v2 := reflect.ValueOf(n.Payload)
// 	if v1.Type() != v2.Type() {
// 		return false
// 	}

// 	switch m.Name {
// 	case RegistrationName:
// 		a, err := m.Registration()
// 		if err != nil {
// 			return false
// 		}
// 		b, err := n.Registration()
// 		if err != nil {
// 			return false
// 		}
// 		return reflect.DeepEqual(a, b)
// 	case LegStartName:
// 		a, err := m.LegStart()
// 		if err != nil {
// 			return false
// 		}
// 		b, err := n.LegStart()
// 		if err != nil {
// 			return false
// 		}
// 		return reflect.DeepEqual(a, b)
// 	case LegEndName:
// 		a, err := m.LegEnd()
// 		if err != nil {
// 			return false
// 		}
// 		b, err := n.LegEnd()
// 		if err != nil {
// 			return false
// 		}
// 		return reflect.DeepEqual(a, b)
// 	case RoundName:
// 		a, err := m.Round()
// 		if err != nil {
// 			return false
// 		}
// 		b, err := n.Round()
// 		if err != nil {
// 			return false
// 		}
// 		return reflect.DeepEqual(a, b)
// 	case ActionName:
// 		a, err := m.Action()
// 		if err != nil {
// 			return false
// 		}
// 		b, err := n.Action()
// 		if err != nil {
// 			return false
// 		}
// 		return reflect.DeepEqual(a, b)
// 	case GameOverName:
// 		a, err := m.GameOver()
// 		if err != nil {
// 			return false
// 		}
// 		b, err := n.GameOver()
// 		if err != nil {
// 			return false
// 		}
// 		return reflect.DeepEqual(a, b)
// 	default:
// 		return false
// 	}
// }

func (m *Message) unmarshalPayload(obj interface{}) error {
	rawMsg := m.Payload.(*json.RawMessage)
	err := json.Unmarshal(*rawMsg, obj)
	if err != nil {
		return err
	}
	m.Payload = obj
	m.Raw = false
	return nil
}

func (m *Message) checkType(typeWanted string) error {
	if m.Name != typeWanted {
		return fmt.Errorf("message error - unmatched type conversion, want %q, actual %q", typeWanted, m.Name)
	}
	return nil
}

// Invitation returns the pointer of the un-marshaled Invitation message.
func (m *Message) Invitation() (*Invitation, error) {
	if err := m.checkType(InvitationName); err != nil {
		return nil, err
	}
	if m.Raw == false {
		return m.Payload.(*Invitation), nil
	}
	obj := new(Invitation)
	return obj, m.unmarshalPayload(obj)
}

// Registration returns the pointer of the un-marshaled Registration message.
func (m *Message) Registration() (*Registration, error) {
	if err := m.checkType(RegistrationName); err != nil {
		return nil, err
	}
	if m.Raw == false {
		return m.Payload.(*Registration), nil
	}
	obj := new(Registration)
	return obj, m.unmarshalPayload(obj)
}

// LegStart returns returns the pointer of the un-marshaled LegStart message.
func (m *Message) LegStart() (*LegStart, error) {
	if err := m.checkType(LegStartName); err != nil {
		return nil, err
	}
	if m.Raw == false {
		return m.Payload.(*LegStart), nil
	}
	obj := new(LegStart)
	return obj, m.unmarshalPayload(obj)
}

// LegEnd returns returns the pointer of the un-marshaled LegEnd message.
func (m *Message) LegEnd() (*LegEnd, error) {
	if err := m.checkType(LegEndName); err != nil {
		return nil, err
	}
	if m.Raw == false {
		return m.Payload.(*LegEnd), nil
	}
	obj := new(LegEnd)
	return obj, m.unmarshalPayload(obj)
}

// Round returns returns the pointer of the un-marshaled Round message.
func (m *Message) Round() (*Round, error) {
	if err := m.checkType(RoundName); err != nil {
		return nil, err
	}
	if m.Raw == false {
		return m.Payload.(*Round), nil
	}
	obj := new(Round)
	return obj, m.unmarshalPayload(obj)
}

// Action returns returns the pointer of the un-marshaled Action message.
func (m *Message) Action() (*Action, error) {
	if err := m.checkType(ActionName); err != nil {
		return nil, err
	}
	if m.Raw == false {
		return m.Payload.(*Action), nil
	}
	obj := new(Action)
	return obj, m.unmarshalPayload(obj)
}

// GameOver returns returns the pointer of the un-marshaled GameOver message.
func (m *Message) GameOver() (*GameOver, error) {
	if err := m.checkType(GameOverName); err != nil {
		return nil, err
	}
	if m.Raw == false {
		return m.Payload.(*GameOver), nil
	}
	obj := new(GameOver)
	return obj, m.unmarshalPayload(obj)
}
