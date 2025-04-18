package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/jtarrio/raildata"
	"github.com/urfave/cli/v2"
)

var cmdGetVehicleData = &cli.Command{
	Name:  "getVehicleData",
	Usage: "gets real-time position data for every active train",
	Action: func(ctx *cli.Context) error {
		return getVehicleData(ctx.Context)
	},
}

func getVehicleData(ctx context.Context) error {
	client := GetClientFromContext(ctx)
	resp, err := client.GetVehicleData(ctx)
	if err != nil {
		return err
	}

	for i := range resp.Vehicles {
		if i > 0 {
			fmt.Println()
		}
		veh := &resp.Vehicles[i]
		dir := "westbound"
		if veh.Direction == raildata.DirectionEastbound {
			dir = "eastbound"
		}
		fmt.Printf("Train %s on %s %s", veh.TrainId, veh.Line.Name, dir)
		if veh.Delay != nil {
			if *veh.Delay > 1*time.Minute {
				delayColor(*veh.Delay).Printf(" (%s late)", *veh.Delay)
			} else if *veh.Delay < -1*time.Minute {
				color.New(color.FgCyan).Printf(" (%s early)", -*veh.Delay)
			}
		}
		fmt.Println()
		if veh.Location != nil {
			fmt.Printf("Last position: %f,%f\n", veh.Location.Latitude, veh.Location.Longitude)
		}
		fmt.Printf("Next stop: %s\n", veh.NextStop.Name)
		fmt.Printf("Departing at %s\n", veh.DepartureTime.Format(time.RFC1123))

	}
	return nil
}

func delayColor(delay time.Duration) *color.Color {
	if delay < 5*time.Minute {
		return color.New(color.Reset)
	}
	if delay < 10*time.Minute {
		return color.New(color.FgYellow)
	}
	return color.New(color.FgHiRed)
}
