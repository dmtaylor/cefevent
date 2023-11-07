package cefevent

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLogger_Log(t *testing.T) {
	type fields struct {
		addSyslogHeader bool
		cefVersion      byte
		getTime         func() time.Time
		getHostname     func() (string, error)
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			_ = &Logger{
				addSyslogHeader: tt.fields.addSyslogHeader,
				cefVersion:      tt.fields.cefVersion,
				out:             buf,
				getTime:         tt.fields.getTime,
				getHostname:     tt.fields.getHostname,
				DeviceVendor:    tt.fields.DeviceVendor,
				DeviceProduct:   tt.fields.DeviceProduct,
				DeviceVersion:   tt.fields.DeviceVersion,
			}
			// TODO write testify code here
			//if err := l.Log(tt.args.deviceEventClassId, tt.args.name, tt.args.severity, tt.args.extensions); (err != nil) != tt.wantErr {
			//	t.Errorf("Log() error = %v, wantErr %v", err, tt.wantErr)
			//}
		})
	}
}

func TestNewLogger(t *testing.T) {
	type args struct {
		deviceVendor  string
		deviceProduct string
		deviceVersion string
		fns           []LoggerConfigOption
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		want    *Logger
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			got := NewLogger(out, tt.args.deviceVendor, tt.args.deviceProduct, tt.args.deviceVersion, tt.args.fns...)
			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("NewLogger() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOmitSyslogHeader(t *testing.T) {
	// TODO implement test
}

func TestWithCefVersion(t *testing.T) {
	tests := []struct {
		name          string
		ver           byte
		want          LoggerConfigOption
		expectedError error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO update test case here
		})
	}
}

func Test_escapeHeaderField(t *testing.T) {
	tests := []struct {
		name  string
		field string
		want  string
	}{
		// TODO: Add test cases.
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
