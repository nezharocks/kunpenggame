package core

import (
	"bufio"
	"bytes"
	"testing"
)

func Test_sendMsg(t *testing.T) {
	type args struct {
		w   *bufio.Writer
		msg *Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sendMsg - succeed to send a registration msg",
			args: args{
				w: bufio.NewWriterSize(new(bytes.Buffer), 1024*10),
				msg: &Message{
					Name: "registration",
					Payload: &Registration{
						TeamID:   1001,
						TeamName: "daolaji",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := sendMsg(tt.args.w, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("sendMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
