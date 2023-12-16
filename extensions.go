package cefevent

import (
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Event type constants for use in the Type field
const (
	BaseEventType        = 0 // Used for base type. Note: will be omitted as per spec
	AggregatedEventType  = 1 // Used for aggregated events
	CorrelationEventType = 2 // Used for correlated events
	ActionEventType      = 3 // Used for action events
)

// Extensions represent the additional fields in a CEF event
type Extensions struct {
	//// General Event Fields

	// Message is an arbitrary message giving more details about the event
	Message string

	// BaseEventCount is the number of times this same event was observed. Omitted if count is less than 2.
	BaseEventCount int

	// ApplicationProtocol application level protocol, example values are HTTP, HTTPS, SSHv2, Telnet, POP, and so on.
	ApplicationProtocol string

	// EndTime is the time at which activity associated with the event ended
	EndTime time.Time

	// ExternalId is an ID used by the originating device.
	// Typically this is a monotonically increasing value, e.g. a Session ID
	ExternalId string

	// Type is the type of event.
	// 0 for base, 1 for aggregated, 2 for correlation, and 3 for action.
	// Base event types will be omitted.
	Type byte

	// BytesIn is the number of incoming bytes transferred from source to destination
	BytesIn *uint

	// BytesOut number of outbound bytes transferred from destination to source.
	BytesOut *uint

	// Outcome is the outcome for the event e.g. "failure"
	Outcome string

	// TransportProtocol identifies the layer 4 protocol used e.g. TCP
	TransportProtocol string

	// Reason is the audit event was generated e.g. "bad password"
	Reason string

	//// Source Fields

	// SourceHostName is the FQDN of the source machine
	SourceHostName string

	// SourceMacAddress is the MAC address of the source machine.
	SourceMacAddress net.HardwareAddr

	// SourceNtDomain is the Windows domain name for the source machine.
	SourceNtDomain string

	// SourceDnsDomain is the DNS domain name portion of the FQDN of the source machine.
	SourceDnsDomain string

	// SourceServiceName is the name of the service generating the event.
	SourceServiceName string

	// SourceTranslatedAddress is the translated IP address of the source machine.
	SourceTranslatedAddress net.IP

	// SourceTranslatedPort is the translated port number of the source machine (e.g. by NAT-ing).
	SourceTranslatedPort *uint

	// SourceProcessId is the PID of the originating process for the event.
	SourceProcessId *int

	//// Destination Fields

	// DestinationDnsDomain the DNS domain part of the complete fully qualified domain name (FQDN).
	DestinationDnsDomain string

	// DestinationServiceName the service targeted by this event. Example "sshd"
	DestinationServiceName string

	// DestinationTranslatedAddress identifies the translated destination that the event refers to in an IP network
	DestinationTranslatedAddress net.IP

	// DestinationTranslatedPort port after it was translated; for example, a firewall. Valid port numbers are 0 to 65535
	DestinationTranslatedPort *uint

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

	//// Device Fields

	// DeviceAction is the action taken by device
	DeviceAction string

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

	// DeviceProcessName process name associated with event e.g. process creating syslog entry.
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

	// DeviceReceiptTime is the time at which the event was received
	DeviceReceiptTime time.Time

	//// File fields

	// FileCreateTime is the time when the file was created
	FileCreateTime time.Time

	// FileHash is a hash of a referenced file
	FileHash string

	// FileId is an ID associated with the file (e.g. inode)
	FileId string

	// FileModificationTime is the time when the file was last modified
	FileModificationTime time.Time

	// FilePath is the absolute path of the file, including the filename.
	FilePath string

	// FilePermission is the permission string for the file
	FilePermission string

	// FileType is the type of file (normal, pipe, socket, etc)
	FileType string

	// FileName is the name of file only (without path)
	FileName string

	// FileSize is the size of the referenced file in bytes
	FileSize *uint

	// OldFileCreateTime is the time when the old file was created
	OldFileCreateTime time.Time

	// OldFileHash is the file hash for the old file
	OldFileHash string

	// OldFileId is the ID for the old file e.g. inode number
	OldFileId string

	// OldFileModificationTime is the time when the old file was last modified
	OldFileModificationTime time.Time

	// OldFileName is the filename for the old file referenced in the event
	OldFileName string

	// OldFilePath is the absolute path to the old file, including file name.
	OldFilePath string

	// OldFilePermission is the permission string for the for the old file
	OldFilePermission string

	// OldFileSize is the size in bytes of the old file.
	OldFileSize *uint

	// OldFileType is the type of the old file (pipe, socket, etc.)
	OldFileType string

	//// HTTP fields

	// RequestUrl is the full URL for an HTTP request, including protocol.
	RequestUrl url.URL

	// RequestClientApplication is the user-agent associated with the request.
	RequestClientApplication string

	// RequestContext is the context for the request origination (e.g. HTTP Referrer)
	RequestContext string

	// RequestCookies is the cookie strings associated with the request
	RequestCookies string

	// RequestMethod is the HTTP verb for the request (e.g. "GET")
	RequestMethod string
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
	if e.BaseEventCount > 1 {
		b.WriteString("cnt=" + strconv.FormatInt(int64(e.BaseEventCount), 10) + " ")
	}
	if !e.EndTime.IsZero() {
		b.WriteString("end=" + strconv.FormatInt(e.EndTime.UnixMilli(), 10) + " ") // Use unix time here
	}
	if e.ExternalId != "" {
		b.WriteString("externalId=" + escapeExtensionField(e.ExternalId) + " ")
	}
	if e.Type != 0 {
		b.WriteString("type=" + strconv.FormatInt(int64(e.Type), 10) + " ")
	}
	if e.BytesIn != nil {
		b.WriteString("in=" + strconv.FormatUint(uint64(*e.BytesIn), 10) + " ")
	}
	if e.BytesOut != nil {
		b.WriteString("out=" + strconv.FormatUint(uint64(*e.BytesOut), 10) + " ")
	}
	if e.Outcome != "" {
		b.WriteString("outcome=" + escapeExtensionField(e.Outcome) + " ")
	}
	if e.TransportProtocol != "" {
		b.WriteString("proto=" + escapeExtensionField(e.TransportProtocol) + " ")
	}
	if e.Reason != "" {
		b.WriteString("reason=" + escapeExtensionField(e.Reason) + " ")
	}
	destinationStr := e.marshalDestinationFields()
	if len(destinationStr) > 0 {
		b.WriteString(destinationStr)
	}
	deviceString := e.marshalDeviceFields()
	if len(deviceString) > 0 {
		b.WriteString(deviceString)
	}
	fileString := e.marshalFileFields()
	if len(fileString) > 0 {
		b.WriteString(fileString)
	}
	// TODO implement

	for k, v := range e.CustomExtensions {
		b.WriteString(escapeExtensionField(k) + "=" + escapeExtensionField(v) + " ")
	}
	return strings.TrimSpace(b.String())
}

func (e Extensions) marshalDeviceFields() string {
	b := strings.Builder{}

	if e.DeviceDirection != nil {
		b.WriteString("deviceDirection=" + strconv.FormatUint(uint64(*e.DeviceDirection), 10) + " ")
	}
	if e.DeviceDnsDomain != "" {
		b.WriteString("deviceDnsDomain=" + escapeExtensionField(e.DeviceDnsDomain) + " ")
	}
	if e.DeviceExternalId != "" {
		b.WriteString("deviceExternalId=" + escapeExtensionField(e.DeviceExternalId) + " ")
	}
	if e.DeviceFacility != "" {
		b.WriteString("deviceFacility=" + escapeExtensionField(e.DeviceFacility) + " ")
	}
	if e.DeviceInboundInterface != "" {
		b.WriteString("deviceInboundInterface=" + escapeExtensionField(e.DeviceInboundInterface) + " ")
	}
	if e.DeviceNtDomain != "" {
		b.WriteString("deviceNtInterface=" + escapeExtensionField(e.DeviceNtDomain) + " ")
	}
	if e.DeviceOutboundInterface != "" {
		b.WriteString("deviceOutboundInterface=" + escapeExtensionField(e.DeviceOutboundInterface) + " ")
	}
	if e.DevicePayloadId != "" {
		b.WriteString("devicePayloadId=" + escapeExtensionField(e.DevicePayloadId) + " ")
	}
	if e.DeviceProcessName != "" {
		b.WriteString("deviceProcessName=" + escapeExtensionField(e.DeviceProcessName) + " ")
	}
	if e.DeviceTimeZone != nil {
		b.WriteString("dtz=" + escapeExtensionField(e.DeviceTimeZone.String()) + " ")
	}
	if str := e.DeviceAddress.String(); str != "<nil>" {
		b.WriteString("dvc=" + str + " ")
	}
	if e.DeviceHostName != "" {
		b.WriteString("dcvhost=" + escapeExtensionField(e.DeviceHostName) + " ")
	}
	if len(e.DeviceMacAddress) != 0 {
		b.WriteString("dvcmac=" + e.DeviceMacAddress.String() + " ")
	}
	if e.DeviceProcessId != nil {
		b.WriteString("dvcpid=" + strconv.FormatUint(uint64(*e.DeviceProcessId), 10) + " ")
	}
	if !e.DeviceReceiptTime.IsZero() {
		b.WriteString("rt=" + strconv.FormatInt(e.DeviceReceiptTime.UnixMilli(), 10) + " ")
	}
	// TODO add custom mapped fields
	return b.String()
}

func (e Extensions) marshalDestinationFields() string {
	b := strings.Builder{}
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
	if e.DestinationHostName != "" {
		b.WriteString("dhost=" + escapeExtensionField(e.DestinationHostName) + " ")
	}
	if len(e.DestinationMacAddress) != 0 {
		b.WriteString("dmac=" + e.DestinationMacAddress.String() + " ")
	}
	if e.DestinationNtDomain != "" {
		b.WriteString("dntdom=" + escapeExtensionField(e.DestinationNtDomain) + " ")
	}
	if e.DestinationProcessId != nil {
		b.WriteString("dpid=" + strconv.FormatUint(uint64(*e.DestinationProcessId), 10) + " ")
	}
	if e.DestinationUserPrivileges != "" {
		b.WriteString("dpriv=" + escapeExtensionField(e.DestinationUserPrivileges) + " ")
	}
	if e.DestinationProcessName != "" {
		b.WriteString("dproc=" + escapeExtensionField(e.DestinationProcessName) + " ")
	}
	if e.DestinationPort != nil {
		b.WriteString("dpt=" + strconv.FormatUint(uint64(*e.DestinationPort), 10) + " ")
	}
	if str := e.DestinationAddress.String(); str != "<nil>" {
		b.WriteString("dst=" + str + " ")
	}
	if e.DestinationUserId != "" {
		b.WriteString("duid=" + escapeExtensionField(e.DestinationUserId) + " ")
	}
	// TODO add destination marshaling

	// TODO add custom mapped fields

	return b.String()
}

func (e Extensions) marshalFileFields() string {
	b := strings.Builder{}

	if !e.FileCreateTime.IsZero() {
		b.WriteString("fileCreateTime=" + strconv.FormatInt(e.FileCreateTime.UnixMilli(), 10) + " ")
	}
	if e.FileHash != "" {
		b.WriteString("fileHash=" + escapeExtensionField(e.FileHash) + " ")
	}
	if e.FileId != "" {
		b.WriteString("fileId=" + escapeExtensionField(e.FileId) + " ")
	}
	if !e.FileModificationTime.IsZero() {
		b.WriteString("fileModificationTime=" + strconv.FormatInt(e.FileModificationTime.UnixMilli(), 10) + " ")
	}
	if e.FilePath != "" {
		b.WriteString("filePath=" + escapeExtensionField(e.FilePath) + " ")
	}
	if e.FilePermission != "" {
		b.WriteString("filePermission=" + escapeExtensionField(e.FilePermission) + " ")
	}
	if e.FileType != "" {
		b.WriteString("fileType=" + escapeExtensionField(e.FileType) + " ")
	}
	if e.FileName != "" {
		b.WriteString("fname=" + escapeExtensionField(e.FileName) + " ")
	}
	if e.FileSize != nil {
		b.WriteString("fsize=" + strconv.FormatUint(uint64(*e.FileSize), 10) + " ")
	}
	if !e.OldFileCreateTime.IsZero() {
		b.WriteString("oldFileCreateTime=" + strconv.FormatInt(e.OldFileCreateTime.UnixMilli(), 10) + " ")
	}
	if e.OldFileHash != "" {
		b.WriteString("oldFileHash=" + escapeExtensionField(e.OldFileHash) + " ")
	}
	if e.OldFileId != "" {
		b.WriteString("oldFileId=" + escapeExtensionField(e.OldFileId) + " ")
	}
	if !e.OldFileModificationTime.IsZero() {
		b.WriteString("oldFileModificationTime=" + strconv.FormatInt(e.OldFileModificationTime.UnixMilli(), 10) + " ")
	}
	if e.OldFileName != "" {
		b.WriteString("oldFileName=" + escapeExtensionField(e.OldFileName) + " ")
	}
	if e.OldFilePath != "" {
		b.WriteString("oldFilePath=" + escapeExtensionField(e.OldFilePath) + " ")
	}
	if e.OldFilePermission != "" {
		b.WriteString("oldFilePermission=" + escapeExtensionField(e.OldFilePermission) + " ")
	}
	if e.OldFileType != "" {
		b.WriteString("oldFileType=" + escapeExtensionField(e.OldFileType) + " ")
	}
	if e.OldFileSize != nil {
		b.WriteString("oldFileSize=" + strconv.FormatUint(uint64(*e.OldFileSize), 10) + " ")
	}

	return b.String()
}

func (e Extensions) marshalHttpFields() string {
	b := strings.Builder{}
	if (url.URL{}) != e.RequestUrl {
		b.WriteString("request=" + escapeExtensionField(e.RequestUrl.String()) + " ")
	}
	if e.RequestClientApplication != "" {
		b.WriteString("requestClientApplication=" + escapeExtensionField(e.RequestClientApplication) + " ")
	}
	if e.RequestContext != "" {
		b.WriteString("requestContext=" + escapeExtensionField(e.RequestContext) + " ")
	}
	if e.RequestCookies != "" {
		b.WriteString("requestCookies=" + escapeExtensionField(e.RequestCookies) + " ")
	}
	if e.RequestMethod != "" {
		b.WriteString("requestMethod=" + escapeExtensionField(e.RequestMethod) + " ")
	}

	return b.String()
}

func (e Extensions) marshalSourceFields() string {
	b := strings.Builder{}
	// TODO implement

	return b.String()
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
