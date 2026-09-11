package config

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseBytes(value string) (int64, error) {
	value = strings.TrimSpace(
		strings.ToUpper(value),
	)

	multiplier := int64(1)

	switch {
	case strings.HasSuffix(value, "KB"):
		multiplier = 1024
		value = strings.TrimSuffix(value, "KB")

	case strings.HasSuffix(value, "MB"):
		multiplier = 1024 * 1024
		value = strings.TrimSuffix(value, "MB")

	case strings.HasSuffix(value, "GB"):
		multiplier = 1024 * 1024 * 1024
		value = strings.TrimSuffix(value, "GB")

	case strings.HasSuffix(value, "B"):
		value = strings.TrimSuffix(value, "B")
	}

	size, err := strconv.ParseInt(
		strings.TrimSpace(value),
		10,
		64,
	)

	if err != nil {
		return 0, fmt.Errorf(
			"invalid size %q",
			value,
		)
	}

	if size < 0 {
		return 0, fmt.Errorf(
			"size cannot be negative",
		)
	}

	return size * multiplier, nil
}

func FormatBytes(value int64) string {
	switch {
	case value >= 1024*1024*1024:
		return fmt.Sprintf(
			"%dGB",
			value/(1024*1024*1024),
		)

	case value >= 1024*1024:
		return fmt.Sprintf(
			"%dMB",
			value/(1024*1024),
		)

	case value >= 1024:
		return fmt.Sprintf(
			"%dKB",
			value/1024,
		)

	default:
		return fmt.Sprintf(
			"%dB",
			value,
		)
	}
}
