package cefevent

import "strings"

// Extensions represent the additional fields in a CEF event
type Extensions struct {
	// TODO add extensions

	// Extras includes non-standard mappings in the extension field. Keys in the map shouldn't overlap with fields in the
	// CEF spec to avoid duplicate values
	CustomExtensions map[string]string
}

func (e Extensions) String() string {
	b := strings.Builder{}
	// TODO implement

	// TODO remove stub return
	for k, v := range e.CustomExtensions {
		b.WriteString(escapeExtensionField(k) + "=" + escapeExtensionField(v) + " ")
	}
	return strings.TrimSpace(b.String())
}

func escapeExtensionField(f string) string {
	b := strings.Builder{}
	for _, r := range []rune(f) {
		switch r {
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '=':
			b.WriteString(`\=`)
		case '\\':
			b.WriteString(`\\`)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
