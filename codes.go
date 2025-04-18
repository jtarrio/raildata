package raildata

import (
	"fmt"
	"strings"
)

// TrainIdPrefix contains the code and the meaning of a special train number prefix.
type TrainIdPrefix struct {
	Prefix      string
	Description string
}

// TrainIdPrefixes contains the list of special train number prefixes.
var TrainIdPrefixes = []TrainIdPrefix{
	{Prefix: "A", Description: "Amtrak Train"},
	{Prefix: "S", Description: "Septa Train"},
	{Prefix: "X", Description: "Non-Revenue train - Does not accept passengers"},
}

// SpecialTrack contains API-to-real-world track number translation data for certain stations.
type SpecialTrack struct {
	Id          string
	StationCode StationCode
	Translation string
}

// SpecialTracks contains a list of track numbers, as provided by the API, that are different in the real world.
// For example, track "Single" in station "UV" is seen as track "1" in the real world.
var SpecialTracks = []SpecialTrack{
	{Id: "Single", StationCode: "ON", Translation: "1"},
	{Id: "2", StationCode: "MP", Translation: "1"},
	{Id: "B", StationCode: "UV", Translation: "2"},
	{Id: "Single", StationCode: "UV", Translation: "1"},
	{Id: "0", StationCode: "NA", Translation: "A"},
	{Id: "4", StationCode: "TS", Translation: "E"},
	{Id: "2", StationCode: "TS", Translation: "F"},
	{Id: "3", StationCode: "TS", Translation: "H"},
	{Id: "1", StationCode: "TS", Translation: "G"},
	{Id: "Single", StationCode: "ST", Translation: "S"},
}

// TranslateTrackNumber takes a track id and a station code and returns the track name used in the real world.
func TranslateTrackNumber(trackId string, stationCode StationCode) string {
	for i := range SpecialTracks {
		if trackId == SpecialTracks[i].Id && stationCode == SpecialTracks[i].StationCode {
			return SpecialTracks[i].Translation
		}
	}
	return trackId
}

// StationPosition contains data about a station's position in relation to a route.
type StationPosition struct {
	Code        string
	Description string
}

// StationPositions contains a list of valid station positions.
var StationPositions = []StationPosition{
	{Code: "0", Description: "First station"},
	{Code: "1", Description: "Intermediate station"},
	{Code: "2", Description: "Final station"},
}

// GetStationPosition returns the [StationPosition] for the given code, or a made-up object for an unknown code.
func GetStationPosition(code string) StationPosition {
	for i := range StationPositions {
		if StationPositions[i].Code == code {
			return StationPositions[i]
		}
	}
	return StationPosition{Code: code, Description: fmt.Sprintf("Unknown station position \"%s\"", code)}
}

// StopCode contains information about a train's behavior at a station
type StopCode struct {
	Code        string
	Description string
}

// StopCodes contains a list of valid stop codes.
var StopCodes = []StopCode{
	{Code: "A", Description: "Arrival time"},
	{Code: "S", Description: "Normal Stop"},
	{Code: "S*", Description: "Normal stop. May leave up to 3 minutes early"},
	{Code: "LV", Description: "Leaves 1 minute after scheduled time"},
	{Code: "L", Description: "Train can leave before scheduled departure. Will hold for connections"},
	{Code: "H", Description: "Will hold for connection unless authorize by dispatcher"},
	{Code: "D", Description: "Stop to discharge passengers only. May leave ahead of schedule"},
	{Code: "R", Description: "Stop to receive passengers only"},
	{Code: "R*", Description: "Stop to receive passengers only. May leave 3 minutes early"},
	{Code: "E", Description: "Employee stop. May leave ahead of schedule"},
}

// GetStopCode returns the [StopCode] for a given code, or a made-up object for an unknown code.
func GetStopCode(code string) StopCode {
	for i := range StopCodes {
		if StopCodes[i].Code == code {
			return StopCodes[i]
		}
	}
	return StopCode{Code: code, Description: fmt.Sprintf("Unknown stop code \"%s\"", code)}
}

// StationCode is a 2-letter identifier for a station.
type StationCode string

// Station contains information about a station.
type Station struct {
	// Code contains the station's 2-letter code.
	Code StationCode
	// Name contains the station's full name.
	Name string
	// ShortName contains a shorter version of the station's name.
	ShortName string
}

var Stations = []Station{
	{Code: "AM", Name: "Aberdeen-Matawan", ShortName: "Matawan"},
	{Code: "AB", Name: "Absecon", ShortName: "Absecon"},
	{Code: "AZ", Name: "Allendale", ShortName: "Allendale"},
	{Code: "AH", Name: "Allenhurst", ShortName: "Allenhurst"},
	{Code: "AS", Name: "Anderson Street", ShortName: "Anderson St."},
	{Code: "AN", Name: "Annandale", ShortName: "Annandale"},
	{Code: "AP", Name: "Asbury Park", ShortName: "Asbury Park"},
	{Code: "AO", Name: "Atco", ShortName: "Atco"},
	{Code: "AC", Name: "Atlantic City Rail Terminal", ShortName: "Atlantic City"},
	{Code: "AV", Name: "Avenel", ShortName: "Avenel"},
	{Code: "BA", Name: "BWI Thurgood Marshall Airport", ShortName: "BWI Airport"},
	{Code: "BL", Name: "Baltimore Station", ShortName: "Baltimore"},
	{Code: "BI", Name: "Basking Ridge", ShortName: "Basking Ridge"},
	{Code: "BH", Name: "Bay Head", ShortName: "Bay Head"},
	{Code: "MC", Name: "Bay Street", ShortName: "Bay Street"},
	{Code: "BS", Name: "Belmar", ShortName: "Belmar"},
	{Code: "BY", Name: "Berkeley Heights", ShortName: "Berkeley Hts"},
	{Code: "BV", Name: "Bernardsville", ShortName: "Bernardsville"},
	{Code: "BM", Name: "Bloomfield", ShortName: "Bloomfield"},
	{Code: "BN", Name: "Boonton", ShortName: "Boonton"},
	{Code: "BK", Name: "Bound Brook", ShortName: "Bound Brook"},
	{Code: "BB", Name: "Bradley Beach", ShortName: "Bradley Beach"},
	{Code: "BU", Name: "Brick Church", ShortName: "Brick Church"},
	{Code: "BW", Name: "Bridgewater", ShortName: "Bridgewater"},
	{Code: "BF", Name: "Broadway Fair Lawn", ShortName: "Broadway-Fl"},
	{Code: "CB", Name: "Campbell Hall", ShortName: "Campbell Hall"},
	{Code: "CM", Name: "Chatham", ShortName: "Chatham"},
	{Code: "CY", Name: "Cherry Hill", ShortName: "Cherry Hill"},
	{Code: "IF", Name: "Clifton", ShortName: "Clifton"},
	{Code: "CN", Name: "Convent Station", ShortName: "Convent Stn"},
	{Code: "XC", Name: "Cranford", ShortName: "Cranford"},
	{Code: "DL", Name: "Delawanna", ShortName: "Delawanna"},
	{Code: "DV", Name: "Denville", ShortName: "Denville"},
	{Code: "DO", Name: "Dover", ShortName: "Dover"},
	{Code: "DN", Name: "Dunellen", ShortName: "Dunellen"},
	{Code: "EO", Name: "East Orange", ShortName: "East Orange"},
	{Code: "ED", Name: "Edison", ShortName: "Edison"},
	{Code: "EH", Name: "Egg Harbor City", ShortName: "Egg Harbor"},
	{Code: "EL", Name: "Elberon", ShortName: "Elberon"},
	{Code: "EZ", Name: "Elizabeth", ShortName: "Elizabeth"},
	{Code: "EN", Name: "Emerson", ShortName: "Emerson"},
	{Code: "EX", Name: "Essex Street", ShortName: "Essex Street"},
	{Code: "FW", Name: "Fanwood", ShortName: "Fanwood"},
	{Code: "FH", Name: "Far Hills", ShortName: "Far Hills"},
	{Code: "FE", Name: "Finderne", ShortName: "Finderne"},
	{Code: "GD", Name: "Garfield", ShortName: "Garfield"},
	{Code: "GW", Name: "Garwood", ShortName: "Garwood"},
	{Code: "GI", Name: "Gillette", ShortName: "Gillette"},
	{Code: "GL", Name: "Gladstone", ShortName: "Gladstone"},
	{Code: "GG", Name: "Glen Ridge", ShortName: "Glen Ridge"},
	{Code: "GK", Name: "Glen Rock Boro Hall", ShortName: "Glen Rock Boro"},
	{Code: "RS", Name: "Glen Rock Main Line", ShortName: "Glen Rock Main"},
	{Code: "GA", Name: "Great Notch", ShortName: "Great Notch"},
	{Code: "HQ", Name: "Hackettstown", ShortName: "Hackettstown"},
	{Code: "HL", Name: "Hamilton", ShortName: "Hamilton"},
	{Code: "HN", Name: "Hammonton", ShortName: "Hammonton"},
	{Code: "RM", Name: "Harriman", ShortName: "Harriman"},
	{Code: "HW", Name: "Hawthorne", ShortName: "Hawthorne"},
	{Code: "HZ", Name: "Hazlet", ShortName: "Hazlet"},
	{Code: "HG", Name: "High Bridge", ShortName: "High Bridge"},
	{Code: "HI", Name: "Highland Avenue", ShortName: "Highland Ave."},
	{Code: "HD", Name: "Hillsdale", ShortName: "Hillsdale"},
	{Code: "HB", Name: "Hoboken", ShortName: "Hoboken"},
	{Code: "UF", Name: "Hohokus", ShortName: "Hohokus"},
	{Code: "JA", Name: "Jersey Avenue", ShortName: "Jersey Ave."},
	{Code: "KG", Name: "Kingsland", ShortName: "Kingsland"},
	{Code: "HP", Name: "Lake Hopatcong", ShortName: "Lake Hopatcong"},
	{Code: "ON", Name: "Lebanon", ShortName: "Lebanon"},
	{Code: "LP", Name: "Lincoln Park", ShortName: "Lincoln Park"},
	{Code: "LI", Name: "Linden", ShortName: "Linden"},
	{Code: "LW", Name: "Lindenwold", ShortName: "Lindenwold"},
	{Code: "FA", Name: "Little Falls", ShortName: "Little Falls"},
	{Code: "LS", Name: "Little Silver", ShortName: "Little Silver"},
	{Code: "LB", Name: "Long Branch", ShortName: "Long Branch"},
	{Code: "LN", Name: "Lyndhurst", ShortName: "Lyndhurst"},
	{Code: "LY", Name: "Lyons", ShortName: "Lyons"},
	{Code: "MA", Name: "Madison", ShortName: "Madison"},
	{Code: "MZ", Name: "Mahwah", ShortName: "Mahwah"},
	{Code: "SQ", Name: "Manasquan", ShortName: "Manasquan"},
	{Code: "MW", Name: "Maplewood", ShortName: "Maplewood"},
	{Code: "XU", Name: "Meadowlands", ShortName: "Meadowlands"},
	{Code: "MP", Name: "Metropark", ShortName: "Metropark"},
	{Code: "MU", Name: "Metuchen", ShortName: "Metuchen"},
	{Code: "MI", Name: "Middletown NJ", ShortName: "Middletown NJ"},
	{Code: "MD", Name: "Middletown NY", ShortName: "Middletown NY"},
	{Code: "MB", Name: "Millburn", ShortName: "Millburn"},
	{Code: "GO", Name: "Millington", ShortName: "Millington"},
	{Code: "MK", Name: "Monmouth Park", ShortName: "Monmouth Park"},
	{Code: "HS", Name: "Montclair Heights", ShortName: "Montclair Hts."},
	{Code: "UV", Name: "Montclair State U", ShortName: "MSU"},
	{Code: "ZM", Name: "Montvale", ShortName: "Montvale"},
	{Code: "MX", Name: "Morris Plains", ShortName: "Morris Plains"},
	{Code: "MR", Name: "Morristown", ShortName: "Morristown"},
	{Code: "HV", Name: "Mount Arlington", ShortName: "Mt. Arlington"},
	{Code: "OL", Name: "Mount Olive", ShortName: "Mount Olive"},
	{Code: "TB", Name: "Mount Tabor", ShortName: "Mount Tabor"},
	{Code: "MS", Name: "Mountain Avenue", ShortName: "Mountain Ave"},
	{Code: "ML", Name: "Mountain Lakes", ShortName: "Mountain Lakes"},
	{Code: "MT", Name: "Mountain Station", ShortName: "Mountain Stn"},
	{Code: "MV", Name: "Mountain View", ShortName: "Mountain View"},
	{Code: "MH", Name: "Murray Hill", ShortName: "Murray Hill"},
	{Code: "NN", Name: "Nanuet", ShortName: "Nanuet"},
	{Code: "NT", Name: "Netcong", ShortName: "Netcong"},
	{Code: "NE", Name: "Netherwood", ShortName: "Netherwood"},
	{Code: "NH", Name: "New Bridge Landing", ShortName: "New Bridge Ldg"},
	{Code: "NB", Name: "New Brunswick", ShortName: "New Brunswick"},
	{Code: "NC", Name: "New Carrollton Station", ShortName: "New Carrollton"},
	{Code: "NV", Name: "New Providence", ShortName: "New Providence"},
	{Code: "NY", Name: "New York Penn Station", ShortName: "New York"},
	{Code: "NA", Name: "Newark Airport", ShortName: "Newark Airport"},
	{Code: "ND", Name: "Newark Broad Street", ShortName: "Newark Broad"},
	{Code: "NP", Name: "Newark Penn Station", ShortName: "Newark Penn"},
	{Code: "OR", Name: "North Branch", ShortName: "North Branch"},
	{Code: "NZ", Name: "North Elizabeth", ShortName: "North Elizab."},
	{Code: "NF", Name: "North Philadelphia", ShortName: ""},
	{Code: "OD", Name: "Oradell", ShortName: "Oradell"},
	{Code: "OG", Name: "Orange", ShortName: "Orange"},
	{Code: "OS", Name: "Otisville", ShortName: "Otisville"},
	{Code: "PV", Name: "Park Ridge", ShortName: "Park Ridge"},
	{Code: "PS", Name: "Passaic", ShortName: "Passaic"},
	{Code: "RN", Name: "Paterson", ShortName: "Paterson"},
	{Code: "PC", Name: "Peapack", ShortName: "Peapack"},
	{Code: "PQ", Name: "Pearl River", ShortName: "Pearl River"},
	{Code: "PN", Name: "Pennsauken", ShortName: "Pennsauken"},
	{Code: "PE", Name: "Perth Amboy", ShortName: "Perth Amboy"},
	{Code: "PH", Name: "Philadelphia", ShortName: "Philadelphia"},
	{Code: "PF", Name: "Plainfield", ShortName: "Plainfield"},
	{Code: "PL", Name: "Plauderville", ShortName: "Plauderville"},
	{Code: "PP", Name: "Point Pleasant Beach", ShortName: "Point Pleasant"},
	{Code: "PO", Name: "Port Jervis", ShortName: "Port Jervis"},
	{Code: "PR", Name: "Princeton", ShortName: "Princeton"},
	{Code: "PJ", Name: "Princeton Junction", ShortName: "Princeton Jct."},
	{Code: "FZ", Name: "Radburn Fair Lawn", ShortName: "Radburn-Fl"},
	{Code: "RH", Name: "Rahway", ShortName: "Rahway"},
	{Code: "RY", Name: "Ramsey Main St", ShortName: "Ramsey"},
	{Code: "17", Name: "Ramsey Route 17", ShortName: "Ramsey Rt 17"},
	{Code: "RA", Name: "Raritan", ShortName: "Raritan"},
	{Code: "RB", Name: "Red Bank", ShortName: "Red Bank"},
	{Code: "RW", Name: "Ridgewood", ShortName: "Ridgewood"},
	{Code: "RG", Name: "River Edge", ShortName: "River Edge"},
	{Code: "RL", Name: "Roselle Park", ShortName: "Roselle Park"},
	{Code: "RF", Name: "Rutherford", ShortName: "Rutherford"},
	{Code: "CW", Name: "Salisbury Mills-Cornwall", ShortName: "Salisbury Mls"},
	{Code: "SC", Name: "Secaucus Concourse", ShortName: ""},
	{Code: "TS", Name: "Secaucus Lower Lvl", ShortName: "Secaucus"},
	{Code: "SE", Name: "Secaucus Upper Lvl", ShortName: "Secaucus"},
	{Code: "RT", Name: "Short Hills", ShortName: "Short Hills"},
	{Code: "XG", Name: "Sloatsburg", ShortName: "Sloatsburg"},
	{Code: "SM", Name: "Somerville", ShortName: "Somerville"},
	{Code: "CH", Name: "South Amboy", ShortName: "South Amboy"},
	{Code: "SO", Name: "South Orange", ShortName: "South Orange"},
	{Code: "LA", Name: "Spring Lake", ShortName: "Spring Lake"},
	{Code: "SV", Name: "Spring Valley", ShortName: "Spring Valley"},
	{Code: "SG", Name: "Stirling", ShortName: "Stirling"},
	{Code: "SF", Name: "Suffern", ShortName: "Suffern"},
	{Code: "ST", Name: "Summit", ShortName: "Summit"},
	{Code: "TE", Name: "Teterboro", ShortName: "Teterboro"},
	{Code: "TO", Name: "Towaco", ShortName: "Towaco"},
	{Code: "TR", Name: "Trenton", ShortName: "Trenton"},
	{Code: "TC", Name: "Tuxedo", ShortName: "Tuxedo"},
	{Code: "US", Name: "Union", ShortName: "Union"},
	{Code: "UM", Name: "Upper Montclair", ShortName: "Upp. Montclair"},
	{Code: "WK", Name: "Waldwick", ShortName: "Waldwick"},
	{Code: "WA", Name: "Walnut Street", ShortName: "Walnut Street"},
	{Code: "WS", Name: "Washington Station", ShortName: "Washington"},
	{Code: "WG", Name: "Watchung Avenue", ShortName: "Watchung Ave."},
	{Code: "WT", Name: "Watsessing Avenue", ShortName: "Watsessing Ave"},
	{Code: "23", Name: "Wayne-Route 23", ShortName: "Wayne Route 23"},
	{Code: "WM", Name: "Wesmont", ShortName: "Wesmont"},
	{Code: "WF", Name: "Westfield", ShortName: "Westfield"},
	{Code: "WW", Name: "Westwood", ShortName: "Westwood"},
	{Code: "WH", Name: "White House", ShortName: "White House"},
	{Code: "WI", Name: "Wilmington Station", ShortName: "Wilmington"},
	{Code: "WR", Name: "Wood Ridge", ShortName: "Wood-Ridge"},
	{Code: "WB", Name: "Woodbridge", ShortName: "Woodbridge"},
	{Code: "WL", Name: "Woodcliff Lake", ShortName: "Woodcliff Lake"},
}

// FindStations returns an object that lets you find a station by code or name.
// If no exact match is found and the name was specified, this function uses fuzzy search to find the closest match.
//
// The [StationFinder.SearchOrSynthesize] method will, if it doesn't find a suitable station, return a synthesized
// [Station] object that uses the provided search data. If no code was specified, "XX" will be used in its place.
// If no name was specified, "Unknown [station code]" will be used in its place, and a shortened version of the name
// will be used in place of the line abbreviation.
func FindStation() StationFinder {
	return finderImpl[Station, StationCode]{
		byCode:        stationsByCode,
		byName:        stationsByName,
		byAbbr:        stationsByShortName,
		list:          Stations,
		getCandidates: func(s *Station) []string { return []string{s.Name, s.ShortName} },
		synthesize: func(code *StationCode, name *string) Station {
			out := Station{}
			if code == nil {
				out.Code = "XX"
			} else {
				out.Code = *code
			}
			if name == nil {
				out.Name = "Unknown " + string(out.Code)
			} else {
				out.Name = *name
			}
			out.ShortName = out.Name[0:min(14, len(out.Name))]
			return out
		},
	}
}

// StationFinder is an object to find stations by code or name.
type StationFinder = Finder[Station, StationCode]

var stationAliases = map[StationCode][]string{
	"AC": {"Atlantic City Terminal"},
	"BA": {"B.W.I. Airport"},
	"MC": {"Bay Street (Montclair)"},
	"BF": {"Broadway"},
	"CN": {"Convent"},
	"ED": {"Edison Station"},
	"FE": {"Manville-Finderne"},
	"GK": {"Glen Rock (Boro Hall)"},
	"RS": {"Glen Rock (Main Line)"},
	"RM": {"Harriman Station"},
	"UF": {"Ho-Ho-Kus"},
	"MI": {"Middletown"},
	"MD": {"Middletown, NY"},
	"UV": {"Montclair State University"},
	"MT": {"Mountain Sta."},
	"ND": {"Newark Broad St.", "Newark Broad St"},
	"NA": {"Newark Int'l Airport", "Newark Airport Railroad Station"},
	"NY": {"Penn Station New York"},
	"PN": {"Pennsauken Transit Center"},
	"PH": {"Philadelphia 30th St.", "30th St. Phl."},
	"PR": {"Princeton Station"},
	"FZ": {"Radburn"},
	"17": {"Route 17 Station", "Ramsey Route 17 Station"},
	"CW": {"Salisbury Mills"},
	"TS": {"Secaucus Junction", "Frank R Lautenberg Secaucus Lower Level"},
	"SE": {"Secaucus Station", "Frank R Lautenberg Secaucus Upper Level"},
	"TR": {"Trenton Station", "Trenton Transit Center"},
	"US": {"Union Station"},
	"WA": {"Walnut Street (Montclair)"},
	"WT": {"Watsessing Avenue (Bloomfield)"},
	"23": {"Wayne/Route 23 Transit Center [RR]"},
}
var stationsByCode = makeMap(Stations, func(s *Station) StationCode { return s.Code })
var stationsByName = makeMmap(Stations, func(s *Station) []string { return append([]string{s.Name}, stationAliases[s.Code]...) })
var stationsByShortName = makeMap(Stations, func(s *Station) string { return s.ShortName })

// LineCode is a 2-letter identifier for a line.
type LineCode string

// Line contains information about a line.
type Line struct {
	// Code contains the line's 2-letter code.
	Code LineCode
	// Name contains the line's full name.
	Name string
	// Abbreviation contains a 3-5 letter abbreviation of the line's name.
	Abbreviation string
	// Color contains the color for this line.
	Color Color
	// OtherAbbrs contains a list with alternative abbreviations for this line.
	OtherAbbrs []string
}

// Lines contains a list of known lines.
var Lines = []Line{
	{Code: "AC", Name: "Atlantic City Line", Abbreviation: "ACRL", Color: MustParseHtmlColor("#075AAA"), OtherAbbrs: []string{"ATLC"}},
	{Code: "MC", Name: "Montclair-Boonton Line", Abbreviation: "MOBO", Color: MustParseHtmlColor("#E66859"), OtherAbbrs: []string{"BNTN", "BNTNM", "MNBTN"}},
	{Code: "BC", Name: "Bergen County Line", Abbreviation: "BERG", Color: MustParseHtmlColor("#FFD411"), OtherAbbrs: []string{"MNBN"}},
	{Code: "ML", Name: "Main Line", Abbreviation: "MAIN", Color: MustParseHtmlColor("#FFD411"), OtherAbbrs: []string{"MNBN"}},
	{Code: "ME", Name: "Morris & Essex Line", Abbreviation: "M&E", Color: MustParseHtmlColor("#08A652"), OtherAbbrs: []string{"MNE"}},
	{Code: "GS", Name: "Gladstone Branch", Abbreviation: "M&E", Color: MustParseHtmlColor("#A4C9AA"), OtherAbbrs: []string{"MNEG"}},
	{Code: "NE", Name: "Northeast Corridor Line", Abbreviation: "NEC", Color: MustParseHtmlColor("#DD3439")},
	{Code: "NC", Name: "North Jersey Coast Line", Abbreviation: "NJCL", Color: MustParseHtmlColor("#03A3DF"), OtherAbbrs: []string{"NJCLL"}},
	{Code: "PV", Name: "Pascack Valley Line", Abbreviation: "PASC", Color: MustParseHtmlColor("#94219A")},
	{Code: "PR", Name: "Princeton Branch", Abbreviation: "PRIN", Color: MustParseHtmlColor("#DD3439")},
	{Code: "RV", Name: "Raritan Valley Line", Abbreviation: "RARV", Color: MustParseHtmlColor("#F2A537")},
	{Code: "SL", Name: "BetMGM Meadowlands", Abbreviation: "BMGM", Color: MustParseHtmlColor("#C1AA72")},
	{Code: "AM", Name: "Amtrak", Abbreviation: "AMTK", Color: MustParseHtmlColor("#FFFF00")},
	{Code: "SP", Name: "Septa", Abbreviation: "SEPTA", Color: MustParseHtmlColor("#1F4FA3")},
}

// FindLine returns an object that lets you find a line by code or name.
// If no exact match is found and the name was specified, this function uses fuzzy search to find the closest match.
//
// The [LineFinder.SearchOrSynthesize] method will, if it doesn't find a suitable line, return a synthesized
// [Line] object that uses the provided search data. If no code was specified, "XX" will be used in its place.
// If no name was specified, "Unknown [line code]" will be used in its place, and "XX" followed by the line code
// will be used in place of the line abbreviation.
func FindLine() LineFinder {
	return finderImpl[Line, LineCode]{
		byCode:        linesByCode,
		byName:        linesByName,
		byAbbr:        linesByAbbreviation,
		list:          Lines,
		getCandidates: func(s *Line) []string { return []string{s.Name, s.Abbreviation} },
		synthesize: func(code *LineCode, name *string) Line {
			out := Line{}
			if code == nil {
				out.Code = "XX"
			} else {
				out.Code = *code
			}
			if name == nil {
				out.Name = "Unknown " + string(out.Code)
			} else {
				out.Name = *name
			}
			out.Abbreviation = "XX" + string(out.Code)
			return out
		},
	}
}

// LineFinder is an object to find lines by code or name.
type LineFinder = Finder[Line, LineCode]

var lineAliases = map[LineCode][]string{
	"AC": {"Atlantic City Rail Line", "Atl. City Line"},
	"MC": {"Montclair-Boonton"},
	"BC": {"Main/Bergen County Line", "Bergen Co. Line"},
	"ML": {"Port Jervis Line"},
	"ME": {"Morristown Line"},
	"NE": {"Northeast Corridor", "Northeast Corrdr"},
	"NC": {"No Jersey Coast"},
	"PV": {"Pascack Valley"},
	"PR": {"Princeton Shuttle"},
	"RV": {"Raritan Valley"},
}
var linesByCode = makeMap(Lines, func(l *Line) LineCode { return l.Code })
var linesByName = makeMmap(Lines, func(s *Line) []string { return append([]string{s.Name}, lineAliases[s.Code]...) })
var linesByAbbreviation = makeMmap(Lines, func(s *Line) []string { return append([]string{s.Abbreviation}, s.OtherAbbrs...) })

func makeMap[I any, C ~string](input []I, getKey func(*I) C) map[string]*I {
	out := map[string]*I{}
	for i := range input {
		v := &input[i]
		out[strings.ToLower(string(getKey(v)))] = v
	}
	return out
}

func makeMmap[I any, C ~string](input []I, getKeys func(*I) []C) map[string]*I {
	out := map[string]*I{}
	for i := range input {
		v := &input[i]
		for _, k := range getKeys(v) {
			out[strings.ToLower(string(k))] = v
		}
	}
	return out
}
