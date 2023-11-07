package cefevent

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateSeverity(t *testing.T) {
	tests := []struct {
		name    string
		sev     string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			"integer_severity",
			"5",
			assert.NoError,
		},
		{
			"medium_severity",
			"Medium",
			assert.NoError,
		},
		{
			"high_severity",
			"High",
			assert.NoError,
		},
		{
			"unknown_severity",
			"Unknown",
			assert.NoError,
		},
		{
			"invalid_int",
			"-5",
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, InvalidSeverityError, i)
			},
		},
		{
			"invalid_adj",
			"something else",
			func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, InvalidSeverityError, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, validateSeverity(tt.sev), fmt.Sprintf("validateSeverity(%v)", tt.sev))
		})
	}
}
