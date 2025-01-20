package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jtarrio/raildata"
	"github.com/jtarrio/raildata/raildata-cli/util"
	"github.com/urfave/cli/v2"
)

var cmdGetStationMsg = &cli.Command{
	Name:  "getStationMSG",
	Usage: "gets station messages",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "station",
			Usage: "code or name of a station to get messages for. When omitted, all stations are queried",
		},
		&cli.StringFlag{
			Name:  "line",
			Usage: "code or name of a line to get messages for. When omitted, all lines are queried",
		},
	},
	Action: func(ctx *cli.Context) error {
		return getStationMsg(ctx.Context, ctx.String("station"), ctx.String("line"))
	},
}

func getStationMsg(ctx context.Context, station string, line string) error {
	req := &raildata.GetStationMsgRequest{}
	if len(station) > 0 {
		stationCode, found := util.FindStation(station)
		if !found {
			return fmt.Errorf("station '%s' unknown", station)
		}
		req.StationCode = stationCode
	}
	if len(line) > 0 {
		lineCode, found := util.FindLine(line)
		if !found {
			return fmt.Errorf("line '%s' unknown", line)
		}
		req.LineCode = lineCode
	}
	client := GetClientFromContext(ctx)
	resp, err := client.GetStationMsg(ctx, req)
	if err != nil {
		return err
	}
	for i := range resp.Messages {
		if i > 0 {
			fmt.Println()
		}
		displayMessage(&resp.Messages[i])
	}
	return nil
}

func displayMessage(msg *raildata.StationMsg) {
	text := strings.TrimSpace(msg.Text)
	if msg.Type == raildata.MsgTypeFullScreen {
		fmt.Fprintf(color.Output, "%s\n", color.HiWhiteString("%s", text))
	} else {
		fmt.Printf("%s\n", text)
	}
	fmt.Printf("Posted on %s", msg.PubDate.Format(time.RFC1123))
	if msg.Agency != nil {
		fmt.Printf(" by %s", *msg.Agency)
	}
	if msg.Source != nil {
		fmt.Printf(" from %s", *msg.Source)
	}
	if msg.Id != nil {
		fmt.Printf(" (id %s)", *msg.Id)
	}
	fmt.Println()
	if len(msg.LineScope) > 0 {
		if len(msg.LineScope) == 1 {
			fmt.Print("For line: ")
		} else {
			fmt.Print("For lines: ")
		}
		for i, line := range msg.LineScope {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(line.Name)
		}
		fmt.Println()
	}
	if len(msg.StationScope) > 0 {
		if len(msg.StationScope) == 1 {
			fmt.Print("For station: ")
		} else {
			fmt.Print("For stations: ")
		}
		for i, station := range msg.StationScope {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(station.Name)
		}
		fmt.Println()
	}
}
