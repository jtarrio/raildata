package raildata_test

import (
	"testing"

	"github.com/jtarrio/raildata"
	"github.com/stretchr/testify/assert"
)

func TestFindStationWithCode(t *testing.T) {
	expected := &raildata.Station{
		Code:      "PJ",
		Name:      "Princeton Junction",
		ShortName: "Princeton Jct.",
	}
	station, found := raildata.FindStation().WithCode("PJ").Search()
	assert.True(t, found)
	assert.Equal(t, expected, station)

	_, found = raildata.FindStation().WithCode("XX").Search()
	assert.False(t, found)
}

func TestFindStationWithName(t *testing.T) {
	// Find by name
	expected := &raildata.Station{
		Code:      "BA",
		Name:      "BWI Thurgood Marshall Airport",
		ShortName: "BWI Airport",
	}
	station, found := raildata.FindStation().WithName("bwi thurgood marshall airport").Search()
	assert.True(t, found)
	assert.Equal(t, expected, station)

	// Find by name (does not exist)
	_, found = raildata.FindStation().WithName("12345678901234567890").Search()
	assert.False(t, found)

	// Find by short name
	expected = &raildata.Station{
		Code:      "BY",
		Name:      "Berkeley Heights",
		ShortName: "Berkeley Hts",
	}
	station, found = raildata.FindStation().WithName("berkeley hts").Search()
	assert.True(t, found)
	assert.Equal(t, expected, station)

	// Find by alias
	expected = &raildata.Station{
		Code:      "UF",
		Name:      "Hohokus",
		ShortName: "Hohokus"}
	station, found = raildata.FindStation().WithName("HO-HO-KUS").Search()
	assert.True(t, found)
	assert.Equal(t, expected, station)
}

func TestFindStationOrSynthesize(t *testing.T) {
	expected := raildata.Station{
		Code:      "XY",
		Name:      "Unknown XY",
		ShortName: "Unknown XY",
	}
	station := raildata.FindStation().WithCode("XY").SearchOrSynthesize()
	assert.Equal(t, expected, station)

	expected = raildata.Station{
		Code:      "XY",
		Name:      "12345678901234567890",
		ShortName: "12345678901234",
	}
	station = raildata.FindStation().WithCode("XY").WithName("12345678901234567890").SearchOrSynthesize()
	assert.Equal(t, expected, station)

	expected = raildata.Station{
		Code:      "XX",
		Name:      "12345678901234567890",
		ShortName: "12345678901234",
	}
	station = raildata.FindStation().WithName("12345678901234567890").SearchOrSynthesize()
	assert.Equal(t, expected, station)
}

func TestFindLineWithCode(t *testing.T) {
	line, found := raildata.FindLine().WithCode("NE").Search()
	assert.True(t, found)
	assert.Equal(t, &raildata.Lines[6], line)

	_, found = raildata.FindLine().WithCode("XX").Search()
	assert.False(t, found)
}

func TestFindLineWithName(t *testing.T) {
	// Find by name
	line, found := raildata.FindLine().WithName("atlantic city line").Search()
	assert.True(t, found)
	assert.Equal(t, &raildata.Lines[0], line)

	// Find by name (does not exist)
	_, found = raildata.FindLine().WithName("12345678901234567890").Search()
	assert.False(t, found)

	// Find by abbreviation
	line, found = raildata.FindLine().WithName("mobo").Search()
	assert.True(t, found)
	assert.Equal(t, &raildata.Lines[1], line)

	// Find by alternative abbreviation
	line, found = raildata.FindLine().WithName("mne").Search()
	assert.True(t, found)
	assert.Equal(t, &raildata.Lines[4], line)

	// Find by alternative name
	line, found = raildata.FindLine().WithName("NORTHEAST CORRDR").Search()
	assert.True(t, found)
	assert.Equal(t, &raildata.Lines[6], line)
}

func TestFindLineOrSynthesize(t *testing.T) {
	expected := raildata.Line{
		Code:         "XY",
		Name:         "Unknown XY",
		Abbreviation: "XXXY",
	}
	line := raildata.FindLine().WithCode("XY").SearchOrSynthesize()
	assert.Equal(t, expected, line)

	expected = raildata.Line{
		Code:         "XY",
		Name:         "12345678901234567890",
		Abbreviation: "XXXY",
	}
	line = raildata.FindLine().WithCode("XY").WithName("12345678901234567890").SearchOrSynthesize()
	assert.Equal(t, expected, line)

	expected = raildata.Line{
		Code:         "XX",
		Name:         "12345678901234567890",
		Abbreviation: "XXXX",
	}
	line = raildata.FindLine().WithName("12345678901234567890").SearchOrSynthesize()
	assert.Equal(t, expected, line)
}

func TestTranslateTrackNumber(t *testing.T) {
	assert.Equal(t, "1", raildata.TranslateTrackNumber("1", "UP"))
	assert.Equal(t, "2", raildata.TranslateTrackNumber("2", "UP"))
	assert.Equal(t, "1", raildata.TranslateTrackNumber("Single", "UV"))
	assert.Equal(t, "2", raildata.TranslateTrackNumber("B", "UV"))
	assert.Equal(t, "C", raildata.TranslateTrackNumber("C", "UV"))
}
