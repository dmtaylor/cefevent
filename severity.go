package cefevent

import (
	"errors"
	"strconv"
)

const UnknownSeverity = "Unknown"
const LowSeverity = "Low"
const MediumSeverity = "Medium"
const HighSeverity = "High"
const VeryHighSeverity = "Very-High"

var InvalidSeverityError = errors.New("invalid severity")

func validateSeverity(sev string) error {
	switch sev {
	case UnknownSeverity:
		fallthrough
	case LowSeverity:
		fallthrough
	case MediumSeverity:
		fallthrough
	case HighSeverity:
		fallthrough
	case VeryHighSeverity:
		return nil
	}
	v, err := strconv.Atoi(sev)
	if err != nil {
		return InvalidSeverityError
	}
	if v > 10 || v < 0 {
		return InvalidSeverityError
	}

	return nil
}
