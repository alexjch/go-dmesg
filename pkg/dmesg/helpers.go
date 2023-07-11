// SPDX-License-Identifier: MIT
package dmesg

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	kmsgLine = `^(?P<priority>\d{1,2}),(?P<sequence>\d+),(?P<monotonic_timestamp>\d+),.*;(?P<message>.*)$`
	kmsgCont = `^\s.*$`
)

var kmsgLineRe = regexp.MustCompile(kmsgLine)
var kmsgContRe = regexp.MustCompile(kmsgCont)

func parseLogPrefix(m []string) (uint8, uint64, int64, error) {
	// parses priority
	p, err := strconv.ParseUint(m[1], 10, 8)
	if err != nil {
		return 0, 0, 0, err
	}
	// parses sequence number
	s, err := strconv.ParseUint(m[2], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	// parses monotonic timestamp
	mon, err := strconv.ParseInt(m[3], 10, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	return uint8(p), s, mon, nil
}

func continuation(line string) bool {
	return kmsgContRe.MatchString(line)
}

func kmsgMatches(line string) []string {
	return kmsgLineRe.FindStringSubmatch(line)
}

// osUptime parses /proc/uptime and returns the uptime as a time.Duration
// notice this "may be inaccurate!". Read more in dmesg source code.
// https://github.com/util-linux/util-linux/blob/master/sys-utils/dmesg.c
func osUptime() (time.Duration, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return time.Duration(0), err
	}
	uptimeArray := strings.Split(string(data), " ")
	if len(uptimeArray) != 2 {
		return time.Duration(0), fmt.Errorf("invalid /proc/uptime: %s", string(data))
	}
	uptime, err := strconv.ParseFloat(uptimeArray[0], 64)
	if err != nil {
		return time.Duration(0), err
	}
	return time.Duration(uptime * float64(time.Second)), nil
}
