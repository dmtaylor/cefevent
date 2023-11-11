package cefevent

import "strings"

// Extensions represent the additional fields in a CEF event
type Extensions struct {
	// DeviceAction is action taken by device
	DeviceAction string

	// ApplicationProtocol Application level protocol, example values are HTTP, HTTPS, SSHv2, Telnet, POP, and so on.
	ApplicationProtocol string

	// DestinationDnsDomain The DNS domain part of the complete fully qualified domain name (FQDN).
	DestinationDnsDomain string

	// DestinationServiceName The service targeted by this event. Example "sshd"
	DestinationServiceName string
	// TODO add all extensions

	// CustomExtensions includes non-standard mappings in the extension field. Keys in the map shouldn't overlap with fields in the
	// CEF spec to avoid duplicate values
	CustomExtensions map[string]string
}

func (e Extensions) String() string {
	b := strings.Builder{}
	if e.DeviceAction != "" {
		b.WriteString("act=" + escapeExtensionField(e.DeviceAction))
	}
	if e.ApplicationProtocol != "" {
		b.WriteString("app=" + escapeExtensionField(e.ApplicationProtocol))
	}
	if e.DestinationDnsDomain != "" {
		b.WriteString("destinationDnsDomain=" + escapeExtensionField(e.DestinationDnsDomain))
	}
	if e.DestinationServiceName != "" {
		b.WriteString("destinationServiceName=" + escapeExtensionField(e.DestinationServiceName))
	}
	// TODO implement

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
