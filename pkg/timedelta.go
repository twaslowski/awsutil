package pkg

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseTimeDelta(input string) (time.Duration, error) {
	re := regexp.MustCompile(`(\d+)([smhdw])`)
	matches := re.FindAllStringSubmatch(input, -1)
	if matches == nil {
		return 0, fmt.Errorf("invalid format")
	}

	var total time.Duration
	for _, match := range matches {
		num, _ := strconv.Atoi(match[1])
		unit := strings.ToLower(match[2])
		var dur time.Duration
		switch unit {
		case "s":
			dur = time.Duration(num) * time.Second
		case "m":
			dur = time.Duration(num) * time.Minute
		case "h":
			dur = time.Duration(num) * time.Hour
		case "d":
			dur = time.Duration(num) * 24 * time.Hour
		case "w":
			dur = time.Duration(num) * 7 * 24 * time.Hour
		default:
			return 0, fmt.Errorf("unknown unit: %s", unit)
		}
		total += dur
	}

	return total, nil
}
