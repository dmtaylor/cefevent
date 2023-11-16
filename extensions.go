package cefevent

import (
	"net"
	"strconv"
	"strings"
	"time"
)

// Extensions represent the additional fields in a CEF event
type Extensions struct {

	// Message is an arbitrary message giving more details about the event
	Message string

	// DeviceAction is the action taken by device
	DeviceAction string

	// ApplicationProtocol application level protocol, example values are HTTP, HTTPS, SSHv2, Telnet, POP, and so on.
	ApplicationProtocol string

	// DestinationDnsDomain the DNS domain part of the complete fully qualified domain name (FQDN).
	DestinationDnsDomain string

	// DestinationServiceName the service targeted by this event. Example "sshd"
	DestinationServiceName string

	// DestinationTranslatedAddress identifies the translated destination that the event refers to in an IP network
	DestinationTranslatedAddress net.IP

	// DestinationTranslatedPort port after it was translated; for example, a firewall. Valid port numbers are 0 to 65535
	DestinationTranslatedPort *uint

	// DeviceDirection any information about what direction the observed communication has taken. 0 for inbound, 1 for outbound
	DeviceDirection *uint8

	// DeviceDnsDomain the DNS domain part of the complete fully qualified domain name (FQDN)
	DeviceDnsDomain string

	// DeviceExternalId a name that uniquely identifies the device generating this event.
	DeviceExternalId string

	// DeviceFacility the facility generating this event
	DeviceFacility string

	// DeviceInboundInterface interface on which the packet or data entered the device
	DeviceInboundInterface string

	// DeviceNtDomain the Windows domain name of the device address.
	DeviceNtDomain string

	// DeviceOutboundInterface interface on which the packet or data left the device
	DeviceOutboundInterface string

	// DevicePayloadId unique identifier for the payload associated with the event
	DevicePayloadId string

	// DeviceProcessName process name associated with event e.g. process creating syslog entry
	DeviceProcessName string

	// DeviceTranslatedAddress identifies the translated device address that the event refers to in an IP network.
	DeviceTranslatedAddress net.IP

	// DeviceTimeZone timezone for device generating event
	DeviceTimeZone *time.Location

	// DeviceAddress identifies the device address that an event refers to
	DeviceAddress net.IP

	// DeviceHostName FQDN associated with device node e.g. "evt.example.com"
	DeviceHostName string

	// DeviceMacAddress MAC address for device in event
	DeviceMacAddress net.HardwareAddr

	// DeviceProcessId is the PID of the process on the device generating the event
	DeviceProcessId *uint

	// DestinationHostName identifies the destination that an event refers to in a network. The format should be a
	// fully qualified domain name associated with the destination node when available. e.g. "sub.example.com" or "example"
	DestinationHostName string

	// DestinationMacAddress MAC address for destination referred to in event
	DestinationMacAddress net.HardwareAddr

	// DestinationNtDomain Windows domain name of destination address
	DestinationNtDomain string

	// DestinationProcessId process ID for destination process associated with event.
	DestinationProcessId *uint

	// DestinationUserPrivileges identify destination user's privileges e.g. "Administrator", "User", "Guest"
	DestinationUserPrivileges string

	// DestinationProcessName name of event's destination process e.g. "ftpd"
	DestinationProcessName string

	// DestinationPort valid port number for destination process. Between 0 & 65535
	DestinationPort *uint

	// DestinationAddress identifies the destination IP address the event refers to.
	DestinationAddress net.IP

	// DestinationUserId identifies the destination user by ID e.g. root is typically "0"
	DestinationUserId string

	// TODO add all extensions

	// CustomExtensions includes non-standard mappings in the extension field. Keys in the map shouldn't overlap with fields in the
	// CEF spec to avoid duplicate values
	CustomExtensions map[string]string
}

// String formats extension for including in CEF event
func (e Extensions) String() string {
	b := strings.Builder{}
	if e.Message != "" {
		b.WriteString("msg=" + escapeExtensionField(e.Message) + " ")
	}
	if e.DeviceAction != "" {
		b.WriteString("act=" + escapeExtensionField(e.DeviceAction) + " ")
	}
	if e.ApplicationProtocol != "" {
		b.WriteString("app=" + escapeExtensionField(e.ApplicationProtocol) + " ")
	}
	fcount, destinationStr := e.marshalDestinationFields()
	if fcount > 0 {
		b.WriteString(destinationStr)
	}
	fcount, deviceString := e.marshalDeviceFields()
	if fcount > 0 {
		b.WriteString(deviceString)
	}
	// TODO implement

	for k, v := range e.CustomExtensions {
		b.WriteString(escapeExtensionField(k) + "=" + escapeExtensionField(v) + " ")
	}
	return strings.TrimSpace(b.String())
}

func (e Extensions) marshalDeviceFields() (int, string) {
	c := 0
	b := strings.Builder{}

	if e.DeviceDirection != nil {
		c += 1
		b.WriteString("deviceDirection=" + strconv.FormatUint(uint64(*e.DeviceDirection), 10) + " ")
	}
	if e.DeviceDnsDomain != "" {
		c += 1
		b.WriteString("deviceDnsDomain=" + escapeExtensionField(e.DeviceDnsDomain) + " ")
	}
	if e.DeviceExternalId != "" {
		c += 1
		b.WriteString("deviceExternalId=" + escapeExtensionField(e.DeviceExternalId) + " ")
	}
	if e.DeviceFacility != "" {
		c += 1
		b.WriteString("deviceFacility=" + escapeExtensionField(e.DeviceFacility) + " ")
	}
	if e.DeviceInboundInterface != "" {
		c += 1
		b.WriteString("deviceInboundInterface=" + escapeExtensionField(e.DeviceInboundInterface) + " ")
	}
	if e.DeviceNtDomain != "" {
		c += 1
		b.WriteString("deviceNtInterface=" + escapeExtensionField(e.DeviceNtDomain) + " ")
	}
	if e.DeviceOutboundInterface != "" {
		c += 1
		b.WriteString("deviceOutboundInterface=" + escapeExtensionField(e.DeviceOutboundInterface) + " ")
	}
	if e.DevicePayloadId != "" {
		c += 1
		b.WriteString("devicePayloadId=" + escapeExtensionField(e.DevicePayloadId) + " ")
	}
	if e.DeviceProcessName != "" {
		c += 1
		b.WriteString("deviceProcessName=" + escapeExtensionField(e.DeviceProcessName) + " ")
	}
	if e.DeviceTimeZone != nil {
		c += 1
		b.WriteString("dtz=" + escapeExtensionField(e.DeviceTimeZone.String()) + " ")
	}
	if str := e.DeviceAddress.String(); str != "<nil>" {
		c += 1
		b.WriteString("dvc=" + str + " ")
	}
	if e.DeviceHostName != "" {
		c += 1
		b.WriteString("dcvhost=" + escapeExtensionField(e.DeviceHostName) + " ")
	}
	if len(e.DeviceMacAddress) != 0 {
		c += 1
		b.WriteString("dvcmac=" + e.DeviceMacAddress.String() + " ")
	}
	if e.DeviceProcessId != nil {
		c += 1
		b.WriteString("dvcpid=" + strconv.FormatUint(uint64(*e.DeviceProcessId), 10) + " ")
	}
	// TODO add custom mapped fields
	return c, b.String()
}

func (e Extensions) marshalDestinationFields() (int, string) {
	b := strings.Builder{}
	c := 0
	if e.DestinationDnsDomain != "" {
		c += 1
		b.WriteString("destinationDnsDomain=" + escapeExtensionField(e.DestinationDnsDomain) + " ")
	}
	if e.DestinationServiceName != "" {
		c += 1
		b.WriteString("destinationServiceName=" + escapeExtensionField(e.DestinationServiceName) + " ")
	}
	if str := e.DestinationTranslatedAddress.String(); str != "<nil>" {
		c += 1
		b.WriteString("destinationTranslatedAddress=" + escapeExtensionField(str) + " ")
	}
	if e.DestinationTranslatedPort != nil {
		c += 1
		b.WriteString("destinationTranslatedPort=" + strconv.FormatUint(uint64(*e.DestinationTranslatedPort), 10) + " ")
	}
	if e.DestinationHostName != "" {
		c += 1
		b.WriteString("dhost=" + escapeExtensionField(e.DestinationHostName) + " ")
	}
	if len(e.DestinationMacAddress) != 0 {
		c += 1
		b.WriteString("dmac=" + e.DestinationMacAddress.String() + " ")
	}
	if e.DestinationNtDomain != "" {
		c += 1
		b.WriteString("dntdom=" + escapeExtensionField(e.DestinationNtDomain) + " ")
	}
	if e.DestinationProcessId != nil {
		c += 1
		b.WriteString("dpid=" + strconv.FormatUint(uint64(*e.DestinationProcessId), 10) + " ")
	}
	if e.DestinationUserPrivileges != "" {
		c += 1
		b.WriteString("dpriv=" + escapeExtensionField(e.DestinationUserPrivileges) + " ")
	}
	if e.DestinationProcessName != "" {
		c += 1
		b.WriteString("dproc=" + escapeExtensionField(e.DestinationProcessName) + " ")
	}
	if e.DestinationPort != nil {
		c += 1
		b.WriteString("dpt=" + strconv.FormatUint(uint64(*e.DestinationPort), 10) + " ")
	}
	if str := e.DestinationAddress.String(); str != "<nil>" {
		c += 1
		b.WriteString("dst=" + str + " ")
	}
	if e.DestinationUserId != "" {
		c += 1
		b.WriteString("duid=" + escapeExtensionField(e.DestinationUserId) + " ")
	}
	// TODO add destination marshaling

	// TODO add custom mapped fields

	return c, b.String()
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

// Ptr is a convenience function to convert literal values to pointers
func Ptr[A any](v A) *A {
	return &v
}
