package cefevent

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func testTime() time.Time {
	return time.Date(2023, 11, 9, 11, 45, 20, 0, time.UTC)
}

func testHostname() (string, error) {
	return "testhost", nil
}

var stubWriterError = errors.New("underlying writer error")

type errorWriter struct{}

func (e errorWriter) Write(_ []byte) (int, error) {
	return 0, stubWriterError
}

func TestLogger_Log(t *testing.T) {
	type fields struct {
		addSyslogHeader bool
		cefVersion      byte
		DeviceVendor    string
		DeviceProduct   string
		DeviceVersion   string
	}
	type args struct {
		deviceEventClassId string
		name               string
		severity           string
		extensions         Extensions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"simple",
			fields{
				addSyslogHeader: true,
				cefVersion:      1,
				DeviceVendor:    "cyberdyne",
				DeviceProduct:   "skynet",
				DeviceVersion:   "0.9.0",
			},
			args{
				deviceEventClassId: "1000",
				name:               "testevent",
				severity:           "Low",
				extensions:         Extensions{},
			},
			"Nov 9 11:45:20 testhost CEF:1|cyberdyne|skynet|0.9.0|1000|testevent|Low|",
			assert.NoError,
		},
		{
			"omit_syslog_and_cef0",
			fields{
				addSyslogHeader: false,
				cefVersion:      0,
				DeviceVendor:    "cyberdyne",
				DeviceProduct:   "skynet",
				DeviceVersion:   "0.9.1",
			},
			args{
				deviceEventClassId: "1001",
				name:               "testeventtofile",
				severity:           LowSeverity,
				extensions:         Extensions{},
			},
			"CEF:0|cyberdyne|skynet|0.9.1|1001|testeventtofile|Low|",
			assert.NoError,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			l := &Logger{
				addSyslogHeader: tt.fields.addSyslogHeader,
				cefVersion:      tt.fields.cefVersion,
				out:             buf,
				getTime:         testTime, // pin time and hostname for tests
				getHostname:     testHostname,
				DeviceVendor:    tt.fields.DeviceVendor,
				DeviceProduct:   tt.fields.DeviceProduct,
				DeviceVersion:   tt.fields.DeviceVersion,
			}
			if tt.wantErr(t, l.Log(tt.args.deviceEventClassId, tt.args.name, tt.args.severity, tt.args.extensions), fmt.Sprintf("Log(%v, %v, %v, %v)", tt.args.deviceEventClassId, tt.args.name, tt.args.severity, tt.args.extensions)) {
				assert.Equal(t, tt.want, buf.String())
			}
		})
	}
}

func TestLogger_LogError(t *testing.T) {
	l := &Logger{
		addSyslogHeader: false,
		cefVersion:      1,
		getTime:         testTime, // pin time and hostname for tests
		getHostname:     testHostname,
		DeviceVendor:    "not",
		DeviceProduct:   "relevant",
		DeviceVersion:   "1",
		out:             &errorWriter{},
	}
	err := l.Log("9001", "scanner", VeryHighSeverity, Extensions{})
	assert.EqualError(t, err, "failed to write log: underlying writer error")
}

func TestNewLogger(t *testing.T) {
	cef0, _ := WithCefVersion(0)
	type args struct {
		deviceVendor  string
		deviceProduct string
		deviceVersion string
		fns           []LoggerConfigOption
	}
	tests := []struct {
		name string
		args args
		want *Logger
	}{
		{
			"basic",
			args{
				"Grand Trunks Semaphore Company",
				"SoftwareClacks",
				"1.0.0",
				[]LoggerConfigOption{},
			},
			&Logger{
				addSyslogHeader: true,
				cefVersion:      1,
				out:             &bytes.Buffer{},
				getTime:         time.Now,
				getHostname:     os.Hostname,
				DeviceVendor:    "Grand Trunks Semaphore Company",
				DeviceProduct:   "SoftwareClacks",
				DeviceVersion:   "1.0.0",
			},
		},
		{
			"with_omit_header",
			args{
				"Daystrom Data Concepts",
				"datalore",
				"1.0.1",
				[]LoggerConfigOption{OmitSyslogHeader()},
			},
			&Logger{
				addSyslogHeader: false,
				cefVersion:      1,
				out:             &bytes.Buffer{},
				getTime:         time.Now,
				getHostname:     os.Hostname,
				DeviceVendor:    "Daystrom Data Concepts",
				DeviceProduct:   "datalore",
				DeviceVersion:   "1.0.1",
			},
		},
		{
			name: "with_cef_version",
			args: args{
				deviceVendor:  "Black Mesa",
				deviceProduct: "Cascade Resonator",
				deviceVersion: "1.0.3",
				fns:           []LoggerConfigOption{cef0},
			},
			want: &Logger{
				addSyslogHeader: true,
				cefVersion:      0,
				out:             &bytes.Buffer{},
				getTime:         time.Now,
				getHostname:     os.Hostname,
				DeviceVendor:    "Black Mesa",
				DeviceProduct:   "Cascade Resonator",
				DeviceVersion:   "1.0.3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			l := NewLogger(out, tt.args.deviceVendor, tt.args.deviceProduct, tt.args.deviceVersion, tt.args.fns...)
			assert.Equal(t, tt.want.addSyslogHeader, l.addSyslogHeader)
			assert.Equal(t, tt.want.cefVersion, l.cefVersion)
			assert.Equal(t, tt.want.DeviceVendor, l.DeviceVendor)
			assert.Equal(t, tt.want.DeviceProduct, l.DeviceProduct)
			assert.Equal(t, tt.want.DeviceVersion, l.DeviceVersion)
		})
	}
}

func TestWithCefVersion_error(t *testing.T) {
	_, err := WithCefVersion(100)
	assert.ErrorIs(t, err, InvalidCefVersionErr)
}

//	func TestOmitSyslogHeader(t *testing.T) {
//		// TODO implement test
//	}
//
//	func TestWithCefVersion(t *testing.T) {
//		tests := []struct {
//			name          string
//			ver           byte
//			want          LoggerConfigOption
//			expectedError error
//		}{
//			// TODO: Add test cases.
//		}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				// TODO update test case here
//			})
//		}
//	}
func Test_escapeHeaderField(t *testing.T) {
	tests := []struct {
		name  string
		field string
		want  string
	}{
		{
			"empty",
			"",
			"",
		},
		{
			"unaffected",
			"vendorheader",
			"vendorheader",
		},
		{
			"spaces",
			"vendor header",
			"vendor header",
		},
		{
			"vertical bars",
			"vendor|header",
			`vendor\|header`,
		},
		{
			"equals",
			"vendor=header",
			"vendor=header",
		},
		{
			"backslash",
			`vendor\header`,
			`vendor\\header`,
		},
		{
			"multi",
			`vendor\header|parttwo`,
			`vendor\\header\|parttwo`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := escapeHeaderField(tt.field)
			assert.Equal(t, tt.want, got)
		})
	}
}
