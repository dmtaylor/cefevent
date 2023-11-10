// Package cefevent provides a 'log' like interface for logging CEF events.
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

var headerEscapeRegex = regexp.MustCompile(`([|\\])`)

// InvalidCefVersionErr error when provided an invalid CEF version. Value should be 0 or 1
var InvalidCefVersionErr = errors.New("invalid cef version")

// LoggerConfigOption is a configuring function for a Logger
type LoggerConfigOption func(l *Logger)

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
func (l *Logger) Log(deviceEventClassId string, name string, severity string, extensions Extensions) error {
	b := strings.Builder{}
	if l.addSyslogHeader {
		b.WriteString(l.getTime().Format(time.Stamp))
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

func escapeHeaderField(field string) string {
	return headerEscapeRegex.ReplaceAllString(field, "\\${1}")
}
