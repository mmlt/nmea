package parser

import (
	"fmt"
	"strings"
)

const (
	// SentenceStart is the token to indicate the start of a sentence.
	SentenceStart = "$"

	// SentenceStartEncapsulated is the token to indicate the start of encapsulated data.
	SentenceStartEncapsulated = "!"

	// FieldSep is the token to delimit fields of a sentence.
	FieldSep = ","

	// ChecksumSep is the token to delimit the checksum of a sentence.
	ChecksumSep = "*"
)

// Sentence interface for all NMEA sentences
type Sentence interface {
	Prefix() string
	DataType() string //TODO rename to Type
	TalkerID() string //TODO rename to Talker
}

// UnkownTypeError is used when a sentence type is encountered that is not implemented by the software.
type UnkownTypeError struct {
	Type string
}

func (e UnkownTypeError) Error() string {
	return fmt.Sprintf("unknown sentence type: %s", e.Type)
}

// Parse parses a NME0183 formmated string and returns a Sentence.
func Parse(s string) (Sentence, error) {
	var err error

	b, err := stringToBase(s)
	if err != nil {
		return nil, err
	}

	var sentence Sentence

	strt := b.Raw[0:1]
	switch strt {
	case "$":
		// MTK message types share the same format
		// so we return the same struct for all types.
		// switch s.Talker {
		// case TypeMTK:
		// 	return newMTK(s)
		// }
		p := parsers[b.Type]
		if p == nil {
			return nil, UnkownTypeError{Type: b.Type}
		}
		sentence, err = p(b)
	case "!":
		// AIVDM/AIVDO encapsulated data
		// 	switch s.Type {
		// 	case TypeVDM, TypeVDO:
		// 		return newVDMVDO(s)
		// 	}
		return nil, fmt.Errorf("AIVDM/AIVDO NIY")
	default:
		return nil, fmt.Errorf("sentence should start with $ or ! but got: %s", strt)
	}

	return sentence, err
}

// stringToBase parses a raw message into it's fields
func stringToBase(raw string) (Base, error) {
	raw = strings.TrimSpace(raw)
	tagBlockParts := strings.SplitN(raw, `\`, 3)

	var (
		tagBlock TagBlock
		err      error
	)
	if len(tagBlockParts) == 3 {
		tags := tagBlockParts[1]
		raw = tagBlockParts[2]
		tagBlock, err = parseTagBlock(tags)
		if err != nil {
			return Base{}, err
		}
	}

	startIndex := strings.IndexAny(raw, SentenceStart+SentenceStartEncapsulated)
	if startIndex != 0 {
		return Base{}, fmt.Errorf("nmea: sentence does not start with a '$' or '!'")
	}
	sumSepIndex := strings.Index(raw, ChecksumSep)
	if sumSepIndex == -1 {
		return Base{}, fmt.Errorf("nmea: sentence does not contain checksum separator")
	}
	var (
		fieldsRaw   = raw[startIndex+1 : sumSepIndex]
		fields      = strings.Split(fieldsRaw, FieldSep)
		checksumRaw = strings.ToUpper(raw[sumSepIndex+1:])
		checksum    = Checksum(fieldsRaw)
	)
	// Validate the checksum
	if checksum != checksumRaw {
		return Base{}, fmt.Errorf(
			"nmea: sentence checksum mismatch [%s != %s]", checksum, checksumRaw)
	}
	talker, typ := parsePrefix(fields[0])
	return Base{
		Talker:   talker,
		Type:     typ,
		Fields:   fields[1:],
		Checksum: checksumRaw,
		Raw:      raw,
		TagBlock: tagBlock,
	}, nil
}

// parsePrefix takes the first field and splits it into a talker id and data type.
func parsePrefix(s string) (string, string) {
	if strings.HasPrefix(s, "PMTK") {
		return "PMTK", s[4:]
	}
	if strings.HasPrefix(s, "P") {
		return "P", s[1:]
	}
	if len(s) < 2 {
		return s, ""
	}
	return s[:2], s[2:]
}

// Checksum xor all the bytes in a string an return it
// as an uppercase hex string
func Checksum(s string) string { //TODO make private
	var checksum uint8
	for i := 0; i < len(s); i++ {
		checksum ^= s[i]
	}
	return fmt.Sprintf("%02X", checksum)
}
