package main

import (
	"context"
	"encoding/json"
	"fmt"
	osc "github.com/outscale/osc-sdk-go/v2"
	cli "github.com/teris-io/cli"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func displayLogs(args []string, options map[string]string) int {
	config := osc.NewConfiguration()
	config.Debug = false
	client := osc.NewAPIClient(config)
	ctx := context.WithValue(context.Background(), osc.ContextAWSv4, osc.AWSv4{
		AccessKey: os.Getenv("OSC_ACCESS_KEY"),
		SecretKey: os.Getenv("OSC_SECRET_KEY"),
	})
	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT)
	logDate := time.Now().UTC().Format("2006-01-02T15:04:05")
	duration := time.Duration(2) * time.Second
	tk := time.NewTicker(duration)
	var file *os.File
	lineBreak := []byte("\n")
	var err error
	if options["write"] != "" {
		file, err = os.OpenFile(options["write"], os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: can not open the file")
			os.Exit(1)
		}
		defer file.Close()
	}
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
			if file == nil {
				fmt.Println(string(jsonLog))
			} else {
				logWriting := []byte(string(jsonLog))
				file.Write(logWriting)
				file.Write(lineBreak)
			}
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
func AddWriteOption() cli.Option {
	return cli.NewOption("write", "Write all traces inside a file instead of writing to standard output").WithChar('w').WithType(cli.TypeString)
}
func main() {
	app := cli.New("osc-logs").
		WithAction(displayLogs).
		WithOption(AddWriteOption())
	ret := app.Run(os.Args, os.Stdout)
	os.Exit(ret)
}
