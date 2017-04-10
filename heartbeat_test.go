package stomp

import (
	"testing"
	"time"
)

func TestParseHeartBeat(t *testing.T) {
	tests := []struct {
		heartBeat string
		x         time.Duration
		y         time.Duration
		ok        bool
	}{
		{"", 0, 0, true},
		{"5000,150", 5 * time.Second, 150 * time.Millisecond, true},
		{",5", 0, 0, false},
		{" 5 , 5", 0, 0, false},
	}
	for _, test := range tests {
		x, y, ok := ParseHeartBeat(test.heartBeat)
		if x != test.x {
			t.Errorf("wrong x")
		}
		if y != test.y {
			t.Errorf("wrong y")
		}
		if ok != test.ok {
			t.Errorf("wrong ok")
		}
	}
}

func TestClientHeartBeat(t *testing.T) {
	tests := []struct {
		ConnectHeartBeat   string
		ConnectedHeartBeat string
		Out                time.Duration
		In                 time.Duration
	}{
		{"", "", 0, 0},
		{"0,0", "0,0", 0, 0},
		{"1000,0", "0,1000", time.Second, 0},
		{"0,1000", "1000,0", 0, time.Second},
		{"0,0", "250,5000", 0, 0},
		{"5000,5000", "0,0", 0, 0},
		{"1000,5000", "250,30000", time.Second * 30, time.Second * 5},
	}
	for _, test := range tests {
		cx, cy, ok := ParseHeartBeat(test.ConnectHeartBeat)
		if !ok {
			t.Fatal("invalid client heart beat")
		}
		sx, sy, ok := ParseHeartBeat(test.ConnectedHeartBeat)
		if !ok {
			t.Fatal("invalid server heart beat")
		}
		out, in := ClientHeartBeat(cx, cy, sx, sy)
		if test.Out != out {
			t.Errorf("out: %s != %s", out, test.Out)
		}
		if test.In != in {
			t.Errorf("in: %s != %s", in, test.In)
		}
	}
}

func TestServerHeartBeat(t *testing.T) {
	tests := []struct {
		ConnectHeartBeat   string
		ConnectedHeartBeat string
		Out                time.Duration
		In                 time.Duration
	}{
		{"", "", 0, 0},
		{"0,0", "0,0", 0, 0},
		{"1000,0", "0,1000", 0, time.Second},
		{"0,1000", "1000,0", time.Second, 0},
		{"0,0", "250,5000", 0, 0},
		{"5000,5000", "0,0", 0, 0},
		{"1000,5000", "250,30000", time.Second * 5, time.Second * 30},
	}
	for _, test := range tests {
		cx, cy, ok := ParseHeartBeat(test.ConnectHeartBeat)
		if !ok {
			t.Fatal("invalid client heart beat")
		}
		sx, sy, ok := ParseHeartBeat(test.ConnectedHeartBeat)
		if !ok {
			t.Fatal("invalid server heart beat")
		}
		out, in := ServerHeartBeat(cx, cy, sx, sy)
		if test.Out != out {
			t.Errorf("out: %s != %s", out, test.Out)
		}
		if test.In != in {
			t.Errorf("in: %s != %s", in, test.In)
		}
	}
}
