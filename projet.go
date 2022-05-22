package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	osc "github.com/outscale/osc-sdk-go/v2"
	cli "github.com/teris-io/cli"
)

func main() {

	app := cli.New("OSC-LOG").
		WithAction(func(args []string, options map[string]string) int {
			config := osc.NewConfiguration()
			config.Debug = false

			client := osc.NewAPIClient(config)

			ctx := context.WithValue(context.Background(), osc.ContextAWSv4, osc.AWSv4{
				AccessKey: os.Getenv("OSC_ACCESS_KEY"),
				SecretKey: os.Getenv("OSC_SECRET_KEY"),
			})
			req := osc.ReadApiLogsRequest{}
			resp, httpRes, err := client.ApiLogApi.ReadApiLogs(ctx).ReadApiLogsRequest(req).Execute()
			var errCode int
			errCode = 0
			if err != nil {
				fmt.Printf("Error %v", err)
				if httpRes != nil {
					fmt.Fprintln(os.Stderr, httpRes.Status)
				}
				errCode = 1
			}

			for _, log := range resp.GetLogs() {
				jsonLog, marshalError := json.Marshal(log)
				fmt.Println(string(jsonLog))
				if marshalError != nil {
					fmt.Printf("Marshal Error: %v", marshalError)
					errCode = 1
				}

			}

			return errCode
		})
	ret := app.Run(os.Args, os.Stdout)
	os.Exit(ret)

}
