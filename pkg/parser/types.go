package parser

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// Primitive types are represented by built in types. TODO for now?
// See spec-add-type.jq for the mapping.

func ParseBoolAV(s string) (bool, error) {
	switch s {
	case "A":
		return true, nil
	case "V":
		return false, nil
	}
	return false, fmt.Errorf("should be one of AV but got: %s", s)
}

func PrintBoolAV(b bool) string {
	if b {
		return "A"
	}
	return "V"
}

func ParseInt(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func PrintInt(i int64) string {
	return strconv.FormatInt(i, 10)
}

func ParseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func PrintFloat(f float64) string {
	fi, ff := math.Modf(f)
	if ff == 0.0 {
		return fmt.Sprintf("%f.0", fi)
	}

	return fmt.Sprintf("%g", f)
}

func ParseString(s string) (string, error) {
	return s, nil
}

func PrintString(s string) string {
	return s
}

// TODO move together to types?
// type FixQuality int64
// // String pretty prints a FixQuality
// func (fq FixQuality) String() {
// 	switch fq {
// 	case 0:
// 		return "NA"
// 	case 1:
// 		return "GPS"
// 	case 2:
// 		return "DGPS"
// 	case 3:
// 		return "PPS"
// 	case 4:
// 		return "RTK"
// 	case 5:
// 		return "FRTK"
// 	case 6:
// 		return "EST" //TODO ?
// 	case 7:
// 		return "M" // TODO ?
// 	case 8:
// 		return "SIM" // TODO ?
// 	}
// 	return "Invalid" //TODO?
// }

func ParseFixQuality(s string) (int64, error) {
	//  0=not available, 1=GPS, 2=Differential GPS, 3=PPS, 4=RealTimeKinematic, 5=FloatRTK, 6=Estimated (dead reckoning), 7=Manual input, 8=Simulation mode
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if v < 0 || v > 8 {
		return 0, fmt.Errorf("should be 0..8 but got: %d", v)
	}
	return v, nil
}

func PrintFixQuality(i int64) string {
	return fmt.Sprint(i)
}

// Date type
type Date struct {
	Valid bool
	DD    int
	MM    int
	YY    int
}

// String pretty prints a Date.
func (d Date) String() string {
	return fmt.Sprintf("%02d/%02d/%02d", d.DD, d.MM, d.YY)
}

// ParseDate parses a string in ddmmyy format and returns a Date.
func ParseDate(ddmmyy string) (Date, error) {
	if ddmmyy == "" {
		return Date{}, nil
	}
	if len(ddmmyy) != 6 {
		return Date{}, fmt.Errorf("should be ddmmyy format but got: %s", ddmmyy)
	}
	dd, err := strconv.Atoi(ddmmyy[0:2])
	if err != nil {
		return Date{}, fmt.Errorf("day in %s: %w", ddmmyy, err)
	}
	mm, err := strconv.Atoi(ddmmyy[2:4])
	if err != nil {
		return Date{}, fmt.Errorf("month in %s: %w", ddmmyy, err)
	}
	yy, err := strconv.Atoi(ddmmyy[4:6])
	if err != nil {
		return Date{}, fmt.Errorf("year in %s: %w", ddmmyy, err)
	}
	return Date{true, dd, mm, yy}, nil
}

// PrintDate prints a Date in ddmmyy format
func PrintDate(d Date) string {
	return fmt.Sprintf("%02d%02d%02d", d.DD, d.MM, d.YY)
}

type Time struct {
	Valid       bool
	Hour        int
	Minute      int
	Second      int
	Millisecond int
}

// String representation of Time
func (t Time) String() string {
	var m string
	if !t.Valid {
		m = "(invalid)"
	}
	return fmt.Sprintf("%02d:%02d:%02d.%03d%s", t.Hour, t.Minute, t.Second, t.Millisecond, m)
}

// timeRe is used to validate time strings
var timeRe = regexp.MustCompile(`^\d{6}(\.\d*)?$`)

// ParseTime parses wall clock time.
// e.g. hhmmss.ssss
// An empty time string will result in an invalid time.
func ParseTime(s string) (Time, error) {
	if s == "" {
		return Time{}, nil
	}
	if !timeRe.MatchString(s) {
		return Time{}, fmt.Errorf("should be hhmmss.ss format but got: %s", s)
	}
	hour, _ := strconv.Atoi(s[:2]) //TODO err ignored?!
	minute, _ := strconv.Atoi(s[2:4])
	second, _ := strconv.ParseFloat(s[4:], 64)
	whole, frac := math.Modf(second)
	return Time{true, hour, minute, int(whole), int(math.Round(frac * 1000))}, nil
}

// PrintTime prints a Time in hhmmss.ss format.
func PrintTime(t Time) string {
	return fmt.Sprintf("%02d%02d%02d.%03d", t.Hour, t.Minute, t.Second, t.Millisecond)
}

type Coordinate struct {
	Val  float64
	Area string
}

func ParseCoordinate(val, area string) (Coordinate, error) {
	v, err := strconv.ParseFloat(val, 64)
	//TODO suport other formats
	if err != nil {
		return Coordinate{}, err
	}
	//TODO validate
	return Coordinate{v, area}, nil
}

func MustParseCoordinate(val, area string) Coordinate {
	r, err := ParseCoordinate(val, area)
	if err != nil {
		panic(err)
	}
	return r
}

// PrintCoordinate prints a Coordinate in val,area format.
func PrintCoordinate(c Coordinate) string {
	return fmt.Sprintf("%0.4f,%s", c.Val, c.Area)
}

type Distance struct {
	Val float64
	// Unit of distance in;
	//  f - Feet (0.3048m)
	//  F - Fathom (1.82m)
	//	K - Kilometer
	//  M - Meter
	//  N - Nautic Mile (1852m)
	//  S - Statue Mile (1609.344m)
	Unit string
}

func ParseDistance(val, unit string) (Distance, error) {
	v, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return Distance{}, err
	}
	u := "fFKMNS"

	if !strings.Contains(u, unit) {
		return Distance{}, fmt.Errorf("unit should be one of %s but got: %s", u, unit)
	}
	return Distance{v, unit}, nil
}

// PrintDistance prints a Distance in val,unit format.
func PrintDistance(d Distance) string {
	return PrintFloat(d.Val) + "," + d.Unit
}
