package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAAM(t *testing.T) {
	var tests = []struct {
		name string
		raw  string
		err  string
		msg  AAM
	}{
		{
			name: "good sentence",
			raw:  "$GPAAM,A,A,0.10,N,WPTNME*32",
			msg: AAM{
				ArrivalCircleEntered: true,
				PerpendicularPassed:  true,
				ArrivalCircleRadius:  Distance{0.1, "N"},
				//TODO remove - ArrivalCircleRadiusUnit:    DistanceUnitNauticalMile,
				DestinationWaypointID: "WPTNME",
			},
		},
		{
			name: "invalid nmea: StatusArrivalCircleEntered",
			raw:  "$GPAAM,x,A,0.10,N,WPTNME*0B",
			err:  "AAM: ArrivalCircleEntered: should be one of AV but got: x",
		},
		{
			name: "invalid nmea: StatusPerpendicularPassed",
			raw:  "$GPAAM,A,x,0.10,N,WPTNME*0B",
			err:  "AAM: PerpendicularPassed: should be one of AV but got: x",
		},
		{
			name: "invalid nmea: DistanceUnitNauticalMile",
			raw:  "$GPAAM,A,A,0.10,x,WPTNME*04",
			err:  "AAM: ArrivalCircleRadius: unit should be one of fFKMNS but got: x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				aam := m.(AAM)
				aam.Base = Base{}
				assert.Equal(t, tt.msg, aam)
			}
		})
	}
}

func TestParseGGA(t *testing.T) {
	var tests = []struct {
		name string
		raw  string
		err  string
		msg  GGA
	}{
		{
			name: "good sentence",
			raw:  "$GNGGA,203415.000,6325.6138,N,01021.4290,E,1,8,2.42,72.5,M,41.5,M,,*7C",
			msg: GGA{
				Time: Time{
					Valid:       true,
					Hour:        20,
					Minute:      34,
					Second:      15,
					Millisecond: 0,
				},
				Latitude:      MustParseCoordinate("6325.6138", "N"),
				Longitude:     MustParseCoordinate("01021.4290", "E"),
				FixQuality:    1,
				NumSatellites: 8,
				HDOP:          2.42,
				Altitude:      Distance{72.5, "M"},
				Separation:    Distance{41.5, "M"},
				DGPSAge:       "",
				DGPSId:        "",
			},
		},
		{
			name: "bad latitude",
			raw:  "$GNGGA,034225.077,x,S,15124.5567,E,1,03,9.7,-25.0,M,21.0,M,,0000*1D",
			err:  "GGA: Latitude: strconv.ParseFloat: parsing \"x\": invalid syntax",
		},
		{
			name: "bad longitude",
			raw:  "$GNGGA,034225.077,3356.4650,S,x,E,1,03,9.7,-25.0,M,21.0,M,,0000*2B",
			err:  "GGA: Longitude: strconv.ParseFloat: parsing \"x\": invalid syntax",
		},
		{
			name: "bad fix quality",
			raw:  "$GNGGA,034225.077,3356.4650,S,15124.5567,E,99,03,9.7,-25.0,M,21.0,M,,0000*7E",
			err:  "GGA: FixQuality: should be 0..8 but got: 99",
		},
		{
			name: "GP talker, good sentence",
			raw:  "$GPGGA,034225.077,3356.4650,S,15124.5567,E,1,03,9.7,-25.0,M,21.0,M,,0000*51",
			msg: GGA{
				Time:          Time{true, 3, 42, 25, 77},
				Latitude:      MustParseCoordinate("3356.4650", "S"),
				Longitude:     MustParseCoordinate("15124.5567", "E"),
				FixQuality:    1,
				NumSatellites: 03,
				HDOP:          9.7,
				Altitude:      Distance{-25.0, "M"},
				Separation:    Distance{21.0, "M"},
				DGPSAge:       "",
				DGPSId:        "0000",
			},
		},
		{
			name: "GP talker, bad latitude",
			raw:  "$GPGGA,034225.077,x,S,15124.5567,E,1,03,9.7,-25.0,M,21.0,M,,0000*03",
			err:  "GGA: Latitude: strconv.ParseFloat: parsing \"x\": invalid syntax",
		},
		{
			name: "GP talker, bad longitude",
			raw:  "$GPGGA,034225.077,3356.4650,S,x,E,1,03,9.7,-25.0,M,21.0,M,,0000*35",
			err:  "GGA: Longitude: strconv.ParseFloat: parsing \"x\": invalid syntax",
		},
		{
			name: "GP talker, bad fix quality",
			raw:  "$GPGGA,034225.077,3356.4650,S,15124.5567,E,99,03,9.7,-25.0,M,21.0,M,,0000*60",
			err:  "GGA: FixQuality: should be 0..8 but got: 99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := Parse(tt.raw)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				gga := m.(GGA)
				gga.Base = Base{}
				assert.Equal(t, tt.msg, gga)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	var tests = []struct {
		name string
		raw  string
		msg  Sentence
		err  string
	}{
		{
			name: "AAM sentence",
			raw:  "$GPAAM,A,A,0.1,N,WPTNME*02",
			msg: AAM{
				Base:                 Base{Talker: "GP", Type: "AAM"},
				ArrivalCircleEntered: true,
				PerpendicularPassed:  true,
				ArrivalCircleRadius:  Distance{0.1, "N"},
				//TODO remove - ArrivalCircleRadiusUnit:    DistanceUnitNauticalMile,
				DestinationWaypointID: "WPTNME",
			},
		},
		{
			name: "GGA sentence",
			raw:  "$GNGGA,203415.000,6325.6138,N,1021.4290,E,1,8,2.42,72.5,M,41.5,M,,*4C",
			msg: GGA{
				Base: Base{Talker: "GN", Type: "GGA"},
				Time: Time{
					Valid:       true,
					Hour:        20,
					Minute:      34,
					Second:      15,
					Millisecond: 0,
				},
				Latitude:      MustParseCoordinate("6325.6138", "N"),
				Longitude:     MustParseCoordinate("01021.4290", "E"),
				FixQuality:    1,
				NumSatellites: 8,
				HDOP:          2.42,
				Altitude:      Distance{72.5, "M"},
				Separation:    Distance{41.5, "M"},
				DGPSAge:       "",
				DGPSId:        "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := Print(tt.msg)
			if tt.err != "" {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.raw, s)
			}
		})
	}
}
