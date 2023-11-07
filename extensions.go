package cefevent

import "strings"

// Extensions represent the additional fields in
type Extensions struct {
	// TODO add extensions

	// Extras includes non-standard mappings in the extension field.
	Extras map[string]string
}

func (e Extensions) String() string {
	// TODO implement

	// TODO remove stub return
	return ""
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
