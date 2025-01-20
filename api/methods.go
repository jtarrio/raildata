package api

var GetToken = method[GetTokenRequest, GetTokenResponse]("getToken")
var IsValidToken = method[TokenRequest, ValidTokenResponse]("isValidToken")
var GetStationList = method[TokenRequest, []GetStations]("getStationList")
var GetStationMSG = method[GetStationMsgRequest, []StationMsgs]("getStationMSG")
var GetStationSchedule = method[GetStationScheduleRequest, []DailyStationInfo]("getStationSchedule")
var GetTrainSchedule = method[GetTrainScheduleRequest, StationInfo]("getTrainSchedule")
var GetTrainSchedule19Rec = method[GetTrainSchedule19RecRequest, StationInfo]("getTrainSchedule19Rec")
var GetTrainStopList = method[GetTrainStopListRequest, Stops]("getTrainStopList")
var GetVehicleData = method[TokenRequest, []VehicleDataInfo]("getVehicleData")

type GetTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenRequest struct {
	Token string `json:"token"`
}

type GetStationMsgRequest struct {
	Token   string `json:"token"`
	Station string `json:"station,omitempty"`
	Line    string `json:"line,omitempty"`
}

type GetStationScheduleRequest struct {
	Token   string `json:"token"`
	Station string `json:"station,omitempty"`
	NjtOnly string `json:"NJTOnly,omitempty"`
}

type GetTrainScheduleRequest struct {
	Token   string `json:"token"`
	Station string `json:"station"`
}

type GetTrainSchedule19RecRequest struct {
	Token   string `json:"token"`
	Station string `json:"station"`
	Line    string `json:"line,omitempty"`
}

type GetTrainStopListRequest struct {
	Token string `json:"token"`
	Train string `json:"train"`
}

type CapacityList struct {
	VEHICLE_NO           string        `json:"VEHICLE_NO"`
	LATITUDE             string        `json:"LATITUDE"`
	LONGITUDE            string        `json:"LONGITUDE"`
	CREATED_TIME         string        `json:"CREATED_TIME"`
	VEHICLE_TYPE         string        `json:"VEHICLE_TYPE"`
	CUR_PERCENTAGE       string        `json:"CUR_PERCENTAGE"`
	CUR_CAPACITY_COLOR   string        `json:"CUR_CAPACITY_COLOR"`
	CUR_PASSENGER_COUNT  string        `json:"CUR_PASSENGER_COUNT"`
	PREV_PERCENTAGE      string        `json:"PREV_PERCENTAGE"`
	PREV_CAPACITY_COLOR  string        `json:"PREV_CAPACITY_COLOR"`
	PREV_PASSENGER_COUNT string        `json:"PREV_PASSENGER_COUNT"`
	SECTIONS             []SectionList `json:"SECTIONS"`
}

type CarList struct {
	CAR_NO              string `json:"CAR_NO"`
	CAR_POSITION        string `json:"CAR_POSITION"`
	CAR_REST            bool   `json:"CAR_REST"`
	CUR_PERCENTAGE      string `json:"CUR_PERCENTAGE"`
	CUR_CAPACITY_COLOR  string `json:"CUR_CAPACITY_COLOR"`
	CUR_PASSENGER_COUNT string `json:"CUR_PASSENGER_COUNT"`
}

type DailyScheduleInfo struct {
	SCHED_DEP_DATE      string `json:"SCHED_DEP_DATE"`
	DESTINATION         string `json:"DESTINATION"`
	TRACK               string `json:"TRACK"`
	LINE                string `json:"LINE"`
	TRAIN_ID            string `json:"TRAIN_ID"`
	CONNECTING_TRAIN_ID string `json:"CONNECTING_TRAIN_ID"`
	STATION_POSITION    string `json:"STATION_POSITION"`
	DIRECTION           string `json:"DIRECTION"`
	DWELL_TIME          string `json:"DWELL_TIME"`
	PERM_PICKUP         string `json:"PERM_PICKUP"`
	PERM_DROPOFF        string `json:"PERM_DROPOFF"`
	STOP_CODE           string `json:"STOP_CODE"`
}

type DailyStationInfo struct {
	STATION_2CHAR string              `json:"STATION_2CHAR"`
	STATIONNAME   string              `json:"STATIONNAME"`
	ITEMS         []DailyScheduleInfo `json:"ITEMS"`
}

type GetTokenResponse struct {
	Authenticated string `json:"Authenticated"`
	UserToken     string `json:"UserToken"`
}

type GetStations struct {
	STATION_2CHAR  string `json:"STATION_2CHAR"`
	STATIONNAME    string `json:"STATIONNAME"`
	STATION_14CHAR string `json:"STATION_14CHAR"`
}

type ScheduleInfo struct {
	SCHED_DEP_DATE      string         `json:"SCHED_DEP_DATE"`
	DESTINATION         string         `json:"DESTINATION"`
	TRACK               string         `json:"TRACK"`
	LINE                string         `json:"LINE"`
	TRAIN_ID            string         `json:"TRAIN_ID"`
	CONNECTING_TRAIN_ID string         `json:"CONNECTING_TRAIN_ID"`
	STATUS              string         `json:"STATUS"`
	SEC_LATE            string         `json:"SEC_LATE"`
	LAST_MODIFIED       string         `json:"LAST_MODIFIED"`
	BACKCOLOR           string         `json:"BACKCOLOR"`
	FORECOLOR           string         `json:"FORECOLOR"`
	SHADOWCOLOR         string         `json:"SHADOWCOLOR"`
	GPSLATITUDE         string         `json:"GPSLATITUDE"`
	GPSLONGITUDE        string         `json:"GPSLONGITUDE"`
	GPSTIME             string         `json:"GPSTIME"`
	STATION_POSITION    string         `json:"STATION_POSITION"`
	LINECODE            string         `json:"LINECODE"`
	LINEABBREVIATION    string         `json:"LINEABBREVIATION"`
	INLINEMSG           string         `json:"INLINEMSG"`
	CAPACITY            []CapacityList `json:"CAPACITY"`
	STOPS               []StopList     `json:"STOPS"`
}

type SectionList struct {
	SECTION_POSITION    string    `json:"SECTION_POSITION"`
	CUR_PERCENTAGE      string    `json:"CUR_PERCENTAGE"`
	CUR_CAPACITY_COLOR  string    `json:"CUR_CAPACITY_COLOR"`
	CUR_PASSENGER_COUNT string    `json:"CUR_PASSENGER_COUNT"`
	CARS                []CarList `json:"CARS"`
}

type StationInfo struct {
	STATION_2CHAR string         `json:"STATION_2CHAR"`
	STATIONNAME   string         `json:"STATIONNAME"`
	STATIONMSGS   []StationMsgs  `json:"STATIONMSGS"`
	ITEMS         []ScheduleInfo `json:"ITEMS"`
}

type StationMsgs struct {
	MSG_TYPE          string `json:"MSG_TYPE"`
	MSG_TEXT          string `json:"MSG_TEXT"`
	MSG_PUBDATE       string `json:"MSG_PUBDATE"`
	MSG_ID            string `json:"MSG_ID"`
	MSG_AGENCY        string `json:"MSG_AGENCY"`
	MSG_SOURCE        string `json:"MSG_SOURCE"`
	MSG_STATION_SCOPE string `json:"MSG_STATION_SCOPE"`
	MSG_LINE_SCOPE    string `json:"MSG_LINE_SCOPE"`
	MSG_PUBDATE_UTC   string `json:"MSG_PUBDATE_UTC"`
}

type StopLines struct {
	LINE_CODE  string `json:"LINE_CODE"`
	LINE_NAME  string `json:"LINE_NAME"`
	LINE_COLOR string `json:"LINE_COLOR"`
}

type StopList struct {
	STATION_2CHAR   string      `json:"STATION_2CHAR"`
	STATIONNAME     string      `json:"STATIONNAME"`
	TIME            string      `json:"TIME"`
	PICKUP          string      `json:"PICKUP"`
	DROPOFF         string      `json:"DROPOFF"`
	DEPARTED        string      `json:"DEPARTED"`
	STOP_STATUS     string      `json:"STOP_STATUS"`
	DEP_TIME        string      `json:"DEP_TIME"`
	TIME_UTC_FORMAT string      `json:"TIME_UTC_FORMAT"`
	STOP_LINES      []StopLines `json:"STOP_LINES"`
}

type Stops struct {
	TRAIN_ID    string         `json:"TRAIN_ID"`
	LINECODE    string         `json:"LINECODE"`
	BACKCOLOR   string         `json:"BACKCOLOR"`
	FORECOLOR   string         `json:"FORECOLOR"`
	SHADOWCOLOR string         `json:"SHADOWCOLOR"`
	DESTINATION string         `json:"DESTINATION"`
	TRANSFERAT  string         `json:"TRANSFERAT"`
	STOPS       []StopList     `json:"STOPS"`
	CAPACITY    []CapacityList `json:"CAPACITY"`
}

type ValidTokenResponse struct {
	ValidToken bool   `json:"validToken"`
	UserID     string `json:"userID"`
}

type VehicleDataInfo struct {
	ID             string `json:"ID"`
	TRAIN_LINE     string `json:"TRAIN_LINE"`
	DIRECTION      string `json:"DIRECTION"`
	ICS_TRACK_CKT  string `json:"ICS_TRACK_CKT"`
	LAST_MODIFIED  string `json:"LAST_MODIFIED"`
	SCHED_DEP_TIME string `json:"SCHED_DEP_TIME"`
	SEC_LATE       string `json:"SEC_LATE"`
	NEXT_STOP      string `json:"NEXT_STOP"`
	LONGITUDE      string `json:"LONGITUDE"`
	LATITUDE       string `json:"LATITUDE"`
}
