package stomp

import (
	"strconv"
	"strings"
	"time"
)

// Parse "heart-beat" header into x and y.
//   x = sender can guarantee sending every x
//   y = would like to receive every y
func ParseHeartBeat(heartBeat string) (time.Duration, time.Duration, bool) {
	if heartBeat == "" {
		return 0, 0, true
	}
	xy := strings.Split(heartBeat, ",")
	if len(xy) != 2 {
		return 0, 0, false
	}
	x, err := strconv.Atoi(xy[0])
	if err != nil {
		return 0, 0, false
	}
	y, err := strconv.Atoi(xy[1])
	if err != nil {
		return 0, 0, false
	}
	return time.Duration(x) * time.Millisecond, time.Duration(y) * time.Millisecond, true
}

// Negotiate client heart beat. Client cx, cy should come from CONNECT frame, and server
// sx, sy should come from CONNECTED frame.
func ClientHeartBeat(cx, cy, sx, sy time.Duration) (time.Duration, time.Duration) {
	var out, in time.Duration
	if cx > 0 && sy > 0 {
		out = maxDuration(cx, sy)
	}
	if cy > 0 && sx > 0 {
		in = maxDuration(cy, sx)
	}
	return out, in
}

// Negotiate server heart beat. Client cx, cy should come from CONNECT frame, server
// sx, sy should come from CONNECTED frame.
func ServerHeartBeat(cx, cy, sx, sy time.Duration) (time.Duration, time.Duration) {
	var out, in time.Duration
	if sx > 0 && cy > 0 {
		out = maxDuration(sx, cy)
	}
	if sy > 0 && cx > 0 {
		in = maxDuration(sy, cx)
	}
	return out, in
}

// Longer duration.
func maxDuration(x, y time.Duration) time.Duration {
	if x > y {
		return x
	}
	return y
}
