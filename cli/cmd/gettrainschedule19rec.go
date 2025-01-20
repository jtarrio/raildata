package cmd

import (
	"context"
	"fmt"

	"github.com/jtarrio/raildata"
	"github.com/jtarrio/raildata/raildata-cli/util"
	"github.com/urfave/cli/v2"
)

var cmdGetTrainSchedule19Rec = &cli.Command{
	Name:  "getTrainSchedule19Rec",
	Usage: "gets schedule data for one station",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "station",
			Usage:    "code or name of a station to get the schedule for",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "line",
			Usage: "code or name of a line to get the schedule for",
		},
	},
	Action: func(ctx *cli.Context) error {
		return getTrainSchedule19Rec(ctx.Context, ctx.String("station"), ctx.String("line"))
	},
}

func getTrainSchedule19Rec(ctx context.Context, station string, line string) error {
	req := &raildata.GetTrainSchedule19RecordsRequest{}
	stationCode, found := util.FindStation(station)
	if !found {
		return fmt.Errorf("station '%s' unknown", station)
	}
	req.StationCode = *stationCode
	if len(line) > 0 {
		lineCode, found := util.FindLine(line)
		if !found {
			return fmt.Errorf("line '%s' unknown", line)
		}
		req.LineCode = lineCode
	}
	client := GetClientFromContext(ctx)
	resp, err := client.GetTrainSchedule19Records(ctx, req)
	if err != nil {
		return err
	}

	fmt.Println(resp.Station.Name)
	for i := range resp.Messages {
		if i > 0 {
			fmt.Println()
		}
		displayMessage(&resp.Messages[i])
	}
	for i := range resp.Entries {
		if i > 0 || len(resp.Messages) > 0 {
			fmt.Println()
		}
		displayTrainScheduleEntry(&resp.Entries[i])
	}

	return nil
}
