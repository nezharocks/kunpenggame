package core

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
)

// todo
func sendMsg(w *bufio.Writer, msg *APIMsg) error {
	msgBytes, err := json.Marshal(&msg)
	if err != nil {
		return err
	}

	msgLen := len(msgBytes)
	if msgLen >= 100000 {
		return errors.New("msgBytes' length is >= 100000 (max msg size)")
	}

	msgLenStr := fmt.Sprintf("%05d", msgLen)
	log.Printf("msg size is %v\n", msgLenStr)
	sendBytes := append([]byte(msgLenStr), msgBytes...)
	fmt.Println("SEND->", string(sendBytes))

	if _, err := w.Write(sendBytes); err != nil {
		log.Println(err)
		return err
	}

	if err := w.Flush(); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func splitMsgSample(data []byte, atEOF bool) (int, []byte, error) {

	if atEOF && len(data) < 5 {
		return 0, nil, nil
	}

	size, err := strconv.ParseInt(string(data[:5]), 10, 32)

	if err != nil {
		return 0, nil, err
	}

	totalSize := int(size + 5)
	if totalSize > len(data) {
		return 0, nil, nil
	}

	msgBytes := make([]byte, size, size)
	copy(msgBytes, data[5:totalSize])

	return totalSize, msgBytes, nil
}

// ErrBrokenMsg is
var ErrBrokenMsg = errors.New("broken message")

func splitMsg(data []byte, atEOF bool) (int, []byte, error) {
	l := len(data)
	if atEOF && l < 5 {
		return 0, nil, ErrBrokenMsg
	}

	size, err := strconv.ParseInt(string(data[:5]), 10, 32)
	if err != nil {
		return 0, nil, errors.New("invalid message: " + err.Error())
	}

	totalSize := int(size + 5)
	if totalSize > l {
		if atEOF {
			return 0, nil, ErrBrokenMsg
		}
		return 0, nil, nil
	}

	msgBytes := make([]byte, size, size)
	copy(msgBytes, data[5:totalSize])

	return totalSize, msgBytes, nil
}
