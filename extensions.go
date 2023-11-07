package cefevent

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
	// TODO implement
	return f
}
