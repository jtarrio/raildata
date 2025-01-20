package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/jtarrio/raildata"
	"github.com/jtarrio/raildata/raildata-cli/util"
	"github.com/urfave/cli/v2"
)

var cmdGetStationSchedule = &cli.Command{
	Name:  "getStationSchedule",
	Usage: "gets schedule data for one stations",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "station",
			Usage:    "code or name of a station to get the schedule for. When omitted, all stations are queried",
			Required: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		return getStationSchedule(ctx.Context, ctx.String("station"))
	},
}

func getStationSchedule(ctx context.Context, station string) error {
	req := &raildata.GetStationScheduleRequest{}
	stationCode, found := util.FindStation(station)
	if !found {
		return fmt.Errorf("station '%s' unknown", station)
	}
	req.StationCode = *stationCode
	client := GetClientFromContext(ctx)
	resp, err := client.RateLimitedMethods().GetStationSchedule(ctx, req)
	if err != nil {
		return err
	}
	for i := range resp.Entries {
		if i > 0 {
			fmt.Println()
		}
		displayStationSchedule(&resp.Entries[i])
	}
	return nil
}

func displayStationSchedule(sched *raildata.StationSchedule) {
	if sched.Station != nil {
		fmt.Printf("Station: %s\n", sched.Station.Name)
	}
	for i := range sched.Entries {
		if i > 0 {
			fmt.Println()
		}
		displayScheduleEntry(&sched.Entries[i])
	}
}

func displayScheduleEntry(entry *raildata.ScheduleEntry) {
	fmt.Printf("%s %s\n", entry.DepartureTime.Format(time.DateTime), entry.Line.Name)
	fmt.Printf("Train %s for %s (", entry.TrainId, entry.Destination)
	if entry.StationPosition.Code != "1" {
		fmt.Printf("%s ", entry.StationPosition.Description)
	}
	if entry.Direction == raildata.DirectionEastbound {
		fmt.Print("eastbound")
	} else {
		fmt.Print("westbound")
	}
	fmt.Print(")")
	if entry.PickupOnly {
		fmt.Print(" (pick-up only)")
	}
	if entry.DropoffOnly {
		fmt.Print(" (drop-off only)")
	}
	fmt.Println()
	if entry.ConnectingTrainId != nil {
		fmt.Printf("Connecting train: %s\n", *entry.ConnectingTrainId)
	}
	if entry.StopCode != nil && entry.StopCode.Code != "S" {
		fmt.Println(entry.StopCode.Description)
	}
}
