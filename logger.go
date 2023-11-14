// Package cefevent provides a 'log'/'slog' like interface for logging CEF events.
package cefevent

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

var defaultLogger *Logger

func init() {
	defaultLogger = NewLogger(os.Stdout, "go", "cefevent", "v0.1")
}

var headerEscapeRegex = regexp.MustCompile(`([|\\])`)

// InvalidCefVersionErr error when provided an invalid CEF version. Value should be 0 or 1
var InvalidCefVersionErr = errors.New("invalid cef version")

// LoggerConfigOption is a configuring function for a Logger
type LoggerConfigOption func(l *Logger)

// MustLoggerConfig panics if err is non-nil. Useful for wrapping WithCefVersion at startup time.
func MustLoggerConfig(l LoggerConfigOption, err error) LoggerConfigOption {
	if err != nil {
		panic(err)
	}
	return l
}

// WithCefVersion overwrite default CEF version. Returns InvalidCefVersionErr if an invalid value is set. Current accepted
// versions are 0 & 1. Most users will not need to overwrite this value.
func WithCefVersion(ver byte) (LoggerConfigOption, error) {
	if ver != 0 && ver != 1 {
		return nil, InvalidCefVersionErr
	}
	return func(l *Logger) {
		l.cefVersion = ver
	}, nil

}

// OmitSyslogHeader omit the syslog header prefix on logged events. Useful if outputting to a file or piping to existing
// syslog implementation
func OmitSyslogHeader() LoggerConfigOption {
	return func(l *Logger) {
		l.addSyslogHeader = false
	}
}

// Logger is a logger for cef events
type Logger struct {
	// addSyslogHeader add syslog style header as per spec. Configurable to allow outputting to file, where that header is omitted
	addSyslogHeader bool
	// cefVersion should be 0 or 1
	cefVersion byte
	// out writer for output
	out io.Writer

	// Manually set time & hostname functions here. This is cursed for testing.
	getTime     func() time.Time       // You basically always want time.Now() for this
	getHostname func() (string, error) // use os.Hostname()

	// DeviceVendor device vendor in CEF header.
	DeviceVendor string

	// DeviceProduct product in CEF header. Ordered pair (DeviceVendor, DeviceProduct) should uniquely identify class of event
	DeviceProduct string

	// DeviceVersion device version in CEF header.
	DeviceVersion string
}

// NewLogger creates a new CEF v1 event logger with default values. This function should be used to create
func NewLogger(out io.Writer, deviceVendor, deviceProduct, deviceVersion string, fns ...LoggerConfigOption) *Logger {
	l := &Logger{
		addSyslogHeader: true,
		cefVersion:      1,
		out:             out,
		getTime:         time.Now,
		getHostname:     os.Hostname,
		DeviceVendor:    deviceVendor,
		DeviceProduct:   deviceProduct,
		DeviceVersion:   deviceVersion,
	}
	for _, fn := range fns {
		fn(l)
	}
	return l
}

// Log logs CEF event to configured writer
func (l *Logger) Log(deviceEventClassId, name, severity string, extensions Extensions) error {
	b := strings.Builder{}
	if l.addSyslogHeader {
		b.WriteString(l.getTime().Format(`Jan 2 15:04:05`))
		hostname, err := l.getHostname()
		if err != nil {
			return fmt.Errorf("failed to get hostname: %w", err)
		}
		b.WriteString(" " + hostname + " ")
	}
	b.WriteString(fmt.Sprintf("CEF:%d|%s|%s|%s|%s|%s|%s|",
		l.cefVersion,
		escapeHeaderField(l.DeviceVendor),
		escapeHeaderField(l.DeviceProduct),
		escapeHeaderField(l.DeviceVersion),
		escapeHeaderField(deviceEventClassId),
		escapeHeaderField(name),
		escapeHeaderField(severity),
	))
	b.WriteString(extensions.String())
	_, err := l.out.Write([]byte(b.String()))
	if err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}
	return nil
}

// Log logs CEF event with default logger
func Log(deviceEventClassId, name, severity string, extensions Extensions) error {
	return defaultLogger.Log(deviceEventClassId, name, severity, extensions)
}

// LogUnknown log CEF event with unknown severity
func (l *Logger) LogUnknown(deviceEventClassId, name string, extensions Extensions) error {
	return l.Log(deviceEventClassId, name, UnknownSeverity, extensions)
}

// LogUnknown log CEF event with unknown severity to default logger
func LogUnknown(deviceEventClassId, name string, extensions Extensions) error {
	return defaultLogger.LogUnknown(deviceEventClassId, name, extensions)
}

// LogLow log CEF event with low severity
func (l *Logger) LogLow(deviceEventClassId, name string, extensions Extensions) error {
	return l.Log(deviceEventClassId, name, LowSeverity, extensions)
}

// LogLow log CEF event with low severity to default logger
func LogLow(deviceEventClassId, name string, extensions Extensions) error {
	return defaultLogger.LogLow(deviceEventClassId, name, extensions)
}

// LogMedium log CEF event with medium severity
func (l *Logger) LogMedium(deviceEventClassId, name string, extensions Extensions) error {
	return l.Log(deviceEventClassId, name, MediumSeverity, extensions)
}

// LogMedium log CEF event with medium severity to default logger
func LogMedium(deviceEventClassId, name string, extensions Extensions) error {
	return defaultLogger.LogMedium(deviceEventClassId, name, extensions)
}

// LogHigh log CEF event with high severity
func (l *Logger) LogHigh(deviceEventClassId, name string, extensions Extensions) error {
	return l.Log(deviceEventClassId, name, HighSeverity, extensions)
}

// LogHigh log CEF event with high severity to default logger
func LogHigh(deviceEventClassId, name string, extensions Extensions) error {
	return defaultLogger.LogHigh(deviceEventClassId, name, extensions)
}

// LogVeryHigh log CEF event with very-high severity
func (l *Logger) LogVeryHigh(deviceEventClassId, name string, extensions Extensions) error {
	return l.Log(deviceEventClassId, name, VeryHighSeverity, extensions)
}

// LogVeryHigh log CEF event with very-high severity to default logger
func LogVeryHigh(deviceEventClassId, name string, extensions Extensions) error {
	return defaultLogger.LogVeryHigh(deviceEventClassId, name, extensions)
}

// SetDefaultLogger sets the default logger to a created one. Useful for using package level functions
func SetDefaultLogger(log *Logger) {
	defaultLogger = log
}

func escapeHeaderField(field string) string {
	return headerEscapeRegex.ReplaceAllString(field, "\\${1}")
}
