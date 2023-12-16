package cefevent

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_escapeExtensionField(t *testing.T) {
	tests := []struct {
		name string
		f    string
		want string
	}{
		{
			"empty",
			"",
			"",
		},
		{
			"base",
			"regular_value",
			"regular_value",
		},
		{
			"newline",
			"regular\nvalue",
			`regular\nvalue`,
		},
		{
			"equals",
			"answer=42",
			"answer\\=42",
		},
		{
			"multi",
			"answer=\r42\\100",
			`answer\=\r42\\100`,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, escapeExtensionField(tt.f), "escapeExtensionField(%v)", tt.f)
		})
	}
}

func TestExtensions_String(t *testing.T) {
	tests := []struct {
		name string
		e    Extensions
		want string
	}{
		{
			"empty",
			Extensions{},
			"",
		},
		{
			"extras",
			Extensions{
				Message:          "value\nwithnewline",
				CustomExtensions: map[string]string{"extra": "value"},
			},
			"msg=value\\nwithnewline extra=value",
		},
		{
			"ip_port_value",
			Extensions{
				DestinationTranslatedPort:    ptr(uint(22)),
				DestinationTranslatedAddress: net.IP{192, 168, 0, 1},
			},
			"destinationTranslatedAddress=192.168.0.1 destinationTranslatedPort=22",
		},
		{
			"file_data_1",
			Extensions{
				FileSize:             ptr(uint(2048)),
				FileType:             "normal",
				FileModificationTime: testTime(),
				FileCreateTime:       testTime(),
				FileId:               "6452",
				FileName:             "example.txt",
			},
			"fileCreateTime=1699530320000 fileId=6452 fileModificationTime=1699530320000 fileType=normal fname=example.txt fsize=2048",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.e.String(), "String()")
		})
	}
}

// ptr is a convenience function to convert literal values to pointers
func ptr[A any](v A) *A {
	return &v
}
