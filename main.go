package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	osc "github.com/outscale/osc-sdk-go/v2"
	cli "github.com/teris-io/cli"
)

func displayLogs(args []string, options map[string]string) int {
	config := osc.NewConfiguration()
	config.Debug = false

	client := osc.NewAPIClient(config)

	ctx := context.WithValue(context.Background(), osc.ContextAWSv4, osc.AWSv4{
		AccessKey: os.Getenv("OSC_ACCESS_KEY"),
		SecretKey: os.Getenv("OSC_SECRET_KEY"),
	})

	logDate := time.Now().UTC().Format("2006-01-02T15:04:05")
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT)
	duration := time.Duration(2) * time.Second
	tk := time.NewTicker(duration)
	for range tk.C {
		req := osc.ReadApiLogsRequest{
			Filters: &osc.FiltersApiLog{
				QueryDateAfter: &logDate,
			},
		}

		resp, httpRes, err := client.ApiLogApi.ReadApiLogs(ctx).ReadApiLogsRequest(req).Execute()
		if err != nil {
			fmt.Printf("Error %v", err)
			if httpRes != nil {
				fmt.Fprintln(os.Stderr, httpRes.Status)
			}
			return 1
		}
		logs := resp.GetLogs()
		if !resp.HasLogs() || len(logs) == 0 {
			continue
		}

		for _, log := range logs {
			jsonLog, marshalError := json.Marshal(log)
			if marshalError != nil {
				fmt.Fprintf(os.Stderr, "Error: can not read log output")
				return 1
			}
			fmt.Println(string(jsonLog))

		}

		lastLog := logs[len(logs)-1]
		logDate = *lastLog.QueryDate

		go func() {
			<-stopSignal
			os.Exit(0)
		}()

	}
	return 0

}

func main() {
	app := cli.New("osc-logs").
		WithAction(displayLogs)
	ret := app.Run(os.Args, os.Stdout)
	os.Exit(ret)
}
