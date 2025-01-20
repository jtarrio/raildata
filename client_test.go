package raildata_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/jtarrio/raildata"
	"github.com/jtarrio/raildata/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testToken = "the-token"

func TestGetToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		assert.Fail(t, "did not expect any API requests")
	}))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)
	actual := client.GetToken()
	assert.Equal(t, testToken, actual)
}

func TestIsValidTokenCanReturnTrue(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "isValidToken", "token", testToken).sendJson("{\n  \"validToken\": true,\n  \"userID\": \"the-user-id\"\n}"))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	actual, err := client.RateLimitedMethods().IsValidToken(context.Background())
	assert.NoError(t, err)
	assert.True(t, actual.ValidToken)
	assert.Equal(t, "the-user-id", *actual.UserId)
}

func TestIsValidTokenCanReturnFalse(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "isValidToken", "token", testToken).sendJson(`{
  "validToken": false,
  "userID": null
}`))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	actual, err := client.RateLimitedMethods().IsValidToken(context.Background())
	assert.NoError(t, err)
	assert.False(t, actual.ValidToken)
	assert.Nil(t, actual.UserId)
}

func TestMissingCredentialsError(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "isValidToken").sendEmpty())

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(""))
	require.NoError(t, err)

	_, err = client.RateLimitedMethods().IsValidToken(context.Background())
	assert.ErrorIs(t, err, errors.MissingCredentialsError)
}

func TestOtherError(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "isValidToken").sendError("some error message"))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	_, err = client.RateLimitedMethods().IsValidToken(context.Background())
	var rderr *errors.RailDataError
	assert.ErrorAs(t, err, &rderr)
	assert.Equal(t, "some error message", rderr.Error())
}

func TestRenewTokenWhenRequired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.NoError(t, req.ParseMultipartForm(5000000))
		switch req.URL.Path {
		case "/isValidToken":
			if req.Form.Get("token") == "oldtoken" {
				expectRequest(t, "isValidToken").sendError("Invalid token.").ServeHTTP(rw, req)
			} else if req.Form.Get("token") == "newtoken" {
				expectRequest(t, "isValidToken").sendJson(`{"validToken":true,"userID":"the-user-id"}`).ServeHTTP(rw, req)
			} else {
				assert.Equal(t, "newtoken", req.Form.Get("token"))
			}
		case "/getToken":
			expectRequest(t, "getToken", "username", "the-user-id", "password", "the-password").sendJson(`{
 "Authenticated": "True",
 "UserToken": "newtoken"
}`).ServeHTTP(rw, req)
		}
	}))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken("oldtoken"), raildata.WithCredentials("the-user-id", "the-password"))
	require.NoError(t, err)

	assert.Equal(t, "oldtoken", client.GetToken())
	actual, err := client.RateLimitedMethods().IsValidToken(context.Background())
	assert.NoError(t, err)
	assert.True(t, actual.ValidToken)
	assert.Equal(t, "the-user-id", *actual.UserId)
	assert.Equal(t, "newtoken", client.GetToken())
}

func TestRenewTokenFailsWithoutCredentials(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "isValidToken").sendError("Invalid token."))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken("oldtoken"))
	require.NoError(t, err)

	assert.Equal(t, "oldtoken", client.GetToken())
	_, err = client.RateLimitedMethods().IsValidToken(context.Background())
	assert.ErrorIs(t, err, errors.MissingCredentialsError)
}

func TestRenewTokenPropagatesErrorOnBadCredentials(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.NoError(t, req.ParseMultipartForm(5000000))
		switch req.URL.Path {
		case "/isValidToken":
			expectRequest(t, "isValidToken").sendError("Invalid token.").ServeHTTP(rw, req)
		case "/getToken":
			expectRequest(t, "getToken", "username", "the-user-id", "password", "the-password").sendJson(`{
 "Authenticated": "False",
 "UserToken": ""
}`).ServeHTTP(rw, req)
		}
	}))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken("oldtoken"), raildata.WithCredentials("the-user-id", "the-password"))
	require.NoError(t, err)

	assert.Equal(t, "oldtoken", client.GetToken())
	_, err = client.RateLimitedMethods().IsValidToken(context.Background())
	assert.ErrorIs(t, err, errors.BadCredentialsError)
}

func TestRenewTokenPropagatesErrorOnLimitExceeded(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.NoError(t, req.ParseMultipartForm(5000000))
		switch req.URL.Path {
		case "/isValidToken":
			expectRequest(t, "isValidToken").sendError("Invalid token.").ServeHTTP(rw, req)
		case "/getToken":
			expectRequest(t, "getToken", "username", "the-user-id", "password", "the-password").sendError("Daily usage limit:10. Your current daily usage: 11").ServeHTTP(rw, req)
		}
	}))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken("oldtoken"), raildata.WithCredentials("the-user-id", "the-password"))
	require.NoError(t, err)

	assert.Equal(t, "oldtoken", client.GetToken())
	_, err = client.RateLimitedMethods().IsValidToken(context.Background())
	var rderr *errors.RailDataError
	assert.ErrorAs(t, err, &rderr)
	assert.Equal(t, "Daily usage limit:10. Your current daily usage: 11", rderr.Error())
}

func TestGetStationList(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "getStationList").sendJson(`[
  {
    "STATION_2CHAR": "AM",
    "STATIONNAME": "Aberdeen-Matawan",
    "STATION_14CHAR": "Matawan"
  },
  {
    "STATION_2CHAR": "AB",
    "STATIONNAME": "Absecon",
    "STATION_14CHAR": "Absecon"
  }
]`))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	actual, err := client.GetStationList(context.Background())
	assert.NoError(t, err)
	expected := &raildata.GetStationListResponse{
		Stations: []raildata.Station{
			{Code: "AM", Name: "Aberdeen-Matawan", ShortName: "Matawan"},
			{Code: "AB", Name: "Absecon", ShortName: "Absecon"},
		},
	}
	assert.Equal(t, expected, actual)
}

func TestGetStationMsgAllStationsAndLines(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "getStationMSG", "station", "", "line", "").sendJson(`[
  {
    "MSG_TYPE": "banner",
    "MSG_TEXT": "Text of the message.",
    "MSG_PUBDATE": "1/17/2025 2:40:00 PM",
    "MSG_ID": "123456789",
    "MSG_AGENCY": "NJT",
    "MSG_SOURCE": "RSS_NJTRailAlerts",
    "MSG_STATION_SCOPE": "*Newark Penn Station,*New York Penn Station",
    "MSG_LINE_SCOPE": "*Montclair-Boonton Line",
    "MSG_PUBDATE_UTC": "1/17/2025 7:40:00 PM"
  },
  {
    "MSG_TYPE": "fullscreen",
    "MSG_TEXT": "Another message.",
    "MSG_PUBDATE": "12/3/2024 10:52:10 AM",
    "MSG_ID": "",
    "MSG_AGENCY": "",
    "MSG_SOURCE": "",
    "MSG_STATION_SCOPE": " ",
    "MSG_LINE_SCOPE": " ",
    "MSG_PUBDATE_UTC": "12/3/2024 3:52:10 PM"
  }
]`))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	req := &raildata.GetStationMsgRequest{}
	actual, err := client.GetStationMsg(context.Background(), req)
	assert.NoError(t, err)
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	expected := &raildata.GetStationMsgResponse{
		Messages: []raildata.StationMsg{
			{
				Type:    raildata.MsgTypeBanner,
				Text:    "Text of the message.",
				PubDate: time.Date(2025, time.January, 17, 14, 40, 0, 0, loc),
				Id:      ptr("123456789"),
				Agency:  ptr("NJT"),
				Source:  ptr("RSS_NJTRailAlerts"),
				StationScope: []raildata.Station{
					{Code: "NP", Name: "Newark Penn Station", ShortName: "Newark Penn"},
					{Code: "NY", Name: "New York Penn Station", ShortName: "New York"},
				},
				LineScope: []raildata.Line{
					{Code: "MC", Name: "Montclair-Boonton Line", Abbreviation: "MOBO"},
				},
			},
			{
				Type:    raildata.MsgTypeFullScreen,
				Text:    "Another message.",
				PubDate: time.Date(2024, time.December, 3, 10, 52, 10, 0, loc),
			},
		},
	}
	assert.Equal(t, expected, actual)
}

func TestGetStationMsgOneStationOneLine(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "getStationMSG", "station", "NY", "line", "MC").sendJson(`[
  {
    "MSG_TYPE": "fullscreen",
    "MSG_TEXT": "A message.",
    "MSG_PUBDATE": "12/3/2024 10:52:10 AM",
    "MSG_ID": "",
    "MSG_AGENCY": "",
    "MSG_SOURCE": "",
    "MSG_STATION_SCOPE": " ",
    "MSG_LINE_SCOPE": " ",
    "MSG_PUBDATE_UTC": "12/3/2024 3:52:10 PM"
  }
]`))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	req := &raildata.GetStationMsgRequest{
		StationCode: (*raildata.StationCode)(ptr("NY")),
		LineCode:    (*raildata.LineCode)(ptr("MC")),
	}
	actual, err := client.GetStationMsg(context.Background(), req)
	assert.NoError(t, err)
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	expected := &raildata.GetStationMsgResponse{
		Messages: []raildata.StationMsg{
			{
				Type:    raildata.MsgTypeFullScreen,
				Text:    "A message.",
				PubDate: time.Date(2024, time.December, 3, 10, 52, 10, 0, loc),
			},
		},
	}
	assert.Equal(t, expected, actual)
}

func TestGetStationMsgNoMessages(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "getStationMSG").sendJson(`[]`))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	req := &raildata.GetStationMsgRequest{}
	actual, err := client.GetStationMsg(context.Background(), req)
	assert.NoError(t, err)
	expected := &raildata.GetStationMsgResponse{}
	assert.Equal(t, expected, actual)
}

func TestGetStationSchedule(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "getStationSchedule", "station", "UM").sendJson(`[
  {
    "STATION_2CHAR": "UM",
    "STATIONNAME": "Upp. Montclair",
    "ITEMS": [
      {
        "SCHED_DEP_DATE": "17-Jan-2025 12:24:45 AM",
        "DESTINATION": "MSU",
        "TRACK": "Montclair-Boonton Line",
        "LINE": "Montclair-Boonton Line",
        "TRAIN_ID": "6299",
        "CONNECTING_TRAIN_ID": "",
        "STATION_POSITION": "1",
        "DIRECTION": "Westbound",
        "DWELL_TIME": "60",
        "PERM_PICKUP": "",
        "PERM_DROPOFF": "",
        "STOP_CODE": "S"
      },
      {
        "SCHED_DEP_DATE": "17-Jan-2025 09:19:00 PM",
        "DESTINATION": "New York -SEC",
        "TRACK": "Montclair-Boonton Line",
        "LINE": "Montclair-Boonton Line",
        "TRAIN_ID": "6274",
        "CONNECTING_TRAIN_ID": "",
        "STATION_POSITION": "0",
        "DIRECTION": "Eastbound",
        "DWELL_TIME": "45",
        "PERM_PICKUP": "True",
        "PERM_DROPOFF": "",
        "STOP_CODE": "S*"
      }
	]
  }
]`))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	req := &raildata.GetStationScheduleRequest{StationCode: "UM"}
	actual, err := client.RateLimitedMethods().GetStationSchedule(context.Background(), req)
	assert.NoError(t, err)
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	expected := &raildata.GetStationScheduleResponse{
		Entries: []raildata.StationSchedule{
			{
				Station: &raildata.Station{Code: "UM", Name: "Upper Montclair", ShortName: "Upp. Montclair"},
				Entries: []raildata.ScheduleEntry{
					{
						DepartureTime:      time.Date(2025, time.January, 17, 0, 24, 45, 0, loc),
						Destination:        "MSU",
						DestinationStation: &raildata.Station{Code: "UV", Name: "Montclair State U", ShortName: "MSU"},
						Line:               raildata.Line{Code: "MC", Name: "Montclair-Boonton Line", Abbreviation: "MOBO"},
						TrainId:            "6299",
						StationPosition:    raildata.StationPositions[1],
						Direction:          raildata.DirectionWestbound,
						DwellTime:          ptr(60 * time.Second),
						StopCode:           &raildata.StopCode{Code: "S", Description: "Normal Stop"},
					},
					{
						DepartureTime:      time.Date(2025, time.January, 17, 21, 19, 0, 0, loc),
						Destination:        "New York -SEC",
						DestinationStation: &raildata.Station{Code: "NY", Name: "New York Penn Station", ShortName: "New York"},
						Line:               raildata.Line{Code: "MC", Name: "Montclair-Boonton Line", Abbreviation: "MOBO"},
						TrainId:            "6274",
						StationPosition:    raildata.StationPositions[0],
						Direction:          raildata.DirectionEastbound,
						DwellTime:          ptr(45 * time.Second),
						PickupOnly:         true,
						StopCode:           &raildata.StopCode{Code: "S*", Description: "Normal stop. May leave up to 3 minutes early"},
					},
				},
			},
		},
	}
	assert.Equal(t, expected, actual)
}

func TestGetStationScheduleEmptyStation(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "getStationSchedule", "station", "").sendEmpty())

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	req := &raildata.GetStationScheduleRequest{}
	_, err = client.RateLimitedMethods().GetStationSchedule(context.Background(), req)
	assert.ErrorIs(t, err, errors.MissingCredentialsError)
}

func TestGetTrainSchedule(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "getTrainSchedule", "station", "NY").sendJson(`{
  "STATION_2CHAR": "NY",
  "STATIONNAME": "New York",
  "STATIONMSGS": [
	{
      "MSG_TYPE": "banner",
      "MSG_TEXT": "A message",
      "MSG_PUBDATE": "12/21/2023 11:13:00 AM",
      "MSG_ID": "",
      "MSG_AGENCY": "",
      "MSG_SOURCE": "",
      "MSG_STATION_SCOPE": "*New York Penn Station,*Newark Airport",
      "MSG_LINE_SCOPE": " ",
      "MSG_PUBDATE_UTC": "12/21/2023 4:13:00 PM"
    }
  ],
  "ITEMS": [
    {
      "SCHED_DEP_DATE": "17-Jan-2025 09:00:00 PM",
      "DESTINATION": "Philadelphia &#9992",
      "TRACK": "9",
      "LINE": "AMTRAK",
      "TRAIN_ID": "A639",
      "CONNECTING_TRAIN_ID": "",
      "STATUS": "BOARDING",
      "SEC_LATE": "-60",
      "LAST_MODIFIED": "17-Jan-2025 08:48:13 PM",
      "BACKCOLOR": "#FFFF00",
      "FORECOLOR": "#000000",
      "SHADOWCOLOR": "#FFFF00",
      "GPSLATITUDE": "",
      "GPSLONGITUDE": "",
      "GPSTIME": "17-Jan-2025 08:03:59 PM",
      "STATION_POSITION": "0",
      "LINECODE": "AM",
      "LINEABBREVIATION": "AMTK",
      "INLINEMSG": "",
      "CAPACITY": [],
      "STOPS": [
        {
          "STATION_2CHAR": "NY",
          "STATIONNAME": "New York Penn Station",
          "TIME": "17-Jan-2025 09:00:00 PM",
          "PICKUP": "",
          "DROPOFF": "",
          "DEPARTED": "NO",
          "STOP_STATUS": "BOARDING",
          "DEP_TIME": "17-Jan-2025 09:00:00 PM",
          "TIME_UTC_FORMAT": "18-Jan-2025 02:00:00 AM",
          "STOP_LINES": []
        },
        {
          "STATION_2CHAR": "NP",
          "STATIONNAME": "Newark Penn Station",
          "TIME": "17-Jan-2025 09:16:00 PM",
          "PICKUP": "",
          "DROPOFF": "",
          "DEPARTED": "NO",
          "STOP_STATUS": "",
          "DEP_TIME": "17-Jan-2025 09:17:00 PM",
          "TIME_UTC_FORMAT": "18-Jan-2025 02:16:00 AM",
          "STOP_LINES": []
        }
      ]
    },
    {
      "SCHED_DEP_DATE": "17-Jan-2025 09:06:00 PM",
      "DESTINATION": "Trenton -SEC &#9992",
      "TRACK": "3",
      "LINE": "Northeast Corrdr",
      "TRAIN_ID": "3887",
      "CONNECTING_TRAIN_ID": "",
      "STATUS": "BOARDING",
      "SEC_LATE": "0",
      "LAST_MODIFIED": "17-Jan-2025 08:57:24 PM",
      "BACKCOLOR": "#F7505E",
      "FORECOLOR": "#FFFFFF",
      "SHADOWCOLOR": "#000000",
      "GPSLATITUDE": "",
      "GPSLONGITUDE": "",
      "GPSTIME": "",
      "STATION_POSITION": "0",
      "LINECODE": "NE",
      "LINEABBREVIATION": "NEC",
      "INLINEMSG": "",
      "CAPACITY": [],
      "STOPS": [
        {
          "STATION_2CHAR": "NY",
          "STATIONNAME": "New York Penn Station",
          "TIME": "17-Jan-2025 09:06:00 PM",
          "PICKUP": "",
          "DROPOFF": "",
          "DEPARTED": "NO",
          "STOP_STATUS": "BOARDING",
          "DEP_TIME": "17-Jan-2025 09:06:00 PM",
          "TIME_UTC_FORMAT": "18-Jan-2025 02:06:00 AM",
          "STOP_LINES": []
        },
        {
          "STATION_2CHAR": "SE",
          "STATIONNAME": "Secaucus Upper Lvl",
          "TIME": "17-Jan-2025 09:15:00 PM",
          "PICKUP": "",
          "DROPOFF": "",
          "DEPARTED": "NO",
          "STOP_STATUS": "",
          "DEP_TIME": "17-Jan-2025 09:15:30 PM",
          "TIME_UTC_FORMAT": "18-Jan-2025 02:15:00 AM",
          "STOP_LINES": []
        }
      ]
    }
  ]
}`))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	req := &raildata.GetTrainScheduleRequest{StationCode: "NY"}
	actual, err := client.GetTrainSchedule(context.Background(), req)
	assert.NoError(t, err)
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	expected := &raildata.GetTrainScheduleResponse{
		Station: raildata.Station{Code: "NY", Name: "New York Penn Station", ShortName: "New York"},
		Messages: []raildata.StationMsg{
			{
				Type:    raildata.MsgTypeBanner,
				Text:    "A message",
				PubDate: time.Date(2023, time.December, 21, 11, 13, 0, 0, loc),
				StationScope: []raildata.Station{
					{Code: "NY", Name: "New York Penn Station", ShortName: "New York"},
					{Code: "NA", Name: "Newark Airport", ShortName: "Newark Airport"},
				},
			},
		},
		Entries: []raildata.TrainScheduleEntry{
			{
				DepartureTime:      time.Date(2025, time.January, 17, 21, 0, 0, 0, loc),
				Destination:        "Philadelphia ✈",
				DestinationStation: &raildata.Station{Code: "PH", Name: "Philadelphia", ShortName: "Philadelphia"},
				Track:              ptr("9"),
				Line:               raildata.Line{Code: "AM", Name: "Amtrak", Abbreviation: "AMTK"},
				TrainId:            "A639",
				Status:             ptr("BOARDING"),
				Delay:              ptr(-60 * time.Second),
				LastUpdated:        ptr(time.Date(2025, time.January, 17, 20, 48, 13, 0, loc)),
				Color:              raildata.ColorSet{Background: color(t, "#FFFF00"), Foreground: color(t, "#000000"), Shadow: color(t, "#FFFF00")},
				GpsTime:            ptr(time.Date(2025, time.January, 17, 20, 3, 59, 0, loc)),
				StationPosition:    raildata.StationPositions[0],
				Stops: []raildata.TrainStop{
					{
						Station:       raildata.Station{Code: "NY", Name: "New York Penn Station", ShortName: "New York"},
						ArrivalTime:   ptr(time.Date(2025, time.January, 17, 21, 0, 0, 0, loc)),
						StopStatus:    ptr("BOARDING"),
						DepartureTime: ptr(time.Date(2025, time.January, 17, 21, 0, 0, 0, loc)),
					},
					{
						Station:       raildata.Station{Code: "NP", Name: "Newark Penn Station", ShortName: "Newark Penn"},
						ArrivalTime:   ptr(time.Date(2025, time.January, 17, 21, 16, 0, 0, loc)),
						DepartureTime: ptr(time.Date(2025, time.January, 17, 21, 17, 0, 0, loc)),
					},
				},
			},
			{
				DepartureTime:      time.Date(2025, time.January, 17, 21, 6, 0, 0, loc),
				Destination:        "Trenton -SEC ✈",
				DestinationStation: &raildata.Station{Code: "TR", Name: "Trenton", ShortName: "Trenton"},
				Track:              ptr("3"),
				Line:               raildata.Line{Code: "NE", Name: "Northeast Corridor Line", Abbreviation: "NEC"},
				TrainId:            "3887",
				Status:             ptr("BOARDING"),
				Delay:              ptr(0 * time.Second),
				LastUpdated:        ptr(time.Date(2025, time.January, 17, 20, 57, 24, 0, loc)),
				Color:              raildata.ColorSet{Background: color(t, "#F7505E"), Foreground: color(t, "#FFFFFF"), Shadow: color(t, "#000000")},
				StationPosition:    raildata.StationPositions[0],
				Stops: []raildata.TrainStop{
					{
						Station:       raildata.Station{Code: "NY", Name: "New York Penn Station", ShortName: "New York"},
						ArrivalTime:   ptr(time.Date(2025, time.January, 17, 21, 6, 0, 0, loc)),
						StopStatus:    ptr("BOARDING"),
						DepartureTime: ptr(time.Date(2025, time.January, 17, 21, 6, 0, 0, loc)),
					},
					{
						Station:       raildata.Station{Code: "SE", Name: "Secaucus Upper Lvl", ShortName: "Secaucus"},
						ArrivalTime:   ptr(time.Date(2025, time.January, 17, 21, 15, 0, 0, loc)),
						DepartureTime: ptr(time.Date(2025, time.January, 17, 21, 15, 30, 0, loc)),
					},
				},
			},
		},
	}
	assert.Equal(t, expected, actual)
}

func TestGetTrainSchedule19Rec(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "getTrainSchedule19Rec", "station", "NY", "line", "NE").sendJson(`{
  "STATION_2CHAR": "NY",
  "STATIONNAME": "New York",
  "STATIONMSGS": [],
  "ITEMS": [
    {
      "SCHED_DEP_DATE": "17-Jan-2025 09:35:00 PM",
      "DESTINATION": "Trenton -SEC &#9992",
      "TRACK": "",
      "LINE": "Northeast Corrdr",
      "TRAIN_ID": "3889",
      "CONNECTING_TRAIN_ID": "",
      "STATUS": " ",
      "SEC_LATE": "-60",
      "LAST_MODIFIED": "17-Jan-2025 08:53:11 PM",
      "BACKCOLOR": "#F7505E",
      "FORECOLOR": "#FFFFFF",
      "SHADOWCOLOR": "#000000",
      "GPSLATITUDE": "",
      "GPSLONGITUDE": "",
      "GPSTIME": "17-Jan-2025 08:51:59 PM",
      "STATION_POSITION": "0",
      "LINECODE": "NE",
      "LINEABBREVIATION": "NEC",
      "INLINEMSG": "",
      "CAPACITY": [],
      "STOPS": null
    },
    {
      "SCHED_DEP_DATE": "17-Jan-2025 10:07:00 PM",
      "DESTINATION": "Jersey Avenue -SEC &#9992",
      "TRACK": "",
      "LINE": "Northeast Corrdr",
      "TRAIN_ID": "3737",
      "CONNECTING_TRAIN_ID": "",
      "STATUS": " ",
      "SEC_LATE": "-60",
      "LAST_MODIFIED": "17-Jan-2025 09:08:29 PM",
      "BACKCOLOR": "#F7505E",
      "FORECOLOR": "#FFFFFF",
      "SHADOWCOLOR": "#000000",
      "GPSLATITUDE": "",
      "GPSLONGITUDE": "",
      "GPSTIME": "17-Jan-2025 09:07:25 PM",
      "STATION_POSITION": "0",
      "LINECODE": "NE",
      "LINEABBREVIATION": "NEC",
      "INLINEMSG": "",
      "CAPACITY": [
        {
          "VEHICLE_NO": "3737",
          "LATITUDE": "40.5734680000",
          "LONGITUDE": "-74.3202910000",
          "CREATED_TIME": "17-Jan-2025 09:26:22 PM",
          "VEHICLE_TYPE": "1",
          "CUR_PERCENTAGE": "1",
          "CUR_CAPACITY_COLOR": "#0B6623",
          "CUR_PASSENGER_COUNT": "6",
          "PREV_PERCENTAGE": "1",
          "PREV_CAPACITY_COLOR": " #0B6623",
          "PREV_PASSENGER_COUNT": "6",
          "SECTIONS": [
            {
              "SECTION_POSITION": "Back",
              "CUR_PERCENTAGE": "0",
              "CUR_CAPACITY_COLOR": "#0B6623",
              "CUR_PASSENGER_COUNT": "2",
              "CARS": [
                {
                  "CAR_NO": "7689",
                  "CAR_POSITION": "2",
                  "CAR_REST": true,
                  "CUR_PERCENTAGE": "0",
                  "CUR_CAPACITY_COLOR": "#0B6623",
                  "CUR_PASSENGER_COUNT": "2"
                }
              ]
            },
            {
              "SECTION_POSITION": "Front",
              "CUR_PERCENTAGE": "1",
              "CUR_CAPACITY_COLOR": "#0B6623",
              "CUR_PASSENGER_COUNT": "4",
              "CARS": [
                {
                  "CAR_NO": "7657",
                  "CAR_POSITION": "1",
                  "CAR_REST": false,
                  "CUR_PERCENTAGE": "1",
                  "CUR_CAPACITY_COLOR": "#0B6623",
                  "CUR_PASSENGER_COUNT": "4"
                }
              ]
            }
          ]
        }
      ],
      "STOPS": null
    }
  ]
}`))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	req := &raildata.GetTrainSchedule19RecordsRequest{StationCode: "NY", LineCode: ptr(raildata.LineCode("NE"))}
	actual, err := client.GetTrainSchedule19Records(context.Background(), req)
	assert.NoError(t, err)
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	expected := &raildata.GetTrainScheduleResponse{
		Station: raildata.Station{Code: "NY", Name: "New York Penn Station", ShortName: "New York"},
		Entries: []raildata.TrainScheduleEntry{
			{
				DepartureTime:      time.Date(2025, time.January, 17, 21, 35, 0, 0, loc),
				Destination:        "Trenton -SEC ✈",
				DestinationStation: &raildata.Station{Code: "TR", Name: "Trenton", ShortName: "Trenton"},
				Line:               raildata.Line{Code: "NE", Name: "Northeast Corridor Line", Abbreviation: "NEC"},
				TrainId:            "3889",
				Delay:              ptr(-60 * time.Second),
				LastUpdated:        ptr(time.Date(2025, time.January, 17, 20, 53, 11, 0, loc)),
				Color:              raildata.ColorSet{Background: color(t, "#F7505E"), Foreground: color(t, "#FFFFFF"), Shadow: color(t, "#000000")},
				GpsTime:            ptr(time.Date(2025, time.January, 17, 20, 51, 59, 0, loc)),
				StationPosition:    raildata.StationPositions[0],
			},
			{
				DepartureTime:      time.Date(2025, time.January, 17, 22, 7, 0, 0, loc),
				Destination:        "Jersey Avenue -SEC ✈",
				DestinationStation: &raildata.Station{Code: "JA", Name: "Jersey Avenue", ShortName: "Jersey Ave."},
				Line:               raildata.Line{Code: "NE", Name: "Northeast Corridor Line", Abbreviation: "NEC"},
				TrainId:            "3737",
				Delay:              ptr(-60 * time.Second),
				LastUpdated:        ptr(time.Date(2025, time.January, 17, 21, 8, 29, 0, loc)),
				Color:              raildata.ColorSet{Background: color(t, "#F7505E"), Foreground: color(t, "#FFFFFF"), Shadow: color(t, "#000000")},
				GpsTime:            ptr(time.Date(2025, time.January, 17, 21, 7, 25, 0, loc)),
				StationPosition:    raildata.StationPositions[0],
				Capacity: []raildata.TrainCapacity{
					{
						Number: "3737",
						Location: raildata.Location{
							Longitude: -74.320291,
							Latitude:  40.573468,
						},
						CreatedTime:     time.Date(2025, time.January, 17, 21, 26, 22, 0, loc),
						Type:            "1",
						CapacityPercent: 1,
						CapacityColor:   color(t, "#0B6623"),
						PassengerCount:  6,
						Sections: []raildata.TrainSection{
							{
								Position:        raildata.SectionPositionBack,
								CapacityPercent: 0,
								CapacityColor:   color(t, "#0B6623"),
								PassengerCount:  2,
								Cars: []raildata.TrainCar{
									{
										TrainId:         "7689",
										Position:        2,
										Restroom:        true,
										CapacityPercent: 0,
										CapacityColor:   color(t, "#0B6623"),
										PassengerCount:  2,
									},
								},
							},
							{
								Position:        raildata.SectionPositionFront,
								CapacityPercent: 1,
								CapacityColor:   color(t, "#0B6623"),
								PassengerCount:  4,
								Cars: []raildata.TrainCar{
									{
										TrainId:         "7657",
										Position:        1,
										Restroom:        false,
										CapacityPercent: 1,
										CapacityColor:   color(t, "#0B6623"),
										PassengerCount:  4,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	assert.Equal(t, expected, actual)
}

func TestGetTrainStopList(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "getTrainStopList", "train", "3737").sendJson(`{
  "TRAIN_ID": "3737",
  "LINECODE": "NE",
  "BACKCOLOR": "#F7505E",
  "FORECOLOR": "#FFFFFF",
  "SHADOWCOLOR": "#000000",
  "DESTINATION": "Jersey Avenue",
  "TRANSFERAT": "",
  "STOPS": [
    {
      "STATION_2CHAR": "NY",
      "STATIONNAME": "New York Penn Station",
      "TIME": "17-Jan-2025 10:07:00 PM",
      "PICKUP": "",
      "DROPOFF": "",
      "DEPARTED": "NO",
      "STOP_STATUS": "OnTime",
      "DEP_TIME": "17-Jan-2025 10:07:00 PM",
      "TIME_UTC_FORMAT": "18-Jan-2025 03:07:00 AM",
      "STOP_LINES": [
        {
          "LINE_CODE": "GS",
          "LINE_NAME": "Gladstone Branch",
          "LINE_COLOR": "#A1D5AE"
        },
        {
          "LINE_CODE": "MC",
          "LINE_NAME": "MontClair-Boonton Line",
          "LINE_COLOR": "#C36366"
        }
      ]
    },
    {
      "STATION_2CHAR": "SE",
      "STATIONNAME": "Secaucus Upper Lvl",
      "TIME": "17-Jan-2025 10:16:00 PM",
      "PICKUP": "",
      "DROPOFF": "",
      "DEPARTED": "NO",
      "STOP_STATUS": "",
      "DEP_TIME": "17-Jan-2025 10:17:00 PM",
      "TIME_UTC_FORMAT": "18-Jan-2025 03:16:00 AM",
      "STOP_LINES": []
    }
  ]
}`))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	req := &raildata.GetTrainStopListRequest{TrainId: "3737"}
	actual, err := client.GetTrainStopList(context.Background(), req)
	assert.NoError(t, err)
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	expected := &raildata.GetTrainStopListResponse{
		TrainId:            "3737",
		Line:               raildata.Line{Code: "NE", Name: "Northeast Corridor Line", Abbreviation: "NEC"},
		Color:              raildata.ColorSet{Background: color(t, "#F7505E"), Foreground: color(t, "#FFFFFF"), Shadow: color(t, "#000000")},
		Destination:        "Jersey Avenue",
		DestinationStation: &raildata.Station{Code: "JA", Name: "Jersey Avenue", ShortName: "Jersey Ave."},
		Stops: []raildata.TrainStop{
			{
				Station:       raildata.Station{Code: "NY", Name: "New York Penn Station", ShortName: "New York"},
				ArrivalTime:   ptr(time.Date(2025, time.January, 17, 22, 7, 0, 0, loc)),
				StopStatus:    ptr("OnTime"),
				DepartureTime: ptr(time.Date(2025, time.January, 17, 22, 7, 0, 0, loc)),
				StopLines: []raildata.StopLine{
					{
						Line:  raildata.Line{Code: "GS", Name: "Gladstone Branch", Abbreviation: "M&E"},
						Color: color(t, "#A1D5AE"),
					},
					{
						Line:  raildata.Line{Code: "MC", Name: "Montclair-Boonton Line", Abbreviation: "MOBO"},
						Color: color(t, "#C36366"),
					},
				},
			},
			{
				Station:       raildata.Station{Code: "SE", Name: "Secaucus Upper Lvl", ShortName: "Secaucus"},
				ArrivalTime:   ptr(time.Date(2025, time.January, 17, 22, 16, 0, 0, loc)),
				DepartureTime: ptr(time.Date(2025, time.January, 17, 22, 17, 0, 0, loc)),
			},
		},
	}
	assert.Equal(t, expected, actual)
}

func TestGetTrainStopListForInvalidTrain(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "getTrainStopList", "train", "3737").sendJson(`{
  "TRAIN_ID": null,
  "LINECODE": null,
  "BACKCOLOR": null,
  "FORECOLOR": null,
  "SHADOWCOLOR": null,
  "DESTINATION": null,
  "TRANSFERAT": null,
  "STOPS": null,
  "CAPACITY": null
}`))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	req := &raildata.GetTrainStopListRequest{TrainId: "3737"}
	actual, err := client.GetTrainStopList(context.Background(), req)
	assert.NoError(t, err)
	assert.Nil(t, actual)
}

func TestGetVehicleData(t *testing.T) {
	server := httptest.NewServer(expectRequest(t, "getVehicleData").sendJson(`[
  {
    "ID": "65",
    "TRAIN_LINE": "Bergen County Line",
    "DIRECTION": "Westbound",
    "ICS_TRACK_CKT": "OV-7112TK",
    "LAST_MODIFIED": "17-Jan-2025 09:49:33 PM",
    "SCHED_DEP_TIME": "17-Jan-2025 09:52:00 PM",
    "SEC_LATE": "98",
    "NEXT_STOP": "Otisville",
    "LONGITUDE": "-74.529233",
    "LATITUDE": "41.471769"
  },
  {
    "ID": "68",
    "TRAIN_LINE": "Bergen County Line",
    "DIRECTION": "Eastbound",
    "ICS_TRACK_CKT": "HO-7021TK",
    "LAST_MODIFIED": "17-Jan-2025 09:49:19 PM",
    "SCHED_DEP_TIME": "17-Jan-2025 09:57:00 PM",
    "SEC_LATE": "-60",
    "NEXT_STOP": "Middletown NY",
    "LONGITUDE": "-74.370446",
    "LATITUDE": "41.457426"
  }
]`))

	client, err := raildata.NewClient(withServerUrl(t, server), raildata.WithToken(testToken))
	require.NoError(t, err)

	actual, err := client.GetVehicleData(context.Background())
	assert.NoError(t, err)
	loc, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)
	expected := &raildata.GetVehicleDataResponse{
		Vehicles: []raildata.VehicleData{
			{
				TrainId:        "65",
				Line:           raildata.Line{Code: "BC", Name: "Bergen County Line", Abbreviation: "BERG"},
				Direction:      raildata.DirectionWestbound,
				TrackCircuitId: "OV-7112TK",
				LastUpdated:    time.Date(2025, time.January, 17, 21, 49, 33, 0, loc),
				DepartureTime:  time.Date(2025, time.January, 17, 21, 52, 0, 0, loc),
				Delay:          ptr(98 * time.Second),
				NextStop:       &raildata.Station{Code: "OS", Name: "Otisville", ShortName: "Otisville"},
				Location:       &raildata.Location{Longitude: -74.529233, Latitude: 41.471769},
			},
			{
				TrainId:        "68",
				Line:           raildata.Line{Code: "BC", Name: "Bergen County Line", Abbreviation: "BERG"},
				Direction:      raildata.DirectionEastbound,
				TrackCircuitId: "HO-7021TK",
				LastUpdated:    time.Date(2025, time.January, 17, 21, 49, 19, 0, loc),
				DepartureTime:  time.Date(2025, time.January, 17, 21, 57, 0, 0, loc),
				Delay:          ptr(-60 * time.Second),
				NextStop:       &raildata.Station{Code: "MD", Name: "Middletown NY", ShortName: "Middletown NY"},
				Location:       &raildata.Location{Longitude: -74.370446, Latitude: 41.457426},
			},
		},
	}
	assert.Equal(t, expected, actual)
}

func expectRequest(t *testing.T, path string, fields ...string) expectingRequest {
	out := expectingRequest{
		t:      t,
		path:   "/" + path,
		fields: map[string]string{},
	}
	for i := 0; i < len(fields); i += 2 {
		out.fields[fields[i]] = fields[i+1]
	}
	return out
}

type expectingRequest struct {
	t      *testing.T
	path   string
	fields map[string]string
}

func (e expectingRequest) sendError(msg string) http.HandlerFunc {
	return e.sendResponse(500, fmt.Sprintf(`{
		"errorMessage": "%s"
	   }`, msg))
}

func (e expectingRequest) sendJson(obj string) http.HandlerFunc {
	return e.sendResponse(200, obj)
}

func (e expectingRequest) sendEmpty() http.HandlerFunc {
	return e.sendResponse(204, "")
}

func (e expectingRequest) sendResponse(statusCode int, body string) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.NoError(e.t, req.ParseMultipartForm(5000000))
		require.Equal(e.t, e.path, req.URL.Path)
		for k, v := range e.fields {
			assert.Equal(e.t, v, req.Form.Get(k))
		}
		rw.WriteHeader(statusCode)
		rw.Write([]byte(body))
	})
}

func withServerUrl(t *testing.T, server *httptest.Server) raildata.Option {
	u, err := url.Parse(server.URL)
	require.NoError(t, err)
	return raildata.WithApiBase(*u)
}

func color(t *testing.T, s string) raildata.Color {
	c, err := raildata.ParseHtmlColor(s)
	require.NoError(t, err)
	return c
}

func ptr[T any](o T) *T {
	return &o
}
