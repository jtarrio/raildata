package cmd

import (
	"context"
	"fmt"

	"github.com/jtarrio/raildata"
	"github.com/jtarrio/raildata/raildata-cli/util"
	"github.com/urfave/cli/v2"
)

var cmdGetTrainStopList = &cli.Command{
	Name:  "getTrainStopList",
	Usage: "gets train stops for a train",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "train",
			Usage:    "number of the train to get the stops for",
			Required: true,
		},
	},
	Action: func(ctx *cli.Context) error {
		return getTrainStopList(ctx.Context, ctx.String("train"))
	},
}

func getTrainStopList(ctx context.Context, train string) error {
	req := &raildata.GetTrainStopListRequest{
		TrainId: train,
	}
	client := GetClientFromContext(ctx)
	resp, err := client.GetTrainStopList(ctx, req)
	if err != nil {
		return err
	}
	if resp == nil {
		fmt.Printf("Train %s not found\n", train)
		return nil
	}

	util.HtmlColors(&resp.Color.Foreground, &resp.Color.Background).Print(resp.Line.Name)
	fmt.Printf("\nTrain %s for %s\n", resp.TrainId, resp.Destination)
	if resp.TransferAt != nil {
		fmt.Printf("Transfer at %s\n", *resp.TransferAt)
	}
	if len(resp.Stops) > 0 {
		displayStops(resp.Stops)
	}
	if len(resp.Capacity) > 0 {
		displayCapacity(resp.Capacity)
	}
	return nil
}
