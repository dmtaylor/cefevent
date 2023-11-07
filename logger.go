package cefevent

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// InvalidCefVersionErr error when provided an invalid CEF version. Value should be 0 or 1
var InvalidCefVersionErr = errors.New("invalid cef version")

// LoggerConfigOption is a configuring function for a Logger
type LoggerConfigOption func(l *Logger)

// WithCefVersion overwrite default CEF version. Returns InvalidCefVersionErr if an invalid value is set
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
	addSyslogHeader bool
	cefVersion      byte
	out             io.Writer
	getTime         func() time.Time // this is cursed for testing. You basically always want time.Now() for this
	DeviceVendor    string
	DeviceProduct   string
	DeviceVersion   string
}

// NewLogger creates a new CEF v1 event logger with default values.
func NewLogger(out io.Writer, deviceVendor, deviceProduct, deviceVersion string, fns ...LoggerConfigOption) *Logger {
	l := &Logger{
		addSyslogHeader: true,
		cefVersion:      1,
		out:             out,
		getTime:         time.Now,
		DeviceVendor:    deviceVendor,
		DeviceProduct:   deviceProduct,
		DeviceVersion:   deviceVersion,
	}
	for _, fn := range fns {
		fn(l)
	}
	return l
}

func (l *Logger) Log(deviceEventClassId string, name string, severity string, extensions Extensions) error {
	b := strings.Builder{}
	if l.addSyslogHeader {
		b.WriteString(l.getTime().Format(time.Stamp))
		hostname, err := os.Hostname()
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
	//TODO add extension formatting here
	_, err := l.out.Write([]byte(b.String()))
	if err != nil {
		return fmt.Errorf("failed to write log: %w", err)
	}
	return nil
}

func escapeHeaderField(field string) string {
	// TODO implement

	return field
}
