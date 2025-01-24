package raildata

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Client is the interface you use to access the RailData server.
type Client interface {
	// GetToken returns the token currently being used by the client.
	GetToken() string
	// GetStationList returns a list of all stations, with their codes, names, and short names.
	GetStationList(context.Context) (*GetStationListResponse, error)
	// GetStationMsg returns a list of messages and alerts, optionally scoped to one station and/or one line.
	GetStationMsg(context.Context, *GetStationMsgRequest) (*GetStationMsgResponse, error)
	// GetTrainSchedule returns the schedule for the next 19 trains departing from a station.
	//
	// This method also returns information about each train's stops.
	GetTrainSchedule(context.Context, *GetTrainScheduleRequest) (*GetTrainScheduleResponse, error)
	// GetTrainSchedule19Records returns the schedule for the next 19 trains departing from a station,
	// optionally scoped to one line.
	//
	// This method does not return information about each train's stops. You can use GetTrainStopList to retrieve it.
	GetTrainSchedule19Records(context.Context, *GetTrainSchedule19RecordsRequest) (*GetTrainScheduleResponse, error)
	// GetTrainStopList returns the list of stops for a train.
	// Returns nil if the provided train id is not valid.
	GetTrainStopList(context.Context, *GetTrainStopListRequest) (*GetTrainStopListResponse, error)
	// GetVehicleData returns the position and status for all active trains.
	//
	// A train appears in this list if it has moved in the last 5 minutes.
	GetVehicleData(context.Context) (*GetVehicleDataResponse, error)
	// RateLimitedMethods returns an interface for rate-limited operations.
	RateLimitedMethods() RateLimitedMethods
}

// RateLimitedMethods contains methods you can only call a few times per day.
type RateLimitedMethods interface {
	// IsValidToken returns whether the current token is valid, and the user name associated with it.
	//
	// Note that NJ Transit sets a limit of 10 calls per day to this endpoint, so don't use it indiscriminately.
	// An alternative is to call GetStationList or GetStationMsg, which will return a result if the token is valid
	// (and will try to autoupdate it if necessary.)
	IsValidToken(context.Context) (*IsValidTokenResponse, error)
	// GetStationSchedule returns the schedule for the next 27 hours for one station.
	//
	// There is a limit of 5 calls per day and, even though the documentation claims you can get a full
	// schedule for all stations omitting the StationCode parameter, it is not true, so you should avoid
	// using this method.
	GetStationSchedule(context.Context, *GetStationScheduleRequest) (*GetStationScheduleResponse, error)
}

// IsValidTokenResponse contains the result of the IsValidToken method.
type IsValidTokenResponse struct {
	// ValidToken contains whether the token is valid.
	ValidToken bool
	// UserId contains the user id for this token, if found.
	UserId *string
}

// GetStationListResponse contains the result of the GetStationList method.
type GetStationListResponse struct {
	// Stations contains the list of stations.
	Stations []Station
}

// GetStationMsgRequest contains the arguments of the GetStationMsg method.
type GetStationMsgRequest struct {
	// StationCode contains an optional station to filter the messages on.
	StationCode *StationCode
	// LineCode contains an optional line to filter the messages on.
	LineCode *LineCode
}

// GetStationMsgResponse contains the result of the GetStationMsg method.
type GetStationMsgResponse struct {
	// Messages contains the list of messages.
	Messages []StationMsg
}

// GetStationScheduleRequest contains the arguments of the GetStationSchedule method.
type GetStationScheduleRequest struct {
	// StationCode contains the station to get the schedule for.
	StationCode StationCode
	// NjtOnly contains whether only NJ Transit trains should be returned.
	// If false, the results may also contain Amtrak trains.
	NjtOnly bool
}

// GetStationScheduleResponse contains the result of the GetStationSchedule method.
type GetStationScheduleResponse struct {
	// Entries contains the schedule entries.
	Entries []StationSchedule
}

// GetTrainScheduleRequest contains the arguments of the GetTrainSchedule method.
type GetTrainScheduleRequest struct {
	// StationCode contains the station to get the schedule for.
	StationCode StationCode
}

// GetTrainSchedule19RecordsRequest contains the arguments of the GetTrainSchedule19Records method.
type GetTrainSchedule19RecordsRequest struct {
	// StationCode contains the station to get the schedule for.
	StationCode StationCode
	// LineCode contains an optional line to filter the schedule on.
	LineCode *LineCode
}

// GetTrainScheduleResponse contains the result of the GetTrainSchedule and GetTrainSchedule19Records methods.
type GetTrainScheduleResponse struct {
	// Station contains the station this schedule belongs to.
	Station Station
	// Messages contains a list of messages for this station.
	Messages []StationMsg
	// Entries contains the schedule entries.
	Entries []TrainScheduleEntry
}

// GetTrainStopListRequest contains the arguments of the GetTrainStopList method.
type GetTrainStopListRequest struct {
	// TrainId contains the train whose list of stops needs to be returned.
	TrainId string
}

// GetTrainStopListResponse contains the result of the GetTrainStopList method.
type GetTrainStopListResponse struct {
	// TrainId contains the train this information belongs to.
	TrainId string
	// Line contains the line this train runs on.
	Line Line
	// Color contains the colors used to render the line name.
	Color ColorSet
	// Destination contains the destination name.
	Destination string
	// DestinationStation contains the destination station, if it could be determined from the destination name.
	DestinationStation *Station
	// TransferAt contains the name of a transfer station. Used for Long Branch connections to Bayhead.
	TransferAt *string
	// Stops contains the list of stops for this train.
	Stops []TrainStop
	// Capacity contains information on how full this train is.
	Capacity []TrainCapacity
}

// GetVehicleDataResponse contains the result of the GetVehicleData method.
type GetVehicleDataResponse struct {
	// Vehicles contains a list of active trains.
	Vehicles []VehicleData
}

// Location contains a vehicle's GPS location.
type Location struct {
	// Longitude contains the vehicle's longitude, in degrees East.
	Longitude float64
	// Latitude contains the vehicle's latitude, in degrees North.
	Latitude float64
}

// StationMsg contains a message or alert.
type StationMsg struct {
	// Type contains the type of message.
	Type MsgType
	// Text contains the message's text. It may contain HTML code and escape sequences.
	Text string
	// PubDate indicates when the message was published.
	PubDate time.Time
	// Id contains an identifier for the message. Typically only used with messages from third-party feeds.
	Id *string
	// Agency identifies the agency that published the message: NJT, AMT, or none.
	Agency *string
	// Source identifies the source of the message: RSS_NJTRailAlerts or none.
	Source *string
	// StationScope contains a list of stations this message pertains to.
	StationScope []Station
	// LineScope contains a list of lines this message pertains to.
	LineScope []Line
}

// StationSchedule contains a station's 27-hour schedule.
type StationSchedule struct {
	// Station identifies the station this schedule belongs to.
	Station *Station
	// Entries contains a list of schedule entries.
	Entries []ScheduleEntry
}

// ScheduleEntry contains an entry in a station's schedule.
type ScheduleEntry struct {
	// DepartureTime contains the scheduled departure date/time for the train.
	DepartureTime time.Time
	// Destination contains the destination name.
	Destination string
	// DestinationStation contains the destination station, if it could be determined from the destination name.
	DestinationStation *Station
	// Line contains the line this train runs on.
	Line Line
	// TrainId contains the train's number.
	TrainId string
	// ConnectingTrainId contains the connecting train's number. Used for Long Branch connections to Bayhead.
	ConnectingTrainId *string
	// StationPosition contains this station's position along the trip.
	StationPosition StationPosition
	// Direction contains the direction of travel.
	Direction Direction
	// DwellTime contains how long the train will wait at the station.
	DwellTime *time.Duration
	// PickupOnly indicates, if true, that the train will only pick up (not discharge) passengers at this stop.
	PickupOnly bool
	// DropoffOnly indicates, if true, that the train will only discharge (not pick up) passengers at this stop.
	DropoffOnly bool
	// StopCode contains a stop code for this station.
	StopCode *StopCode
}

// TrainScheduleEntry contains an entry in a train's schedule at a station.
type TrainScheduleEntry struct {
	// DepartureTime contains the scheduled departure date/time for the train.
	DepartureTime time.Time
	// Destination contains the destination name.
	Destination string
	// Track contains the name of the track this train will leave from, if known.
	Track *string
	// Line contains the line this train runs on.
	Line Line
	// LineName contains the display name for the line. For example, the [Line] may be Amtrak, but [LineName] may contain "Acela Express".
	LineName string
	// TrainId contains the train's number.
	TrainId string
	// ConnectingTrainId contains the connecting train's number. Used for Long Branch connections to Bayhead.
	ConnectingTrainId *string
	// Status contains the train's current status.
	Status *string
	// Delay contains the train's current delay.
	Delay *time.Duration
	// LastUpdated contains the date/time this entry was updated.
	LastUpdated *time.Time
	// Color contains the colors used to render the line name.
	Color ColorSet
	// GpsLocation contains the train's last known position.
	GpsLocation *Location
	// GpsTime contains the date/time the GpsLocation was captured.
	GpsTime *time.Time
	// StationPosition contains this station's position along the trip.
	StationPosition StationPosition
	// InlineMessage contains an in-line message for the train at the station.
	InlineMessage *string
	// Capacity contains information on how full this train is.
	Capacity []TrainCapacity
	// Stops contains the list of stops for this train.
	Stops []TrainStop
}

// TrainCapacity contains information on how full a train is.
type TrainCapacity struct {
	// Number contains the train's number.
	Number string
	// Location contains the train's last known position.
	Location Location
	// CreatedTime contains the time this record was creted.
	CreatedTime time.Time
	// Type contains the vehicle type.
	Type string
	// CapacityPercent contains the percentage of capacity used.
	CapacityPercent int
	// CapacityColor contains a color that represents how full the train is.
	CapacityColor Color
	// PassengerCount contains the number of passengers on board the train.
	PassengerCount int
	// Sections contains capacity information for each train section.
	Sections []TrainSection
}

// TrainSection contains information on how full a train section is.
type TrainSection struct {
	// Position contains the section's position on the train.
	Position SectionPosition
	// CapacityPercent contains the percentage of capacity used.
	CapacityPercent int
	// CapacityColor contains a color that represents how full the section is.
	CapacityColor Color
	// PassengerCount contains the number of passengers on board the section.
	PassengerCount int
	// Cars contains capacity information for each car in this section.
	Cars []TrainCar
}

// TrainCar contains information on how full a train car is.
type TrainCar struct {
	// TrainId contains the train car's number.
	TrainId string
	// Position contains the car's position on the train, 1 being the front.
	Position int
	// Restroom contains whether this car has a restroom.
	Restroom bool
	// CapacityPercent contains the percentage of capacity used.
	CapacityPercent int
	// CapacityColor contains a color that represents how full the car is.
	CapacityColor Color
	// PassengerCount contains the number of passengers on board the car.
	PassengerCount int
}

// TrainStop contains information about a train's stop.
type TrainStop struct {
	// Station contains the station where this train stops.
	Station Station
	// ArrivalTime contains the expected arrival time.
	ArrivalTime *time.Time
	// PickupOnly indicates, if true, that the train will only pick up (not discharge) passengers at this stop.
	PickupOnly bool
	// DropoffOnly indicates, if true, that the train will only discharge (not pick up) passengers at this stop.
	DropoffOnly bool
	// Departed indicates, if true, that the train has already left this station.
	Departed bool
	// StopStatus contains an optional status at the stop: OnTime, Delayed, Cancelled, or none.
	StopStatus *string
	// DepartureTime contains the expected departure time.
	DepartureTime *time.Time
	// StopLines contains a list of lines that connect at this stop.
	StopLines []StopLine
}

// StopLine contains information about a connecting line.
type StopLine struct {
	// Line contain's the line that connects at this stop.
	Line Line
	// Color contains a color used to represent the line.
	Color Color
}

// VehicleData contains innformation about an active train.
type VehicleData struct {
	// TrainId contains the train's number.
	TrainId string
	// Line contains the line this train is running on.
	Line Line
	// Direction contains this train's direction of travel.
	Direction Direction
	// TrackCircuitId contains the last identified circuit id for this train.
	TrackCircuitId string
	// LastUpdated contains the time this information was last refreshed.
	LastUpdated time.Time
	// DepartureTime contains the expected departure time at the next station.
	DepartureTime time.Time
	// Delay contains the train's current delay.
	Delay *time.Duration
	// NextStop contains the train's next stop.
	NextStop *Station
	// Location contains the train's GPS location.
	Location *Location
}

// MsgType represents the type of a message.
type MsgType int

const (
	MsgTypeBanner     MsgType = iota // a "banner" style message, displayed along other information.
	MsgTypeFullScreen                // a message that takes over the screen.
)

// Direction represents the direction of travel of a trip.
type Direction int

const (
	DirectionEastbound Direction = iota // eastbound.
	DirectionWestbound                  // westbound.
)

// SectionPosition designates a section in a train.
type SectionPosition int

const (
	SectionPositionFront  SectionPosition = iota // the first cars in a train.
	SectionPositionMiddle                        // the middle cars in a train.
	SectionPositionBack                          // the last cars in a train.
)

// ColorSet contains colors used to render a line name.
type ColorSet struct {
	// Foreground contains the color for the text.
	Foreground Color
	// Background contains the color for the background.
	Background Color
	// Shadow contains the color for the text's shadow.
	Shadow Color
}

// Color contains a color specification to render a line name.
type Color struct {
	rgb [3]byte
}

// NewColor returns a Color object given its red, green, and blue components (from 0 to 255).
func NewColor(r, g, b int) (Color, error) {
	if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
		return Color{}, errors.New("invalid color")
	}
	return Color{rgb: [3]byte{byte(r), byte(g), byte(b)}}, nil
}

// ParseHtmlColor parses an HTML color specification into a Color object.
func ParseHtmlColor(s string) (Color, error) {
	if len(s) == 0 || s[0] != '#' {
		return Color{}, errors.New("color specification does not start with #")
	}
	var hexBytes []byte
	if len(s) == 4 {
		hexBytes = []byte{s[1], s[1], s[2], s[2], s[3], s[3]}
	} else if len(s) == 7 {
		hexBytes = []byte(s[1:7])
	} else {
		return Color{}, errors.New("color specification does not have the correct length")
	}
	rgb := make([]byte, 3)
	_, err := hex.Decode(rgb, hexBytes)
	if err != nil {
		return Color{}, err
	}
	return Color{rgb: [3]byte(rgb)}, nil
}

// Html returns an HTML color specification for this color.
func (c Color) Html() string {
	dst := make([]byte, 1, 7)
	dst[0] = '#'
	dst = hex.AppendEncode(dst, c.rgb[:])
	return fmt.Sprintf("#%s%s%s", string(dst[0:2]), string(dst[2:4]), string(dst[4:6]))
}

// RGB returns the red, green, and blue components of this color (from 0 to 255).
func (c Color) RGB() (r, g, b int) {
	return int(c.rgb[0]), int(c.rgb[1]), int(c.rgb[2])
}
