package core

import (
	"fmt"
	"testing"
)

func TestNewMapFromStrTing(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name    string
		args    args
		wantM   *Map
		wantErr bool
	}{
		{
			name: "TestNewMapFromString - map1",
			args: args{
				data: map1,
			},
			wantErr: false,
		},
		{
			name: "TestNewMapFromString - map2",
			args: args{
				data: map2,
			},
			wantErr: false,
		},
		{
			name: "TestNewMapFromString - map3",
			args: args{
				data: map3,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotM, err := NewMapFromString(tt.args.data)
			fmt.Printf("%v\n\n", gotM.JSON())
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMapFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(gotM, tt.wantM) {
			// 	t.Errorf("NewMapFromString() = %v, want %v", gotM, tt.wantM)
			// }
		})
	}
}
