package cefevent

import (
	"net"
	"strconv"
	"strings"
)

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

	// DestinationTranslatedAddress Identifies the translated destination that the event refers to in an IP network
	DestinationTranslatedAddress net.IP

	// DestinationTranslatedPort Port after it was translated; for example, a firewall. Valid port numbers are 0 to 65535
	DestinationTranslatedPort *uint

	// DeviceDirection Any information about what direction the observed communication has taken. 0 for inbound, 1 for outbound
	DeviceDirection *uint8

	// DeviceDnsDomain The DNS domain part of the complete fully qualified domain name (FQDN)
	DeviceDnsDomain string
	// TODO add all extensions

	// CustomExtensions includes non-standard mappings in the extension field. Keys in the map shouldn't overlap with fields in the
	// CEF spec to avoid duplicate values
	CustomExtensions map[string]string
}

// String formats extension for including in CEF event
func (e Extensions) String() string {
	b := strings.Builder{}
	if e.DeviceAction != "" {
		b.WriteString("act=" + escapeExtensionField(e.DeviceAction) + " ")
	}
	if e.ApplicationProtocol != "" {
		b.WriteString("app=" + escapeExtensionField(e.ApplicationProtocol) + " ")
	}
	if e.DestinationDnsDomain != "" {
		b.WriteString("destinationDnsDomain=" + escapeExtensionField(e.DestinationDnsDomain) + " ")
	}
	if e.DestinationServiceName != "" {
		b.WriteString("destinationServiceName=" + escapeExtensionField(e.DestinationServiceName) + " ")
	}
	if str := e.DestinationTranslatedAddress.String(); str != "<nil>" {
		b.WriteString("destinationTranslatedAddress=" + escapeExtensionField(str) + " ")
	}
	if e.DestinationTranslatedPort != nil {
		b.WriteString("destinationTranslatedPort=" + strconv.FormatUint(uint64(*e.DestinationTranslatedPort), 10) + " ")
	}
	if e.DeviceDirection != nil {
		b.WriteString("deviceDirection=" + strconv.FormatUint(uint64(*e.DeviceDirection), 10) + " ")
	}
	if e.DeviceDnsDomain != "" {
		b.WriteString("deviceDnsDomain=" + e.DeviceDnsDomain + " ")
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
