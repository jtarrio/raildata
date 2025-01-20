package cmd

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/fatih/color"
	"github.com/jtarrio/raildata"
	"github.com/jtarrio/raildata/raildata-cli/util"
	"github.com/urfave/cli/v2"
)

var cmdGetTrainSchedule = &cli.Command{
	Name:  "getTrainSchedule",
	Usage: "gets schedule data for one station",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "station",
			Usage:    "code or name of a station to get the schedule for",
			Required: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		return getTrainSchedule(ctx.Context, ctx.String("station"))
	},
}

func getTrainSchedule(ctx context.Context, station string) error {
	req := &raildata.GetTrainScheduleRequest{}
	stationCode, found := util.FindStation(station)
	if !found {
		return fmt.Errorf("station '%s' unknown", station)
	}
	req.StationCode = *stationCode
	client := GetClientFromContext(ctx)
	resp, err := client.GetTrainSchedule(ctx, req)
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

func displayTrainScheduleEntry(entry *raildata.TrainScheduleEntry) {
	fmt.Printf("%s ", entry.DepartureTime.Format(time.DateTime))
	util.HtmlColors(&entry.Color.Foreground, &entry.Color.Background).Print(entry.Line.Name)
	if entry.Status != nil {
		fmt.Printf(" %s", *entry.Status)
	}
	fmt.Printf("\nTrain %s for %s", entry.TrainId, entry.Destination)
	if entry.Track != nil {
		fmt.Printf(" on track %s", *entry.Track)
	}
	if entry.StationPosition.Code != "1" {
		fmt.Printf(" (%s)", entry.StationPosition.Description)
	}
	fmt.Println()
	if entry.ConnectingTrainId != nil {
		fmt.Printf("Connecting train: %s\n", *entry.ConnectingTrainId)
	}
	if entry.GpsLocation != nil {
		fmt.Printf("Train at %f,%f", entry.GpsLocation.Latitude, entry.GpsLocation.Longitude)
		if entry.GpsTime != nil {
			fmt.Printf(" (%s)", entry.GpsTime.Format(time.RFC1123))
		}
		fmt.Println()
	}
	if entry.Delay != nil {
		if *entry.Delay > 1*time.Minute {
			delayColor(*entry.Delay).Printf("Running late %s\n", *entry.Delay)
		} else if *entry.Delay < -1*time.Minute {
			color.New(color.FgCyan).Printf("Running early %s\n", -*entry.Delay)
		}
	}
	if len(entry.Stops) > 0 {
		displayStops(entry.Stops)
	}
	if len(entry.Capacity) > 0 {
		displayCapacity(entry.Capacity)
	}
	if entry.LastUpdated != nil {
		fmt.Printf("Last updated: %s\n", entry.LastUpdated.Format(time.RFC1123))
	}
}

func displayStops(stops []raildata.TrainStop) {
	fmt.Print("Stops: ")
	for i := range stops {
		if i > 0 {
			fmt.Print(", ")
		}
		stop := &stops[i]
		stopColor := color.New(color.Reset)
		if stop.Departed {
			stopColor = color.New(color.CrossedOut)
		}
		stopColor.Print(stop.Station.Name)
		if !stop.Departed {
			if stop.ArrivalTime != nil {
				fmt.Printf(" %s", stop.ArrivalTime.Format(time.TimeOnly))
			}
			if stop.DropoffOnly {
				fmt.Print(" (drop-off only)")
			}
			if stop.PickupOnly {
				fmt.Print(" (pick-up only)")
			}
			if stop.StopStatus != nil && *stop.StopStatus != "OnTime" {
				fmt.Printf(" %s", *stop.StopStatus)
			}
			if len(stop.StopLines) > 0 {
				for l := range stop.StopLines {
					fmt.Print(" ")
					sl := &stop.StopLines[l]
					util.HtmlColors(&sl.Color, nil).Print(sl.Line.Abbreviation)
				}
			}
		}
	}
	fmt.Println()
}

func displayCapacity(capacity []raildata.TrainCapacity) {
	fmt.Print("Capacity: ")
	for i := range capacity {
		if i > 0 {
			fmt.Print(", ")
		}
		cap := &capacity[i]
		fmt.Printf("Vehicle %s (%f, %f) %s\n",
			cap.Number, cap.Location.Latitude, cap.Location.Longitude,
			util.HtmlColors(&cap.CapacityColor, nil).Sprintf("%d%% full (%d pax)", cap.CapacityPercent, cap.PassengerCount))
		cars := []*raildata.TrainCar{}
		for s := range cap.Sections {
			sec := &cap.Sections[s]
			for c := range sec.Cars {
				cars = append(cars, &sec.Cars[c])
			}
		}
		slices.SortFunc(cars, func(a, b *raildata.TrainCar) int {
			return a.Position - b.Position
		})
		if len(cars) > 0 {
			fmt.Print("\tFront: ")
			for c, car := range cars {
				if c > 0 {
					fmt.Print(", ")
				}
				restroom := ""
				if car.Restroom {
					restroom = "🚻 "
				}
				fmt.Printf("%s %s%s", car.TrainId, restroom,
					util.HtmlColors(&car.CapacityColor, nil).Sprintf("%d%% full (%d pax)", car.CapacityPercent, car.PassengerCount))
			}
			fmt.Println(" :Back")
		}
	}
}
