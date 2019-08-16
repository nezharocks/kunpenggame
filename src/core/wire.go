package core

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
)

// Wire wraps
type Wire struct {
	ReadWriter io.ReadWriter
	BufferSize int
	Reader     *bufio.Reader
	Writer     *bufio.Writer
	MsgCh      chan Message
	ErrCh      chan error
}

// NewWire creates wire between game server and team client.
func NewWire(rw io.ReadWriter, bufferSize int) *Wire {
	return &Wire{
		Reader: bufio.NewReaderSize(rw, bufferSize),
		Writer: bufio.NewWriterSize(rw, bufferSize),
		MsgCh:  make(chan Message, 2),
		ErrCh:  make(chan error, 10),
	}
}

// Send sends message to the wire
func (w *Wire) Send(msg *Message) error {
	wireBytes, err := msg.BytesOnWire()
	if err != nil {
		return err
	}

	if _, err := w.Writer.Write(wireBytes); err != nil {
		return fmt.Errorf("wire error - fail to write message to wire, error: %v", err)
	}

	if err := w.Writer.Flush(); err != nil {
		return fmt.Errorf("wire error - fail to write message to wire, error: %v", err)
	}

	if debugMessage {
		log.Printf("message sent ->\n%v", string(wireBytes))
	}
	return nil
}

// Receive receives message from the wire
func (w *Wire) Receive() {
	scanner := bufio.NewScanner(w.Reader)
	scanner.Split(msgSplit)
	for scanner.Scan() {
		msg, err := ParseMessageOnWire(scanner.Bytes())
		if err != nil {
			w.ErrCh <- err
			continue
		}
		if debugMessage {
			log.Printf("message received ->\n%v", string(scanner.Bytes()))
		}
		w.MsgCh <- *msg
	}
}

func msgSplit(data []byte, atEOF bool) (int, []byte, error) {
	l := len(data)
	if atEOF && l < 5 {
		return 0, nil, fmt.Errorf("message error - message is broken, msg:\n%v", string(data))
	}

	size, err := strconv.ParseInt(string(data[:5]), 10, 32)
	if err != nil {
		return 0, nil, fmt.Errorf("message error - message is invalid, error: %v", err)
	}

	totalSize := int(size + 5)
	if totalSize > l {
		if atEOF {
			return 0, nil, fmt.Errorf("message error - message is broken, msg:\n%v", string(data))
		}
		return 0, nil, nil
	}

	msgBytes := make([]byte, totalSize, totalSize)
	copy(msgBytes, data[:totalSize])
	return totalSize, msgBytes, nil
}
