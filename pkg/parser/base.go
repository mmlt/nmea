package parser

// Base represents the NMEA sentence as textual fields.
type Base struct {
	Talker   string   // The talker id (e.g GP)
	Type     string   // The data type (e.g GSA)
	Fields   []string // Array of fields
	Checksum string   // The Checksum
	Raw      string   // The raw NMEA sentence received //TODO Needed for troubleshooting?
	TagBlock TagBlock // NMEA tagblock
}

// Prefix returns the talker and type of message
func (b Base) Prefix() string {
	return b.Talker + b.Type
}

// DataType returns the type of the message
func (b Base) DataType() string {
	return b.Type
}

// TalkerID returns the talker of the message
func (b Base) TalkerID() string {
	return b.Talker
}

// Field types and parsers
