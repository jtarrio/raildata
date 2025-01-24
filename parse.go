package raildata

import (
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/jtarrio/raildata/api"
)

var njLocation = func() *time.Location {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		panic(err)
	}
	return loc
}()

const msgDateTimeFormat = "1/2/2006 3:04:05 PM"
const dateTimeFormat = "02-Jan-2006 03:04:05 PM"

func ParseValidTokenResponse(input *api.ValidTokenResponse) (*IsValidTokenResponse, error) {
	response := &IsValidTokenResponse{
		ValidToken: input.ValidToken,
		UserId:     strToPtr(input.UserID),
	}
	return response, nil
}

func ParseGetStationsList(input []api.GetStations) (*GetStationListResponse, error) {
	response := &GetStationListResponse{}
	for _, item := range input {
		response.Stations = append(response.Stations, *ParseGetStations(&item))
	}
	return response, nil
}

func ParseGetStations(input *api.GetStations) *Station {
	stationEntry := &Station{
		Code:      StationCode(input.STATION_2CHAR),
		Name:      input.STATIONNAME,
		ShortName: input.STATION_14CHAR,
	}
	return stationEntry
}

func ParseStationMsgsList(input []api.StationMsgs) *GetStationMsgResponse {
	response := &GetStationMsgResponse{}
	for _, item := range input {
		response.Messages = append(response.Messages, *ParseStationMsgs(&item))
	}
	return response
}

func ParseStationMsgs(input *api.StationMsgs) *StationMsg {
	stationMsg := &StationMsg{
		Type:         strToMsgType(input.MSG_TYPE),
		Text:         input.MSG_TEXT,
		PubDate:      *strToLocalTime(input.MSG_PUBDATE, msgDateTimeFormat),
		Id:           strToPtr(input.MSG_ID),
		Agency:       strToPtr(input.MSG_AGENCY),
		Source:       strToPtr(input.MSG_SOURCE),
		StationScope: decodeStationScope(input.MSG_STATION_SCOPE),
		LineScope:    decodeLineScope(input.MSG_LINE_SCOPE),
	}
	return stationMsg
}

func ParseDailyStationInfoList(input []api.DailyStationInfo) (*GetStationScheduleResponse, error) {
	response := &GetStationScheduleResponse{}
	for _, item := range input {
		response.Entries = append(response.Entries, *ParseDailyStationInfo(&item))
	}
	return response, nil
}

func ParseDailyStationInfo(input *api.DailyStationInfo) *StationSchedule {
	stationSchedule := &StationSchedule{
		Station: strToStation(input.STATION_2CHAR, input.STATIONNAME),
	}
	for _, item := range input.ITEMS {
		stationSchedule.Entries = append(stationSchedule.Entries, *ParseDailyScheduleInfo(&item))
	}
	return stationSchedule
}

func ParseDailyScheduleInfo(input *api.DailyScheduleInfo) *ScheduleEntry {
	destination := strUnquote(input.DESTINATION)
	scheduleEntry := &ScheduleEntry{
		DepartureTime:      *strToLocalTime(input.SCHED_DEP_DATE, dateTimeFormat),
		Destination:        destination,
		DestinationStation: strToStation("", destination),
		Line:               *strToLine("", input.LINE),
		TrainId:            input.TRAIN_ID,
		ConnectingTrainId:  strToPtr(input.CONNECTING_TRAIN_ID),
		StationPosition:    GetStationPosition(input.STATION_POSITION),
		Direction:          strToDirection(input.DIRECTION),
		DwellTime:          strToDurationSeconds(input.DWELL_TIME),
		PickupOnly:         strToBool(input.PERM_PICKUP),
		DropoffOnly:        strToBool(input.PERM_DROPOFF),
		StopCode:           strToStopCode(input.STOP_CODE),
	}
	return scheduleEntry
}

func ParseStationInfo(input *api.StationInfo) *GetTrainScheduleResponse {
	response := &GetTrainScheduleResponse{
		Station: *strToStation(input.STATION_2CHAR, input.STATIONNAME),
	}
	for _, item := range input.STATIONMSGS {
		response.Messages = append(response.Messages, *ParseStationMsgs(&item))
	}
	for _, item := range input.ITEMS {
		response.Entries = append(response.Entries, *ParseScheduleInfo(&item, &response.Station))
	}
	return response
}

func ParseScheduleInfo(input *api.ScheduleInfo, station *Station) *TrainScheduleEntry {
	destination := strUnquote(input.DESTINATION)
	scheduleEntry := &TrainScheduleEntry{
		DepartureTime:     *strToLocalTime(input.SCHED_DEP_DATE, dateTimeFormat),
		Destination:       destination,
		Track:             strToTrackName(input.TRACK, station),
		Line:              *strToLine(input.LINECODE, input.LINE),
		LineName:          input.LINE,
		TrainId:           input.TRAIN_ID,
		ConnectingTrainId: strToPtr(input.CONNECTING_TRAIN_ID),
		Status:            strToPtr(input.STATUS),
		Delay:             strToDurationSeconds(input.SEC_LATE),
		LastUpdated:       strToLocalTime(input.LAST_MODIFIED, dateTimeFormat),
		Color:             *strsToColorSet(input.FORECOLOR, input.BACKCOLOR, input.SHADOWCOLOR),
		GpsLocation:       strsToLocation(input.GPSLONGITUDE, input.GPSLATITUDE),
		GpsTime:           strToLocalTime(input.GPSTIME, dateTimeFormat),
		StationPosition:   GetStationPosition(input.STATION_POSITION),
		InlineMessage:     strToPtr(input.INLINEMSG),
	}
	for _, item := range input.CAPACITY {
		scheduleEntry.Capacity = append(scheduleEntry.Capacity, *ParseCapacityList(&item))
	}
	for _, item := range input.STOPS {
		scheduleEntry.Stops = append(scheduleEntry.Stops, *ParseStopList(&item))
	}
	return scheduleEntry
}

func ParseCapacityList(input *api.CapacityList) *TrainCapacity {
	response := &TrainCapacity{
		Number:          *strToPtr(input.VEHICLE_NO),
		Location:        *strsToLocation(input.LONGITUDE, input.LATITUDE),
		CreatedTime:     *strToLocalTime(input.CREATED_TIME, dateTimeFormat),
		Type:            *strToPtr(input.VEHICLE_TYPE),
		CapacityPercent: *strToInt(input.CUR_PERCENTAGE),
		CapacityColor:   *strToColor(input.CUR_CAPACITY_COLOR),
		PassengerCount:  *strToInt(input.CUR_PASSENGER_COUNT),
	}
	for _, item := range input.SECTIONS {
		response.Sections = append(response.Sections, *ParseSectionList(&item))
	}
	return response
}

func ParseSectionList(input *api.SectionList) *TrainSection {
	response := &TrainSection{
		Position:        strToSectionPosition(input.SECTION_POSITION),
		CapacityPercent: *strToInt(input.CUR_PERCENTAGE),
		CapacityColor:   *strToColor(input.CUR_CAPACITY_COLOR),
		PassengerCount:  *strToInt(input.CUR_PASSENGER_COUNT),
	}
	for _, item := range input.CARS {
		response.Cars = append(response.Cars, *ParseCarList(&item))
	}
	return response
}

func ParseCarList(input *api.CarList) *TrainCar {
	response := &TrainCar{
		TrainId:         *strToPtr(input.CAR_NO),
		Position:        *strToInt(input.CAR_POSITION),
		Restroom:        input.CAR_REST,
		CapacityPercent: *strToInt(input.CUR_PERCENTAGE),
		CapacityColor:   *strToColor(input.CUR_CAPACITY_COLOR),
		PassengerCount:  *strToInt(input.CUR_PASSENGER_COUNT),
	}
	return response
}

func ParseStopList(input *api.StopList) *TrainStop {
	response := &TrainStop{
		Station:       *strToStation(input.STATION_2CHAR, input.STATIONNAME),
		ArrivalTime:   strToLocalTime(input.TIME, dateTimeFormat),
		PickupOnly:    strToBool(input.PICKUP),
		DropoffOnly:   strToBool(input.DROPOFF),
		Departed:      strToBool(input.DEPARTED),
		StopStatus:    strToPtr(input.STOP_STATUS),
		DepartureTime: strToLocalTime(input.DEP_TIME, dateTimeFormat),
	}
	for _, item := range input.STOP_LINES {
		response.StopLines = append(response.StopLines, *ParseStopLines(&item))
	}
	return response
}

func ParseStopLines(input *api.StopLines) *StopLine {
	response := &StopLine{
		Line:  *strToLine(input.LINE_CODE, input.LINE_NAME),
		Color: *strToColor(input.LINE_COLOR),
	}
	return response
}

func ParseStops(input *api.Stops) *GetTrainStopListResponse {
	trainidp := strToPtr(input.TRAIN_ID)
	if trainidp == nil {
		return nil
	}
	destination := strUnquote(input.DESTINATION)
	response := &GetTrainStopListResponse{
		TrainId:            *trainidp,
		Line:               *strToLine(input.LINECODE, ""),
		Color:              *strsToColorSet(input.FORECOLOR, input.BACKCOLOR, input.SHADOWCOLOR),
		Destination:        destination,
		DestinationStation: strToStation("", destination),
		TransferAt:         strToPtr(input.TRANSFERAT),
	}
	for _, item := range input.STOPS {
		response.Stops = append(response.Stops, *ParseStopList(&item))
	}
	for _, item := range input.CAPACITY {
		response.Capacity = append(response.Capacity, *ParseCapacityList(&item))
	}
	return response
}

func ParseVehicleDataInfoList(input []api.VehicleDataInfo) *GetVehicleDataResponse {
	response := &GetVehicleDataResponse{}
	for _, item := range input {
		response.Vehicles = append(response.Vehicles, *ParseVehicleDataInfo(&item))
	}
	return response
}

func ParseVehicleDataInfo(input *api.VehicleDataInfo) *VehicleData {
	response := &VehicleData{
		TrainId:        input.ID,
		Line:           *strToLine("", input.TRAIN_LINE),
		Direction:      strToDirection(input.DIRECTION),
		TrackCircuitId: input.ICS_TRACK_CKT,
		LastUpdated:    *strToLocalTime(input.LAST_MODIFIED, dateTimeFormat),
		DepartureTime:  *strToLocalTime(input.SCHED_DEP_TIME, dateTimeFormat),
		Delay:          strToDurationSeconds(input.SEC_LATE),
		NextStop:       strToStation("", input.NEXT_STOP),
		Location:       strsToLocation(input.LONGITUDE, input.LATITUDE),
	}
	return response
}

func strUnquote(s string) string {
	return html.UnescapeString(s)
}

func strToPtr(s string) *string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return nil
	}
	return &s
}

func strToBool(s string) bool {
	s = strings.ToLower(s)
	return s == "true" || s == "yes"
}

func strToFloat(s string) *float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &f
}

func strToInt(s string) *int {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	r := int(i)
	return &r
}

func strToColor(s string) *Color {
	p := strToPtr(s)
	if p == nil {
		return nil
	}
	c, err := ParseHtmlColor(*p)
	if err != nil {
		return nil
	}
	return &c
}

func strsToColorSet(fg, bg, shadow string) *ColorSet {
	fgc := strToColor(fg)
	bgc := strToColor(bg)
	shadowc := strToColor(shadow)
	if fgc == nil || bgc == nil {
		return nil
	}
	if shadowc == nil {
		shadowc = &Color{}
	}
	return &ColorSet{
		Foreground: *fgc,
		Background: *bgc,
		Shadow:     *shadowc,
	}
}

func strToLocalTime(s string, format string) *time.Time {
	t, err := time.ParseInLocation(format, s, njLocation)
	if err != nil {
		return nil
	}
	return &t
}

func strToDurationSeconds(s string) *time.Duration {
	secs, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil
	}
	duration := time.Duration(secs) * time.Second
	return &duration
}

func strsToLocation(lon string, lat string) *Location {
	lonf := strToFloat(lon)
	latf := strToFloat(lat)
	if lonf == nil || latf == nil {
		return nil
	}
	return &Location{
		Longitude: *lonf,
		Latitude:  *latf,
	}
}

func strToStation(code string, name string) *Station {
	fs := FindStation()
	if codep := (*StationCode)(strToPtr(code)); codep != nil {
		fs = fs.WithCode(*codep)
	}
	if namep := strToPtr(name); namep != nil {
		fs = fs.WithName(*namep)
	}
	return fs.SearchOrSynthesize()
}

func strToLine(code string, name string) *Line {
	fs := FindLine()
	if codep := (*LineCode)(strToPtr(code)); codep != nil {
		fs = fs.WithCode(*codep)
	}
	if namep := strToPtr(name); namep != nil {
		fs = fs.WithName(*namep)
	}
	return fs.SearchOrSynthesize()
}

func strToTrackName(track string, station *Station) *string {
	trackp := strToPtr(track)
	if trackp == nil || station == nil {
		return trackp
	}
	translation := TranslateTrackNumber(track, station.Code)
	return &translation
}

func strToMsgType(msgType string) MsgType {
	switch msgType {
	case "fullscreen":
		return MsgTypeFullScreen
	default:
		return MsgTypeBanner
	}
}

func strToDirection(direction string) Direction {
	switch direction {
	case "Eastbound":
		return DirectionEastbound
	default:
		return DirectionWestbound
	}
}

func strToSectionPosition(position string) SectionPosition {
	switch position {
	case "Front":
		return SectionPositionFront
	case "Back":
		return SectionPositionBack
	default:
		return SectionPositionMiddle
	}
}

func strToStopCode(code string) *StopCode {
	codep := strToPtr(code)
	if codep == nil {
		return nil
	}
	stopCode := GetStopCode(*codep)
	return &stopCode
}

func decodeStationScope(s string) []Station {
	scope := decodeScope(s)
	var out []Station
	for _, stationName := range scope {
		out = append(out, *FindStation().WithName(stationName).SearchOrSynthesize())
	}
	return out
}

func decodeLineScope(s string) []Line {
	scope := decodeScope(s)
	var out []Line
	for _, lineName := range scope {
		if line, found := linesByName[lineName]; found {
			out = append(out, *line)
		} else if line, found := linesByAbbreviation[lineName]; found {
			out = append(out, *line)
		}
	}
	return out
}

func decodeScope(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		for _, part := range strings.Split(part, "*") {
			part = strings.TrimSpace(part)
			if len(part) > 0 {
				out = append(out, part)
			}
		}
	}
	return out
}
